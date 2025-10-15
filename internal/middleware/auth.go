package middleware

import (
	"net/http"
	"strings"

	"advanced-user-api/internal/config"
	"advanced-user-api/internal/pkg/jwt"

	"github.com/gin-gonic/gin" // Gin фреймворк
)

// ================================================================
// AUTH MIDDLEWARE - Проверка JWT токена
// ================================================================

// AuthMiddleware создаёт middleware для проверки JWT токена
// Параметры:
//   - cfg: конфигурация (для получения JWT secret)
// Возвращает:
//   - gin.HandlerFunc: middleware функцию
//
// Использование:
//   authorized := r.Group("/api/v1/users")
//   authorized.Use(middleware.AuthMiddleware(cfg))
//   {
//       authorized.GET("", handler.GetAll) // Требует токен
//   }
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Возвращаем функцию-обработчик
	// Эта функция будет вызываться для каждого запроса к защищённым routes
	return func(c *gin.Context) {
		// === ШАГ 1: ИЗВЛЕЧЕНИЕ ТОКЕНА ИЗ ЗАГОЛОВКА ===
		// Токен должен быть в заголовке Authorization
		// Формат: "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
		authHeader := c.GetHeader("Authorization")
		
		// Проверяем наличие заголовка
		if authHeader == "" {
			// Заголовок отсутствует - отклоняем запрос
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "отсутствует токен аутентификации",
			})
			c.Abort() // Прерываем обработку запроса (не вызываем следующие handlers)
			return
		}

		// === ШАГ 2: ПАРСИНГ ЗАГОЛОВКА ===
		// Формат: "Bearer TOKEN"
		// Разделяем по пробелу
		parts := strings.Split(authHeader, " ")
		
		// Проверяем формат
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Неправильный формат заголовка
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "неверный формат токена (используйте: Bearer TOKEN)",
			})
			c.Abort()
			return
		}

		// Извлекаем сам токен (вторая часть после "Bearer ")
		tokenString := parts[1]

		// === ШАГ 3: ВАЛИДАЦИЯ ТОКЕНА ===
		// Проверяем подпись и срок действия токена
		claims, err := jwt.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			// Токен невалиден (истёк, неправильная подпись, повреждён)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "невалидный или истёкший токен",
			})
			c.Abort()
			return
		}

		// === ШАГ 4: СОХРАНЕНИЕ ДАННЫХ В КОНТЕКСТ ===
		// Gin Context - хранилище данных для текущего запроса
		// Сохраняем данные из токена, чтобы handlers могли их использовать
		
		// Сохраняем ID пользователя
		c.Set("userID", claims.UserID)
		
		// Сохраняем email
		c.Set("userEmail", claims.Email)
		
		// Сохраняем роль (для проверки прав доступа)
		c.Set("userRole", claims.Role)

		// === ШАГ 5: ПРОДОЛЖЕНИЕ ОБРАБОТКИ ===
		// c.Next() - вызывает следующий handler в цепочке
		// Если не вызвать Next(), запрос остановится здесь
		c.Next()
	}
}

// ================================================================
// HELPER FUNCTIONS - Вспомогательные функции
// ================================================================

// GetUserIDFromContext - извлекает ID пользователя из Gin контекста
// Используется в handlers для получения ID текущего пользователя
//
// Пример использования:
//   userID := middleware.GetUserIDFromContext(c)
//   user, _ := service.GetUser(userID)
func GetUserIDFromContext(c *gin.Context) uint {
	// c.Get() - получает значение из контекста
	// Возвращает (interface{}, bool)
	userID, exists := c.Get("userID")
	if !exists {
		// Если userID нет в контексте - возвращаем 0
		return 0
	}
	
	// Приводим interface{} к типу uint
	if id, ok := userID.(uint); ok {
		return id
	}
	
	return 0
}

// GetUserRoleFromContext - извлекает роль пользователя из контекста
func GetUserRoleFromContext(c *gin.Context) string {
	role, exists := c.Get("userRole")
	if !exists {
		return ""
	}
	
	if r, ok := role.(string); ok {
		return r
	}
	
	return ""
}

// RequireRole - middleware для проверки роли пользователя
// Используется для ограничения доступа (например, только для admin)
//
// Пример:
//   admin := r.Group("/api/v1/admin")
//   admin.Use(middleware.AuthMiddleware(cfg))
//   admin.Use(middleware.RequireRole("admin"))
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем роль из контекста (установлена в AuthMiddleware)
		userRole := GetUserRoleFromContext(c)
		
		// Проверяем роль
		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "недостаточно прав доступа",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

