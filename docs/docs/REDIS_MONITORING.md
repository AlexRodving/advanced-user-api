# 🔴 Redis & Monitoring Guide

## Часть 1: Redis - Зачем и как использовать

### 🎯 Что такое Redis?

**Redis** (Remote Dictionary Server) - in-memory хранилище данных типа ключ-значение.

**Особенность:** Данные хранятся в **RAM** (оперативной памяти), а не на диске.

**Скорость:**
- PostgreSQL (диск): ~5-10ms запрос
- Redis (RAM): ~0.1-1ms запрос ⚡ **в 10-100 раз быстрее!**

---

### 💡 Зачем нужен Redis?

#### 1. Кеширование (Cache)

**Проблема без Redis:**
```go
// Каждый запрос идёт в PostgreSQL
func GetPopularPosts() {
    posts := db.Find(&Post{}).Order("view_count DESC").Limit(10)
    // Запрос к БД: ~10ms
}

// Если 1000 пользователей в секунду:
// 1000 × 10ms = 10,000ms = 10 секунд нагрузки на БД!
```

**С Redis:**
```go
func GetPopularPosts() {
    // Проверяем кеш
    cached := redis.Get("popular_posts")
    if cached != nil {
        return cached  // Из кеша: ~0.1ms ⚡
    }
    
    // Если нет в кеше - берём из БД
    posts := db.Find(&Post{}).Order("view_count DESC").Limit(10)
    
    // Кешируем на 5 минут
    redis.Set("popular_posts", posts, 5*time.Minute)
    
    return posts
}

// Первый запрос: 10ms (БД)
// Остальные 999: 0.1ms (Redis) ⚡
// Нагрузка на БД: 1 запрос вместо 1000!
```

---

#### 2. Session Storage (Сессии)

**Проблема с JWT:**
- Нельзя отозвать (revoke) токен до истечения
- Нужно ждать 24 часа пока токен истечёт

**С Redis:**
```go
// При логине - сохраняем session
sessionID := generateUUID()
redis.Set("session:"+sessionID, userID, 24*time.Hour)

// При каждом запросе
userID := redis.Get("session:" + sessionID)

// Logout - мгновенная инвалидация
redis.Del("session:" + sessionID)  // Сессия удалена ⚡
```

---

#### 3. Rate Limiting (Ограничение запросов)

**Защита от DDoS и спама:**

```go
// Middleware для rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        
        // Счётчик запросов от IP
        key := "rate_limit:" + ip
        count := redis.Incr(key)
        
        // Первый запрос - устанавливаем TTL 1 минута
        if count == 1 {
            redis.Expire(key, time.Minute)
        }
        
        // Максимум 100 запросов в минуту
        if count > 100 {
            c.AbortWithStatusJSON(429, gin.H{
                "error": "too many requests",
            })
            return
        }
        
        c.Next()
    }
}
```

---

#### 4. Counters (Счётчики)

**View counter для постов:**

```go
// Каждый просмотр поста
func ViewPost(postID uint) {
    // Увеличиваем счётчик в Redis
    redis.Incr("post:" + strconv.Itoa(postID) + ":views")
    
    // Раз в минуту синхронизируем с БД (background worker)
    // Вместо UPDATE при каждом просмотре!
}

// Background worker (запускается раз в минуту)
func SyncViewCounters() {
    keys := redis.Keys("post:*:views")
    
    for _, key := range keys {
        postID := extractPostID(key)
        count := redis.Get(key)
        
        // Обновляем БД
        db.Model(&Post{}).Where("id = ?", postID).
            UpdateColumn("view_count", gorm.Expr("view_count + ?", count))
        
        // Очищаем Redis
        redis.Del(key)
    }
}
```

**Преимущества:**
- 1000 просмотров → 1000 операций в Redis (быстро)
- 1 UPDATE в PostgreSQL раз в минуту (вместо 1000!)

---

#### 5. Pub/Sub (Сообщения в реальном времени)

**Real-time уведомления:**

```go
// Подписка на уведомления
func SubscribeToNotifications(userID uint) {
    pubsub := redis.Subscribe("notifications:" + userID)
    
    for msg := range pubsub.Channel() {
        fmt.Println("Новое уведомление:", msg.Payload)
        // Отправить в WebSocket клиенту
    }
}

// Отправка уведомления
func NotifyUser(userID uint, message string) {
    redis.Publish("notifications:"+userID, message)
}
```

---

### 💻 Как добавить Redis в проект

#### Шаг 1: Раскомментировать в docker-compose.yml

```yaml
# Было закомментировано:
# redis:
#   image: redis:7-alpine

# Раскомментируйте:
redis:
  image: redis:7-alpine
  container_name: advanced-api-redis
  ports:
    - "6379:6379"
  networks:
    - app-network
```

---

#### Шаг 2: Установить Redis client

```bash
go get github.com/redis/go-redis/v9
```

---

#### Шаг 3: Создать Redis connection

**Создайте:** `internal/repository/redis.go`

```go
package repository

import (
    "context"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// InitRedis - подключение к Redis
func InitRedis(host, port string) *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     host + ":" + port,  // "redis:6379"
        Password: "",                  // Без пароля
        DB:       0,                   // База данных 0
    })
    
    // Проверка подключения
    _, err := client.Ping(ctx).Result()
    if err != nil {
        panic("Failed to connect to Redis: " + err.Error())
    }
    
    return client
}
```

---

#### Шаг 4: Использовать в Service

**Обновите:** `internal/service/post_service.go`

```go
type postService struct {
    postRepo    repository.PostRepository
    redisClient *redis.Client  // НОВОЕ
}

func NewPostService(postRepo repository.PostRepository, redisClient *redis.Client) PostService {
    return &postService{
        postRepo:    postRepo,
        redisClient: redisClient,
    }
}

// GetPopularPosts с кешированием
func (s *postService) GetPopularPosts() ([]domain.Post, error) {
    // 1. Проверяем кеш
    cached, err := s.redisClient.Get(ctx, "popular_posts").Result()
    if err == nil {
        var posts []domain.Post
        json.Unmarshal([]byte(cached), &posts)
        return posts, nil
    }
    
    // 2. Если нет в кеше - берём из БД
    posts, err := s.postRepo.FindPopular(10)
    if err != nil {
        return nil, err
    }
    
    // 3. Кешируем на 5 минут
    data, _ := json.Marshal(posts)
    s.redisClient.Set(ctx, "popular_posts", data, 5*time.Minute)
    
    return posts, nil
}
```

---

#### Шаг 5: Инициализация в main.go

```go
// cmd/api/main.go

func main() {
    cfg := config.Load()
    
    // PostgreSQL
    db := repository.InitDB(cfg)
    
    // Redis
    redisClient := repository.InitRedis(cfg.RedisHost, cfg.RedisPort)
    
    // Services с Redis
    postService := service.NewPostService(postRepo, redisClient)
    
    // ...
}
```

---

### 📊 Redis Data Types

#### 1. String (ключ-значение)

```go
// Установить
redis.Set(ctx, "user:1:name", "Alice", 0)

// Получить
name := redis.Get(ctx, "user:1:name").Val()  // "Alice"

// С TTL (время жизни)
redis.Set(ctx, "session:abc", "user_id:123", 24*time.Hour)
```

---

#### 2. Hash (хеш-таблица)

```go
// Сохранить пользователя
redis.HSet(ctx, "user:1", map[string]interface{}{
    "name":  "Alice",
    "email": "alice@example.com",
    "role":  "admin",
})

// Получить поле
name := redis.HGet(ctx, "user:1", "name").Val()

// Получить всё
user := redis.HGetAll(ctx, "user:1").Val()
// map[string]string{"name": "Alice", "email": "alice@example.com", ...}
```

---

#### 3. List (список)

```go
// Добавить в начало
redis.LPush(ctx, "notifications:user:1", "New comment on your post")

// Получить первые 10
notifications := redis.LRange(ctx, "notifications:user:1", 0, 9).Val()
```

---

#### 4. Set (множество)

```go
// Добавить в множество
redis.SAdd(ctx, "post:1:likers", "user:10", "user:20", "user:30")

// Проверить есть ли в множестве
exists := redis.SIsMember(ctx, "post:1:likers", "user:10").Val()  // true

// Количество элементов
count := redis.SCard(ctx, "post:1:likers").Val()  // 3
```

---

#### 5. Sorted Set (отсортированное множество)

```go
// Leaderboard (топ пользователей)
redis.ZAdd(ctx, "leaderboard", redis.Z{Score: 1500, Member: "user:1"})
redis.ZAdd(ctx, "leaderboard", redis.Z{Score: 2000, Member: "user:2"})

// Топ 10
top10 := redis.ZRevRange(ctx, "leaderboard", 0, 9).Val()
```

---

## Часть 2: Мониторинг приложения

### 🔍 Что мониторить?

1. **Здоровье сервиса** (работает ли?)
2. **Производительность** (скорость ответов)
3. **Ошибки** (что ломается?)
4. **Ресурсы** (CPU, RAM, диск)
5. **Бизнес метрики** (регистрации, посты, и т.д.)

---

### 📊 Уровни мониторинга

```
┌─────────────────────────────────────┐
│  1. Health Checks (жив ли сервис?)  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  2. Logs (что происходит?)          │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  3. Metrics (сколько, как быстро?)  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  4. Traces (где тормозит?)          │
└─────────────────────────────────────┘
```

---

## 1️⃣ Health Checks (уже есть!)

### Endpoint /health

**Файл:** `internal/handler/routes.go`

```go
router.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "service": "advanced-user-api",
        "status":  "ok",
    })
})
```

### Расширенный health check

```go
router.GET("/health", func(c *gin.Context) {
    // Проверка БД
    dbHealth := "ok"
    if err := db.DB().Ping(); err != nil {
        dbHealth = "error"
    }
    
    // Проверка Redis
    redisHealth := "ok"
    if _, err := redisClient.Ping(ctx).Result(); err != nil {
        redisHealth = "error"
    }
    
    c.JSON(200, gin.H{
        "service":  "advanced-user-api",
        "status":   "ok",
        "database": dbHealth,
        "redis":    redisHealth,
        "uptime":   time.Since(startTime).String(),
    })
})
```

**Мониторинг:**
```bash
# Проверка каждые 30 секунд
watch -n 30 curl http://localhost:8080/health
```

---

## 2️⃣ Structured Logging

### Zap Logger (уже используем!)

**Файл:** `internal/middleware/logger.go`

**Что логируем:**
```go
logger.Info("HTTP Request",
    zap.Int("status", 200),
    zap.String("method", "GET"),
    zap.String("path", "/api/v1/users"),
    zap.Duration("latency", 15*time.Millisecond),
    zap.String("ip", "192.168.1.1"),
)
```

**Вывод (JSON):**
```json
{
  "level": "info",
  "ts": 1634567890,
  "msg": "HTTP Request",
  "status": 200,
  "method": "GET",
  "path": "/api/v1/users",
  "latency": "15ms",
  "ip": "192.168.1.1"
}
```

### Логирование в файл

```go
// cmd/api/main.go

func setupLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    
    // Логи в файл
    config.OutputPaths = []string{
        "stdout",              // Консоль
        "/var/log/api/app.log", // Файл
    }
    
    logger, _ := config.Build()
    return logger
}
```

### Log Rotation (ротация логов)

```bash
# /etc/logrotate.d/advanced-api
/var/log/api/*.log {
    daily           # Ротация каждый день
    rotate 7        # Хранить 7 файлов
    compress        # Сжимать старые
    delaycompress   # Сжимать со следующей ротации
    notifempty      # Не ротировать пустые
    create 0640 appuser appgroup
}
```

---

## 3️⃣ Metrics - Prometheus

### Установка Prometheus client

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

---

### Создание метрик

**Создайте:** `internal/pkg/metrics/metrics.go`

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP запросы
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Общее количество HTTP запросов",
        },
        []string{"method", "endpoint", "status"},
    )
    
    // Длительность запросов
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Длительность HTTP запросов",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    // Активные пользователи
    ActiveUsers = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users_count",
            Help: "Количество активных пользователей",
        },
    )
    
    // Размер БД
    DatabaseSize = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "database_size_bytes",
            Help: "Размер базы данных в байтах",
        },
    )
)
```

---

### Middleware для метрик

```go
// internal/middleware/metrics.go

func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        // Записываем метрики
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        
        metrics.HTTPRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
            status,
        ).Inc()
        
        metrics.HTTPRequestDuration.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
        ).Observe(duration)
    }
}
```

---

### Endpoint для метрик

```go
// cmd/api/main.go

import "github.com/prometheus/client_golang/prometheus/promhttp"

func main() {
    // ...
    
    // Prometheus metrics endpoint
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))
    
    // ...
}
```

**Проверка:**
```bash
curl http://localhost:8080/metrics
```

**Вывод:**
```
# HELP http_requests_total Общее количество HTTP запросов
# TYPE http_requests_total counter
http_requests_total{endpoint="/api/v1/users",method="GET",status="200"} 42
http_requests_total{endpoint="/api/v1/auth/login",method="POST",status="200"} 15

# HELP http_request_duration_seconds Длительность HTTP запросов
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{endpoint="/api/v1/users",method="GET",le="0.005"} 30
http_request_duration_seconds_bucket{endpoint="/api/v1/users",method="GET",le="0.01"} 40
```

---

### Docker Compose с Prometheus

```yaml
# docker-compose.monitoring.yml

services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    networks:
      - app-network

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
    networks:
      - app-network
    depends_on:
      - prometheus

volumes:
  prometheus_data:
  grafana_data:
```

---

### Конфигурация Prometheus

**Создайте:** `prometheus.yml`

```yaml
global:
  scrape_interval: 15s  # Собирать метрики каждые 15 секунд

scrape_configs:
  - job_name: 'advanced-user-api'
    static_configs:
      - targets: ['api:8080']  # Наш API
    metrics_path: '/metrics'
```

---

### Запуск мониторинга

```bash
# Запустить API + PostgreSQL + Prometheus + Grafana
docker compose -f docker-compose.yml -f docker-compose.monitoring.yml up -d

# Откройте в браузере:
# Prometheus: http://localhost:9090
# Grafana:    http://localhost:3000 (admin/admin)
```

---

### Grafana Dashboard

**В Grafana:**

1. Add Data Source → Prometheus → URL: `http://prometheus:9090`
2. Create Dashboard
3. Add Panel → Metrics:

```promql
# Запросы в секунду
rate(http_requests_total[1m])

# Средняя длительность
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# Количество ошибок (5xx)
sum(rate(http_requests_total{status=~"5.."}[5m]))

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

**Визуализация:**
- Графики
- Счётчики
- Таблицы
- Алерты

---

## 4️⃣ APM (Application Performance Monitoring)

### Sentry - Error Tracking

**Установка:**
```bash
go get github.com/getsentry/sentry-go
```

**Инициализация:**
```go
import "github.com/getsentry/sentry-go"

func main() {
    sentry.Init(sentry.ClientOptions{
        Dsn: "https://your-dsn@sentry.io/project-id",
        Environment: "production",
    })
    defer sentry.Flush(2 * time.Second)
    
    // ...
}
```

**Отлов ошибок:**
```go
func (h *Handler) CreateUser(c *gin.Context) {
    user, err := h.service.CreateUser(req)
    if err != nil {
        // Отправляем ошибку в Sentry
        sentry.CaptureException(err)
        
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
}
```

---

## 5️⃣ Uptime Monitoring

### Внешний мониторинг

**Сервисы:**
- **UptimeRobot** (бесплатно 50 мониторов)
- **Pingdom**
- **StatusCake**
- **Better Uptime**

**Настройка UptimeRobot:**
1. Создать монитор HTTP(s)
2. URL: `https://your-api.com/health`
3. Интервал: 5 минут
4. Email alert при downtime

---

## 6️⃣ Dashboard пример

### Простой dashboard с Grafana

**Панели:**

1. **Requests per second (RPS)**
   ```promql
   sum(rate(http_requests_total[1m]))
   ```

2. **Error rate**
   ```promql
   sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
   ```

3. **Average response time**
   ```promql
   rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])
   ```

4. **Database connections**
   ```promql
   go_sql_open_connections
   ```

5. **Memory usage**
   ```promql
   go_memstats_alloc_bytes
   ```

6. **Goroutines**
   ```promql
   go_goroutines
   ```

---

## 🚨 Alerting (Алерты)

### Prometheus Alertmanager

**Создайте:** `alerts.yml`

```yaml
groups:
  - name: api_alerts
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }}"
      
      # High latency
      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
      
      # Service down
      - alert: ServiceDown
        expr: up{job="advanced-user-api"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service is down!"
```

---

## 📈 Метрики которые стоит собирать

### Application Metrics

```go
// Регистрации
metrics.RegistrationsTotal.Inc()

// Логины
metrics.LoginsTotal.Inc()

// Созданные посты
metrics.PostsCreated.Inc()

// Активные пользователи (обновляется каждую минуту)
func UpdateActiveUsers() {
    count := db.Model(&User{}).Where("last_seen > ?", time.Now().Add(-5*time.Minute)).Count()
    metrics.ActiveUsers.Set(float64(count))
}
```

### Database Metrics

```go
sqlDB, _ := db.DB()
stats := sqlDB.Stats()

metrics.DBOpenConnections.Set(float64(stats.OpenConnections))
metrics.DBInUse.Set(float64(stats.InUse))
metrics.DBIdle.Set(float64(stats.Idle))
metrics.DBWaitCount.Set(float64(stats.WaitCount))
```

### System Metrics (автоматически от Prometheus)

- CPU usage
- Memory usage
- Goroutines count
- GC stats

---

## 🔔 Notification Channels

### Slack

```yaml
# alertmanager.yml
receivers:
  - name: 'slack'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/YOUR/WEBHOOK/URL'
        channel: '#alerts'
        text: '{{ .CommonAnnotations.summary }}'
```

### Telegram

```yaml
receivers:
  - name: 'telegram'
    telegram_configs:
      - bot_token: 'YOUR_BOT_TOKEN'
        chat_id: YOUR_CHAT_ID
```

### Email

```yaml
receivers:
  - name: 'email'
    email_configs:
      - to: 'admin@example.com'
        from: 'alerts@example.com'
        smarthost: 'smtp.gmail.com:587'
        auth_username: 'your-email@gmail.com'
        auth_password: 'app-password'
```

---

## 🎯 Полный стек мониторинга

```yaml
# docker-compose.monitoring.yml

services:
  # Prometheus - сбор метрик
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
  
  # Grafana - визуализация
  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
  
  # AlertManager - алерты
  alertmanager:
    image: prom/alertmanager
    ports:
      - "9093:9093"
    volumes:
      - ./alertmanager.yml:/etc/alertmanager/alertmanager.yml
  
  # Loki - логи
  loki:
    image: grafana/loki
    ports:
      - "3100:3100"
  
  # Promtail - сбор логов
  promtail:
    image: grafana/promtail
    volumes:
      - /var/log:/var/log
      - ./promtail.yml:/etc/promtail/promtail.yml
```

---

## 📊 Что получаем в итоге

### Grafana Dashboard покажет:

```
┌────────────────────────────────────────────┐
│  Advanced User API - Dashboard             │
├────────────────────────────────────────────┤
│                                            │
│  📈 Requests/sec:  150 req/s               │
│  ⏱️  Avg latency:   25ms                    │
│  ❌ Error rate:     0.1%                    │
│  👥 Active users:   1,234                   │
│                                            │
│  ┌──────────────────────────────────────┐ │
│  │     RPS Graph (last 24h)             │ │
│  │                    /\                 │ │
│  │                   /  \      /\        │ │
│  │         /\       /    \    /  \       │ │
│  │        /  \     /      \  /    \      │ │
│  └──────────────────────────────────────┘ │
│                                            │
│  ┌──────────────────────────────────────┐ │
│  │     Top Endpoints                    │ │
│  │  1. GET /api/v1/posts     - 45%      │ │
│  │  2. POST /api/v1/auth/... - 30%      │ │
│  │  3. GET /api/v1/users     - 15%      │ │
│  └──────────────────────────────────────┘ │
│                                            │
│  🚨 Alerts: None ✅                        │
└────────────────────────────────────────────┘
```

---

## 🎓 Резюме

### Redis используется для:
✅ Кеширование (популярные посты, пользователи)  
✅ Sessions (отзыв токенов)  
✅ Rate limiting (защита от DDoS)  
✅ Counters (просмотры, лайки)  
✅ Pub/Sub (real-time уведомления)  
✅ Leaderboards (топы)  

### Мониторинг состоит из:
✅ Health checks (/health endpoint)  
✅ Structured logging (Zap)  
✅ Metrics (Prometheus)  
✅ Visualization (Grafana)  
✅ Error tracking (Sentry)  
✅ Uptime monitoring (UptimeRobot)  
✅ Alerting (Slack, Email, Telegram)  

---

## 📖 См. также

- [Docker Guide](./DOCKER.md) - Docker Compose для мониторинга
- [Architecture](./ARCHITECTURE.md) - Где вписать Redis
- [Blog Extension](./EXTEND_TO_BLOG.md) - Примеры использования Redis

---

## 📚 Дополнительные ресурсы

### Redis:
- [Redis Documentation](https://redis.io/docs/)
- [go-redis Client](https://redis.uptrace.dev/)
- [Redis Best Practices](https://redis.io/docs/manual/patterns/)

### Monitoring:
- [Prometheus](https://prometheus.io/docs/)
- [Grafana](https://grafana.com/docs/)
- [Sentry](https://docs.sentry.io/)
- [OpenTelemetry](https://opentelemetry.io/) - современный стандарт

---

**Redis + Мониторинг = Production-ready приложение!** 🚀

