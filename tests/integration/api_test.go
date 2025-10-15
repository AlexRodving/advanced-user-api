package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"advanced-user-api/internal/config"
	"advanced-user-api/internal/domain"
	"advanced-user-api/internal/handler"
	"advanced-user-api/internal/repository"
	"advanced-user-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ================================================================
// INTEGRATION TESTS - Интеграционные тесты с реальной БД
// ================================================================

// setupTestDB - создаёт тестовую БД
// В реальных проектах используйте testcontainers-go
func setupTestDB() *gorm.DB {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=advanced_api_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Auto Migrate
	db.AutoMigrate(&domain.User{})

	return db
}

// cleanupTestDB - очищает тестовую БД
func cleanupTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM users")
}

// TestFullAuthFlow - полный цикл: регистрация → вход → получение данных
func TestFullAuthFlow(t *testing.T) {
	// === SETUP ===
	db := setupTestDB()
	defer cleanupTestDB(db)

	// Создаём слои приложения
	userRepo := repository.NewUserRepository(db)
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: "24h",
	}
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)

	// Создаём handlers
	authHandler := handler.NewAuthHandler(authService, userService)

	// Создаём роутер
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler.SetupRoutes(router, authHandler, nil, cfg)

	// === TEST: РЕГИСТРАЦИЯ ===
	registerReq := domain.RegisterRequest{
		Email:    "integration@test.com",
		Name:     "Integration Test",
		Password: "password123",
	}
	body, _ := json.Marshal(registerReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Проверяем статус
	assert.Equal(t, http.StatusCreated, w.Code)

	// Парсим ответ
	var registerResp domain.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &registerResp)

	// Проверяем токен
	assert.NotEmpty(t, registerResp.Token)
	assert.Equal(t, registerReq.Email, registerResp.User.Email)

	// === TEST: ВХОД ===
	loginReq := domain.LoginRequest{
		Email:    "integration@test.com",
		Password: "password123",
	}
	body, _ = json.Marshal(loginReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp domain.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NotEmpty(t, loginResp.Token)

	// === TEST: ЗАЩИЩЁННЫЙ ENDPOINT С ТОКЕНОМ ===
	token := loginResp.Token

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var user domain.User
	json.Unmarshal(w.Body.Bytes(), &user)
	assert.Equal(t, "integration@test.com", user.Email)
}

