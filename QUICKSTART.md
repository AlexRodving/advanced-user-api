# ⚡ Быстрый старт

## 🚀 Запуск за 3 команды

### Вариант 1: С Docker Compose (рекомендуется)

```bash
# 1. Клонируйте репозиторий
git clone <your-repo-url>
cd 08_advanced_api

# 2. Запустите всё (API + PostgreSQL + Redis)
docker-compose up -d

# 3. Готово! API работает на http://localhost:8080
```

### Вариант 2: Локально

```bash
# 1. Запустите PostgreSQL
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
docker exec -it postgres psql -U postgres -c "CREATE DATABASE advanced_api;"

# 2. Установите зависимости
go mod tidy

# 3. Запустите API
go run cmd/api/main.go
```

---

## 🧪 Тестирование API

### 1. Регистрация
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test","password":"password123"}'
```

**Получите:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {"id": 1, "email": "test@example.com", ...}
}
```

### 2. Вход
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 3. Получить список пользователей (с токеном)
```bash
TOKEN="ваш_токен_из_register"

curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Текущий пользователь
```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📊 Архитектура

```
┌──────────────┐
│   Gin Router │  ← HTTP запросы
└──────┬───────┘
       ↓
┌──────────────┐
│   Middleware │  ← JWT, CORS, Logger
└──────┬───────┘
       ↓
┌──────────────┐
│   Handlers   │  ← HTTP обработчики
└──────┬───────┘
       ↓
┌──────────────┐
│   Services   │  ← Бизнес-логика, JWT, bcrypt
└──────┬───────┘
       ↓
┌──────────────┐
│  Repository  │  ← GORM операции
└──────┬───────┘
       ↓
┌──────────────┐
│  PostgreSQL  │  ← База данных
└──────────────┘
```

---

## 🔥 Ключевые особенности

✅ **JWT аутентификация** - безопасная авторизация  
✅ **Bcrypt** - хеширование паролей  
✅ **GORM** - ORM с Auto Migration  
✅ **Gin** - быстрый веб-фреймворк  
✅ **Clean Architecture** - 3 слоя  
✅ **Docker** - полная контейнеризация  
✅ **CORS** - поддержка frontend  
✅ **Graceful Shutdown** - корректная остановка  
✅ **Подробные комментарии** - весь код объяснён  

---

## 📖 Документация

- **README.md** - полное описание проекта
- **PLAN.md** - план разработки
- **Код с комментариями** - каждый файл подробно объяснён

---

## 🎓 Для обучения

Этот проект идеален для изучения:
- Go веб-разработки
- JWT аутентификации
- GORM ORM
- Docker контейнеризации
- Clean Architecture

**Каждый файл содержит подробные комментарии!** 📚

---

## 💼 Для портфолио

- Production-ready код
- Современный стек технологий
- Best practices
- Готов к деплою

---

**Приятного изучения!** 🚀

