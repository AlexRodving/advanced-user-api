package handler

import (
	"net/http"

	"advanced-user-api/internal/domain"
	"advanced-user-api/internal/middleware"
	"advanced-user-api/internal/service"

	"github.com/gin-gonic/gin"
)

// ================================================================
// AUTH HANDLER - HTTP обработчики для аутентификации
// ================================================================

// AuthHandler - структура для обработки auth запросов
type AuthHandler struct {
	authService service.AuthService // Зависимость от Auth Service
	userService service.UserService // Зависимость от User Service (для /me)
}

// NewAuthHandler - конструктор
func NewAuthHandler(authService service.AuthService, userService service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// ================================================================
// REGISTER - POST /auth/register
// ================================================================

// Register обрабатывает регистрацию нового пользователя
// Endpoint: POST /api/v1/auth/register
// Body: {"email": "...", "name": "...", "password": "..."}
// Response: {"token": "...", "user": {...}}
func (h *AuthHandler) Register(c *gin.Context) {
	// === ШАГ 1: ПАРСИНГ И ВАЛИДАЦИЯ JSON ===
	// Создаём структуру для приёма данных
	var req domain.RegisterRequest

	// c.ShouldBindJSON() - парсит JSON и валидирует по binding тегам
	// Автоматически проверяет:
	//   - required: поля обязательны
	//   - email: валидный email
	//   - min=6: минимальная длина пароля
	if err := c.ShouldBindJSON(&req); err != nil {
		// Если валидация не прошла - возвращаем 400 Bad Request
		// err.Error() содержит описание ошибки валидации
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// === ШАГ 2: ВЫЗОВ SERVICE ===
	// Передаём данные в Auth Service для регистрации
	// Service:
	//   - Проверит уникальность email
	//   - Захеширует пароль
	//   - Создаст пользователя в БД
	//   - Сгенерирует JWT токен
	authResponse, err := h.authService.Register(&req)
	if err != nil {
		// Ошибка регистрации (email уже существует, ошибка БД, etc.)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// === ШАГ 3: ОТПРАВКА ОТВЕТА ===
	// Возвращаем 201 Created с токеном и данными пользователя
	c.JSON(http.StatusCreated, authResponse)
}

// ================================================================
// LOGIN - POST /auth/login
// ================================================================

// Login обрабатывает вход пользователя
// Endpoint: POST /api/v1/auth/login
// Body: {"email": "...", "password": "..."}
// Response: {"token": "...", "user": {...}}
func (h *AuthHandler) Login(c *gin.Context) {
	// === ШАГ 1: ПАРСИНГ И ВАЛИДАЦИЯ ===
	var req domain.LoginRequest

	// Парсим и валидируем JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// === ШАГ 2: АУТЕНТИФИКАЦИЯ ===
	// Service:
	//   - Найдёт пользователя по email
	//   - Проверит пароль (bcrypt)
	//   - Сгенерирует JWT токен
	authResponse, err := h.authService.Login(&req)
	if err != nil {
		// Неверный email или пароль
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// === ШАГ 3: ОТПРАВКА ОТВЕТА ===
	// Возвращаем 200 OK с токеном
	c.JSON(http.StatusOK, authResponse)
}

// ================================================================
// ME - GET /auth/me (защищённый endpoint)
// ================================================================

// Me возвращает данные текущего аутентифицированного пользователя
// Endpoint: GET /api/v1/auth/me
// Headers: Authorization: Bearer TOKEN
// Response: {"id": 1, "email": "...", "name": "..."}
//
// Этот endpoint ТРЕБУЕТ JWT токен (защищён AuthMiddleware)
func (h *AuthHandler) Me(c *gin.Context) {
	// === ШАГ 1: ПОЛУЧЕНИЕ ID ИЗ КОНТЕКСТА ===
	// AuthMiddleware уже проверил токен и сохранил userID в контекст
	// Извлекаем ID текущего пользователя
	userID := middleware.GetUserIDFromContext(c)

	// Проверяем, что ID получен
	if userID == 0 {
		// Если userID = 0, значит middleware не установил его
		// Это не должно произойти, если AuthMiddleware работает правильно
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "не удалось определить пользователя",
		})
		return
	}

	// === ШАГ 2: ПОЛУЧЕНИЕ ДАННЫХ ПОЛЬЗОВАТЕЛЯ ===
	// Загружаем полные данные пользователя из БД
	user, err := h.userService.GetCurrentUser(userID)
	if err != nil {
		// Пользователь не найден (маловероятно, но возможно если удалён)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "пользователь не найден",
		})
		return
	}

	// === ШАГ 3: ОТПРАВКА ОТВЕТА ===
	// Возвращаем данные пользователя (без пароля - json:"-")
	c.JSON(http.StatusOK, user)
}
