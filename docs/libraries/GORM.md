# 🗄️ GORM - Go ORM Library

## Overview

GORM - мощная ORM библиотека для Go с полной поддержкой функций.

**Официальная документация:** https://gorm.io/docs/

**Версия в проекте:** v1.25.5

**Driver:** PostgreSQL (gorm.io/driver/postgres v1.5.4)

---

## Подключение к БД

### gorm.Open()

**Описание:** Открывает подключение к БД

**Сигнатура:**
```go
func Open(dialector Dialector, opts ...Option) (db *DB, err error)
```

**Параметры:**
- `dialector Dialector` - драйвер БД (postgres.Open(dsn))
- `opts ...Option` - опции GORM (логирование, и т.д.)

**Возвращает:**
- `*DB` - подключение к БД
- `error` - ошибка подключения

**Пример использования:**
```go
dsn := "host=localhost user=postgres password=postgres dbname=mydb port=5432"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

**Где используем:**
- `internal/repository/database.go:41`

---

## CRUD Operations

### 1. db.Create()

**Описание:** Создаёт новую запись в БД

**Сигнатура:**
```go
func (db *DB) Create(value interface{}) *DB
```

**Параметры:**
- `value interface{}` - указатель на структуру для создания

**Возвращает:** `*DB` - для цепочки вызовов и проверки ошибок

**Пример использования:**
```go
user := &domain.User{
    Email: "test@example.com",
    Name: "Test User",
    Password: hashedPassword,
}

result := db.Create(user)
if result.Error != nil {
    return result.Error
}

// После Create, user.ID заполнится автоматически
fmt.Println(user.ID)  // 1, 2, 3...
```

**Генерируемый SQL:**
```sql
INSERT INTO users (email, name, password, created_at, updated_at)
VALUES ('test@example.com', 'Test User', '$2a$10...', NOW(), NOW())
RETURNING id
```

**Где используем:**
- `internal/repository/user_repository.go:25`

---

### 2. db.First()

**Описание:** Находит первую запись по условию

**Сигнатура:**
```go
func (db *DB) First(dest interface{}, conds ...interface{}) *DB
```

**Параметры:**
- `dest interface{}` - указатель на структуру результата
- `conds ...interface{}` - условия поиска (WHERE)

**Возвращает:** `*DB`

**Ошибки:**
- `gorm.ErrRecordNotFound` - если запись не найдена

**Пример использования:**
```go
var user domain.User

// По первичному ключу
db.First(&user, 1)  // WHERE id = 1

// По условию
db.First(&user, "email = ?", "test@example.com")
```

**Генерируемый SQL:**
```sql
SELECT * FROM users WHERE id = 1 AND deleted_at IS NULL LIMIT 1
```

**Где используем:**
- `internal/repository/user_repository.go:32`

---

### 3. db.Where()

**Описание:** Добавляет WHERE условие

**Сигнатура:**
```go
func (db *DB) Where(query interface{}, args ...interface{}) *DB
```

**Параметры:**
- `query interface{}` - SQL условие или struct/map
- `args ...interface{}` - параметры для SQL (защита от injection)

**Возвращает:** `*DB` - для цепочки

**Пример использования:**
```go
// String условие
db.Where("email = ?", "test@example.com").First(&user)

// Struct условие
db.Where(&User{Email: "test@example.com"}).First(&user)

// Map условие
db.Where(map[string]interface{}{"email": "test@example.com"}).First(&user)

// Множественные условия
db.Where("email = ? AND role = ?", "test@example.com", "admin").First(&user)

// IN условие
db.Where("id IN ?", []int{1, 2, 3}).Find(&users)

// LIKE условие
db.Where("name LIKE ?", "%John%").Find(&users)
```

**Генерируемый SQL:**
```sql
SELECT * FROM users WHERE email = 'test@example.com' AND deleted_at IS NULL
```

**Где используем:**
- `internal/repository/user_repository.go:32`

---

### 4. db.Find()

**Описание:** Находит все записи по условию

**Сигнатура:**
```go
func (db *DB) Find(dest interface{}, conds ...interface{}) *DB
```

**Параметры:**
- `dest interface{}` - указатель на slice для результатов
- `conds ...interface{}` - условия (опционально)

**Возвращает:** `*DB`

**Пример использования:**
```go
var users []domain.User

// Все записи
db.Find(&users)

// С условием
db.Where("role = ?", "admin").Find(&users)

// По ID
db.Find(&users, []int{1, 2, 3})
```

**Генерируемый SQL:**
```sql
SELECT * FROM users WHERE deleted_at IS NULL
```

**Где используем:**
- `internal/repository/user_repository.go:47`

---

### 5. db.Save()

**Описание:** Обновляет ВСЕ поля или создаёт если нет ID

**Сигнатура:**
```go
func (db *DB) Save(value interface{}) *DB
```

**Параметры:**
- `value interface{}` - указатель на структуру

**Возвращает:** `*DB`

**Поведение:**
- Если ID = 0 → CREATE
- Если ID > 0 → UPDATE всех полей

**Пример использования:**
```go
user.Name = "Updated Name"
user.Email = "new@example.com"
db.Save(&user)  // UPDATE всех полей
```

**Генерируемый SQL:**
```sql
UPDATE users SET email='new@example.com', name='Updated Name', updated_at=NOW()
WHERE id = 1 AND deleted_at IS NULL
```

**Где используем:**
- `internal/repository/user_repository.go:57` (через Update с Select)

---

### 6. db.Updates()

**Описание:** Обновляет только заполненные поля

**Сигнатура:**
```go
func (db *DB) Updates(values interface{}) *DB
```

**Параметры:**
- `values interface{}` - struct или map с новыми значениями

**Возвращает:** `*DB`

**Пример использования:**
```go
// Struct (игнорирует zero values!)
db.Model(&user).Updates(User{Name: "New Name"})

// Map (обновляет все, даже zero values)
db.Model(&user).Updates(map[string]interface{}{
    "name": "New Name",
    "age": 0,  // Обновит на 0
})
```

**Где используем:**
- `internal/repository/user_repository.go:57`

---

### 7. db.Delete()

**Описание:** Удаляет запись (soft delete если есть DeletedAt)

**Сигнатура:**
```go
func (db *DB) Delete(value interface{}, conds ...interface{}) *DB
```

**Параметры:**
- `value interface{}` - модель для удаления
- `conds ...interface{}` - дополнительные условия

**Возвращает:** `*DB`

**Пример использования:**
```go
// Soft delete (если есть gorm.DeletedAt)
db.Delete(&user, 1)  // SET deleted_at = NOW()

// Permanent delete
db.Unscoped().Delete(&user, 1)  // DELETE FROM users
```

**Генерируемый SQL (soft delete):**
```sql
UPDATE users SET deleted_at = NOW() WHERE id = 1 AND deleted_at IS NULL
```

**Где используем:**
- `internal/repository/user_repository.go:71`

---

## Model Methods

### db.Model()

**Описание:** Указывает модель для операции

**Сигнатура:**
```go
func (db *DB) Model(value interface{}) *DB
```

**Параметры:**
- `value interface{}` - модель

**Возвращает:** `*DB`

**Пример использования:**
```go
db.Model(&User{}).Where("role = ?", "admin").Count(&count)
db.Model(&user).Update("name", "New Name")
```

**Где используем:**
- `internal/repository/user_repository.go:57`

---

### db.Select()

**Описание:** Выбирает только определённые поля

**Сигнатура:**
```go
func (db *DB) Select(query interface{}, args ...interface{}) *DB
```

**Параметры:**
- `query interface{}` - поля для выборки
- `args ...interface{}` - параметры

**Возвращает:** `*DB`

**Пример использования:**
```go
// Строка
db.Select("name, email").Find(&users)

// Slice
db.Select([]string{"name", "email"}).Find(&users)

// Все поля кроме
db.Omit("password").Find(&users)
```

**Где используем:**
- `internal/repository/user_repository.go:57` - выбор полей для обновления

---

## Auto Migration

### db.AutoMigrate()

**Описание:** Автоматически создаёт/обновляет таблицы

**Сигнатура:**
```go
func (db *DB) AutoMigrate(dst ...interface{}) error
```

**Параметры:**
- `dst ...interface{}` - модели для миграции

**Возвращает:** `error`

**Что делает:**
- Создаёт таблицу если не существует
- Добавляет новые колонки
- Добавляет индексы
- **НЕ** удаляет колонки
- **НЕ** меняет типы

**Пример использования:**
```go
db.AutoMigrate(&domain.User{}, &domain.Post{}, &domain.Comment{})
```

**Где используем:**
- `internal/repository/database.go:67`

---

## Struct Tags

### GORM теги для моделей

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Email     string         `gorm:"uniqueIndex;not null"`
    Name      string         `gorm:"not null"`
    Age       int            `gorm:"default:18"`
    Role      string         `gorm:"type:varchar(50);default:'user'"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### Основные теги:

| Тег | Описание |
|-----|----------|
| `primaryKey` | PRIMARY KEY |
| `unique` | UNIQUE constraint |
| `uniqueIndex` | UNIQUE INDEX |
| `not null` | NOT NULL |
| `default:value` | DEFAULT value |
| `type:varchar(100)` | Тип колонки |
| `index` | Создать индекс |
| `autoCreateTime` | Авто заполнение при создании |
| `autoUpdateTime` | Авто обновление |
| `-` | Игнорировать поле |
| `column:custom_name` | Имя колонки |

**Где используем:**
- `internal/domain/user.go:18-33`

---

## Error Handling

### Проверка ошибок

```go
result := db.First(&user, id)

// Способ 1: Проверка result.Error
if result.Error != nil {
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, errors.New("user not found")
    }
    return nil, result.Error
}

// Способ 2: Inline
if err := db.First(&user, id).Error; err != nil {
    return nil, err
}
```

### Частые ошибки:

```go
gorm.ErrRecordNotFound      // Запись не найдена
gorm.ErrInvalidTransaction  // Некорректная транзакция
gorm.ErrNotImplemented      // Не реализовано
gorm.ErrMissingWhereClause  // Отсутствует WHERE (защита)
gorm.ErrUnsupportedRelation // Неподдерживаемая связь
gorm.ErrPrimaryKeyRequired  // Требуется primary key
```

**Где используем:**
- Все методы в `internal/repository/user_repository.go`

---

## Advanced Features

### 1. RowsAffected

**Описание:** Количество затронутых строк

**Пример:**
```go
result := db.Delete(&User{}, id)
if result.RowsAffected == 0 {
    return errors.New("user not found")
}
```

**Где используем:**
- `internal/repository/user_repository.go:74`

---

### 2. Transaction

**Описание:** Выполнение в транзакции

**Пример:**
```go
db.Transaction(func(tx *gorm.DB) error {
    // Все операции в транзакции
    if err := tx.Create(&user1).Error; err != nil {
        return err  // Rollback
    }
    
    if err := tx.Create(&user2).Error; err != nil {
        return err  // Rollback
    }
    
    return nil  // Commit
})
```

---

### 3. Preload

**Описание:** Загрузка связанных данных

**Пример:**
```go
type User struct {
    ID    uint
    Posts []Post
}

// Загрузить пользователя с его постами
db.Preload("Posts").Find(&users)
```

---

### 4. Scopes

**Описание:** Переиспользуемые query условия

**Пример:**
```go
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

db.Scopes(ActiveUsers).Find(&users)
```

---

### 5. Raw SQL

**Описание:** Выполнение сырого SQL

**Пример:**
```go
db.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&user)

db.Exec("UPDATE users SET name = ? WHERE id = ?", name, id)
```

---

## Connection Pool

### Настройка пула подключений

```go
sqlDB, err := db.DB()

// Максимум idle connections
sqlDB.SetMaxIdleConns(10)

// Максимум open connections
sqlDB.SetMaxOpenConns(100)

// Максимальное время жизни connection
sqlDB.SetConnMaxLifetime(time.Hour)
```

**Где используем:**
- `internal/repository/database.go:74-76`

---

## Soft Delete

### gorm.DeletedAt

**Описание:** Автоматический soft delete

```go
type User struct {
    ID        uint
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**Поведение:**
- `Delete()` → SET deleted_at = NOW()
- Все queries автоматически добавляют `WHERE deleted_at IS NULL`
- `Unscoped()` → игнорирует deleted_at

**Примеры:**
```go
// Soft delete
db.Delete(&user)  // deleted_at = NOW()

// Найти всё (включая удалённое)
db.Unscoped().Find(&users)

// Permanent delete
db.Unscoped().Delete(&user)  // DELETE FROM users
```

**Где используем:**
- `internal/domain/user.go:31`

---

## Best Practices

### 1. Всегда проверяйте ошибки

```go
if err := db.Create(&user).Error; err != nil {
    return err
}
```

### 2. Используйте указатели

```go
db.First(&user, id)  // ✅ Правильно
db.First(user, id)   // ❌ Неправильно
```

### 3. Избегайте N+1 проблемы

```go
// ❌ N+1 queries
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts)
}

// ✅ 1 query
db.Preload("Posts").Find(&users)
```

### 4. Используйте плейсхолдеры (?) для защиты от SQL injection

```go
// ✅ Безопасно
db.Where("email = ?", userInput).Find(&user)

// ❌ Опасно!
db.Where(fmt.Sprintf("email = '%s'", userInput)).Find(&user)
```

---

## См. также

- [Gin Framework](./GIN.md) - HTTP web framework
- [PostgreSQL](https://www.postgresql.org/docs/) - База данных
- [Architecture Guide](../ARCHITECTURE.md) - Архитектура проекта
- [Testing Guide](../TESTING.md) - Тестирование с GORM

