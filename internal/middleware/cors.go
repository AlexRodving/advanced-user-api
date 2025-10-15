package middleware

import (
	"github.com/gin-gonic/gin"
)

// ================================================================
// CORS MIDDLEWARE - Cross-Origin Resource Sharing
// ================================================================

// CORS разрешает запросы с других доменов
// Без CORS браузер блокирует запросы с frontend (например, React app на localhost:3000)
// к API на localhost:8080

// CORSMiddleware настраивает CORS заголовки
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// === CORS ЗАГОЛОВКИ ===
		
		// Access-Control-Allow-Origin - какие домены могут делать запросы
		// "*" - разрешить всем (для development)
		// В production лучше указать конкретные домены:
		// c.Header("Access-Control-Allow-Origin", "https://myapp.com")
		c.Header("Access-Control-Allow-Origin", "*")
		
		// Access-Control-Allow-Methods - какие HTTP методы разрешены
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// Access-Control-Allow-Headers - какие заголовки может отправлять клиент
		// Authorization - для JWT токена
		// Content-Type - для JSON
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		
		// Access-Control-Allow-Credentials - разрешить отправку cookies
		c.Header("Access-Control-Allow-Credentials", "true")

		// === PREFLIGHT REQUEST ===
		// Браузер отправляет OPTIONS запрос перед реальным запросом
		// Это называется "preflight request"
		// Нужно ответить 200 OK без обработки
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 204 No Content
			return
		}

		// Продолжаем обработку запроса
		c.Next()
	}
}

