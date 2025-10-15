# 🐳 Docker & Docker Compose Guide

## Обзор

Полное руководство по контейнеризации проекта с Docker и оркестрации с Docker Compose.

**Файлы в проекте:**
- `docker/Dockerfile` - образ приложения
- `docker-compose.yml` - оркестрация всех сервисов
- `.dockerignore` - исключения при сборке

---

## 🎯 Зачем Docker?

### Проблемы без Docker

```
Developer A:
- Go 1.23
- PostgreSQL 14
- Redis 6
- ✅ Работает

Developer B:
- Go 1.21
- PostgreSQL 15
- Нет Redis
- ❌ Не работает!
```

### С Docker

```
Все разработчики:
- Docker 20+
- ✅ Работает одинаково!

Production:
- Docker 20+
- ✅ Работает так же!
```

**Преимущества:**
- ✅ Одинаковое окружение для всех
- ✅ Легко деплоить
- ✅ Изоляция приложений
- ✅ Версионирование окружения

---

## 📄 Dockerfile - Образ приложения

**Файл:** `docker/Dockerfile`

### Multi-Stage Build

Наш Dockerfile использует **multi-stage build** - два этапа сборки для оптимизации размера образа.

```dockerfile
# ================================================================
# STAGE 1: BUILDER - Компиляция приложения
# ================================================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Установка системных зависимостей
RUN apk add --no-cache git ca-certificates

# Копирование go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование всего кода
COPY . .

# Компиляция
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/api cmd/api/main.go

# ================================================================
# STAGE 2: RUNTIME - Минимальный образ для запуска
# ================================================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Создание непривилегированного пользователя
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Копирование ТОЛЬКО бинарника из builder
COPY --from=builder /app/api .

RUN chown -R appuser:appgroup /home/appuser

USER appuser

EXPOSE 8080

CMD ["./api"]
```

---

## 🔍 Разбор Dockerfile по секциям

### Stage 1: Builder (компиляция)

#### 1. Base Image

```dockerfile
FROM golang:1.25-alpine AS builder
```

**Что происходит:**
- `FROM` - базовый образ
- `golang:1.25-alpine` - официальный Go образ на Alpine Linux
- `AS builder` - имя этапа (для ссылки из stage 2)

**Alpine Linux:**
- Минимальный Linux (~5MB vs ~500MB Ubuntu)
- Быстрая установка пакетов через `apk`

**Альтернативы:**
```dockerfile
FROM golang:1.25         # Debian-based (~800MB)
FROM golang:1.25-alpine  # Alpine-based (~300MB) ✅
```

---

#### 2. Working Directory

```dockerfile
WORKDIR /app
```

**Что делает:**
- Создаёт директорию `/app` в контейнере
- Устанавливает её как текущую
- Все последующие команды выполняются в `/app`

---

#### 3. System Dependencies

```dockerfile
RUN apk add --no-cache git ca-certificates
```

**Параметры:**
- `apk` - пакетный менеджер Alpine Linux
- `add` - установить пакеты
- `--no-cache` - не сохранять кеш (экономит место)

**Пакеты:**
- `git` - нужен для `go mod download` (некоторые модули используют git)
- `ca-certificates` - SSL сертификаты для HTTPS запросов

---

#### 4. Dependencies Caching

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
```

**Почему в таком порядке?**

Docker кеширует каждый слой (RUN, COPY, etc.):
1. Если `go.mod` не изменился → используется кеш → **быстро!**
2. Если `go.mod` изменился → перескачиваем зависимости

**Без кеширования:**
```dockerfile
COPY . .           # Копируем весь код
RUN go mod download  # Скачиваем ВСЕГДА (медленно!)
```

**С кешированием (наш вариант):**
```dockerfile
COPY go.mod go.sum ./  # Копируем только go.mod
RUN go mod download    # Скачиваем ТОЛЬКО если go.mod изменился ✅
COPY . .               # Затем копируем код
```

**Результат:**
- Первая сборка: ~2 минуты
- Повторная сборка (без изменения зависимостей): ~10 секунд!

---

#### 5. Build Application

```dockerfile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/api cmd/api/main.go
```

**Флаги компиляции:**

| Флаг | Значение | Зачем |
|------|----------|-------|
| `CGO_ENABLED=0` | Отключить CGO | Статическая сборка (не зависит от libc) |
| `GOOS=linux` | Целевая ОС | Linux (для Alpine) |
| `-a` | Пересобрать все пакеты | Полная сборка |
| `-installsuffix cgo` | Суффикс для разделения | Изоляция от CGO версий |
| `-o /app/api` | Выходной файл | Где сохранить бинарник |

**Результат:** Один исполняемый файл `api` (~15-20MB)

---

### Stage 2: Runtime (выполнение)

#### 1. Minimal Base Image

```dockerfile
FROM alpine:latest
```

**Почему Alpine?**
- Размер: ~5MB (vs 130MB Ubuntu)
- Безопасность: минимум пакетов = меньше уязвимостей
- Скорость: быстрое скачивание и запуск

---

#### 2. CA Certificates

```dockerfile
RUN apk --no-cache add ca-certificates
```

**Зачем:**
- Для HTTPS запросов к внешним API
- Без них: `x509: certificate signed by unknown authority`

---

#### 3. Security: Non-root User

```dockerfile
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
```

**Зачем:**
- ❌ Запуск от root = небезопасно
- ✅ Запуск от appuser = безопасно

**Параметры:**
- `-S` - system user/group
- `-G appgroup` - добавить в группу

**Потом переключаемся:**
```dockerfile
USER appuser
```

Все команды после этого выполняются от `appuser`, не от `root`!

---

#### 4. Copy Binary from Builder

```dockerfile
COPY --from=builder /app/api .
```

**Магия Multi-Stage:**
- `--from=builder` - копируем из **Stage 1** (builder)
- Копируем **ТОЛЬКО** бинарник `api`
- **НЕ** копируем исходный код, зависимости, и т.д.

**Результат:**
- Builder stage: ~500MB (не входит в финальный образ!)
- Runtime stage: ~15-20MB ✅

**Сравнение:**

| Подход | Размер образа |
|--------|---------------|
| Без multi-stage | ~500MB |
| С multi-stage ✅ | ~15-20MB |

**Выигрыш:** 25x меньше!

---

#### 5. Permissions

```dockerfile
RUN chown -R appuser:appgroup /home/appuser
```

**Даём права пользователю `appuser` на все файлы**

---

#### 6. Expose Port

```dockerfile
EXPOSE 8080
```

**Документирует** что приложение слушает порт 8080

**Важно:** Это НЕ публикует порт! Нужен `-p` флаг при `docker run`

---

#### 7. Start Command

```dockerfile
CMD ["./api"]
```

**Запускает** бинарник при старте контейнера

**Альтернативы:**
```dockerfile
CMD ["./api"]              # ✅ exec form (PID 1, корректный SIGTERM)
CMD ./api                  # ❌ shell form (sh -c, проблемы с сигналами)
ENTRYPOINT ["./api"]       # Для фиксированной команды
```

---

## 📦 .dockerignore

**Файл:** `.dockerignore`

**Назначение:** Исключает файлы из контекста сборки (ускоряет сборку)

```dockerignore
# Git
.git
.gitignore

# Environment
.env
env.example

# IDE
.vscode
.idea

# Documentation
*.md

# Tests
*_test.go
tests/

# Build artifacts
bin/
*.out

# Logs
*.log
```

**Почему важно:**
- Меньше файлов → быстрее `docker build`
- Не копируем лишнее в образ
- Безопасность (не копируем `.env`)

---

## 🎼 Docker Compose - Оркестрация

**Файл:** `docker-compose.yml`

### Что это?

Docker Compose позволяет управлять **несколькими контейнерами** одной командой.

**Наши сервисы:**
1. `postgres` - PostgreSQL база данных
2. `redis` - Redis кеш
3. `api` - наше приложение
4. `pgadmin` - веб-интерфейс для PostgreSQL (опционально)

---

### Полный разбор docker-compose.yml

```yaml
version: '3.8'

# Определение сервисов
services:
  
  # ================================================================
  # POSTGRES - База данных
  # ================================================================
  postgres:
    image: postgres:15-alpine           # Образ из Docker Hub
    container_name: advanced-api-postgres  # Имя контейнера
    environment:                        # Переменные окружения
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: advanced_api
    ports:
      - "5432:5432"                     # Порт хоста:порт контейнера
    volumes:
      - postgres_data:/var/lib/postgresql/data  # Persistent storage
    networks:
      - app-network                     # Внутренняя сеть
    restart: unless-stopped             # Политика перезапуска
    healthcheck:                        # Проверка здоровья
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
```

---

### Разбор каждой секции

#### 1. Image

```yaml
image: postgres:15-alpine
```

**Формат:** `<repository>:<tag>`

- `postgres` - официальный образ PostgreSQL
- `15` - версия PostgreSQL
- `alpine` - на базе Alpine Linux (меньше размер)

**Где берётся:** Docker Hub (https://hub.docker.com/_/postgres)

---

#### 2. Container Name

```yaml
container_name: advanced-api-postgres
```

**Зачем:**
- Удобное имя вместо случайного
- Используется в логах и команд

ах
- Можно ссылаться в других сервисах

**Без container_name:** `advanced-user-api_postgres_1` (автогенерация)
**С container_name:** `advanced-api-postgres` ✅

---

#### 3. Environment Variables

```yaml
environment:
  POSTGRES_USER: postgres
  POSTGRES_PASSWORD: postgres
  POSTGRES_DB: advanced_api
```

**Настройки PostgreSQL:**
- `POSTGRES_USER` - имя пользователя БД
- `POSTGRES_PASSWORD` - пароль (в production используйте secrets!)
- `POSTGRES_DB` - имя базы данных (создаётся автоматически)

**В production:**
```yaml
environment:
  POSTGRES_PASSWORD: ${DB_PASSWORD}  # Из .env файла
```

---

#### 4. Ports

```yaml
ports:
  - "5432:5432"
```

**Формат:** `"HOST_PORT:CONTAINER_PORT"`

- `5432` (слева) - порт на **вашем компьютере**
- `5432` (справа) - порт **внутри контейнера**

**Примеры:**
```yaml
ports:
  - "5432:5432"   # Стандартный PostgreSQL
  - "5433:5432"   # Если 5432 уже занят на хосте
  - "8080:8080"   # API приложение
```

**Как работает:**
- Приложение на хосте → `localhost:5432` → контейнер PostgreSQL
- Другие контейнеры → `postgres:5432` (через network)

---

#### 5. Volumes

```yaml
volumes:
  - postgres_data:/var/lib/postgresql/data
```

**Формат:** `VOLUME_NAME:CONTAINER_PATH`

**Зачем:**
- Сохранение данных при перезапуске контейнера
- Данные **не удаляются** при `docker-compose down`

**Без volume:**
```bash
docker-compose down
# ❌ Все данные PostgreSQL потеряны!
```

**С volume:**
```bash
docker-compose down
docker-compose up
# ✅ Данные на месте!
```

**Определение volumes:**
```yaml
volumes:
  postgres_data:  # Именованный volume (управляется Docker)
  redis_data:
```

---

#### 6. Networks

```yaml
networks:
  - app-network
```

**Зачем:**
- Изоляция сервисов
- Контейнеры в одной сети могут общаться между собой

**Как работает:**
```
api контейнер → postgres:5432 → postgres контейнер ✅
api контейнер → localhost:5432 → ❌ (другая сеть)
```

**Внутри сети:**
- Используем имена сервисов как hostnames
- `postgres` вместо `localhost`

**В коде приложения:**
```bash
# .env
DB_HOST=postgres   # ✅ Имя сервиса
DB_PORT=5432       # Внутренний порт контейнера
```

**Определение network:**
```yaml
networks:
  app-network:
    driver: bridge
```

---

#### 7. Restart Policy

```yaml
restart: unless-stopped
```

**Политики:**
- `no` - никогда не перезапускать
- `always` - всегда перезапускать
- `on-failure` - только при ошибке
- `unless-stopped` - всегда, кроме явной остановки ✅

**Поведение:**
```bash
# Контейнер упал
# → Docker автоматически перезапускает

# docker-compose stop
# → Docker НЕ перезапускает (явная остановка)

# Перезагрузка сервера
# → Docker автоматически запускает контейнеры
```

---

#### 8. Health Check

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres"]
  interval: 10s
  timeout: 5s
  retries: 5
```

**Параметры:**
- `test` - команда для проверки (pg_isready для PostgreSQL)
- `interval` - как часто проверять (каждые 10 секунд)
- `timeout` - максимальное время выполнения (5 секунд)
- `retries` - сколько раз повторить при неудаче (5)

**Статусы:**
- `starting` - ещё не проверен
- `healthy` - проверка прошла ✅
- `unhealthy` - проверка не прошла ❌

**Зачем:**
- API не запустится пока PostgreSQL не готов
- `depends_on` с `condition: service_healthy`

---

### API Service

```yaml
api:
  build:
    context: .                    # Где искать Dockerfile
    dockerfile: docker/Dockerfile # Путь к Dockerfile
  container_name: advanced-api-app
  environment:
    DB_HOST: postgres             # ✅ Имя сервиса, не localhost!
    DB_PORT: 5432
    DB_USER: postgres
    DB_PASSWORD: postgres
    DB_NAME: advanced_api
    JWT_SECRET: your-secret-key-change-in-production
    SERVER_PORT: 8080
    GIN_MODE: release
  ports:
    - "8080:8080"
  depends_on:
    postgres:
      condition: service_healthy  # Ждём пока PostgreSQL не станет healthy
    redis:
      condition: service_healthy
  networks:
    - app-network
  restart: unless-stopped
```

#### depends_on с condition

```yaml
depends_on:
  postgres:
    condition: service_healthy
```

**Что делает:**
- API **не запустится** пока PostgreSQL не пройдёт healthcheck
- Избегаем ошибку "connection refused" при старте

**Без condition:**
```yaml
depends_on:
  - postgres  # Запускается одновременно, API может стартовать раньше БД
```

---

### Redis Service

```yaml
redis:
  image: redis:7-alpine
  container_name: advanced-api-redis
  ports:
    - "6379:6379"
  volumes:
    - redis_data:/data
  networks:
    - app-network
  restart: unless-stopped
  healthcheck:
    test: ["CMD", "redis-cli", "ping"]
    interval: 10s
    timeout: 5s
    retries: 5
```

**Redis healthcheck:**
- `redis-cli ping` → `PONG` = healthy ✅

---

### pgAdmin (опционально)

```yaml
pgadmin:
  image: dpage/pgadmin4:latest
  container_name: advanced-api-pgadmin
  environment:
    PGADMIN_DEFAULT_EMAIL: admin@admin.com
    PGADMIN_DEFAULT_PASSWORD: admin
  ports:
    - "5050:80"
  networks:
    - app-network
  depends_on:
    - postgres
```

**Веб-интерфейс для управления PostgreSQL:**
- URL: http://localhost:5050
- Email: admin@admin.com
- Password: admin

**Подключение к БД:**
- Host: `postgres` (имя сервиса)
- Port: `5432`
- Username: `postgres`
- Password: `postgres`

---

## 🚀 Docker Compose команды

### Основные команды

```bash
# Запустить все сервисы
docker compose up

# В фоновом режиме (detached)
docker compose up -d

# Пересобрать образы и запустить
docker compose up -d --build

# Остановить все сервисы
docker compose down

# Остановить И удалить volumes (данные БД!)
docker compose down -v

# Посмотреть статус сервисов
docker compose ps

# Логи всех сервисов
docker compose logs

# Логи конкретного сервиса
docker compose logs api
docker compose logs postgres

# Следить за логами в реальном времени
docker compose logs -f api

# Перезапустить сервис
docker compose restart api

# Остановить конкретный сервис
docker compose stop api

# Запустить один сервис
docker compose up postgres
```

---

### Выполнение команд в контейнере

```bash
# Зайти в контейнер
docker compose exec api sh

# Выполнить команду
docker compose exec api ls -la

# Выполнить в PostgreSQL
docker compose exec postgres psql -U postgres -d advanced_api

# Выполнить SQL
docker compose exec postgres psql -U postgres -d advanced_api -c "SELECT * FROM users;"

# Выполнить в Redis
docker compose exec redis redis-cli ping
```

---

### Управление volumes

```bash
# Список volumes
docker volume ls

# Информация о volume
docker volume inspect advanced-user-api_postgres_data

# Удалить неиспользуемые volumes
docker volume prune

# Бэкап PostgreSQL volume
docker compose exec postgres pg_dump -U postgres advanced_api > backup.sql

# Восстановление
docker compose exec -T postgres psql -U postgres advanced_api < backup.sql
```

---

### Просмотр ресурсов

```bash
# Использование ресурсов (CPU, RAM)
docker stats

# Только для наших контейнеров
docker stats advanced-api-app advanced-api-postgres advanced-api-redis

# Размер образов
docker images | grep advanced

# Размер volumes
docker system df -v
```

---

## 🔧 Troubleshooting

### Проблема: Port already in use

```
Error: bind: address already in use
```

**Решение:**
```bash
# Найти процесс на порту 8080
lsof -i :8080

# Убить процесс
kill -9 <PID>

# Или изменить порт в docker-compose.yml
ports:
  - "8081:8080"  # Теперь доступен на 8081
```

---

### Проблема: Cannot connect to database

```
Error: dial tcp: lookup postgres: no such host
```

**Решение:**
```bash
# Проверьте что services в одной сети
docker compose ps

# Проверьте DB_HOST в переменных
environment:
  DB_HOST: postgres  # ✅ Имя сервиса, не localhost!
```

---

### Проблема: Permission denied

```
Error: permission denied while trying to connect to Docker daemon
```

**Решение:**
```bash
# Добавьте пользователя в группу docker
sudo usermod -aG docker $USER

# Перелогиньтесь
newgrp docker

# Или используйте sudo
sudo docker compose up
```

---

### Проблема: Image build failed

```
Error: failed to solve: failed to compute cache key
```

**Решение:**
```bash
# Очистите кеш
docker builder prune

# Пересоберите без кеша
docker compose build --no-cache

# Или пересоберите с pull новых образов
docker compose build --pull
```

---

## 🎯 Оптимизация Dockerfile

### 1. Layer Caching

```dockerfile
# ❌ Плохо - копируем всё сразу
COPY . .
RUN go mod download  # Всегда перекачивает

# ✅ Хорошо - отдельные слои
COPY go.mod go.sum ./
RUN go mod download    # Кешируется!
COPY . .
```

### 2. .dockerignore

```dockerfile
# Исключите ненужное
.git
tests/
*.md
```

### 3. Multi-stage build

```dockerfile
FROM golang:1.25-alpine AS builder  # Stage 1: сборка
FROM alpine:latest                   # Stage 2: запуск (без Go!)
COPY --from=builder /app/api .       # Только бинарник
```

### 4. Порядок команд

```dockerfile
# От менее изменяемого к более изменяемому:
FROM ...           # Редко меняется
WORKDIR ...        # Редко меняется
RUN apk add ...    # Редко меняется
COPY go.mod ...    # Иногда меняется
COPY . .           # Часто меняется
RUN go build ...   # Часто пересобирается
```

---

## 🔒 Security Best Practices

### 1. Не используйте root

```dockerfile
# ❌ Плохо
USER root
CMD ["./api"]

# ✅ Хорошо
RUN adduser -S appuser
USER appuser
CMD ["./api"]
```

### 2. Не храните секреты в образе

```dockerfile
# ❌ Плохо
ENV JWT_SECRET=hardcoded-secret

# ✅ Хорошо
# Передавайте через environment в runtime
```

```yaml
# docker-compose.yml
environment:
  JWT_SECRET: ${JWT_SECRET}  # Из .env файла
```

### 3. Используйте specific tags

```dockerfile
# ❌ Плохо
FROM golang:latest      # Может внезапно измениться!

# ✅ Хорошо
FROM golang:1.25-alpine  # Предсказуемо
```

### 4. Scan для уязвимостей

```bash
# Docker Scout (встроен в Docker Desktop)
docker scout cves advanced-user-api:latest

# Trivy
trivy image advanced-user-api:latest
```

---

## 📊 Development vs Production

### Development (docker-compose.yml)

```yaml
services:
  api:
    build: .
    volumes:
      - .:/app  # Live reload - изменения сразу видны
    environment:
      GIN_MODE: debug
    ports:
      - "8080:8080"
```

### Production (docker-compose.prod.yml)

```yaml
services:
  api:
    image: registry.example.com/advanced-user-api:v1.0.0  # Pre-built image
    environment:
      GIN_MODE: release
      JWT_SECRET: ${JWT_SECRET}  # Из secrets
    restart: always
    deploy:
      replicas: 3              # Несколько инстансов
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
```

**Запуск:**
```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

---

## 🛠️ Makefile Integration

**Файл:** `Makefile`

```makefile
# Docker команды в Makefile
.PHONY: docker-build docker-up docker-down docker-logs

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api

docker-clean:
	docker compose down -v
	docker system prune -f
```

**Использование:**
```bash
make docker-build
make docker-up
make docker-logs
```

---

## 📖 Примеры реальных сценариев

### Сценарий 1: Первый запуск

```bash
# 1. Клонировать репозиторий
git clone https://github.com/AlexRodving/advanced-user-api.git
cd advanced-user-api

# 2. Создать .env
cp env.example .env

# 3. Запустить всё
docker compose up -d

# 4. Проверить статус
docker compose ps

# 5. Посмотреть логи
docker compose logs -f api

# 6. Проверить API
curl http://localhost:8080/health
```

---

### Сценарий 2: Обновление кода

```bash
# 1. Остановить API
docker compose stop api

# 2. Пересобрать образ
docker compose build api

# 3. Запустить обновлённый
docker compose up -d api

# Или всё в одной команде
docker compose up -d --build api
```

---

### Сценарий 3: Бэкап БД

```bash
# Создать бэкап
docker compose exec postgres pg_dump -U postgres advanced_api > backup_$(date +%Y%m%d).sql

# Восстановить
docker compose exec -T postgres psql -U postgres advanced_api < backup_20251015.sql
```

---

### Сценарий 4: Очистка

```bash
# Остановить всё
docker compose down

# Удалить volumes (данные!)
docker compose down -v

# Удалить образы
docker rmi advanced-user-api:latest

# Полная очистка Docker
docker system prune -a --volumes
```

---

## 🔗 Docker Network подробнее

### Внутреннее общение сервисов

```
┌─────────────────────────────────────────┐
│         Docker Network: app-network     │
│                                         │
│  ┌──────────┐  ┌──────────┐           │
│  │   API    │→ │ postgres │           │
│  │ :8080    │  │ :5432    │           │
│  └──────────┘  └──────────┘           │
│       ↓                                │
│  ┌──────────┐                         │
│  │  redis   │                         │
│  │ :6379    │                         │
│  └──────────┘                         │
└─────────────────────────────────────────┘
       ↕ (через порты)
    Host Machine
    localhost:8080
    localhost:5432
```

**DNS внутри сети:**
- `postgres` → IP адрес postgres контейнера
- `redis` → IP адрес redis контейнера
- `api` → IP адрес api контейнера

---

## 📦 Docker Images Best Practices

### Размер образа

```bash
# Проверить размер
docker images advanced-user-api

# Наш образ: ~15-20MB ✅
# Без multi-stage: ~500MB ❌
```

### Слои (Layers)

```bash
# Посмотреть слои образа
docker history advanced-user-api:latest

# Каждая команда в Dockerfile = новый слой
RUN apk add git        # Слой 1
RUN go mod download    # Слой 2
RUN go build           # Слой 3
```

**Оптимизация:**
```dockerfile
# ❌ Много слоёв
RUN apk add git
RUN apk add ca-certificates
RUN apk add curl

# ✅ Один слой
RUN apk add --no-cache git ca-certificates curl
```

---

## 🌍 Environment Variables

### Способы передачи

#### 1. Прямо в docker-compose.yml

```yaml
environment:
  JWT_SECRET: hardcoded-secret  # ❌ Плохо для production
```

#### 2. Из .env файла

```yaml
environment:
  JWT_SECRET: ${JWT_SECRET}
```

```bash
# .env
JWT_SECRET=my-secret-key
```

#### 3. Через env_file

```yaml
env_file:
  - .env
  - .env.production
```

---

## 🧪 Тестирование с Docker

```bash
# Запустить тестовую БД
docker compose -f docker-compose.test.yml up -d postgres

# Выполнить тесты
go test ./...

# Остановить тестовую БД
docker compose -f docker-compose.test.yml down -v
```

---

## 📖 См. также

- [Deployment Guide](docs/DEPLOY.md) - Деплой с Docker
- [GitHub Actions](docs/GITHUB_ACTIONS.md) - CI/CD с Docker
- [Architecture](docs/ARCHITECTURE.md) - Как Docker вписывается в архитектуру

---

## 📚 Дополнительные ресурсы

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Reference](https://docs.docker.com/compose/compose-file/)
- [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Alpine Linux](https://alpinelinux.org/)

---

**Общая экономия:** 500MB → 15MB образ = **97% меньше!** 🎉

