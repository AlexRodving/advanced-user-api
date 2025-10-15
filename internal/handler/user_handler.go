package handler

import (
	"net/http"
	"strconv"

	"advanced-user-api/internal/domain"
	"advanced-user-api/internal/service"

	"github.com/gin-gonic/gin"
)

// ================================================================
// USER HANDLER - HTTP обработчики для пользователей
// ================================================================

// UserHandler - структура для обработки user запросов
type UserHandler struct {
	userService service.UserService // Зависимость от User Service
}

// NewUserHandler - конструктор
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// ================================================================
// HANDLERS - Обработчики HTTP запросов
// ================================================================

// GetAll получает список всех пользователей
// Endpoint: GET /api/v1/users
// Headers: Authorization: Bearer TOKEN (защищён!)
// Response: [{"id": 1, "email": "...", "name": "..."}, ...]
func (h *UserHandler) GetAll(c *gin.Context) {
	// === ШАГ 1: ВЫЗОВ SERVICE ===
	// Получаем всех пользователей из service
	users, err := h.userService.GetAllUsers()
	if err != nil {
		// Ошибка БД
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка получения пользователей",
		})
		return
	}

	// === ШАГ 2: ОТПРАВКА ОТВЕТА ===
	// c.JSON() - автоматически устанавливает Content-Type: application/json
	// и кодирует данные в JSON
	c.JSON(http.StatusOK, users)
}

// GetByID получает одного пользователя по ID
// Endpoint: GET /api/v1/users/:id
// Headers: Authorization: Bearer TOKEN (защищён!)
// Response: {"id": 1, "email": "...", "name": "..."}
func (h *UserHandler) GetByID(c *gin.Context) {
	// === ШАГ 1: ИЗВЛЕЧЕНИЕ ID ИЗ URL ===
	// c.Param("id") - получает параметр из URL
	// Пример: /users/42 → c.Param("id") = "42"
	idStr := c.Param("id")
	
	// Конвертируем строку в uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// ID невалиден (например, /users/abc)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "невалидный ID",
		})
		return
	}

	// === ШАГ 2: ВЫЗОВ SERVICE ===
	// Получаем пользователя по ID
	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		// Пользователь не найден
		c.JSON(http.StatusNotFound, gin.H{
			"error": "пользователь не найден",
		})
		return
	}

	// === ШАГ 3: ОТПРАВКА ОТВЕТА ===
	c.JSON(http.StatusOK, user)
}

// Update обновляет данные пользователя
// Endpoint: PUT /api/v1/users/:id
// Headers: Authorization: Bearer TOKEN (защищён!)
// Body: {"name": "...", "email": "..."}
// Response: {"id": 1, "email": "...", "name": "..."}
func (h *UserHandler) Update(c *gin.Context) {
	// === ШАГ 1: ИЗВЛЕЧЕНИЕ ID ===
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "невалидный ID",
		})
		return
	}

	// === ШАГ 2: ПАРСИНГ И ВАЛИДАЦИЯ JSON ===
	var req domain.UpdateUserRequest
	
	// ShouldBindJSON() парсит и валидирует
	// Проверяет binding теги (omitempty, email, min=2)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// === ШАГ 3: ВЫЗОВ SERVICE ===
	// Service обновит пользователя в БД
	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// === ШАГ 4: ОТПРАВКА ОТВЕТА ===
	c.JSON(http.StatusOK, user)
}

// Delete удаляет пользователя (soft delete)
// Endpoint: DELETE /api/v1/users/:id
// Headers: Authorization: Bearer TOKEN (защищён!)
// Response: {"message": "пользователь удалён"}
func (h *UserHandler) Delete(c *gin.Context) {
	// === ШАГ 1: ИЗВЛЕЧЕНИЕ ID ===
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "невалидный ID",
		})
		return
	}

	// === ШАГ 2: ВЫЗОВ SERVICE ===
	// Service удалит пользователя (soft delete)
	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// === ШАГ 3: ОТПРАВКА ОТВЕТА ===
	// Возвращаем подтверждение удаления
	c.JSON(http.StatusOK, gin.H{
		"message": "пользователь удалён",
	})
}

