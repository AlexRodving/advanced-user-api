package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ================================================================
// LOGGER MIDDLEWARE - Логирование HTTP запросов
// ================================================================

// LoggerMiddleware логирует каждый HTTP запрос
// Записывает:
// - Метод (GET, POST, etc.)
// - Path (/api/v1/users)
// - Статус код (200, 404, etc.)
// - Время выполнения
// - IP адрес клиента
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Запоминаем время начала обработки запроса
		start := time.Now()
		
		// Запоминаем путь (может измениться в handlers)
		path := c.Request.URL.Path
		
		// Обрабатываем запрос (вызываем следующие handlers)
		c.Next()
		
		// === ЛОГИРОВАНИЕ ПОСЛЕ ОБРАБОТКИ ===
		// Вычисляем время выполнения
		latency := time.Since(start)
		
		// Получаем информацию о запросе
		statusCode := c.Writer.Status()      // HTTP статус код
		method := c.Request.Method           // HTTP метод
		clientIP := c.ClientIP()             // IP адрес клиента
		errorMessage := c.Errors.String()    // Ошибки (если были)
		
		// Логируем с разным уровнем в зависимости от статуса
		if statusCode >= 500 {
			// 5xx - ошибки сервера (Error level)
			logger.Error("HTTP Request",
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.Duration("latency", latency),
				zap.String("ip", clientIP),
				zap.String("error", errorMessage),
			)
		} else if statusCode >= 400 {
			// 4xx - ошибки клиента (Warn level)
			logger.Warn("HTTP Request",
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.Duration("latency", latency),
				zap.String("ip", clientIP),
			)
		} else {
			// 2xx, 3xx - успех (Info level)
			logger.Info("HTTP Request",
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.Duration("latency", latency),
				zap.String("ip", clientIP),
			)
		}
	}
}

