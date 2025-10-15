package repository

import (
	"errors"

	"advanced-user-api/internal/domain"

	"gorm.io/gorm" // GORM ORM
)

// ================================================================
// REPOSITORY INTERFACE - Контракт для работы с пользователями
// ================================================================

// UserRepository - интерфейс для работы с пользователями в БД
// Определяет ЧТО можно делать с пользователями (контракт)
// Реализация (КАК это делается) находится в userRepository ниже
type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}

// ================================================================
// REPOSITORY IMPLEMENTATION - Реализация с GORM
// ================================================================

// userRepository - приватная структура, реализующая UserRepository
// Содержит подключение к БД через GORM
type userRepository struct {
	db *gorm.DB // Подключение к БД
}

// NewUserRepository - конструктор для создания нового repository
// Принимает GORM подключение как зависимость (Dependency Injection)
// Возвращает интерфейс (не конкретный тип!)
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// ================================================================
// CRUD OPERATIONS - Операции с базой данных
// ================================================================

// Create - создаёт нового пользователя в БД
// Параметры:
//   - user: указатель на структуру User для создания
// Возвращает:
//   - error: ошибка создания (если есть)
//
// GORM автоматически:
// - Генерирует ID (автоинкремент)
// - Устанавливает CreatedAt и UpdatedAt
// - Возвращает созданную запись с ID в ту же структуру user
func (r *userRepository) Create(user *domain.User) error {
	// db.Create() - вставляет новую запись в БД
	// Генерирует SQL: INSERT INTO users (email, name, password, ...) VALUES (?, ?, ?, ...)
	// После выполнения user.ID будет содержать ID из БД!
	// .Error - возвращает ошибку (если есть)
	return r.db.Create(user).Error
}

// FindByID - ищет пользователя по ID
// Параметры:
//   - id: ID пользователя для поиска
// Возвращает:
//   - *domain.User: найденный пользователь
//   - error: ошибка поиска или gorm.ErrRecordNotFound если не найден
func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	// Создаём пустую структуру для заполнения данными из БД
	var user domain.User
	
	// db.First() - находит ПЕРВУЮ запись по условию
	// Генерирует SQL: SELECT * FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1
	// &user - указатель на структуру, куда GORM запишет результат
	// id - значение для поиска (подставится вместо ?)
	// .Error - ошибка выполнения
	err := r.db.First(&user, id).Error
	
	// Проверяем специальную ошибку "запись не найдена"
	if err == gorm.ErrRecordNotFound {
		// Возвращаем более понятную ошибку для пользователя
		return nil, errors.New("пользователь не найден")
	}
	
	// Если другая ошибка (например, ошибка БД) - возвращаем её
	if err != nil {
		return nil, err
	}
	
	// Возвращаем найденного пользователя
	return &user, nil
}

// FindByEmail - ищет пользователя по email адресу
// Используется для аутентификации (login)
// Параметры:
//   - email: email для поиска
// Возвращает:
//   - *domain.User: найденный пользователь
//   - error: ошибка поиска
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	// Создаём пустую структуру
	var user domain.User
	
	// db.Where() - добавляет условие WHERE в запрос
	// Генерирует SQL: SELECT * FROM users WHERE email = ? AND deleted_at IS NULL
	// "email = ?" - условие (? заменится на значение email)
	// .First(&user) - выполняет запрос и сканирует результат в user
	err := r.db.Where("email = ?", email).First(&user).Error
	
	// Проверяем, найден ли пользователь
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("пользователь с таким email не найден")
	}
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// FindAll - получает всех пользователей из БД
// Возвращает:
//   - []domain.User: slice всех пользователей
//   - error: ошибка выполнения запроса
func (r *userRepository) FindAll() ([]domain.User, error) {
	// Создаём пустой slice для заполнения
	var users []domain.User
	
	// db.Find() - находит ВСЕ записи
	// Генерирует SQL: SELECT * FROM users WHERE deleted_at IS NULL
	// &users - указатель на slice, куда GORM запишет все найденные записи
	// GORM автоматически исключает "удалённые" записи (где deleted_at не NULL)
	err := r.db.Find(&users).Error
	
	// Если ошибка - возвращаем nil и ошибку
	if err != nil {
		return nil, err
	}
	
	// Возвращаем slice пользователей (может быть пустым, если пользователей нет)
	return users, nil
}

// Update - обновляет данные пользователя в БД
// Параметры:
//   - user: указатель на User с обновлёнными данными
//           ВАЖНО: user.ID должен быть установлен!
// Возвращает:
//   - error: ошибка обновления
func (r *userRepository) Update(user *domain.User) error {
	// db.Save() - обновляет ВСЕ поля записи
	// Генерирует SQL: UPDATE users SET email=?, name=?, updated_at=? WHERE id=?
	// GORM автоматически:
	// 1. Обновляет UpdatedAt на текущее время
	// 2. Использует user.ID для поиска записи
	// 3. Обновляет все поля (кроме ID, CreatedAt, DeletedAt)
	//
	// Альтернатива - db.Updates() для обновления только изменённых полей:
	// r.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(user)
	return r.db.Save(user).Error
}

// Delete - "мягко" удаляет пользователя (soft delete)
// Параметры:
//   - id: ID пользователя для удаления
// Возвращает:
//   - error: ошибка удаления
//
// ВАЖНО: Это НЕ физическое удаление!
// GORM просто устанавливает deleted_at = NOW()
// Запись остаётся в БД, но игнорируется во всех запросах
func (r *userRepository) Delete(id uint) error {
	// db.Delete() - "мягкое" удаление (soft delete)
	// Генерирует SQL: UPDATE users SET deleted_at = NOW() WHERE id = ?
	// &domain.User{} - пустая структура (нужна только для определения таблицы)
	// id - ID для удаления
	//
	// Если нужно ФИЗИЧЕСКОЕ удаление (hard delete):
	// r.db.Unscoped().Delete(&domain.User{}, id)
	result := r.db.Delete(&domain.User{}, id)
	
	// Проверяем ошибку выполнения
	if result.Error != nil {
		return result.Error
	}
	
	// Проверяем, была ли затронута хотя бы одна строка
	// RowsAffected - количество затронутых строк
	// Если 0 - пользователь с таким ID не существует (или уже удалён)
	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}
	
	// Успешно удалили
	return nil
}

// ================================================================
// ДОПОЛНИТЕЛЬНЫЕ МЕТОДЫ (примеры для расширения)
// ================================================================

// FindByRole - находит всех пользователей с определённой ролью
// Пример расширения функциональности
func (r *userRepository) FindByRole(role string) ([]domain.User, error) {
	var users []domain.User
	
	// WHERE с условием по role
	err := r.db.Where("role = ?", role).Find(&users).Error
	
	return users, err
}

// CountAll - подсчитывает общее количество пользователей
func (r *userRepository) CountAll() (int64, error) {
	var count int64
	
	// db.Model() - указываем модель
	// Count() - подсчитывает количество записей
	// Генерирует SQL: SELECT COUNT(*) FROM users WHERE deleted_at IS NULL
	err := r.db.Model(&domain.User{}).Count(&count).Error
	
	return count, err
}

