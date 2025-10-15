# 🏗️ Architecture Documentation

## Clean Architecture Overview

Проект следует принципам **Clean Architecture** с разделением на слои:

```
┌─────────────────────────────────────────────┐
│            HTTP Request (Gin)               │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│         Handler Layer (HTTP)                │
│  • Парсинг JSON                             │
│  • Валидация входных данных                 │
│  • HTTP статус коды                         │
│  • Преобразование в DTO                     │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│         Service Layer (Business Logic)      │
│  • Бизнес-правила                           │
│  • JWT генерация/валидация                  │
│  • Bcrypt хеширование                       │
│  • Обработка ошибок                         │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│       Repository Layer (Data Access)        │
│  • GORM операции                            │
│  • SQL запросы                              │
│  • Транзакции                               │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│            PostgreSQL Database              │
└─────────────────────────────────────────────┘
```

---

## 📂 Project Structure

```
advanced-user-api/
│
├── cmd/
│   └── api/
│       └── main.go                 # Entry point
│
├── internal/                       # Private application code
│   ├── config/
│   │   └── config.go              # Configuration (Viper)
│   │
│   ├── domain/
│   │   └── user.go                # Domain models & DTOs
│   │
│   ├── handler/                   # HTTP Layer
│   │   ├── auth_handler.go        # Auth endpoints
│   │   ├── user_handler.go        # User endpoints
│   │   └── routes.go              # Route setup
│   │
│   ├── service/                   # Business Logic Layer
│   │   ├── auth_service.go        # Authentication logic
│   │   └── user_service.go        # User business logic
│   │
│   ├── repository/                # Data Access Layer
│   │   ├── database.go            # DB connection & migrations
│   │   └── user_repository.go     # User CRUD operations
│   │
│   ├── middleware/                # Middleware components
│   │   ├── auth.go                # JWT authentication
│   │   ├── cors.go                # CORS headers
│   │   └── logger.go              # Request logging
│   │
│   └── pkg/                       # Shared utilities
│       ├── jwt/                   # JWT utilities
│       ├── password/              # Password hashing
│       └── validator/             # Custom validators
│
├── tests/
│   ├── unit/                      # Unit tests with mocks
│   └── integration/               # Integration tests
│
├── docs/                          # Documentation
├── docker/                        # Docker files
├── migrations/                    # SQL migrations
└── .github/workflows/             # CI/CD pipelines
```

---

## 🔄 Request Flow

### Example: User Registration

```
1. HTTP Request
   POST /api/v1/auth/register
   Body: {"email": "...", "name": "...", "password": "..."}
   
   ↓
   
2. Middleware Chain
   • CORS
   • Logger
   
   ↓
   
3. Handler (auth_handler.go)
   • Parse JSON → RegisterRequest DTO
   • Validate via Gin binding tags
   
   ↓
   
4. Service (auth_service.go)
   • Business validation (email format, password length)
   • Check if email already exists (via Repository)
   • Hash password with bcrypt
   • Create user (via Repository)
   • Generate JWT token
   
   ↓
   
5. Repository (user_repository.go)
   • GORM Create operation
   • Auto-fill timestamps (created_at, updated_at)
   • Return user with ID
   
   ↓
   
6. Database (PostgreSQL)
   • INSERT INTO users ...
   • Return generated ID
   
   ↓
   
7. Response Flow (back up)
   Repository → Service → Handler
   
   ↓
   
8. HTTP Response
   Status: 201 Created
   Body: {"token": "...", "user": {...}}
```

---

## 🎯 Layer Responsibilities

### 1. Handler Layer (`internal/handler/`)

**Responsibilities:**
- HTTP request/response handling
- JSON parsing & encoding
- HTTP status codes
- Input validation (basic)
- DTOs (Data Transfer Objects)

**Does NOT:**
- Contain business logic
- Access database directly
- Know about other handlers

**Example:**
```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req domain.RegisterRequest
    
    // Parse & validate
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Call service
    response, err := h.authService.Register(&req)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Return response
    c.JSON(201, response)
}
```

---

### 2. Service Layer (`internal/service/`)

**Responsibilities:**
- Business logic & rules
- Data validation (business rules)
- Orchestration between repositories
- Error handling
- JWT generation
- Password hashing

**Does NOT:**
- Know about HTTP (no `gin.Context`)
- Access database directly (uses Repository)
- Know about other services directly

**Example:**
```go
func (s *authService) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
    // Business validation
    if len(req.Password) < 6 {
        return nil, errors.New("password too short")
    }
    
    // Check if exists
    existing, _ := s.userRepo.FindByEmail(req.Email)
    if existing != nil {
        return nil, errors.New("email already exists")
    }
    
    // Hash password
    hashedPassword, err := password.Hash(req.Password)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user := &domain.User{
        Email:    req.Email,
        Name:     req.Name,
        Password: hashedPassword,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    // Generate token
    token, err := jwt.Generate(user)
    
    return &domain.AuthResponse{
        Token: token,
        User:  user,
    }, nil
}
```

---

### 3. Repository Layer (`internal/repository/`)

**Responsibilities:**
- Database operations (CRUD)
- GORM queries
- Transactions
- Data mapping

**Does NOT:**
- Contain business logic
- Validate business rules
- Know about HTTP

**Example:**
```go
func (r *userRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

---

## 🔌 Dependency Injection

Зависимости внедряются через конструкторы:

```go
// main.go
func main() {
    // 1. Database
    db := repository.InitDB(cfg)
    
    // 2. Repositories (depend on DB)
    userRepo := repository.NewUserRepository(db)
    
    // 3. Services (depend on Repositories)
    authService := service.NewAuthService(userRepo, cfg)
    userService := service.NewUserService(userRepo)
    
    // 4. Handlers (depend on Services)
    authHandler := handler.NewAuthHandler(authService, userService)
    userHandler := handler.NewUserHandler(userService)
    
    // 5. Routes (depend on Handlers)
    handler.SetupRoutes(router, authHandler, userHandler, cfg)
}
```

**Преимущества:**
- ✅ Легко тестировать (моки)
- ✅ Слабая связанность
- ✅ Можно заменить реализацию
- ✅ Понятный flow зависимостей

---

## 🧩 Key Components

### Middleware

```go
// Middleware chain
router.Use(
    middleware.CORSMiddleware(),      // CORS headers
    gin.Logger(),                     // Request logging
    gin.Recovery(),                   // Panic recovery
)

// Protected routes
protected := router.Group("/api/v1")
protected.Use(middleware.AuthMiddleware(cfg))  // JWT validation
```

### Configuration

```go
// Viper for configuration management
type Config struct {
    ServerPort    string
    DBHost        string
    DBPort        string
    JWTSecret     string
    JWTExpiration string
}

// Load from .env or environment variables
cfg := config.Load()
```

### DTOs (Data Transfer Objects)

```go
// Request DTOs
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Name     string `json:"name" binding:"required,min=2"`
    Password string `json:"password" binding:"required,min=6"`
}

// Response DTOs
type AuthResponse struct {
    Token string `json:"token"`
    User  *User  `json:"user"`
}
```

---

## 🔐 Security

### Password Hashing
```go
// Never store plain passwords
hashedPassword, _ := bcrypt.GenerateFromPassword(
    []byte(password), 
    bcrypt.DefaultCost,
)
```

### JWT Authentication
```go
// Generate token
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, _ := token.SignedString([]byte(secret))

// Validate token
token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(secret), nil
})
```

### CORS
```go
// Allow cross-origin requests
c.Header("Access-Control-Allow-Origin", "*")
c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
```

---

## 📊 Database Schema

### Users Table
```sql
CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(255) UNIQUE NOT NULL,
    name       VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL        -- Soft delete
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

### GORM Auto Migration
```go
db.AutoMigrate(&domain.User{})
// Automatically creates/updates table based on struct
```

---

## 🧪 Testing Strategy

### Unit Tests
- Test each layer independently
- Use mocks for dependencies
- Focus on business logic

### Integration Tests
- Test full HTTP flow
- Real database (test DB)
- All layers together

Подробнее в [`TESTING.md`](./TESTING.md)

---

## 🚀 Deployment

### Docker Multi-Stage Build
```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api cmd/api/main.go

# Stage 2: Runtime
FROM alpine:latest
COPY --from=builder /app/api .
CMD ["./api"]
```

Подробнее в [`DEPLOY.md`](./DEPLOY.md)

---

## 📖 Additional Resources

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

