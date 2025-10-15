# 🔒 Bcrypt & Viper

## Bcrypt - Password Hashing

### Overview

Bcrypt - алгоритм хеширования паролей с автоматической солью.

**Пакет:** `golang.org/x/crypto/bcrypt`

**Версия в проекте:** v0.18.0

---

### Основные методы

#### 1. bcrypt.GenerateFromPassword()

**Описание:** Хеширует пароль с солью

**Сигнатура:**
```go
func GenerateFromPassword(password []byte, cost int) ([]byte, error)
```

**Параметры:**
- `password []byte` - пароль для хеширования
- `cost int` - сложность (4-31, рекомендуется 10-14)

**Возвращает:**
- `[]byte` - хеш пароля (60 символов)
- `error` - ошибка хеширования

**Пример использования:**
```go
password := "mySecretPassword123"
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    bcrypt.DefaultCost,  // 10
)
if err != nil {
    return err
}

// hashedPassword: "$2a$10$N9qo8uLOickgx2ZMRZoMye.6IrYtIB7LhGbp3bLMqGPH..."
```

**Где используем:**
- `internal/pkg/password/password.go:16`

**Cost values:**
```go
bcrypt.MinCost     = 4
bcrypt.MaxCost     = 31
bcrypt.DefaultCost = 10  // ✅ Рекомендуется
```

**Время хеширования:**
- Cost 10: ~100ms
- Cost 12: ~400ms
- Cost 14: ~1.6s

---

#### 2. bcrypt.CompareHashAndPassword()

**Описание:** Сравнивает хеш и пароль

**Сигнатура:**
```go
func CompareHashAndPassword(hashedPassword, password []byte) error
```

**Параметры:**
- `hashedPassword []byte` - хеш из БД
- `password []byte` - пароль от пользователя

**Возвращает:**
- `nil` - пароли совпадают
- `error` - пароли не совпадают или ошибка

**Пример использования:**
```go
// Получаем хеш из БД
user := &domain.User{
    Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye...",
}

// Проверяем пароль
err := bcrypt.CompareHashAndPassword(
    []byte(user.Password),
    []byte("userInputPassword"),
)

if err != nil {
    // Пароль неверный
    return errors.New("invalid password")
}

// Пароль верный
```

**Где используем:**
- `internal/pkg/password/password.go:25`

---

### Наша реализация

#### Hash() - Хеширование пароля

**Файл:** `internal/pkg/password/password.go`

**Сигнатура:**
```go
func Hash(password string) (string, error)
```

**Параметры:**
- `password string` - пароль в открытом виде

**Возвращает:**
- `string` - хешированный пароль
- `error` - ошибка хеширования

**Пример:**
```go
hashedPassword, err := password.Hash("myPassword123")
if err != nil {
    return err
}

// Сохраняем в БД
user.Password = hashedPassword
```

**Где используем:**
- `internal/service/auth_service.go:67` - при регистрации

---

#### Verify() - Проверка пароля

**Сигнатура:**
```go
func Verify(hashedPassword, password string) error
```

**Параметры:**
- `hashedPassword string` - хеш из БД
- `password string` - пароль от пользователя

**Возвращает:**
- `nil` - пароль верный
- `error` - пароль неверный

**Пример:**
```go
// Получаем пользователя из БД
user, _ := repo.FindByEmail(email)

// Проверяем пароль
err := password.Verify(user.Password, requestPassword)
if err != nil {
    return errors.New("invalid credentials")
}

// Пароль верный, продолжаем
```

**Где используем:**
- `internal/service/auth_service.go:120` - при входе

---

### Security Best Practices

#### 1. Никогда не храните пароли в открытом виде

```go
// ❌ НИКОГДА!
user.Password = "password123"
db.Create(&user)

// ✅ Всегда хешируйте
hashedPassword, _ := password.Hash("password123")
user.Password = hashedPassword
db.Create(&user)
```

#### 2. Используйте подходящий cost

```go
// ❌ Слишком низкий (быстро взламывается)
bcrypt.GenerateFromPassword(password, 4)

// ✅ Оптимально для production
bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)  // 10

// ✅ Для высокой безопасности
bcrypt.GenerateFromPassword(password, 12)
```

#### 3. Не возвращайте пароль в API

```go
type User struct {
    ID       uint   `json:"id"`
    Email    string `json:"email"`
    Password string `json:"-"`  // ✅ НЕ возвращаем в JSON!
}
```

**Где используем:**
- `internal/domain/user.go:25` - `json:"-"` тег

#### 4. Одинаковые ошибки для несуществующего пользователя и неверного пароля

```go
// ❌ Атакующий узнаёт, что email существует
user, err := repo.FindByEmail(email)
if err != nil {
    return errors.New("user not found")  // Плохо!
}
if !password.Verify(user.Password, pass) {
    return errors.New("invalid password")  // Плохо!
}

// ✅ Одинаковое сообщение
user, err := repo.FindByEmail(email)
if err != nil || !password.Verify(user.Password, pass) {
    return errors.New("invalid credentials")  // Хорошо!
}
```

**Где используем:**
- `internal/service/auth_service.go:122` - общая ошибка "неверные данные"

---

## Viper - Configuration Management

### Overview

Viper - библиотека для управления конфигурацией приложения.

**Официальная документация:** https://github.com/spf13/viper

**Версия в проекте:** v1.21.0

---

### Основные методы

#### 1. viper.SetConfigFile()

**Описание:** Устанавливает путь к конфигурационному файлу

**Сигнатура:**
```go
func SetConfigFile(in string)
```

**Параметры:**
- `in string` - путь к файлу (например, ".env")

**Пример:**
```go
viper.SetConfigFile(".env")
```

**Где используем:**
- `internal/config/config.go:42`

---

#### 2. viper.ReadInConfig()

**Описание:** Читает конфигурацию из файла

**Сигнатура:**
```go
func ReadInConfig() error
```

**Возвращает:** `error` - ошибка чтения файла

**Пример:**
```go
if err := viper.ReadInConfig(); err != nil {
    log.Println("Config file not found, using env variables")
}
```

**Где используем:**
- `internal/config/config.go:45`

---

#### 3. viper.AutomaticEnv()

**Описание:** Автоматически читает environment variables

**Сигнатура:**
```go
func AutomaticEnv()
```

**Пример:**
```go
viper.AutomaticEnv()

// Теперь можно читать:
// export DB_HOST=localhost
dbHost := viper.GetString("DB_HOST")
```

**Где используем:**
- `internal/config/config.go:53`

---

#### 4. viper.GetString()

**Описание:** Получает строковое значение конфигурации

**Сигнатура:**
```go
func GetString(key string) string
```

**Параметры:**
- `key string` - ключ конфигурации

**Возвращает:** `string` - значение

**Пример:**
```go
dbHost := viper.GetString("DB_HOST")
```

**Где используем:**
- `internal/config/config.go:58-66`

---

#### 5. viper.SetDefault()

**Описание:** Устанавливает значение по умолчанию

**Сигнатура:**
```go
func SetDefault(key string, value interface{})
```

**Параметры:**
- `key string` - ключ
- `value interface{}` - значение по умолчанию

**Пример:**
```go
viper.SetDefault("SERVER_PORT", "8080")
viper.SetDefault("JWT_EXPIRATION", "24h")

// Если не установлено в .env или env variables, используется default
port := viper.GetString("SERVER_PORT")  // "8080"
```

**Где используем:**
- `internal/config/config.go:70-76`

---

### Приоритет источников конфигурации

Viper использует следующий порядок приоритета (от высшего к низшему):

1. **Set()** - явная установка в коде
2. **Environment variables** - переменные окружения
3. **Config file** (.env)
4. **Defaults** - значения по умолчанию

```go
// Установлены все источники:
viper.SetDefault("DB_HOST", "default.com")    // Приоритет 4
// .env: DB_HOST=env-file.com                 // Приоритет 3
// export DB_HOST=environment.com             // Приоритет 2
viper.Set("DB_HOST", "code.com")              // Приоритет 1

result := viper.GetString("DB_HOST")  // "code.com" (самый высокий приоритет)
```

---

### Наша конфигурация

**Файл:** `internal/config/config.go`

#### Структура Config

```go
type Config struct {
    // Server
    ServerPort string
    GinMode    string
    
    // Database
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    
    // JWT
    JWTSecret     string
    JWTExpiration string
    
    // Redis
    RedisHost string
    RedisPort string
}
```

#### Load() - Загрузка конфигурации

**Что делает:**
1. Пытается прочитать `.env` файл
2. Включает чтение environment variables
3. Устанавливает defaults
4. Возвращает заполненную структуру Config

**Пример использования:**
```go
cfg := config.Load()

fmt.Println(cfg.ServerPort)   // "8080"
fmt.Println(cfg.DBHost)       // "localhost"
fmt.Println(cfg.JWTSecret)    // "your-secret"
```

**Где используем:**
- `cmd/api/main.go:21` - при запуске приложения

---

### Формат .env файла

```env
# Server
SERVER_PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=advanced_api

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

**Пример:**
- `env.example` в корне проекта

---

### Best Practices

#### 1. Никогда не коммитьте .env с секретами

```gitignore
# .gitignore
.env
```

#### 2. Предоставляйте .env.example

```bash
cp env.example .env
# Затем редактируйте .env с реальными значениями
```

#### 3. Используйте осмысленные defaults

```go
viper.SetDefault("SERVER_PORT", "8080")      // ✅
viper.SetDefault("JWT_EXPIRATION", "24h")    // ✅
viper.SetDefault("DB_HOST", "localhost")     // ✅
```

#### 4. Валидируйте критичные настройки

```go
if cfg.JWTSecret == "" {
    log.Fatal("JWT_SECRET must be set")
}
```

---

## См. также

- [JWT](./JWT.md) - Использование JWTSecret из конфигурации
- [GORM](./GORM.md) - Использование DB настроек
- [Architecture Guide](../ARCHITECTURE.md) - Конфигурация в архитектуре

