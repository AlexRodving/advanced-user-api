# 🎯 Этап 1: GORM Модели и настройка

## 📋 Задание

Создайте модели данных с использованием GORM и настройте подключение к БД.

---

## ✅ Задача 1: Создайте GORM модель User

Файл: `internal/domain/user.go`

```go
package domain

import (
	"time"
	"gorm.io/gorm"
)

// User - модель пользователя для GORM
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Name      string         `gorm:"not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"` // "-" = не показывать в JSON
	Role      string         `gorm:"default:'user'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// TableName - имя таблицы в БД (опционально, по умолчанию "users")
func (User) TableName() string {
	return "users"
}

// DTO (Data Transfer Objects) - для запросов/ответов

// RegisterRequest - данные для регистрации
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=2"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest - данные для входа
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest - данные для обновления
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2"`
	Email string `json:"email" binding:"omitempty,email"`
}

// AuthResponse - ответ с токеном
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
```

### 📝 Что нового в GORM:

#### GORM теги:
- `gorm:"primaryKey"` - первичный ключ
- `gorm:"uniqueIndex"` - уникальный индекс
- `gorm:"not null"` - обязательное поле
- `gorm:"default:'user'"` - значение по умолчанию
- `gorm:"index"` - создать индекс

#### Gin binding теги:
- `binding:"required"` - обязательное поле
- `binding:"email"` - валидация email
- `binding:"min=6"` - минимальная длина
- `binding:"omitempty"` - может быть пустым

#### Soft Delete:
- `DeletedAt gorm.DeletedAt` - "мягкое" удаление
- Строка не удаляется физически, просто помечается
- GORM автоматически игнорирует "удалённые" строки

---

## ✅ Задача 2: Создайте Config с Viper

Файл: `internal/config/config.go`

```go
package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	// Database
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	
	// Redis
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	
	// JWT
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiration string `mapstructure:"JWT_EXPIRATION"`
	
	// Server
	ServerPort string `mapstructure:"SERVER_PORT"`
	GinMode    string `mapstructure:"GIN_MODE"`
	
	// Logging
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

// Load загружает конфигурацию из .env или environment variables
func Load() *Config {
	// TODO: Реализуйте загрузку конфигурации
	
	// 1. Установите путь к .env файлу
	viper.SetConfigFile(".env")
	
	// 2. Установите defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("GIN_MODE", "debug")
	
	// 3. Читайте environment variables
	viper.AutomaticEnv()
	
	// 4. Читайте .env файл (если есть)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️  .env файл не найден, используем environment variables")
	}
	
	// 5. Unmarshal в структуру Config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Ошибка чтения конфигурации:", err)
	}
	
	return &cfg
}
```

---

## ✅ Задача 3: Подключение к БД с GORM

Файл: `internal/repository/database.go`

```go
package repository

import (
	"fmt"
	"log"
	
	"advanced-user-api/internal/config"
	"advanced-user-api/internal/domain"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB инициализирует подключение к БД с GORM
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// TODO: Реализуйте подключение
	
	// 1. Формируем DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	
	// 2. Подключаемся через GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Логирование SQL запросов
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// 3. Auto Migration - автоматическое создание таблиц!
	// GORM сам создаст таблицу на основе struct
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}
	
	log.Println("✅ База данных подключена и migrations применены")
	
	return db, nil
}
```

### 🔥 Преимущества GORM:

**1. Auto Migration:**
```go
// Вместо SQL миграций вручную:
db.AutoMigrate(&User{})
// GORM сам создаёт таблицу!
```

**2. Автоматические timestamps:**
```go
// CreatedAt и UpdatedAt обновляются автоматически!
```

**3. Soft Delete:**
```go
// db.Delete(&user) - не удаляет физически
// Просто устанавливает DeletedAt
```

---

## ✅ Задача 4: User Repository с GORM

Файл: `internal/repository/user_repository.go`

```go
package repository

import (
	"advanced-user-api/internal/domain"
	"gorm.io/gorm"
)

// UserRepository - интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}

// userRepository - реализация с GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository - конструктор
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create - создаёт пользователя
func (r *userRepository) Create(user *domain.User) error {
	// TODO: Реализуйте
	// ПОДСКАЗКА: return r.db.Create(user).Error
	
	return nil
}

// FindByID - ищет по ID
func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	// TODO: Реализуйте
	// ПОДСКАЗКА: 
	// var user domain.User
	// err := r.db.First(&user, id).Error
	// return &user, err
	
	return nil, nil
}

// FindByEmail - ищет по email
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	// TODO: Реализуйте
	// ПОДСКАЗКА: r.db.Where("email = ?", email).First(&user).Error
	
	return nil, nil
}

// FindAll - получает всех пользователей
func (r *userRepository) FindAll() ([]domain.User, error) {
	// TODO: Реализуйте
	// ПОДСКАЗКА: r.db.Find(&users).Error
	
	return nil, nil
}

// Update - обновляет пользователя
func (r *userRepository) Update(user *domain.User) error {
	// TODO: Реализуйте
	// ПОДСКАЗКА: r.db.Save(user).Error
	
	return nil
}

// Delete - удаляет пользователя (soft delete)
func (r *userRepository) Delete(id uint) error {
	// TODO: Реализуйте
	// ПОДСКАЗКА: r.db.Delete(&domain.User{}, id).Error
	
	return nil
}
```

### 🎯 Сравнение: SQL vs GORM

| Операция | database/sql | GORM |
|----------|--------------|------|
| **Создать** | `db.QueryRow("INSERT...")` | `db.Create(&user)` |
| **Найти по ID** | `db.QueryRow("SELECT... WHERE id=?")` | `db.First(&user, id)` |
| **Найти все** | `db.Query("SELECT...")` | `db.Find(&users)` |
| **Обновить** | `db.Exec("UPDATE...")` | `db.Save(&user)` |
| **Удалить** | `db.Exec("DELETE...")` | `db.Delete(&user)` |

**GORM проще!** Но нужно понимать, какой SQL генерируется.

---

## 🧪 Тестирование

После реализации запустите:

```bash
# 1. Установите зависимости
go mod tidy

# 2. Запустите PostgreSQL
docker-compose up -d postgres

# 3. Создайте тестовый файл для проверки
```

---

## 📝 Подсказки для реализации

### Create:
```go
return r.db.Create(user).Error
```

### Find:
```go
var user domain.User
err := r.db.First(&user, id).Error
if err == gorm.ErrRecordNotFound {
    return nil, errors.New("пользователь не найден")
}
return &user, err
```

### Update:
```go
return r.db.Save(user).Error
```

### Delete (soft delete):
```go
return r.db.Delete(&domain.User{}, id).Error
```

---

## 🎓 Критерии проверки

Я проверю:
- ✅ GORM модель создана с правильными тегами
- ✅ Config загружается через Viper
- ✅ БД подключается и Auto Migration работает
- ✅ Все Repository методы реализованы
- ✅ Код читаемый с комментариями

---

**Создайте эти файлы и реализуйте TODO методы!**

Когда закончите - напишите **"готово"**! 💪

