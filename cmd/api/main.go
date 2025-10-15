package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"advanced-user-api/internal/config"
	"advanced-user-api/internal/handler"
	"advanced-user-api/internal/repository"
	"advanced-user-api/internal/service"

	"github.com/gin-gonic/gin"
)

// ================================================================
// MAIN - Точка входа в приложение
// ================================================================

func main() {
	fmt.Println("=== 🚀 Advanced User Management API ===\n")

	// === ШАГ 1: ЗАГРУЗКА КОНФИГУРАЦИИ ===
	// Загружаем настройки из .env и environment variables
	cfg := config.Load()
	log.Println("✅ Конфигурация загружена")

	// === ШАГ 2: ПОДКЛЮЧЕНИЕ К БД ===
	// InitDB() подключается к PostgreSQL и выполняет Auto Migration
	// GORM автоматически создаст таблицу users на основе struct
	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}
	
	// defer - закроем подключение при завершении main()
	defer func() {
		if err := repository.CloseDB(db); err != nil {
			log.Println("❌ Ошибка закрытия БД:", err)
		}
	}()

	// === ШАГ 3: СОЗДАНИЕ СЛОЁВ (Dependency Injection) ===
	// Создаём слои приложения снизу вверх
	
	// 3.1: Repository (работа с БД)
	userRepo := repository.NewUserRepository(db)
	
	// 3.2: Services (бизнес-логика)
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	
	// 3.3: Handlers (HTTP обработчики)
	authHandler := handler.NewAuthHandler(authService, userService)
	userHandler := handler.NewUserHandler(userService)
	
	log.Println("✅ Все слои приложения инициализированы")

	// === ШАГ 4: НАСТРОЙКА GIN ===
	// Устанавливаем режим Gin (debug, release, test)
	gin.SetMode(cfg.GinMode)
	
	// Создаём новый Gin роутер
	// gin.Default() включает:
	//   - Logger middleware (логирование запросов)
	//   - Recovery middleware (восстановление после panic)
	router := gin.Default()
	
	log.Println("✅ Gin роутер создан")

	// === ШАГ 5: РЕГИСТРАЦИЯ МАРШРУТОВ ===
	// Настраиваем все HTTP endpoints
	handler.SetupRoutes(router, authHandler, userHandler, cfg)
	log.Println("✅ Маршруты зарегистрированы")

	// === ШАГ 6: СОЗДАНИЕ HTTP СЕРВЕРА ===
	// Создаём HTTP сервер с настройками
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,           // Адрес и порт (например, ":8080")
		Handler: router,                          // Gin роутер как handler
		
		// Таймауты для защиты от медленных клиентов
		ReadTimeout:  10 * time.Second,  // Максимальное время чтения запроса
		WriteTimeout: 10 * time.Second,  // Максимальное время записи ответа
		IdleTimeout:  120 * time.Second, // Максимальное время простоя соединения
	}

	// === ШАГ 7: ЗАПУСК СЕРВЕРА В ГОРУТИНЕ ===
	// Запускаем сервер в отдельной горутине
	// Это позволяет обрабатывать graceful shutdown
	go func() {
		fmt.Printf("\n🚀 Сервер запущен на http://localhost:%s\n\n", cfg.ServerPort)
		fmt.Println("📡 Доступные endpoints:")
		fmt.Println("   PUBLIC (без токена):")
		fmt.Println("     POST   /api/v1/auth/register  - Регистрация")
		fmt.Println("     POST   /api/v1/auth/login     - Вход")
		fmt.Println("     GET    /health                - Health check")
		fmt.Println("\n   PROTECTED (требуют JWT токен):")
		fmt.Println("     GET    /api/v1/auth/me        - Текущий пользователь")
		fmt.Println("     GET    /api/v1/users          - Список пользователей")
		fmt.Println("     GET    /api/v1/users/:id      - Получить пользователя")
		fmt.Println("     PUT    /api/v1/users/:id      - Обновить пользователя")
		fmt.Println("     DELETE /api/v1/users/:id      - Удалить пользователя")
		fmt.Println("\n💡 Нажмите Ctrl+C для остановки\n")
		
		// ListenAndServe() - запускает HTTP сервер
		// Блокирующая функция - работает до ошибки или остановки
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("❌ Ошибка запуска сервера:", err)
		}
	}()

	// === ШАГ 8: GRACEFUL SHUTDOWN ===
	// Ожидаем сигнал остановки (Ctrl+C или kill)
	// Это позволяет корректно завершить все запросы перед остановкой
	
	// Создаём канал для получения сигналов ОС
	quit := make(chan os.Signal, 1)
	
	// signal.Notify() - подписываемся на сигналы
	// SIGINT - Ctrl+C
	// SIGTERM - kill команда
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// Блокируемся до получения сигнала
	<-quit
	log.Println("\n🛑 Получен сигнал остановки...")

	// Создаём контекст с таймаутом для graceful shutdown
	// Даём серверу 5 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown() - корректно останавливает сервер
	// Ждёт завершения всех активных запросов (но не более 5 секунд)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("❌ Ошибка при остановке сервера:", err)
	}

	log.Println("✅ Сервер корректно остановлен")
}


