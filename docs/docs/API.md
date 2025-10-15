# 📡 API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

---

## 🔓 Public Endpoints (без токена)

### 1. Health Check
Проверка работоспособности API

**Endpoint:** `GET /health`

**Response:**
```json
{
  "service": "advanced-user-api",
  "status": "ok"
}
```

**Example:**
```bash
curl http://localhost:8080/health
```

---

### 2. Register
Регистрация нового пользователя

**Endpoint:** `POST /api/v1/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "name": "User Name",
  "password": "password123"
}
```

**Validation:**
- `email` - обязательно, валидный email
- `name` - обязательно, минимум 2 символа
- `password` - обязательно, минимум 6 символов

**Response 201 Created:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "User Name",
    "role": "user",
    "created_at": "2025-10-15T10:00:00Z",
    "updated_at": "2025-10-15T10:00:00Z"
  }
}
```

**Errors:**
- `400 Bad Request` - невалидные данные
- `409 Conflict` - email уже зарегистрирован

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "name": "Alice",
    "password": "secret123"
  }'
```

---

### 3. Login
Вход в систему

**Endpoint:** `POST /api/v1/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response 200 OK:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "User Name",
    "role": "user",
    "created_at": "2025-10-15T10:00:00Z",
    "updated_at": "2025-10-15T10:00:00Z"
  }
}
```

**Errors:**
- `401 Unauthorized` - неверный email или пароль

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "secret123"
  }'
```

---

## 🔒 Protected Endpoints (требуют JWT токен)

### Аутентификация
Все защищённые endpoints требуют JWT токен в заголовке:

```
Authorization: Bearer <your-jwt-token>
```

**Errors:**
- `401 Unauthorized` - токен отсутствует, невалиден или истёк

---

### 4. Get Current User
Получить данные текущего пользователя

**Endpoint:** `GET /api/v1/auth/me`

**Headers:**
```
Authorization: Bearer <token>
```

**Response 200 OK:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "User Name",
  "role": "user",
  "created_at": "2025-10-15T10:00:00Z",
  "updated_at": "2025-10-15T10:00:00Z"
}
```

**Example:**
```bash
TOKEN="your-jwt-token"

curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

---

### 5. Get All Users
Получить список всех пользователей

**Endpoint:** `GET /api/v1/users`

**Headers:**
```
Authorization: Bearer <token>
```

**Response 200 OK:**
```json
[
  {
    "id": 1,
    "email": "alice@example.com",
    "name": "Alice",
    "role": "user",
    "created_at": "2025-10-15T10:00:00Z",
    "updated_at": "2025-10-15T10:00:00Z"
  },
  {
    "id": 2,
    "email": "bob@example.com",
    "name": "Bob",
    "role": "user",
    "created_at": "2025-10-15T11:00:00Z",
    "updated_at": "2025-10-15T11:00:00Z"
  }
]
```

**Example:**
```bash
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

---

### 6. Get User by ID
Получить пользователя по ID

**Endpoint:** `GET /api/v1/users/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` - ID пользователя (integer)

**Response 200 OK:**
```json
{
  "id": 1,
  "email": "alice@example.com",
  "name": "Alice",
  "role": "user",
  "created_at": "2025-10-15T10:00:00Z",
  "updated_at": "2025-10-15T10:00:00Z"
}
```

**Errors:**
- `404 Not Found` - пользователь не найден

**Example:**
```bash
curl http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

### 7. Update User
Обновить данные пользователя

**Endpoint:** `PUT /api/v1/users/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Path Parameters:**
- `id` - ID пользователя (integer)

**Request Body:**
```json
{
  "name": "New Name",
  "email": "newemail@example.com"
}
```

**Validation:**
- `name` - опционально, минимум 2 символа если указано
- `email` - опционально, валидный email если указан

**Response 200 OK:**
```json
{
  "id": 1,
  "email": "newemail@example.com",
  "name": "New Name",
  "role": "user",
  "created_at": "2025-10-15T10:00:00Z",
  "updated_at": "2025-10-15T12:00:00Z"
}
```

**Errors:**
- `404 Not Found` - пользователь не найден
- `400 Bad Request` - невалидные данные

**Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Updated",
    "email": "alice.new@example.com"
  }'
```

---

### 8. Delete User
Удалить пользователя (Soft Delete)

**Endpoint:** `DELETE /api/v1/users/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` - ID пользователя (integer)

**Response 200 OK:**
```json
{
  "message": "пользователь удалён"
}
```

**Errors:**
- `404 Not Found` - пользователь не найден

**Note:** Используется soft delete - запись не удаляется физически, а помечается как удалённая (поле `deleted_at`)

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🔑 JWT Token

### Структура токена
JWT токен содержит:
- `user_id` - ID пользователя
- `email` - Email пользователя
- `role` - Роль пользователя
- `exp` - Время истечения (24 часа)
- `iat` - Время создания
- `iss` - Издатель (advanced-user-api)

### Пример использования
```bash
# 1. Получите токен через login или register
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

# 2. Извлеките токен
TOKEN=$(echo $RESPONSE | jq -r '.token')

# 3. Используйте для защищённых endpoints
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📋 HTTP Status Codes

| Code | Значение | Когда используется |
|------|----------|-------------------|
| 200 | OK | Успешный GET, PUT, DELETE |
| 201 | Created | Успешный POST (создание) |
| 400 | Bad Request | Невалидные данные |
| 401 | Unauthorized | Нет токена или токен невалиден |
| 404 | Not Found | Ресурс не найден |
| 409 | Conflict | Email уже существует |
| 500 | Internal Server Error | Ошибка сервера |

---

## 🧪 Тестирование API

### Postman Collection
Импортируйте эту коллекцию в Postman:

```json
{
  "info": {
    "name": "Advanced User API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Auth",
      "item": [
        {
          "name": "Register",
          "request": {
            "method": "POST",
            "header": [{"key": "Content-Type", "value": "application/json"}],
            "body": {
              "mode": "raw",
              "raw": "{\"email\":\"test@example.com\",\"name\":\"Test\",\"password\":\"password123\"}"
            },
            "url": "{{base_url}}/api/v1/auth/register"
          }
        },
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "header": [{"key": "Content-Type", "value": "application/json"}],
            "body": {
              "mode": "raw",
              "raw": "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
            },
            "url": "{{base_url}}/api/v1/auth/login"
          }
        }
      ]
    }
  ],
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080"
    }
  ]
}
```

### cURL Examples
Полные примеры в [`docs/QUICKSTART.md`](docs/QUICKSTART.md)

---

## 🔐 Безопасность

### Best Practices
1. **HTTPS в production** - всегда используйте HTTPS
2. **Сильный JWT_SECRET** - минимум 32 случайных символа
3. **Короткий срок жизни токенов** - 24 часа или меньше
4. **Rate limiting** - ограничьте количество запросов
5. **Валидация** - все входные данные валидируются

### Генерация безопасного секрета
```bash
openssl rand -base64 32
```

---

## 📖 Дополнительная документация

- [Quick Start](docs/QUICKSTART.md) - Быстрый старт
- [Deployment](docs/DEPLOY.md) - Деплой в production
- [Project Summary](docs/PROJECT_SUMMARY.md) - Сводка проекта

