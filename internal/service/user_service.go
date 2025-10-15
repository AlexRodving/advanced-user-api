package service

import (
	"advanced-user-api/internal/domain"
	"advanced-user-api/internal/repository"
)

// ================================================================
// USER SERVICE - Бизнес-логика для пользователей
// ================================================================

// UserService - интерфейс для работы с пользователями
type UserService interface {
	GetUser(id uint) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	UpdateUser(id uint, req *domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(id uint) error
	GetCurrentUser(id uint) (*domain.User, error)
}

// userService - реализация сервиса
type userService struct {
	userRepo repository.UserRepository // Зависимость от Repository
}

// NewUserService - конструктор
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// ================================================================
// SERVICE METHODS
// ================================================================

// GetUser - получает пользователя по ID
func (s *userService) GetUser(id uint) (*domain.User, error) {
	// Просто вызываем repository
	// Service слой здесь не добавляет логики (но может добавить в будущем)
	return s.userRepo.FindByID(id)
}

// GetAllUsers - получает всех пользователей
func (s *userService) GetAllUsers() ([]domain.User, error) {
	// Получаем всех пользователей из repository
	return s.userRepo.FindAll()
}

// UpdateUser - обновляет данные пользователя
// Параметры:
//   - id: ID пользователя для обновления
//   - req: новые данные (email, name)
// Возвращает:
//   - *domain.User: обновлённый пользователь
//   - error: ошибка обновления
func (s *userService) UpdateUser(id uint, req *domain.UpdateUserRequest) (*domain.User, error) {
	// === ШАГ 1: ПРОВЕРКА СУЩЕСТВОВАНИЯ ===
	// Находим пользователя по ID
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err // Пользователь не найден
	}

	// === ШАГ 2: ОБНОВЛЕНИЕ ПОЛЕЙ ===
	// Обновляем только те поля, которые переданы
	
	// Если передан новый email - обновляем
	if req.Email != "" {
		user.Email = req.Email
	}
	
	// Если передано новое имя - обновляем
	if req.Name != "" {
		user.Name = req.Name
	}

	// === ШАГ 3: СОХРАНЕНИЕ В БД ===
	// Save() обновит запись в БД
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Возвращаем обновлённого пользователя
	return user, nil
}

// DeleteUser - удаляет пользователя (soft delete)
func (s *userService) DeleteUser(id uint) error {
	// Проверяем существование пользователя
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Удаляем через repository
	return s.userRepo.Delete(id)
}

// GetCurrentUser - получает данные текущего аутентифицированного пользователя
// Используется для endpoint GET /auth/me
func (s *userService) GetCurrentUser(id uint) (*domain.User, error) {
	// Находим пользователя по ID из JWT токена
	return s.userRepo.FindByID(id)
}

