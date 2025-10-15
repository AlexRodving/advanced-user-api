package config

import (
	"log"

	"github.com/spf13/viper" // Viper - библиотека для работы с конфигурацией
)

// ================================================================
// CONFIG - Конфигурация приложения
// ================================================================

// Config - структура, содержащая все настройки приложения
// Настройки загружаются из:
// 1. .env файла (если есть)
// 2. Environment variables (переменные окружения)
// 3. Defaults (значения по умолчанию)
type Config struct {
	// === DATABASE SETTINGS ===
	// Настройки подключения к PostgreSQL
	
	// DBHost - адрес сервера БД (localhost, db.example.com, etc.)
	// mapstructure:"DB_HOST" - имя переменной окружения
	DBHost string `mapstructure:"DB_HOST"`
	
	// DBPort - порт PostgreSQL (обычно 5432)
	DBPort string `mapstructure:"DB_PORT"`
	
	// DBUser - имя пользователя БД
	DBUser string `mapstructure:"DB_USER"`
	
	// DBPassword - пароль для подключения к БД
	DBPassword string `mapstructure:"DB_PASSWORD"`
	
	// DBName - имя базы данных
	DBName string `mapstructure:"DB_NAME"`

	// === REDIS SETTINGS ===
	// Настройки для Redis (кеш, сессии)
	
	// RedisHost - адрес Redis сервера
	RedisHost string `mapstructure:"REDIS_HOST"`
	
	// RedisPort - порт Redis (обычно 6379)
	RedisPort string `mapstructure:"REDIS_PORT"`

	// === JWT SETTINGS ===
	// Настройки для JSON Web Tokens (аутентификация)
	
	// JWTSecret - секретный ключ для подписи токенов
	// ВАЖНО: В production используйте длинную случайную строку!
	JWTSecret string `mapstructure:"JWT_SECRET"`
	
	// JWTExpiration - время жизни токена (например, "24h", "7d")
	JWTExpiration string `mapstructure:"JWT_EXPIRATION"`

	// === SERVER SETTINGS ===
	// Настройки HTTP сервера
	
	// ServerPort - порт, на котором запускается API (например, "8080")
	ServerPort string `mapstructure:"SERVER_PORT"`
	
	// GinMode - режим работы Gin ("debug", "release", "test")
	// debug - подробные логи, release - production режим
	GinMode string `mapstructure:"GIN_MODE"`

	// === LOGGING SETTINGS ===
	// Настройки логирования
	
	// LogLevel - уровень логирования ("debug", "info", "warn", "error")
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

// ================================================================
// LOAD CONFIGURATION
// ================================================================

// Load загружает конфигурацию из .env файла и environment variables
// Приоритет (от высшего к низшему):
// 1. Environment variables (переменные окружения)
// 2. .env файл
// 3. Defaults (значения по умолчанию)
func Load() *Config {
	// === ШАГ 1: НАСТРОЙКА VIPER ===
	
	// SetConfigFile - указываем путь к .env файлу
	// Viper будет искать файл .env в текущей директории
	viper.SetConfigFile(".env")

	// === ШАГ 2: УСТАНОВКА ЗНАЧЕНИЙ ПО УМОЛЧАНИЮ ===
	// Если переменная не найдена ни в .env, ни в environment - используется default
	
	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "advanced_api")
	
	// Redis defaults
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	
	// JWT defaults
	viper.SetDefault("JWT_SECRET", "change-this-secret-in-production")
	viper.SetDefault("JWT_EXPIRATION", "24h")
	
	// Server defaults
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("GIN_MODE", "debug")
	
	// Logging defaults
	viper.SetDefault("LOG_LEVEL", "debug")

	// === ШАГ 3: ЧТЕНИЕ ENVIRONMENT VARIABLES ===
	// AutomaticEnv() - автоматически читает переменные окружения
	// Например, если установлена DB_HOST=production.db.com,
	// она перезапишет значение из .env и default
	viper.AutomaticEnv()

	// === ШАГ 4: ЧТЕНИЕ .env ФАЙЛА ===
	// Пытаемся прочитать .env файл
	// Если файла нет - не проблема, используем environment variables и defaults
	if err := viper.ReadInConfig(); err != nil {
		// Файл не найден - это нормально для production (там используются env vars)
		log.Println("⚠️  .env файл не найден, используем environment variables и defaults")
	} else {
		// Файл найден и прочитан
		log.Println("✅ Конфигурация загружена из .env")
	}

	// === ШАГ 5: UNMARSHAL В СТРУКТУРУ ===
	// Viper автоматически заполняет структуру Config значениями
	// Использует mapstructure теги для сопоставления
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		// Критическая ошибка - не можем загрузить конфигурацию
		log.Fatal("❌ Ошибка чтения конфигурации:", err)
	}

	// Возвращаем заполненную конфигурацию
	return &cfg
}

