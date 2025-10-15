# 📋 План разработки Advanced API

## 🎯 Цель

Создать production-ready API с:
- JWT аутентификацией
- Docker контейнеризацией
- Автоматическими тестами
- CI/CD pipeline
- Swagger документацией

---

## 📅 Этапы разработки (5-7 дней)

### ✅ Этап 1: Настройка проекта (2-3 часа)

**Цель:** Создать базовую структуру с зависимостями

**Задачи:**
1. Создать структуру директорий
2. Инициализировать `go.mod` с зависимостями:
   - `github.com/gin-gonic/gin` - веб-фреймворк
   - `gorm.io/gorm` - ORM
   - `gorm.io/driver/postgres` - PostgreSQL драйвер для GORM
   - `github.com/golang-jwt/jwt/v5` - JWT
   - `golang.org/x/crypto/bcrypt` - хеширование паролей
   - `github.com/go-playground/validator/v10` - валидация
   - `go.uber.org/zap` - логирование
   - `github.com/spf13/viper` - конфигурация
3. Создать `docker-compose.yml`
4. Создать `Makefile` для команд

**Результат:** Рабочее окружение для разработки

---

### ✅ Этап 2: GORM модели и миграции (2-3 часа)

**Цель:** Настроить работу с БД через GORM

**Задачи:**
1. Создать GORM модели:
   ```go
   type User struct {
       ID        uint      `gorm:"primaryKey"`
       Email     string    `gorm:"uniqueIndex;not null"`
       Name      string    `gorm:"not null"`
       Password  string    `gorm:"not null"`
       Role      string    `gorm:"default:'user'"`
       CreatedAt time.Time
       UpdatedAt time.Time
   }
   ```

2. Настроить Auto Migration
3. Создать Repository с GORM методами:
   - `Create`, `FindByID`, `FindAll`, `Update`, `Delete`
   - `FindByEmail` (для аутентификации)

**Результат:** Работа с БД через GORM

---

### ✅ Этап 3: Gin роутер и базовые endpoints (2-3 часа)

**Цель:** Настроить Gin и создать CRUD endpoints

**Задачи:**
1. Настроить Gin роутер
2. Создать группу routes `/api/v1`
3. Реализовать User handlers:
   - GET `/users` - список
   - GET `/users/:id` - один пользователь
   - PUT `/users/:id` - обновление
   - DELETE `/users/:id` - удаление
4. Добавить базовые middleware:
   - CORS
   - Logger
   - Recovery

**Результат:** Рабочий REST API на Gin

---

### ✅ Этап 4: JWT Аутентификация (3-4 часа)

**Цель:** Реализовать регистрацию, вход и защиту endpoints

**Задачи:**
1. Создать JWT утилиты:
   - `GenerateToken(userID, email)` - генерация токена
   - `ValidateToken(tokenString)` - валидация токена
   - `ExtractClaims(token)` - извлечение данных

2. Реализовать Auth Service:
   - `Register(email, name, password)` - регистрация с хешированием пароля
   - `Login(email, password)` - вход с проверкой пароля
   - `GetCurrentUser(userID)` - получение текущего пользователя

3. Создать Auth Handlers:
   - POST `/auth/register`
   - POST `/auth/login`
   - GET `/auth/me` (требует токен)

4. Реализовать Auth Middleware:
   - Проверка наличия токена в заголовке
   - Валидация токена
   - Извлечение userID в контекст

**Результат:** Полная JWT аутентификация

---

### ✅ Этап 5: Middleware и улучшения (2-3 часа)

**Цель:** Добавить production-ready middleware

**Задачи:**
1. **Logging Middleware:**
   - Логирование каждого запроса
   - Время выполнения
   - Статус код
   - Использование zap logger

2. **CORS Middleware:**
   - Настройка разрешённых origins
   - Поддержка preflight запросов

3. **Rate Limiting:**
   - Ограничение запросов (например, 100/мин)
   - Защита от DDoS

4. **Request ID:**
   - Уникальный ID для каждого запроса
   - Для трейсинга

5. **Validator:**
   - Кастомные правила валидации
   - Красивые сообщения об ошибках

**Результат:** Production-ready middleware

---

### ✅ Этап 6: Тестирование (4-5 часов)

**Цель:** Покрыть код тестами

**Задачи:**
1. **Unit тесты для Service:**
   ```go
   func TestUserService_CreateUser(t *testing.T) {
       // Arrange
       mockRepo := &MockUserRepository{}
       service := NewUserService(mockRepo)
       
       // Act
       user, err := service.CreateUser(...)
       
       // Assert
       assert.NoError(t, err)
       assert.NotNil(t, user)
   }
   ```

2. **Unit тесты для Handlers:**
   - Используя `httptest`
   - Моки для service

3. **Integration тесты:**
   - Тестовая БД
   - Полный цикл (создание → чтение → обновление → удаление)

4. **JWT тесты:**
   - Генерация токена
   - Валидация
   - Истёкший токен

**Результат:** > 80% code coverage

---

### ✅ Этап 7: Swagger документация (1-2 часа)

**Цель:** Автоматическая документация API

**Задачи:**
1. Добавить swagger аннотации:
   ```go
   // @Summary Create user
   // @Tags users
   // @Accept json
   // @Produce json
   // @Param user body CreateUserRequest true "User data"
   // @Success 201 {object} User
   // @Router /users [post]
   func (h *UserHandler) CreateUser(c *gin.Context) {
       // ...
   }
   ```

2. Генерация документации: `swag init`
3. Swagger UI на `/swagger/*`

**Результат:** Интерактивная документация API

---

### ✅ Этап 8: Docker и деплой (2-3 часа)

**Цель:** Контейнеризация и автоматический деплой

**Задачи:**
1. **Dockerfile (multi-stage build):**
   - Builder stage (компиляция)
   - Runtime stage (минимальный образ)

2. **docker-compose.yml:**
   - API сервис
   - PostgreSQL
   - Redis (для кеша)
   - pgAdmin (для управления БД)

3. **GitHub Actions CI/CD:**
   - Запуск тестов
   - Сборка Docker образа
   - Деплой на сервер

4. **Деплой на VPS:**
   - Docker Compose на сервере
   - Nginx reverse proxy
   - SSL сертификат (Let's Encrypt)

**Результат:** Автоматический деплой в production

---

## 🛠️ Инструменты разработки

### Makefile команды:

```bash
make run          # Запустить локально
make test         # Запустить тесты
make test-coverage # Тесты с coverage
make docker-build # Собрать Docker образ
make docker-up    # Запустить через Docker Compose
make docker-down  # Остановить контейнеры
make migrate-up   # Применить миграции
make migrate-down # Откатить миграции
make swagger      # Генерация Swagger документации
make lint         # Проверка кода
make fmt          # Форматирование
```

---

## 📊 Метрики проекта (после завершения)

- **Строк кода:** ~2000-2500
- **Файлов:** ~25-30
- **Покрытие тестами:** >80%
- **Endpoints:** ~10-12
- **Middleware:** ~5-6
- **Docker services:** 3-4

---

## 🎯 Готовы начать?

**Порядок работы:**
1. Я создам базовую структуру и файлы
2. Вы будете реализовывать по шагам
3. Я проверяю и объясняю
4. Двигаемся дальше

**Начинаем с Этапа 1?** 🚀

