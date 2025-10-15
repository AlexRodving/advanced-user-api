# 🔐 JWT - JSON Web Tokens

## Overview

JWT (JSON Web Token) - стандарт для безопасной передачи информации между сторонами.

**Официальная документация:** https://github.com/golang-jwt/jwt

**Версия в проекте:** v5.3.0 (github.com/golang-jwt/jwt/v5)

---

## Структура JWT

### Формат токена

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MDAwMDAwMDB9.signature
│                                     │                                     │
│        Header (base64)              │       Payload (base64)              │  Signature
```

### Компоненты:

1. **Header** - алгоритм и тип
2. **Payload** - данные (claims)
3. **Signature** - подпись для верификации

---

## Claims (Данные в токене)

### Наша структура Claims

```go
type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
```

**Где определена:**
- `internal/pkg/jwt/jwt.go:14-19`

### RegisteredClaims (стандартные поля)

```go
type RegisteredClaims struct {
    Issuer    string             // "iss" - издатель токена
    Subject   string             // "sub" - субъект
    Audience  jwt.ClaimStrings   // "aud" - аудитория
    ExpiresAt *jwt.NumericDate   // "exp" - время истечения
    NotBefore *jwt.NumericDate   // "nbf" - не раньше чем
    IssuedAt  *jwt.NumericDate   // "iat" - время создания
    ID        string             // "jti" - уникальный ID
}
```

**Где используем:**
- `internal/pkg/jwt/jwt.go:36-42` - установка exp, iat, iss

---

## Методы которые используем

### 1. jwt.NewWithClaims()

**Описание:** Создаёт новый JWT токен с данными

**Сигнатура:**
```go
func NewWithClaims(method SigningMethod, claims Claims) *Token
```

**Параметры:**
- `method SigningMethod` - алгоритм подписи (HS256, RS256, etc.)
- `claims Claims` - данные для токена

**Возвращает:** `*Token` - объект токена

**Пример использования:**
```go
claims := &Claims{
    UserID: user.ID,
    Email:  user.Email,
    Role:   user.Role,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        Issuer:    "advanced-user-api",
    },
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
```

**Где используем:**
- `internal/pkg/jwt/jwt.go:35`

---

### 2. token.SignedString()

**Описание:** Подписывает токен и возвращает строку

**Сигнатура:**
```go
func (t *Token) SignedString(key interface{}) (string, error)
```

**Параметры:**
- `key interface{}` - секретный ключ для подписи ([]byte для HS256)

**Возвращает:**
- `string` - подписанный JWT токен
- `error` - ошибка подписи

**Пример использования:**
```go
secretKey := []byte("your-secret-key")
tokenString, err := token.SignedString(secretKey)
if err != nil {
    return "", err
}

// tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Где используем:**
- `internal/pkg/jwt/jwt.go:45`

---

### 3. jwt.ParseWithClaims()

**Описание:** Парсит и валидирует JWT токен

**Сигнатура:**
```go
func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc, options ...ParserOption) (*Token, error)
```

**Параметры:**
- `tokenString string` - JWT токен для парсинга
- `claims Claims` - структура для заполнения данными
- `keyFunc Keyfunc` - функция получения ключа для верификации

**Возвращает:**
- `*Token` - объект токена
- `error` - ошибка парсинга или валидации

**Пример использования:**
```go
claims := &Claims{}
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    // Проверяем алгоритм
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    
    // Возвращаем ключ для верификации
    return []byte(secretKey), nil
})

if err != nil {
    return nil, err
}

if !token.Valid {
    return nil, errors.New("invalid token")
}

// Теперь claims содержит данные из токена
fmt.Println(claims.UserID, claims.Email)
```

**Где используем:**
- `internal/pkg/jwt/jwt.go:58`

---

### 4. jwt.SigningMethodHS256

**Описание:** Алгоритм подписи HMAC-SHA256

**Тип:** `*SigningMethodHMAC`

**Пример использования:**
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
```

**Альтернативы:**
- `jwt.SigningMethodHS384` - HMAC-SHA384
- `jwt.SigningMethodHS512` - HMAC-SHA512
- `jwt.SigningMethodRS256` - RSA-SHA256 (требует приватный/публичный ключ)

**Где используем:**
- `internal/pkg/jwt/jwt.go:35`

---

### 5. jwt.NewNumericDate()

**Описание:** Создаёт NumericDate для exp, iat, nbf

**Сигнатура:**
```go
func NewNumericDate(t time.Time) *NumericDate
```

**Параметры:**
- `t time.Time` - время

**Возвращает:** `*NumericDate`

**Пример использования:**
```go
now := time.Now()
expirationTime := now.Add(24 * time.Hour)

claims := jwt.RegisteredClaims{
    ExpiresAt: jwt.NewNumericDate(expirationTime),  // Истекает через 24 часа
    IssuedAt:  jwt.NewNumericDate(now),             // Создан сейчас
    NotBefore: jwt.NewNumericDate(now),             // Валиден с текущего момента
}
```

**Где используем:**
- `internal/pkg/jwt/jwt.go:37-39`

---

## Наша реализация

### Generate() - Генерация токена

**Файл:** `internal/pkg/jwt/jwt.go`

**Сигнатура:**
```go
func Generate(user *domain.User, secret string, expiration time.Duration) (string, error)
```

**Параметры:**
- `user *domain.User` - данные пользователя
- `secret string` - секретный ключ
- `expiration time.Duration` - время жизни токена (например, 24 * time.Hour)

**Возвращает:**
- `string` - JWT токен
- `error` - ошибка генерации

**Что делает:**
1. Создаёт Claims с данными пользователя
2. Устанавливает время истечения
3. Подписывает токен секретным ключом
4. Возвращает строку токена

**Пример использования:**
```go
tokenString, err := jwt.Generate(user, cfg.JWTSecret, 24*time.Hour)
if err != nil {
    return nil, err
}

// tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Где используем:**
- `internal/service/auth_service.go:83` - после регистрации
- `internal/service/auth_service.go:131` - после входа

---

### Validate() - Валидация токена

**Сигнатура:**
```go
func Validate(tokenString string, secret string) (*Claims, error)
```

**Параметры:**
- `tokenString string` - JWT токен
- `secret string` - секретный ключ

**Возвращает:**
- `*Claims` - данные из токена
- `error` - ошибка валидации

**Что проверяет:**
1. Формат токена
2. Подпись (signature)
3. Алгоритм подписи
4. Время истечения (exp)
5. Валидность токена

**Пример использования:**
```go
claims, err := jwt.Validate(tokenString, cfg.JWTSecret)
if err != nil {
    return nil, errors.New("invalid token")
}

// Используем данные из токена
userID := claims.UserID
email := claims.Email
role := claims.Role
```

**Где используем:**
- `internal/middleware/auth.go:55` - в auth middleware

---

## Работа с токеном в HTTP

### Формат заголовка Authorization

```
Authorization: Bearer <token>
```

### Извлечение токена

```go
// Получаем заголовок
authHeader := c.GetHeader("Authorization")

// Проверяем формат
const bearerPrefix = "Bearer "
if !strings.HasPrefix(authHeader, bearerPrefix) {
    return errors.New("invalid authorization header")
}

// Извлекаем токен
tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
```

**Где используем:**
- `internal/middleware/auth.go:39-52`

---

## Ошибки валидации

### Частые ошибки JWT

```go
jwt.ErrTokenMalformed       // Токен неправильного формата
jwt.ErrTokenUnverifiable    // Не удаётся верифицировать
jwt.ErrTokenSignatureInvalid // Неверная подпись
jwt.ErrTokenExpired         // Токен истёк
jwt.ErrTokenNotValidYet     // Токен ещё не валиден (nbf)
```

### Проверка типа ошибки

```go
if errors.Is(err, jwt.ErrTokenExpired) {
    return errors.New("token expired, please login again")
}
```

---

## Security Best Practices

### 1. Используйте сильный секретный ключ

```go
// ❌ Плохо
secret := "secret"

// ✅ Хорошо (минимум 32 символа)
secret := "a-very-long-and-random-secret-key-with-at-least-32-characters"

// ✅ Ещё лучше - генерация
secret := generateRandomString(64)
```

### 2. Устанавливайте разумное время истечения

```go
// ❌ Слишком долго
expiration := 365 * 24 * time.Hour  // 1 год

// ✅ Оптимально
expiration := 24 * time.Hour        // 1 день
expiration := 15 * time.Minute      // 15 минут для sensitive operations
```

**Где настраиваем:**
- `internal/config/config.go:30` - JWT_EXPIRATION из env
- По умолчанию: `24h`

### 3. Проверяйте алгоритм подписи

```go
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    // ВАЖНО: Проверяем алгоритм!
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(secret), nil
})
```

**Защищает от:** Algorithm confusion attacks

**Где используем:**
- `internal/pkg/jwt/jwt.go:60-63`

### 4. Не храните чувствительные данные в токене

```go
// ❌ Плохо - пароль в токене!
type BadClaims struct {
    UserID   uint
    Password string  // НИКОГДА!
}

// ✅ Хорошо - только ID и публичные данные
type GoodClaims struct {
    UserID uint
    Email  string
    Role   string
}
```

### 5. Используйте HTTPS в production

```
// ❌ HTTP - токен виден всем!
http://api.example.com/users

// ✅ HTTPS - зашифрованное соединение
https://api.example.com/users
```

---

## Refresh Tokens (опционально)

### Концепция

- **Access Token** - короткий срок (15 мин - 1 час)
- **Refresh Token** - длинный срок (7-30 дней)

### Workflow

```
1. Login → получаем access + refresh tokens
2. Запросы → используем access token
3. Access истёк → используем refresh для получения нового access
4. Refresh истёк → нужен повторный login
```

**Не реализовано в проекте** (можно добавить в future enhancements)

---

## Примеры использования

### Полный flow аутентификации

```go
// 1. Регистрация/Вход
user := &domain.User{...}
token, err := jwt.Generate(user, "secret", 24*time.Hour)

// 2. Клиент получает токен
response := domain.AuthResponse{
    Token: token,
    User:  user,
}

// 3. Клиент отправляет токен в последующих запросах
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

// 4. Сервер валидирует токен
claims, err := jwt.Validate(tokenString, "secret")

// 5. Используем данные из токена
userID := claims.UserID
```

---

## См. также

- [Gin Framework](./GIN.md) - Получение токена из headers
- [Bcrypt](./BCRYPT.md) - Хеширование паролей
- [Architecture Guide](../ARCHITECTURE.md) - Auth flow
- [API Documentation](../API.md) - Использование JWT в API

