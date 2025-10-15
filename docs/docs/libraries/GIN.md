# 🍸 Gin Web Framework

## Overview

Gin - высокопроизводительный HTTP веб-фреймворк для Go.

**Официальная документация:** https://gin-gonic.com/docs/

**Версия в проекте:** v1.11.0

---

## Основные компоненты

### gin.Engine

Главный роутер приложения.

```go
// Создание
router := gin.Default()  // С Logger и Recovery middleware
router := gin.New()      // Без middleware
```

**Где используем:**
- [`cmd/api/main.go`](../../cmd/api/main.go) - создание роутера

---

## Методы которые используем

### 1. gin.Default()

**Описание:** Создаёт Engine с предустановленными middleware (Logger, Recovery)

**Сигнатура:**
```go
func Default() *Engine
```

**Параметры:** нет

**Возвращает:** `*Engine` - экземпляр роутера

**Пример использования:**
```go
router := gin.Default()
```

**Где используем:**
- `cmd/api/main.go:46`

---

### 2. router.Use()

**Описание:** Добавляет middleware в цепочку обработки

**Сигнатура:**
```go
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes
```

**Параметры:**
- `middleware ...HandlerFunc` - один или несколько middleware

**Возвращает:** `IRoutes` - для цепочки вызовов

**Пример использования:**
```go
router.Use(middleware.CORSMiddleware())
router.Use(gin.Logger(), gin.Recovery())
```

**Где используем:**
- `internal/handler/routes.go:27` - CORS middleware

---

### 3. router.Group()

**Описание:** Создаёт группу маршрутов с общим префиксом

**Сигнатура:**
```go
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup
```

**Параметры:**
- `relativePath string` - префикс пути (например, "/api/v1")
- `handlers ...HandlerFunc` - middleware для группы (опционально)

**Возвращает:** `*RouterGroup` - новая группа

**Пример использования:**
```go
v1 := router.Group("/api/v1")
{
    v1.POST("/auth/register", authHandler.Register)
}

// С middleware
protected := router.Group("/api/v1")
protected.Use(middleware.AuthMiddleware(cfg))
```

**Где используем:**
- `internal/handler/routes.go:34` - публичные routes
- `internal/handler/routes.go:45` - защищённые routes

---

### 4. group.POST() / GET() / PUT() / DELETE()

**Описание:** Регистрирует handler для HTTP метода

**Сигнатура:**
```go
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes
func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes
func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes
```

**Параметры:**
- `relativePath string` - путь endpoint (например, "/users/:id")
- `handlers ...HandlerFunc` - один или несколько handlers

**Возвращает:** `IRoutes`

**Пример использования:**
```go
v1.POST("/auth/register", authHandler.Register)
v1.GET("/users/:id", userHandler.GetByID)
v1.PUT("/users/:id", userHandler.Update)
v1.DELETE("/users/:id", userHandler.Delete)
```

**Где используем:**
- `internal/handler/routes.go:36-38` - auth endpoints
- `internal/handler/routes.go:47-51` - user endpoints

---

### 5. c.ShouldBindJSON()

**Описание:** Парсит JSON из request body и валидирует

**Сигнатура:**
```go
func (c *Context) ShouldBindJSON(obj interface{}) error
```

**Параметры:**
- `obj interface{}` - указатель на структуру для парсинга

**Возвращает:** 
- `error` - nil если успешно, иначе ошибка валидации

**Пример использования:**
```go
var req domain.RegisterRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

**Валидация по тегам:**
```go
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

**Где используем:**
- `internal/handler/auth_handler.go:30` - Register
- `internal/handler/auth_handler.go:70` - Login
- `internal/handler/user_handler.go:81` - Update

---

### 6. c.JSON()

**Описание:** Отправляет JSON ответ

**Сигнатура:**
```go
func (c *Context) JSON(code int, obj interface{})
```

**Параметры:**
- `code int` - HTTP status code
- `obj interface{}` - объект для сериализации в JSON

**Возвращает:** ничего

**Пример использования:**
```go
c.JSON(200, user)
c.JSON(201, gin.H{"token": token, "user": user})
c.JSON(400, gin.H{"error": "invalid input"})
```

**Где используем:**
- Все handlers в `internal/handler/`

---

### 7. c.Param()

**Описание:** Получает параметр из URL path

**Сигнатура:**
```go
func (c *Context) Param(key string) string
```

**Параметры:**
- `key string` - имя параметра (без ':')

**Возвращает:** `string` - значение параметра

**Пример использования:**
```go
// Route: /users/:id
id := c.Param("id")  // "42" для /users/42

// Конвертация в число
idUint, err := strconv.ParseUint(id, 10, 32)
if err != nil {
    c.JSON(400, gin.H{"error": "invalid ID"})
    return
}
```

**Где используем:**
- `internal/handler/user_handler.go:35` - GetByID
- `internal/handler/user_handler.go:63` - Update
- `internal/handler/user_handler.go:105` - Delete

---

### 8. c.Set() / c.Get()

**Описание:** Устанавливает/получает значения в контексте запроса

**Сигнатура:**
```go
func (c *Context) Set(key string, value interface{})
func (c *Context) Get(key string) (value interface{}, exists bool)
func (c *Context) MustGet(key string) interface{}
```

**Параметры:**
- `key string` - ключ
- `value interface{}` - значение (для Set)

**Возвращает:**
- `Get()` - значение и флаг существования
- `MustGet()` - значение (panic если нет)

**Пример использования:**
```go
// В middleware
c.Set("user_id", uint(42))

// В handler
userID, exists := c.Get("user_id")
if !exists {
    c.JSON(401, gin.H{"error": "unauthorized"})
    return
}

// Или с приведением типа
userID := c.MustGet("user_id").(uint)
```

**Где используем:**
- `internal/middleware/auth.go:70` - Set user_id
- `internal/handler/auth_handler.go:120` - Get user_id

---

### 9. c.GetHeader()

**Описание:** Получает значение заголовка

**Сигнатура:**
```go
func (c *Context) GetHeader(key string) string
```

**Параметры:**
- `key string` - имя заголовка

**Возвращает:** `string` - значение заголовка

**Пример использования:**
```go
authHeader := c.GetHeader("Authorization")
// "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Где используем:**
- `internal/middleware/auth.go:39` - получение JWT token

---

### 10. c.AbortWithStatusJSON()

**Описание:** Останавливает цепочку handlers и отправляет JSON

**Сигнатура:**
```go
func (c *Context) AbortWithStatusJSON(code int, jsonObj interface{})
```

**Параметры:**
- `code int` - HTTP status code
- `jsonObj interface{}` - JSON объект

**Возвращает:** ничего

**Пример использования:**
```go
c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
// Дальнейшие handlers не выполнятся
```

**Где используем:**
- `internal/middleware/auth.go:43,52,61` - ошибки аутентификации

---

### 11. c.Next()

**Описание:** Передаёт управление следующему handler в цепочке

**Сигнатура:**
```go
func (c *Context) Next()
```

**Параметры:** нет

**Возвращает:** ничего

**Пример использования:**
```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // До обработки запроса
        log.Println("Before")
        
        c.Next()  // Вызываем следующий handler
        
        // После обработки запроса
        log.Println("After")
    }
}
```

**Где используем:**
- `internal/middleware/auth.go:75` - после успешной аутентификации
- `internal/middleware/cors.go:34` - после CORS
- `internal/middleware/logger.go:24` - в logger

---

### 12. gin.H{}

**Описание:** Shortcut для `map[string]interface{}`

**Сигнатура:**
```go
type H map[string]interface{}
```

**Пример использования:**
```go
c.JSON(200, gin.H{
    "message": "success",
    "data": user,
    "count": 10,
})
```

**Где используем:**
- Все handlers для error responses и простых JSON объектов

---

### 13. c.Request

**Описание:** Доступ к оригинальному `*http.Request`

**Тип:** `*http.Request`

**Пример использования:**
```go
method := c.Request.Method        // "GET", "POST", etc.
path := c.Request.URL.Path        // "/api/v1/users"
body := c.Request.Body            // io.ReadCloser
```

**Где используем:**
- `internal/middleware/logger.go:22` - URL.Path для логирования

---

### 14. c.Writer

**Описание:** Доступ к `http.ResponseWriter`

**Тип:** `ResponseWriter`

**Пример использования:**
```go
status := c.Writer.Status()       // HTTP status code
c.Writer.Header().Set("X-Custom", "value")
```

**Где используем:**
- `internal/middleware/logger.go:28` - получение status code

---

### 15. c.ClientIP()

**Описание:** Получает IP адрес клиента

**Сигнатура:**
```go
func (c *Context) ClientIP() string
```

**Параметры:** нет

**Возвращает:** `string` - IP адрес

**Пример использования:**
```go
ip := c.ClientIP()  // "192.168.1.1"
```

**Где используем:**
- `internal/middleware/logger.go:29` - логирование IP

---

### 16. c.Header()

**Описание:** Устанавливает заголовок ответа

**Сигнатура:**
```go
func (c *Context) Header(key, value string)
```

**Параметры:**
- `key string` - имя заголовка
- `value string` - значение

**Возвращает:** ничего

**Пример использования:**
```go
c.Header("Content-Type", "application/json")
c.Header("X-Total-Count", "100")
```

**Где используем:**
- `internal/middleware/cors.go:16-23` - CORS заголовки

---

## Validation Tags

Gin использует validator/v10 для валидации.

### Часто используемые теги:

| Tag | Описание | Пример |
|-----|----------|--------|
| `required` | Обязательное поле | `binding:"required"` |
| `email` | Валидный email | `binding:"required,email"` |
| `min=N` | Минимальная длина | `binding:"min=6"` |
| `max=N` | Максимальная длина | `binding:"max=100"` |
| `omitempty` | Опциональное поле | `binding:"omitempty,email"` |
| `oneof=a b` | Одно из значений | `binding:"oneof=admin user"` |

**Где используем:**
- `internal/domain/user.go:54-60` - RegisterRequest
- `internal/domain/user.go:67-72` - LoginRequest
- `internal/domain/user.go:80-85` - UpdateUserRequest

---

## Middleware Pattern

```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Предобработка
        
        // 2. Проверки
        if shouldAbort {
            c.AbortWithStatusJSON(400, gin.H{"error": "..."})
            return
        }
        
        // 3. Установка данных в контекст
        c.Set("key", value)
        
        // 4. Вызов следующего handler
        c.Next()
        
        // 5. Постобработка (после всех handlers)
    }
}
```

**Примеры:**
- [`internal/middleware/auth.go`](../../internal/middleware/auth.go) - JWT auth
- [`internal/middleware/cors.go`](../../internal/middleware/cors.go) - CORS
- [`internal/middleware/logger.go`](../../internal/middleware/logger.go) - Logger

---

## Error Handling

### Способ 1: Прямой ответ
```go
if err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### Способ 2: Abort
```go
if err != nil {
    c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
    return
}
```

**Где используем:**
- Все handlers в `internal/handler/`

---

## Best Practices

### 1. Используйте binding теги
```go
type Request struct {
    Email string `json:"email" binding:"required,email"`
}
```

### 2. Проверяйте ошибки binding
```go
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### 3. Используйте группы для версионирования
```go
v1 := router.Group("/api/v1")
v2 := router.Group("/api/v2")
```

### 4. Middleware для общей логики
```go
protected := router.Group("/api/v1")
protected.Use(middleware.AuthMiddleware(cfg))
```

---

## См. также

- [GORM](docs/libraries/GORM.md) - ORM для работы с БД
- [JWT](docs/libraries/JWT.md) - JWT токены
- [Bcrypt](docs/libraries/BCRYPT_VIPER.md) - Хеширование паролей
- [Architecture Guide](docs/ARCHITECTURE.md) - Архитектура проекта

