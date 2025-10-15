package repository

import (
	"fmt"
	"log"

	"advanced-user-api/internal/config"
	"advanced-user-api/internal/domain"

	"gorm.io/driver/postgres" // PostgreSQL драйвер для GORM
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Логирование SQL запросов
)

// ================================================================
// DATABASE CONNECTION - Подключение к PostgreSQL через GORM
// ================================================================

// InitDB инициализирует подключение к PostgreSQL базе данных
// Параметры:
//   - cfg: конфигурация с настройками подключения
// Возвращает:
//   - *gorm.DB: подключение к БД
//   - error: ошибка подключения (если есть)
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// === ШАГ 1: ФОРМИРОВАНИЕ DSN (Data Source Name) ===
	// DSN - строка подключения к PostgreSQL
	// Формат: "host=localhost port=5432 user=postgres password=secret dbname=mydb sslmode=disable"
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,     // Адрес сервера БД
		cfg.DBPort,     // Порт (обычно 5432)
		cfg.DBUser,     // Имя пользователя
		cfg.DBPassword, // Пароль
		cfg.DBName,     // Имя базы данных
	)

	// === ШАГ 2: ПОДКЛЮЧЕНИЕ К БД ЧЕРЕЗ GORM ===
	// gorm.Open() - открывает подключение к БД
	// Параметры:
	//   - postgres.Open(dsn): PostgreSQL драйвер с DSN
	//   - &gorm.Config{...}: настройки GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger - настройка логирования SQL запросов
		// logger.Info - логировать все SQL запросы (полезно для отладки)
		// В production лучше использовать logger.Warn или logger.Error
		Logger: logger.Default.LogMode(logger.Info),
		
		// Другие полезные настройки (опционально):
		// NowFunc: func() time.Time { return time.Now().UTC() }, // Использовать UTC время
		// PrepareStmt: true, // Кешировать prepared statements (быстрее)
	})
	
	// Проверяем ошибку подключения
	if err != nil {
		// Если не удалось подключиться - возвращаем ошибку с контекстом
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// === ШАГ 3: AUTO MIGRATION ===
	// AutoMigrate - автоматическое создание/обновление таблиц
	// GORM анализирует структуру User и:
	// 1. Создаёт таблицу users (если не существует)
	// 2. Добавляет недостающие колонки (если структура изменилась)
	// 3. Создаёт индексы (uniqueIndex, index)
	// 4. НЕ удаляет существующие колонки (безопасно)
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Логируем успешное подключение
	log.Println("✅ База данных подключена")
	log.Println("✅ Auto Migration выполнен (таблица users создана/обновлена)")

	// === ШАГ 4: НАСТРОЙКА CONNECTION POOL (опционально) ===
	// Получаем базовый sql.DB для тонкой настройки
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// SetMaxIdleConns - максимальное количество неактивных соединений в пуле
	// Неактивные соединения остаются открытыми для переиспользования
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns - максимальное количество открытых соединений
	// Ограничивает нагрузку на БД
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime - максимальное время жизни соединения
	// После этого времени соединение закрывается и создаётся новое
	// Полезно для балансировки нагрузки и обновления соединений
	// sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Connection pool настроен (10 idle, 100 max connections)")

	// Возвращаем готовое подключение к БД
	return db, nil
}

// ================================================================
// ДОПОЛНИТЕЛЬНЫЕ ФУНКЦИИ (опционально)
// ================================================================

// CloseDB - корректное закрытие подключения к БД
// Вызывайте при остановке приложения (в main с defer)
func CloseDB(db *gorm.DB) error {
	// Получаем базовый sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	
	// Закрываем все соединения
	return sqlDB.Close()
}

