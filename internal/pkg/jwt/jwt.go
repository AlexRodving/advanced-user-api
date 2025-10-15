package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5" // JWT библиотека
)

// ================================================================
// JWT UTILITIES - Работа с JSON Web Tokens
// ================================================================

// Claims - структура данных, которые хранятся в JWT токене
// Содержит информацию о пользователе и метаданные токена
type Claims struct {
	// UserID - ID пользователя (для идентификации)
	UserID uint `json:"user_id"`
	
	// Email - email пользователя (дополнительная информация)
	Email string `json:"email"`
	
	// Role - роль пользователя (для проверки прав доступа)
	Role string `json:"role"`
	
	// RegisteredClaims - стандартные JWT claims (exp, iat, iss, etc.)
	// Включает:
	//   - ExpiresAt: время истечения токена
	//   - IssuedAt: время создания токена
	//   - NotBefore: токен не валиден до этого времени
	//   - Issuer: кто выдал токен
	//   - Subject: для кого токен
	jwt.RegisteredClaims
}

// ================================================================
// GENERATE TOKEN - Создание JWT токена
// ================================================================

// GenerateToken создаёт новый JWT токен для пользователя
// Параметры:
//   - userID: ID пользователя
//   - email: email пользователя
//   - role: роль пользователя
//   - secret: секретный ключ для подписи токена
//   - expiration: время жизни токена (например, "24h")
// Возвращает:
//   - string: JWT токен (строка вида "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
//   - error: ошибка генерации
func GenerateToken(userID uint, email, role, secret string, expiration time.Duration) (string, error) {
	// === ШАГ 1: СОЗДАНИЕ CLAIMS ===
	// Claims - данные, которые будут закодированы в токене
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// ExpiresAt - время истечения токена
			// time.Now().Add(expiration) - текущее время + 24 часа (например)
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			
			// IssuedAt - время создания токена
			IssuedAt: jwt.NewNumericDate(time.Now()),
			
			// Issuer - кто выдал токен (название вашего приложения)
			Issuer: "advanced-user-api",
		},
	}

	// === ШАГ 2: СОЗДАНИЕ ТОКЕНА ===
	// jwt.NewWithClaims() - создаёт новый токен с указанными claims
	// Параметры:
	//   - jwt.SigningMethodHS256: алгоритм подписи (HMAC-SHA256)
	//   - claims: данные для кодирования
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// === ШАГ 3: ПОДПИСЬ ТОКЕНА ===
	// SignedString() - подписывает токен секретным ключом
	// Возвращает строку вида: "header.payload.signature"
	// Только тот, кто знает secret, может создать валидный токен!
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	// Возвращаем готовый JWT токен
	return tokenString, nil
}

// ================================================================
// VALIDATE TOKEN - Проверка JWT токена
// ================================================================

// ValidateToken проверяет валидность JWT токена
// Параметры:
//   - tokenString: JWT токен от клиента
//   - secret: секретный ключ для проверки подписи
// Возвращает:
//   - *Claims: данные из токена (если токен валиден)
//   - error: ошибка валидации
func ValidateToken(tokenString, secret string) (*Claims, error) {
	// === ШАГ 1: ПАРСИНГ ТОКЕНА ===
	// jwt.ParseWithClaims() - парсит и валидирует токен
	// Параметры:
	//   - tokenString: токен для проверки
	//   - &Claims{}: структура для заполнения данными из токена
	//   - keyFunc: функция для получения секретного ключа
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		// Защита от атаки "algorithm confusion"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неожиданный алгоритм подписи")
		}
		
		// Возвращаем секретный ключ для проверки подписи
		return []byte(secret), nil
	})

	// === ШАГ 2: ПРОВЕРКА ОШИБОК ===
	if err != nil {
		// Возможные ошибки:
		// - Токен истёк (expired)
		// - Неправильная подпись (invalid signature)
		// - Токен повреждён (malformed)
		return nil, err
	}

	// === ШАГ 3: ИЗВЛЕЧЕНИЕ CLAIMS ===
	// token.Claims - данные из токена
	// Приводим к типу *Claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		// Если не удалось привести к Claims или токен невалиден
		return nil, errors.New("невалидный токен")
	}

	// Возвращаем claims с данными пользователя
	return claims, nil
}

// ================================================================
// ДОПОЛНИТЕЛЬНЫЕ ФУНКЦИИ
// ================================================================

// ExtractUserID - извлекает ID пользователя из токена
// Удобная функция для быстрого получения только ID
func ExtractUserID(tokenString, secret string) (uint, error) {
	// Валидируем токен
	claims, err := ValidateToken(tokenString, secret)
	if err != nil {
		return 0, err
	}
	
	// Возвращаем только UserID
	return claims.UserID, nil
}

// IsTokenExpired - проверяет, истёк ли токен
func IsTokenExpired(tokenString, secret string) bool {
	claims, err := ValidateToken(tokenString, secret)
	if err != nil {
		return true // Если ошибка - считаем истёкшим
	}
	
	// Проверяем ExpiresAt
	return claims.ExpiresAt.Before(time.Now())
}

