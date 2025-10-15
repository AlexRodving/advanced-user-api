package domain

import (
	"time"

	"gorm.io/gorm" // GORM ORM библиотека
)

// ================================================================
// DOMAIN MODEL - Основная модель пользователя
// ================================================================

// User - модель пользователя для GORM
// Эта структура определяет:
// 1. Как данные хранятся в БД (gorm теги)
// 2. Как данные отправляются в JSON (json теги)
type User struct {
	// ID - уникальный идентификатор пользователя
	// uint - беззнаковое целое число (1, 2, 3...)
	// gorm:"primaryKey" - это PRIMARY KEY в БД (уникальный, автоинкремент)
	// json:"id" - в JSON будет поле "id"
	ID uint `gorm:"primaryKey" json:"id"`

	// Email - электронная почта пользователя
	// gorm:"uniqueIndex" - создаёт уникальный индекс (два одинаковых email быть не может)
	// gorm:"not null" - поле обязательно (не может быть NULL в БД)
	// json:"email" - в JSON будет поле "email"
	Email string `gorm:"uniqueIndex;not null" json:"email"`

	// Name - имя пользователя
	// gorm:"not null" - обязательное поле
	// json:"name" - в JSON будет поле "name"
	Name string `gorm:"not null" json:"name"`

	// Password - хеш пароля (НЕ сам пароль!)
	// gorm:"not null" - обязательное поле
	// json:"-" - ВАЖНО! Пароль НЕ отправляется в JSON (безопасность!)
	Password string `gorm:"not null" json:"-"`

	// Role - роль пользователя (user, admin, etc.)
	// gorm:"default:'user'" - по умолчанию "user" при создании
	// json:"role" - в JSON будет поле "role"
	Role string `gorm:"default:'user'" json:"role"`

	// CreatedAt - время создания записи
	// GORM автоматически устанавливает при Create()
	// json:"created_at" - в JSON будет поле "created_at"
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt - время последнего обновления
	// GORM автоматически обновляет при Save() или Update()
	// json:"updated_at" - в JSON будет поле "updated_at"
	UpdatedAt time.Time `json:"updated_at"`

	// DeletedAt - время "мягкого" удаления
	// gorm.DeletedAt - специальный тип для soft delete
	// gorm:"index" - создать индекс для быстрого поиска
	// json:"-" - не показывать в JSON
	// Когда вызываем Delete(), GORM не удаляет строку физически,
	// а просто устанавливает DeletedAt = текущее время
	// Все последующие запросы автоматически игнорируют "удалённые" строки
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName - переопределение имени таблицы в БД
// По умолчанию GORM использует множественное число от имени структуры: User → "users"
// Эта функция явно указывает имя таблицы (опционально, если нужно другое имя)
func (User) TableName() string {
	return "users"
}

// ================================================================
// DTO (Data Transfer Objects) - структуры для HTTP запросов/ответов
// ================================================================
// DTO используются для:
// 1. Приёма данных от клиента (Request)
// 2. Отправки данных клиенту (Response)
// Это отделяет внутреннюю модель БД от внешнего API

// RegisterRequest - данные для регистрации нового пользователя
type RegisterRequest struct {
	// Email - электронная почта
	// json:"email" - название поля в JSON
	// binding:"required,email" - валидация Gin:
	//   - required: поле обязательно
	//   - email: должно быть валидным email адресом
	Email string `json:"email" binding:"required,email"`

	// Name - имя пользователя
	// binding:"required,min=2" - валидация:
	//   - required: обязательно
	//   - min=2: минимум 2 символа
	Name string `json:"name" binding:"required,min=2"`

	// Password - пароль (будет хеширован перед сохранением)
	// binding:"required,min=6" - минимум 6 символов
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest - данные для входа (аутентификации)
type LoginRequest struct {
	// Email - для идентификации пользователя
	// binding:"required,email" - обязательно и валидный email
	Email string `json:"email" binding:"required,email"`

	// Password - пароль для проверки
	// binding:"required" - обязательно
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest - данные для обновления профиля пользователя
type UpdateUserRequest struct {
	// Name - новое имя
	// binding:"omitempty,min=2" - валидация:
	//   - omitempty: поле опционально (может не присутствовать)
	//   - min=2: если присутствует, то минимум 2 символа
	Name string `json:"name" binding:"omitempty,min=2"`

	// Email - новый email
	// binding:"omitempty,email" - опционально, но если есть - валидный email
	Email string `json:"email" binding:"omitempty,email"`
}

// AuthResponse - ответ после успешной регистрации или входа
// Содержит JWT токен и данные пользователя
type AuthResponse struct {
	// Token - JWT токен для аутентификации последующих запросов
	// Клиент должен отправлять этот токен в заголовке Authorization
	Token string `json:"token"`

	// User - данные пользователя (без пароля!)
	// Указатель *User позволяет вернуть nil если нужно
	User *User `json:"user"`
}
