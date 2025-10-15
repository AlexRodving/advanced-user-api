package service

import (
	"errors"
	"time"

	"advanced-user-api/internal/config"
	"advanced-user-api/internal/domain"
	"advanced-user-api/internal/pkg/jwt"
	"advanced-user-api/internal/pkg/password"
	"advanced-user-api/internal/repository"
)

// ================================================================
// AUTH SERVICE - Сервис аутентификации
// ================================================================

// AuthService - интерфейс для аутентификации пользователей
type AuthService interface {
	Register(req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(req *domain.LoginRequest) (*domain.AuthResponse, error)
}

// authService - реализация сервиса аутентификации
type authService struct {
	userRepo repository.UserRepository // Зависимость от Repository
	cfg      *config.Config            // Конфигурация (для JWT secret)
}

// NewAuthService - конструктор для создания Auth Service
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// ================================================================
// REGISTER - Регистрация нового пользователя
// ================================================================

// Register регистрирует нового пользователя в системе
// Параметры:
//   - req: данные для регистрации (email, name, password)
// Возвращает:
//   - *domain.AuthResponse: JWT токен и данные пользователя
//   - error: ошибка регистрации
//
// Процесс:
// 1. Проверяем, не существует ли уже пользователь с таким email
// 2. Хешируем пароль (bcrypt)
// 3. Создаём пользователя в БД
// 4. Генерируем JWT токен
// 5. Возвращаем токен и данные пользователя
func (s *authService) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// === ШАГ 1: ПРОВЕРКА СУЩЕСТВОВАНИЯ ПОЛЬЗОВАТЕЛЯ ===
	// Проверяем, не зарегистрирован ли уже пользователь с таким email
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		// Пользователь с таким email уже существует
		return nil, errors.New("пользователь с таким email уже зарегистрирован")
	}

	// === ШАГ 2: ХЕШИРОВАНИЕ ПАРОЛЯ ===
	// НИКОГДА не сохраняйте пароли в открытом виде!
	// Хешируем пароль с помощью bcrypt
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, errors.New("ошибка хеширования пароля")
	}

	// === ШАГ 3: СОЗДАНИЕ ПОЛЬЗОВАТЕЛЯ ===
	// Создаём структуру User для сохранения в БД
	user := &domain.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword, // Сохраняем ХЕШ, не сам пароль!
		Role:     "user",         // По умолчанию роль "user"
	}

	// Сохраняем пользователя в БД через repository
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("ошибка создания пользователя")
	}

	// === ШАГ 4: ГЕНЕРАЦИЯ JWT ТОКЕНА ===
	// Парсим время жизни токена из конфигурации
	// "24h" → 24 часа
	expiration, err := time.ParseDuration(s.cfg.JWTExpiration)
	if err != nil {
		expiration = 24 * time.Hour // Если ошибка парсинга - используем 24 часа
	}

	// Генерируем JWT токен с данными пользователя
	token, err := jwt.GenerateToken(
		user.ID,            // ID пользователя
		user.Email,         // Email
		user.Role,          // Роль
		s.cfg.JWTSecret,    // Секретный ключ из конфигурации
		expiration,         // Время жизни токена
	)
	if err != nil {
		return nil, errors.New("ошибка генерации токена")
	}

	// === ШАГ 5: ФОРМИРОВАНИЕ ОТВЕТА ===
	// Возвращаем токен и данные пользователя (без пароля!)
	return &domain.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// ================================================================
// LOGIN - Вход пользователя
// ================================================================

// Login аутентифицирует пользователя и выдаёт JWT токен
// Параметры:
//   - req: данные для входа (email, password)
// Возвращает:
//   - *domain.AuthResponse: JWT токен и данные пользователя
//   - error: ошибка аутентификации
//
// Процесс:
// 1. Находим пользователя по email
// 2. Проверяем пароль (bcrypt.Compare)
// 3. Генерируем JWT токен
// 4. Возвращаем токен и данные пользователя
func (s *authService) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// === ШАГ 1: ПОИСК ПОЛЬЗОВАТЕЛЯ ===
	// Ищем пользователя по email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Пользователь не найден
		// ВАЖНО: Не говорим "email не найден" - это утечка информации
		// Говорим общее "неверные credentials"
		return nil, errors.New("неверный email или пароль")
	}

	// === ШАГ 2: ПРОВЕРКА ПАРОЛЯ ===
	// Сравниваем хеш из БД с введённым паролем
	// password.Verify() использует bcrypt.CompareHashAndPassword()
	if !password.Verify(user.Password, req.Password) {
		// Пароль неправильный
		return nil, errors.New("неверный email или пароль")
	}

	// === ШАГ 3: ГЕНЕРАЦИЯ JWT ТОКЕНА ===
	// Парсим время жизни токена
	expiration, err := time.ParseDuration(s.cfg.JWTExpiration)
	if err != nil {
		expiration = 24 * time.Hour
	}

	// Генерируем токен для аутентифицированного пользователя
	token, err := jwt.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
		s.cfg.JWTSecret,
		expiration,
	)
	if err != nil {
		return nil, errors.New("ошибка генерации токена")
	}

	// === ШАГ 4: ВОЗВРАТ ОТВЕТА ===
	// Возвращаем токен и данные пользователя
	return &domain.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

