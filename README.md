# 🚀 Advanced User Management API

> Production-ready REST API с JWT аутентификацией, Gin фреймворком, GORM ORM, Docker и полной документацией

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?logo=postgresql)](https://www.postgresql.org/)
[![Gin](https://img.shields.io/badge/Gin-Web%20Framework-00ADD8)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-ORM-00ADD8)](https://gorm.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 🌐 [📖 Документация на сайте](https://alexrodving.github.io/advanced-user-api/)

**Красивый веб-сайт с полной документацией проекта!**

✨ Gradient дизайн | 🔍 Поиск по документации | 📱 Адаптивная вёрстка | 🎨 Syntax highlighting

## 📚 Документация

> **💡 Совет:** Вся документация также доступна как [красивый веб-сайт](https://alexrodving.github.io/advanced-user-api/) с поиском и syntax highlighting!

### Основная документация
- 📡 **[API Documentation](docs/docs/API.md)** - Полное описание всех endpoints с примерами
- ⚡ **[Quick Start](docs/docs/QUICKSTART.md)** - Запуск за 3 команды
- 🏗️ **[Architecture Guide](docs/docs/ARCHITECTURE.md)** - Дизайн системы и структура кода
- 🧪 **[Testing Guide](docs/docs/TESTING.md)** - Unit и Integration тесты
- 🌿 **[Git Workflow](docs/docs/GIT_WORKFLOW.md)** - Работа с Git в команде (branches, PR, commits)
- 🚀 **[Deployment Guide](docs/docs/DEPLOY.md)** - Деплой в production (VPS, K8s, Heroku, AWS)
- 📊 **[Project Summary](docs/docs/PROJECT_SUMMARY.md)** - Статистика и технический обзор

### Библиотеки (детальная документация)
- 📖 **[Libraries Overview](docs/docs/libraries/)** - Обзор всех библиотек
- 🍸 **[Gin Framework](docs/docs/libraries/GIN.md)** - HTTP router, handlers, middleware
- 🗄️ **[GORM ORM](docs/docs/libraries/GORM.md)** - Database operations, models, migrations
- 🔐 **[JWT](docs/docs/libraries/JWT.md)** - Token generation and validation
- 🔒 **[Bcrypt & Viper](docs/docs/libraries/BCRYPT_VIPER.md)** - Password hashing & configuration

### Веб-сайт документации
- 🌐 **[Documentation Website Guide](docs/docs/WEBSITE.md)** - Как работает сайт документации (Docsify, CSS, GitHub Pages)

---

## 🎯 Что нового по сравнению с базовой версией?

### ✨ Новые технологии:
- ✅ **Gin** - популярный веб-фреймворк (быстрее net/http)
- ✅ **GORM** - ORM для упрощения работы с БД
- ✅ **JWT** - аутентификация с токенами
- ✅ **Docker & Docker Compose** - контейнеризация
- ✅ **Middleware** - логирование, CORS, recovery
- ✅ **Unit & Integration тесты** - полное покрытие
- ✅ **Graceful Shutdown** - корректная остановка
- ✅ **Structured Logging** - zap logger
- ✅ **Swagger** - документация API
- ✅ **Validator** - валидация запросов

---

## 📂 Структура проекта

```
advanced-user-api/
├── cmd/
│   └── api/
│       └── main.go                  # Точка входа
│
├── internal/
│   ├── config/
│   │   └── config.go               # Конфигурация (Viper)
│   │
│   ├── domain/
│   │   ├── user.go                 # User модель
│   │   └── auth.go                 # Auth модели (Login, Register)
│   │
│   ├── repository/
│   │   ├── user_repository.go      # GORM repository
│   │   └── interfaces.go           # Интерфейсы
│   │
│   ├── service/
│   │   ├── user_service.go         # Бизнес-логика
│   │   ├── auth_service.go         # JWT аутентификация
│   │   └── interfaces.go           # Интерфейсы
│   │
│   ├── handler/
│   │   ├── user_handler.go         # User endpoints (Gin)
│   │   ├── auth_handler.go         # Auth endpoints
│   │   └── routes.go               # Настройка маршрутов
│   │
│   ├── middleware/
│   │   ├── auth.go                 # JWT middleware
│   │   ├── logger.go               # Logging middleware
│   │   ├── cors.go                 # CORS middleware
│   │   └── recovery.go             # Panic recovery
│   │
│   └── pkg/
│       ├── jwt/
│       │   └── jwt.go              # JWT утилиты
│       ├── password/
│       │   └── password.go         # Bcrypt хеширование
│       └── validator/
│           └── validator.go        # Custom validators
│
├── tests/
│   ├── unit/
│   │   └── auth_service_test.go    # Unit тесты с моками
│   └── integration/
│       └── api_test.go             # Integration тесты
│
├── docs/                           # 🌐 GitHub Pages сайт документации
│   ├── index.html                  # Docsify + Custom CSS
│   ├── README.md                   # Главная страница сайта
│   ├── _sidebar.md                 # Боковое меню
│   ├── _navbar.md                  # Верхнее меню
│   ├── .nojekyll                   # Отключение Jekyll
│   └── docs/                       # Markdown документация
│       ├── API.md                  # 📡 API endpoints
│       ├── ARCHITECTURE.md         # 🏗️ Архитектура
│       ├── TESTING.md              # 🧪 Тестирование
│       ├── GIT_WORKFLOW.md         # 🌿 Git workflow
│       ├── QUICKSTART.md           # ⚡ Быстрый старт
│       ├── DEPLOY.md               # 🚀 Деплой
│       ├── PROJECT_SUMMARY.md      # 📊 Сводка
│       ├── WEBSITE.md              # 🌐 Как работает сайт
│       └── libraries/              # 📚 Документация библиотек
│           ├── README.md           # Обзор
│           ├── GIN.md              # Gin (16 методов)
│           ├── GORM.md             # GORM ORM
│           ├── JWT.md              # JWT auth
│           └── BCRYPT_VIPER.md     # Bcrypt & Viper
│
├── .github/
│   └── workflows/
│       └── ci.yml                  # GitHub Actions CI/CD
│
├── docker/
│   └── Dockerfile                  # Multi-stage Dockerfile (15MB)
│
├── docker-compose.yml              # Весь стек (API + PostgreSQL + Redis)
├── Makefile                        # Команды для разработки
├── env.example                     # Пример переменных окружения
├── .dockerignore
├── .gitignore
├── LICENSE                         # MIT License
├── go.mod
├── go.sum
└── README.md
```

---

## 🛠️ Технологический стек

| Технология | Зачем нужна | Альтернативы |
|------------|-------------|--------------|
| **Gin** | Веб-фреймворк | Echo, Fiber, Chi |
| **GORM** | ORM для БД | sqlx, pgx |
| **JWT** | Аутентификация | Session-based, OAuth |
| **Zap** | Структурированные логи | Logrus, Zerolog |
| **Viper** | Конфигурация | envconfig |
| **Testify** | Assertions для тестов | Standard testing |
| **Swagger** | API документация | - |
| **Docker** | Контейнеризация | - |
| **PostgreSQL** | База данных | MySQL, MongoDB |
| **Redis** | Кеш, сессии | Memcached |

---

## 🚀 Быстрый старт

### С Docker Compose (рекомендуется):

```bash
# Запустить всё (API + PostgreSQL + Redis)
docker-compose up -d

# Применить миграции
make migrate-up

# API доступен на http://localhost:8080
```

### Без Docker:

```bash
# Установить зависимости
go mod tidy

# Запустить PostgreSQL
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

# Применить миграции
make migrate-up

# Запустить приложение
go run cmd/api/main.go
```

---

## 🔐 Аутентификация (JWT)

### 1. Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "name": "Alice",
    "password": "secret123"
  }'
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "alice@example.com",
    "name": "Alice"
  }
}
```

### 2. Вход (Login)
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "secret123"
  }'
```

### 3. Защищённые endpoints
```bash
# Без токена - ошибка 401
curl http://localhost:8080/api/v1/users

# С токеном - успех
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🧪 Тестирование

### Запуск тестов:

```bash
# Все тесты
make test

# С покрытием
make test-coverage

# Только unit тесты
go test ./internal/service/... -v

# Только integration тесты
go test ./tests/integration/... -v
```

---

## 🐳 Docker

### Сборка образа:

```bash
# Production образ
docker build -t user-api:latest -f docker/Dockerfile .

# Development образ
docker build -t user-api:dev -f docker/Dockerfile.dev .
```

### Запуск с Docker Compose:

```bash
# Запустить
docker-compose up -d

# Посмотреть логи
docker-compose logs -f api

# Остановить
docker-compose down

# Остановить и удалить данные
docker-compose down -v
```

---

## 📡 API Endpoints

### Публичные (без авторизации):

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/v1/auth/register` | Регистрация |
| POST | `/api/v1/auth/login` | Вход |
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/docs/*` | Swagger UI |

### Защищённые (требуют JWT токен):

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/api/v1/users` | Список пользователей |
| GET | `/api/v1/users/:id` | Получить пользователя |
| PUT | `/api/v1/users/:id` | Обновить пользователя |
| DELETE | `/api/v1/users/:id` | Удалить пользователя |
| GET | `/api/v1/users/me` | Текущий пользователь |

---

## 🎯 План реализации (пошаговый)

### Этап 1: Настройка проекта (День 1)
- [x] Структура директорий
- [ ] go.mod с зависимостями
- [ ] Docker и docker-compose
- [ ] Makefile для команд

### Этап 2: GORM и базовые модели (День 1-2)
- [ ] GORM модели
- [ ] Миграции через GORM
- [ ] Repository с GORM

### Этап 3: Gin роутер (День 2)
- [ ] Настройка Gin
- [ ] Базовые routes
- [ ] Middleware (CORS, Logger, Recovery)

### Этап 4: JWT Аутентификация (День 2-3)
- [ ] JWT генерация и валидация
- [ ] Register endpoint
- [ ] Login endpoint
- [ ] Auth middleware

### Этап 5: Улучшенные handlers (День 3)
- [ ] Все CRUD операции
- [ ] Валидация с validator
- [ ] Error handling

### Этап 6: Тестирование (День 4)
- [ ] Unit тесты для service
- [ ] Unit тесты для handlers
- [ ] Integration тесты
- [ ] Моки с testify

### Этап 7: Документация (День 4)
- [ ] Swagger интеграция
- [ ] API документация
- [ ] Postman коллекция

### Этап 8: Деплой (День 5)
- [ ] Dockerfile оптимизация
- [ ] CI/CD (GitHub Actions)
- [ ] Деплой на VPS/Cloud

---

## 💡 Ключевые улучшения

### 1. Gin вместо net/http
```go
// Было (net/http):
http.HandleFunc("/users", handler)

// Стало (Gin):
r := gin.Default()
r.GET("/users", handler)
r.POST("/users", handler)
```

### 2. GORM вместо database/sql
```go
// Было (SQL):
db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)

// Стало (GORM):
db.First(&user, id)
```

### 3. JWT Middleware
```go
// Защищённые routes
authorized := r.Group("/api/v1/users")
authorized.Use(middleware.AuthMiddleware())
{
    authorized.GET("", userHandler.GetAll)
    authorized.GET("/:id", userHandler.GetByID)
}
```

---

## 🎓 Что изучите

1. **Gin framework** - роутинг, middleware, параметры
2. **GORM** - ORM, миграции, ассоциации
3. **JWT** - создание, валидация, refresh tokens
4. **bcrypt** - хеширование паролей
5. **Docker** - multi-stage builds, compose
6. **Testing** - моки, assertions, table-driven tests
7. **Swagger** - автоматическая документация
8. **Makefile** - автоматизация задач
9. **Graceful Shutdown** - корректная остановка
10. **CI/CD** - автоматический деплой

---

## 🔥 Это реальный production проект!

После завершения у вас будет:
- ✅ Полноценный API готовый к деплою
- ✅ Автоматические тесты
- ✅ Docker контейнеры
- ✅ CI/CD pipeline
- ✅ Swagger документация
- ✅ **Production-ready код для портфолио**

---

## 📖 Additional Resources

- [Git Workflow Guide](docs/GIT_WORKFLOW.md) - Работа с Git в команде
- [Go Best Practices](https://go.dev/doc/effective_go)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

