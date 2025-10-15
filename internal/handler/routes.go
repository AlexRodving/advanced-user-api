package handler

import (
	"advanced-user-api/internal/config"
	"advanced-user-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// ================================================================
// ROUTES SETUP - Настройка всех маршрутов приложения
// ================================================================

// SetupRoutes настраивает все HTTP маршруты (endpoints)
// Параметры:
//   - router: Gin роутер
//   - authHandler: обработчик auth запросов
//   - userHandler: обработчик user запросов
//   - cfg: конфигурация (для JWT secret в middleware)
func SetupRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	cfg *config.Config,
) {
	// Применяем глобальные middleware
	router.Use(middleware.CORSMiddleware())
	// ================================================================
	// API VERSION 1 - Группа маршрутов /api/v1
	// ================================================================
	// Группировка позволяет:
	// 1. Версионировать API (/api/v1, /api/v2)
	// 2. Применять middleware к группе
	// 3. Организовать код
	api := router.Group("/api/v1")
	{
		// ============================================================
		// PUBLIC ROUTES - Публичные маршруты (без аутентификации)
		// ============================================================
		
		// --- AUTH ROUTES ---
		// Группа для аутентификации
		auth := api.Group("/auth")
		{
			// POST /api/v1/auth/register - Регистрация
			// Любой может зарегистрироваться (публичный endpoint)
			auth.POST("/register", authHandler.Register)
			
			// POST /api/v1/auth/login - Вход
			// Любой может войти (публичный endpoint)
			auth.POST("/login", authHandler.Login)
			
			// --- PROTECTED AUTH ROUTES ---
			// GET /api/v1/auth/me - Текущий пользователь
			// ТРЕБУЕТ JWT токен (защищён AuthMiddleware)
			auth.GET("/me", middleware.AuthMiddleware(cfg), authHandler.Me)
		}

		// ============================================================
		// PROTECTED ROUTES - Защищённые маршруты (требуют JWT)
		// ============================================================
		
		// --- USER ROUTES ---
		// Группа для работы с пользователями
		// ВСЕ endpoints в этой группе требуют JWT токен!
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg)) // Применяем middleware ко всей группе
		{
			// GET /api/v1/users - Список всех пользователей
			// Требует: Authorization: Bearer TOKEN
			users.GET("", userHandler.GetAll)
			
			// GET /api/v1/users/:id - Получить пользователя по ID
			// Пример: GET /api/v1/users/42
			// Требует: Authorization: Bearer TOKEN
			users.GET("/:id", userHandler.GetByID)
			
			// PUT /api/v1/users/:id - Обновить пользователя
			// Пример: PUT /api/v1/users/42
			// Body: {"name": "New Name", "email": "new@email.com"}
			// Требует: Authorization: Bearer TOKEN
			users.PUT("/:id", userHandler.Update)
			
			// DELETE /api/v1/users/:id - Удалить пользователя
			// Пример: DELETE /api/v1/users/42
			// Требует: Authorization: Bearer TOKEN
			users.DELETE("/:id", userHandler.Delete)
		}
	}

	// ================================================================
	// HEALTH CHECK - Проверка здоровья сервиса
	// ================================================================
	// GET /health - Публичный endpoint для мониторинга
	// Используется для проверки, что сервис работает
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "advanced-user-api",
		})
	})
}

// ================================================================
// ВИЗУАЛИЗАЦИЯ МАРШРУТОВ
// ================================================================
//
// PUBLIC (без токена):
//   POST   /api/v1/auth/register
//   POST   /api/v1/auth/login
//   GET    /health
//
// PROTECTED (требуют JWT токен):
//   GET    /api/v1/auth/me
//   GET    /api/v1/users
//   GET    /api/v1/users/:id
//   PUT    /api/v1/users/:id
//   DELETE /api/v1/users/:id
//
// ================================================================

