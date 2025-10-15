# 🧪 Testing Guide

## Обзор

Проект включает два типа тестов:
- **Unit тесты** - тестирование отдельных компонентов с моками
- **Integration тесты** - тестирование полного flow с реальной БД

---

## 📁 Структура тестов

```
tests/
├── unit/
│   └── auth_service_test.go      # Unit тесты AuthService с моками
└── integration/
    └── api_test.go                # Integration тесты HTTP API
```

---

## 🔧 Unit тесты

### Что тестируется
- `AuthService.Register()` - регистрация пользователя
- `AuthService.Login()` - вход пользователя
- Бизнес-логика валидации
- Обработка ошибок

### Технологии
- **testify/assert** - assertions
- **testify/mock** - моки для repository

### Запуск
```bash
# Все unit тесты
go test -v ./tests/unit/...

# С покрытием
go test -v -cover ./tests/unit/...

# Детальное покрытие
go test -v -coverprofile=coverage.out ./tests/unit/...
go tool cover -html=coverage.out
```

### Пример теста
```go
func TestRegister_Success(t *testing.T) {
    // Arrange - подготовка моков
    mockRepo := new(MockUserRepository)
    authService := service.NewAuthService(mockRepo, cfg)
    
    // Act - выполнение действия
    response, err := authService.Register(req)
    
    // Assert - проверка результатов
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.NotEmpty(t, response.Token)
}
```

---

## 🌐 Integration тесты

### Что тестируется
- Полный HTTP flow: Register → Login → Get User
- Работа с реальной PostgreSQL БД
- Middleware (JWT authentication)
- Валидация запросов
- HTTP статус коды

### Требования
- **PostgreSQL** база данных для тестов
- Database: `advanced_api_test`

### Настройка БД для тестов

#### Вариант 1: Docker
```bash
docker run --name postgres-test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=advanced_api_test \
  -p 5433:5432 -d postgres:15-alpine
```

#### Вариант 2: Локальная PostgreSQL
```bash
psql -U postgres -c "CREATE DATABASE advanced_api_test;"
```

### Запуск
```bash
# Установите переменные окружения
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=advanced_api_test
export JWT_SECRET=test-secret

# Запустите тесты
go test -v ./tests/integration/...
```

### Пример теста
```go
func TestFullAuthFlow(t *testing.T) {
    // 1. Регистрация
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/auth/register", body)
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusCreated, w.Code)
    
    // 2. Вход
    // ...
    
    // 3. Защищённый endpoint
    req.Header.Set("Authorization", "Bearer "+token)
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

---

## 📊 Покрытие кода

### Генерация отчёта
```bash
# Общее покрытие
go test -v -coverprofile=coverage.out ./...

# HTML отчёт
go tool cover -html=coverage.out -o coverage.html

# Открыть в браузере
xdg-open coverage.html  # Linux
open coverage.html      # macOS
```

### Цель покрытия
- **Минимум**: 60%
- **Хорошо**: 80%
- **Отлично**: 90%+

---

## 🚀 CI/CD тесты

Тесты автоматически запускаются в GitHub Actions при каждом push:

```yaml
# .github/workflows/ci.yml
- name: Run unit tests
  run: go test -v -race -coverprofile=coverage.txt ./tests/unit/...

- name: Run integration tests
  run: go test -v ./tests/integration/...
```

---

## 🛠️ Написание новых тестов

### Unit тест шаблон
```go
func TestFeature_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)
    mockRepo.On("Method", args).Return(result, nil)
    
    // Act
    result, err := service.DoSomething(input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
    mockRepo.AssertExpectations(t)
}
```

### Integration тест шаблон
```go
func TestEndpoint_Success(t *testing.T) {
    // Setup
    router := setupTestRouter()
    
    // Request
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/endpoint", nil)
    
    // Execute
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    // ... проверка response body
}
```

---

## 🔍 Best Practices

### 1. Изоляция тестов
- Каждый тест должен быть независимым
- Используйте `cleanup` функции
- Не полагайтесь на порядок выполнения

### 2. Моки
- Мокируйте внешние зависимости (БД, HTTP)
- Используйте `testify/mock` для структурированных моков
- Проверяйте вызовы моков с `AssertExpectations`

### 3. Именование
- `Test<Function>_<Scenario>` - например `TestLogin_InvalidPassword`
- Описательные имена для понимания что тестируется

### 4. Arrange-Act-Assert
```go
// Arrange - подготовка
// Act - действие
// Assert - проверка
```

### 5. Table-Driven Tests
```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## 📚 Полезные команды

```bash
# Запустить все тесты
go test ./...

# Подробный вывод
go test -v ./...

# Только быстрые (unit) тесты
go test -short ./...

# Проверка race conditions
go test -race ./...

# Бенчмарки
go test -bench=. ./...

# Конкретный тест
go test -run TestRegister_Success ./tests/unit/...

# С таймаутом
go test -timeout 30s ./...
```

---

## 🐛 Отладка тестов

```bash
# Вывод всех логов
go test -v ./... 2>&1 | tee test.log

# Использование delve debugger
dlv test ./tests/unit -- -test.run TestRegister

# Печать в тестах
func TestSomething(t *testing.T) {
    t.Log("Debug message")
    fmt.Printf("Value: %+v\n", obj)
}
```

---

## 📖 Дополнительные ресурсы

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Advanced Testing Patterns](https://go.dev/blog/subtests)

