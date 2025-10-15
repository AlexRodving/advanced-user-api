# 📝 Расширение проекта: От User API к Blog Platform

## Обзор

Этот гайд показывает как расширить текущий User Management API в полноценную **Blog Platform** с постами, комментариями и системой ролей.

**Что добавим:**
- 📝 Posts (посты)
- 💬 Comments (комментарии)
- 👥 Roles (роли: admin, moderator, user)
- 🔒 Permissions (разрешения по ролям)
- 📁 Categories (категории постов)
- 🏷️ Tags (теги)
- ❤️ Likes (лайки)

---

## 🎯 Архитектура блога

```
Users (пользователи)
  ↓ has many
Posts (посты)
  ↓ has many
Comments (комментарии)
  ↓ belongs to
Users (автор комментария)
```

### Связи между моделями:

```
User (1) ←→ (N) Posts       # Один пользователь → много постов
Post (1) ←→ (N) Comments    # Один пост → много комментариев
User (1) ←→ (N) Comments    # Один пользователь → много комментариев
Post (N) ←→ (N) Tags        # Многие-ко-многим (через join таблицу)
Post (N) ←→ (1) Category    # Много постов → одна категория
```

---

## 📊 Новые модели данных

### 1. Role Enum (роли пользователей)

**Обновляем:** `internal/domain/user.go`

```go
package domain

// Константы для ролей
const (
    RoleUser      = "user"       // Обычный пользователь (по умолчанию)
    RoleModerator = "moderator"  // Модератор (может удалять комментарии)
    RoleAdmin     = "admin"      // Администратор (полный доступ)
)

// User - обновлённая модель
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"uniqueIndex;not null" json:"email"`
    Name      string         `gorm:"not null" json:"name"`
    Password  string         `gorm:"not null" json:"-"`
    Role      string         `gorm:"type:varchar(20);default:'user';check:role IN ('user','moderator','admin')" json:"role"`
    Bio       string         `gorm:"type:text" json:"bio"`                    // НОВОЕ: Биография
    Avatar    string         `gorm:"type:varchar(500)" json:"avatar"`         // НОВОЕ: URL аватара
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    // НОВОЕ: Связи
    Posts    []Post    `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`    // Посты пользователя
    Comments []Comment `gorm:"foreignKey:AuthorID" json:"comments,omitempty"` // Комментарии
}

// HasRole - проверка роли
func (u *User) HasRole(role string) bool {
    return u.Role == role
}

// IsAdmin - проверка админских прав
func (u *User) IsAdmin() bool {
    return u.Role == RoleAdmin
}

// IsModerator - проверка прав модератора (или выше)
func (u *User) IsModerator() bool {
    return u.Role == RoleModerator || u.Role == RoleAdmin
}
```

---

### 2. Post (пост)

**Создаём:** `internal/domain/post.go`

```go
package domain

import (
    "time"
    "gorm.io/gorm"
)

// Post - модель поста в блоге
type Post struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Title       string         `gorm:"type:varchar(200);not null;index" json:"title"`
    Slug        string         `gorm:"type:varchar(250);uniqueIndex;not null" json:"slug"` // URL-friendly название
    Content     string         `gorm:"type:text;not null" json:"content"`
    Excerpt     string         `gorm:"type:varchar(500)" json:"excerpt"`                   // Краткое описание
    CoverImage  string         `gorm:"type:varchar(500)" json:"cover_image"`               // URL обложки
    Status      string         `gorm:"type:varchar(20);default:'draft';check:status IN ('draft','published','archived')" json:"status"`
    ViewCount   int            `gorm:"default:0" json:"view_count"`                        // Счётчик просмотров
    
    // Связи
    AuthorID    uint           `gorm:"not null;index" json:"author_id"`
    Author      User           `gorm:"foreignKey:AuthorID" json:"author"`                  // Автор поста
    
    CategoryID  *uint          `gorm:"index" json:"category_id"`                           // Может быть NULL
    Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    
    Comments    []Comment      `gorm:"foreignKey:PostID" json:"comments,omitempty"`        // Комментарии
    Tags        []Tag          `gorm:"many2many:post_tags;" json:"tags,omitempty"`         // Многие-ко-многим
    Likes       []Like         `gorm:"foreignKey:PostID" json:"likes,omitempty"`           // Лайки
    
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    PublishedAt *time.Time     `json:"published_at,omitempty"`                             // Время публикации
}

// TableName - имя таблицы
func (Post) TableName() string {
    return "posts"
}

// IsPublished - проверка опубликован ли пост
func (p *Post) IsPublished() bool {
    return p.Status == "published"
}

// CanEdit - может ли пользователь редактировать пост
func (p *Post) CanEdit(user *User) bool {
    // Автор или админ могут редактировать
    return p.AuthorID == user.ID || user.IsAdmin()
}

// CanDelete - может ли пользователь удалить пост
func (p *Post) CanDelete(user *User) bool {
    // Автор, модератор или админ могут удалять
    return p.AuthorID == user.ID || user.IsModerator()
}
```

**DTOs для постов:**

```go
// CreatePostRequest - создание поста
type CreatePostRequest struct {
    Title      string  `json:"title" binding:"required,min=3,max=200"`
    Content    string  `json:"content" binding:"required,min=10"`
    Excerpt    string  `json:"excerpt" binding:"omitempty,max=500"`
    CoverImage string  `json:"cover_image" binding:"omitempty,url"`
    CategoryID *uint   `json:"category_id" binding:"omitempty"`
    TagIDs     []uint  `json:"tag_ids" binding:"omitempty"`               // Массив ID тегов
    Status     string  `json:"status" binding:"omitempty,oneof=draft published"`
}

// UpdatePostRequest - обновление поста
type UpdatePostRequest struct {
    Title      string  `json:"title" binding:"omitempty,min=3,max=200"`
    Content    string  `json:"content" binding:"omitempty,min=10"`
    Excerpt    string  `json:"excerpt" binding:"omitempty,max=500"`
    CoverImage string  `json:"cover_image" binding:"omitempty,url"`
    CategoryID *uint   `json:"category_id"`
    TagIDs     []uint  `json:"tag_ids"`
    Status     string  `json:"status" binding:"omitempty,oneof=draft published archived"`
}

// PostListResponse - список постов с пагинацией
type PostListResponse struct {
    Posts      []Post `json:"posts"`
    TotalCount int64  `json:"total_count"`
    Page       int    `json:"page"`
    PageSize   int    `json:"page_size"`
    TotalPages int    `json:"total_pages"`
}
```

---

### 3. Comment (комментарий)

**Создаём:** `internal/domain/comment.go`

```go
package domain

import (
    "time"
    "gorm.io/gorm"
)

// Comment - модель комментария
type Comment struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Content   string         `gorm:"type:text;not null" json:"content"`
    
    // Связи
    PostID    uint           `gorm:"not null;index" json:"post_id"`
    Post      Post           `gorm:"foreignKey:PostID" json:"post,omitempty"`
    
    AuthorID  uint           `gorm:"not null;index" json:"author_id"`
    Author    User           `gorm:"foreignKey:AuthorID" json:"author"`
    
    // Вложенные комментарии (ответы на комментарии)
    ParentID  *uint          `gorm:"index" json:"parent_id,omitempty"`                    // NULL если корневой
    Parent    *Comment       `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Replies   []Comment      `gorm:"foreignKey:ParentID" json:"replies,omitempty"`        // Ответы
    
    // Модерация
    IsEdited  bool           `gorm:"default:false" json:"is_edited"`
    IsBlocked bool           `gorm:"default:false" json:"is_blocked"`                     // Заблокирован модератором
    
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CanEdit - может ли пользователь редактировать комментарий
func (c *Comment) CanEdit(user *User) bool {
    return c.AuthorID == user.ID
}

// CanDelete - может ли пользователь удалить комментарий
func (c *Comment) CanDelete(user *User) bool {
    // Автор, модератор или админ
    return c.AuthorID == user.ID || user.IsModerator()
}

// CanBlock - может ли пользователь заблокировать комментарий
func (c *Comment) CanBlock(user *User) bool {
    // Только модераторы и админы
    return user.IsModerator()
}
```

**DTOs:**

```go
// CreateCommentRequest - создание комментария
type CreateCommentRequest struct {
    Content  string `json:"content" binding:"required,min=1,max=5000"`
    PostID   uint   `json:"post_id" binding:"required"`
    ParentID *uint  `json:"parent_id" binding:"omitempty"`  // Для ответов на комментарии
}

// UpdateCommentRequest - обновление комментария
type UpdateCommentRequest struct {
    Content string `json:"content" binding:"required,min=1,max=5000"`
}
```

---

### 4. Category (категория)

**Создаём:** `internal/domain/category.go`

```go
package domain

import (
    "time"
    "gorm.io/gorm"
)

// Category - категория постов
type Category struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
    Slug        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
    Description string         `gorm:"type:text" json:"description"`
    
    Posts       []Post         `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
    
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Примеры категорий: "Технологии", "Новости", "Туториалы"
```

---

### 5. Tag (тег)

**Создаём:** `internal/domain/tag.go`

```go
package domain

import "time"

// Tag - тег для постов
type Tag struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
    Slug      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
    
    Posts     []Post    `gorm:"many2many:post_tags;" json:"posts,omitempty"`
    
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Примеры тегов: "golang", "api", "tutorial", "docker"
```

---

### 6. Like (лайк)

**Создаём:** `internal/domain/like.go`

```go
package domain

import "time"

// Like - лайк поста
type Like struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    
    PostID    uint      `gorm:"not null;index:idx_post_user" json:"post_id"`
    Post      Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
    
    UserID    uint      `gorm:"not null;index:idx_post_user" json:"user_id"`
    User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    
    CreatedAt time.Time `json:"created_at"`
}

// TableName - имя таблицы
func (Like) TableName() string {
    return "likes"
}

// Уникальный индекс: один пользователь может лайкнуть пост только раз
// В GORM: `gorm:"uniqueIndex:idx_post_user_unique"`
// Или в миграции: CREATE UNIQUE INDEX idx_post_user_unique ON likes(post_id, user_id)
```

---

## 📝 API Endpoints для блога

### Posts

```
GET    /api/v1/posts              # Список постов (с пагинацией, фильтрами)
GET    /api/v1/posts/:id          # Получить пост
POST   /api/v1/posts              # Создать пост (auth required)
PUT    /api/v1/posts/:id          # Обновить пост (auth + ownership)
DELETE /api/v1/posts/:id          # Удалить пост (auth + ownership/moderator)
POST   /api/v1/posts/:id/publish  # Опубликовать пост
GET    /api/v1/posts/slug/:slug   # Получить пост по slug
```

### Comments

```
GET    /api/v1/posts/:id/comments        # Комментарии к посту
POST   /api/v1/posts/:id/comments        # Добавить комментарий (auth)
PUT    /api/v1/comments/:id              # Редактировать комментарий (auth + ownership)
DELETE /api/v1/comments/:id              # Удалить комментарий (auth + ownership/moderator)
POST   /api/v1/comments/:id/block        # Заблокировать комментарий (moderator)
```

### Likes

```
POST   /api/v1/posts/:id/like    # Лайкнуть пост (auth)
DELETE /api/v1/posts/:id/like    # Убрать лайк (auth)
GET    /api/v1/posts/:id/likes   # Количество лайков
```

### Categories & Tags

```
GET    /api/v1/categories         # Список категорий
POST   /api/v1/categories         # Создать категорию (admin)
GET    /api/v1/tags               # Список тегов
POST   /api/v1/tags               # Создать тег (admin)
```

---

## 🗄️ Database Schema (PostgreSQL)

### Posts Table

```sql
CREATE TABLE posts (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(200) NOT NULL,
    slug         VARCHAR(250) UNIQUE NOT NULL,
    content      TEXT NOT NULL,
    excerpt      VARCHAR(500),
    cover_image  VARCHAR(500),
    status       VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    view_count   INTEGER DEFAULT 0,
    author_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id  INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW(),
    deleted_at   TIMESTAMP NULL,
    published_at TIMESTAMP NULL
);

CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_category_id ON posts(category_id);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at);
```

### Comments Table

```sql
CREATE TABLE comments (
    id         SERIAL PRIMARY KEY,
    content    TEXT NOT NULL,
    post_id    INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id  INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    is_edited  BOOLEAN DEFAULT FALSE,
    is_blocked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
```

### Likes Table

```sql
CREATE TABLE likes (
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(post_id, user_id)  -- Один пользователь = один лайк
);

CREATE INDEX idx_likes_post_id ON likes(post_id);
CREATE INDEX idx_likes_user_id ON likes(user_id);
```

### Categories Table

```sql
CREATE TABLE categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) UNIQUE NOT NULL,
    slug        VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP NULL
);
```

### Tags Table

```sql
CREATE TABLE tags (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(50) UNIQUE NOT NULL,
    slug       VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Post_Tags (Many-to-Many)

```sql
CREATE TABLE post_tags (
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);
```

---

## 🔒 Middleware для проверки ролей

**Создаём:** `internal/middleware/role.go`

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "advanced-user-api/internal/domain"
    "advanced-user-api/internal/service"
)

// RequireRole - middleware для проверки роли
func RequireRole(userService service.UserService, roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Получаем user_id из контекста (установлен AuthMiddleware)
        userID, exists := c.Get("user_id")
        if !exists {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        
        // Получаем пользователя из БД
        user, err := userService.GetUser(userID.(uint))
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "user not found"})
            return
        }
        
        // Проверяем роль
        hasRole := false
        for _, role := range roles {
            if user.Role == role {
                hasRole = true
                break
            }
        }
        
        if !hasRole {
            c.AbortWithStatusJSON(403, gin.H{"error": "insufficient permissions"})
            return
        }
        
        // Сохраняем user в контекст для использования в handlers
        c.Set("user", user)
        c.Next()
    }
}

// RequireAdmin - только для админов
func RequireAdmin(userService service.UserService) gin.HandlerFunc {
    return RequireRole(userService, domain.RoleAdmin)
}

// RequireModerator - для модераторов и админов
func RequireModerator(userService service.UserService) gin.HandlerFunc {
    return RequireRole(userService, domain.RoleModerator, domain.RoleAdmin)
}
```

**Использование:**

```go
// internal/handler/routes.go

// Только админы могут создавать категории
adminRoutes := v1.Group("/admin")
adminRoutes.Use(middleware.RequireAdmin(userService))
{
    adminRoutes.POST("/categories", categoryHandler.Create)
    adminRoutes.DELETE("/users/:id", userHandler.Delete)
}

// Модераторы могут блокировать комментарии
moderatorRoutes := v1.Group("/moderator")
moderatorRoutes.Use(middleware.RequireModerator(userService))
{
    moderatorRoutes.POST("/comments/:id/block", commentHandler.Block)
    moderatorRoutes.DELETE("/posts/:id", postHandler.Delete)
}
```

---

## 🔧 Repository Layer

### PostRepository

**Создаём:** `internal/repository/post_repository.go`

```go
package repository

import (
    "advanced-user-api/internal/domain"
    "gorm.io/gorm"
)

type PostRepository interface {
    Create(post *domain.Post) error
    FindByID(id uint) (*domain.Post, error)
    FindBySlug(slug string) (*domain.Post, error)
    FindAll(page, pageSize int, filters PostFilters) ([]domain.Post, int64, error)
    Update(post *domain.Post) error
    Delete(id uint) error
    IncrementViewCount(id uint) error
    FindByAuthor(authorID uint, page, pageSize int) ([]domain.Post, int64, error)
    FindByCategory(categoryID uint, page, pageSize int) ([]domain.Post, int64, error)
    FindByTag(tagID uint, page, pageSize int) ([]domain.Post, int64, error)
}

type PostFilters struct {
    Status     string  // "published", "draft", "archived"
    AuthorID   *uint
    CategoryID *uint
    Search     string  // Поиск по title и content
}

type postRepository struct {
    db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
    return &postRepository{db: db}
}

// Create - создаёт новый пост
func (r *postRepository) Create(post *domain.Post) error {
    return r.db.Create(post).Error
}

// FindByID - получает пост с предзагрузкой связей
func (r *postRepository) FindByID(id uint) (*domain.Post, error) {
    var post domain.Post
    err := r.db.
        Preload("Author").           // Загружаем автора
        Preload("Category").         // Загружаем категорию
        Preload("Tags").             // Загружаем теги
        Preload("Comments.Author").  // Загружаем комментарии с авторами
        First(&post, id).Error
    
    if err != nil {
        return nil, err
    }
    return &post, nil
}

// FindAll - список постов с фильтрами и пагинацией
func (r *postRepository) FindAll(page, pageSize int, filters PostFilters) ([]domain.Post, int64, error) {
    var posts []domain.Post
    var totalCount int64
    
    query := r.db.Model(&domain.Post{})
    
    // Применяем фильтры
    if filters.Status != "" {
        query = query.Where("status = ?", filters.Status)
    }
    
    if filters.AuthorID != nil {
        query = query.Where("author_id = ?", *filters.AuthorID)
    }
    
    if filters.CategoryID != nil {
        query = query.Where("category_id = ?", *filters.CategoryID)
    }
    
    if filters.Search != "" {
        search := "%" + filters.Search + "%"
        query = query.Where("title LIKE ? OR content LIKE ?", search, search)
    }
    
    // Подсчёт общего количества
    query.Count(&totalCount)
    
    // Пагинация
    offset := (page - 1) * pageSize
    err := query.
        Preload("Author").
        Preload("Category").
        Preload("Tags").
        Order("created_at DESC").
        Limit(pageSize).
        Offset(offset).
        Find(&posts).Error
    
    return posts, totalCount, err
}

// IncrementViewCount - увеличивает счётчик просмотров
func (r *postRepository) IncrementViewCount(id uint) error {
    return r.db.Model(&domain.Post{}).
        Where("id = ?", id).
        UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).
        Error
}

// FindBySlug - поиск по slug (для SEO-friendly URLs)
func (r *postRepository) FindBySlug(slug string) (*domain.Post, error) {
    var post domain.Post
    err := r.db.
        Preload("Author").
        Preload("Category").
        Preload("Tags").
        Where("slug = ?", slug).
        First(&post).Error
    
    if err != nil {
        return nil, err
    }
    return &post, nil
}
```

---

### CommentRepository

**Создаём:** `internal/repository/comment_repository.go`

```go
package repository

import (
    "advanced-user-api/internal/domain"
    "gorm.io/gorm"
)

type CommentRepository interface {
    Create(comment *domain.Comment) error
    FindByID(id uint) (*domain.Comment, error)
    FindByPost(postID uint) ([]domain.Comment, error)
    Update(comment *domain.Comment) error
    Delete(id uint) error
    Block(id uint) error
    Unblock(id uint) error
}

type commentRepository struct {
    db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
    return &commentRepository{db: db}
}

// FindByPost - получает комментарии поста с вложенной структурой
func (r *commentRepository) FindByPost(postID uint) ([]domain.Comment, error) {
    var comments []domain.Comment
    
    // Получаем только корневые комментарии (parent_id IS NULL)
    err := r.db.
        Where("post_id = ? AND parent_id IS NULL", postID).
        Preload("Author").                    // Автор
        Preload("Replies.Author").            // Ответы с авторами
        Preload("Replies.Replies.Author").    // Вложенные ответы (3 уровня)
        Order("created_at DESC").
        Find(&comments).Error
    
    return comments, err
}

// Block - блокирует комментарий (модератор)
func (r *commentRepository) Block(id uint) error {
    return r.db.Model(&domain.Comment{}).
        Where("id = ?", id).
        Update("is_blocked", true).Error
}
```

---

## 💼 Service Layer

### PostService

**Создаём:** `internal/service/post_service.go`

```go
package service

import (
    "errors"
    "time"
    "advanced-user-api/internal/domain"
    "advanced-user-api/internal/repository"
)

type PostService interface {
    CreatePost(authorID uint, req *domain.CreatePostRequest) (*domain.Post, error)
    GetPost(id uint) (*domain.Post, error)
    GetPostBySlug(slug string) (*domain.Post, error)
    GetAllPosts(page, pageSize int, filters repository.PostFilters) (*domain.PostListResponse, error)
    UpdatePost(id uint, userID uint, req *domain.UpdatePostRequest) (*domain.Post, error)
    DeletePost(id uint, user *domain.User) error
    PublishPost(id uint, userID uint) (*domain.Post, error)
}

type postService struct {
    postRepo repository.PostRepository
    userRepo repository.UserRepository
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) PostService {
    return &postService{
        postRepo: postRepo,
        userRepo: userRepo,
    }
}

// CreatePost - создание поста
func (s *postService) CreatePost(authorID uint, req *domain.CreatePostRequest) (*domain.Post, error) {
    // Валидация
    if req.Title == "" {
        return nil, errors.New("title обязателен")
    }
    
    // Генерация slug из title
    slug := generateSlug(req.Title)
    
    post := &domain.Post{
        Title:      req.Title,
        Slug:       slug,
        Content:    req.Content,
        Excerpt:    req.Excerpt,
        CoverImage: req.CoverImage,
        Status:     req.Status,
        AuthorID:   authorID,
        CategoryID: req.CategoryID,
    }
    
    // Если статус "published", устанавливаем время публикации
    if post.Status == "published" {
        now := time.Now()
        post.PublishedAt = &now
    }
    
    // Создаём пост
    if err := s.postRepo.Create(post); err != nil {
        return nil, err
    }
    
    // Если указаны теги, добавляем их
    if len(req.TagIDs) > 0 {
        // TODO: Добавить теги через TagRepository
    }
    
    return post, nil
}

// UpdatePost - обновление поста с проверкой прав
func (s *postService) UpdatePost(id uint, userID uint, req *domain.UpdatePostRequest) (*domain.Post, error) {
    // Получаем пост
    post, err := s.postRepo.FindByID(id)
    if err != nil {
        return nil, errors.New("пост не найден")
    }
    
    // Получаем пользователя
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("пользователь не найден")
    }
    
    // Проверяем права на редактирование
    if !post.CanEdit(user) {
        return nil, errors.New("нет прав на редактирование")
    }
    
    // Обновляем поля
    if req.Title != "" {
        post.Title = req.Title
        post.Slug = generateSlug(req.Title)
    }
    if req.Content != "" {
        post.Content = req.Content
    }
    if req.Status != "" {
        post.Status = req.Status
    }
    
    // Обновляем в БД
    if err := s.postRepo.Update(post); err != nil {
        return nil, err
    }
    
    return post, nil
}

// DeletePost - удаление поста с проверкой прав
func (s *postService) DeletePost(id uint, user *domain.User) error {
    post, err := s.postRepo.FindByID(id)
    if err != nil {
        return errors.New("пост не найден")
    }
    
    // Проверяем права на удаление
    if !post.CanDelete(user) {
        return errors.New("нет прав на удаление")
    }
    
    return s.postRepo.Delete(id)
}

// Вспомогательная функция для генерации slug
func generateSlug(title string) string {
    // Упрощённая версия:
    // В production используйте библиотеку: github.com/gosimple/slug
    
    slug := strings.ToLower(title)
    slug = strings.ReplaceAll(slug, " ", "-")
    // Удаляем специальные символы
    slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")
    
    return slug
}
```

---

## 🎨 Handler Layer

### PostHandler

**Создаём:** `internal/handler/post_handler.go`

```go
package handler

import (
    "strconv"
    "github.com/gin-gonic/gin"
    "advanced-user-api/internal/domain"
    "advanced-user-api/internal/service"
)

type PostHandler struct {
    postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
    return &PostHandler{postService: postService}
}

// GetAll - список постов с пагинацией
func (h *PostHandler) GetAll(c *gin.Context) {
    // Query параметры
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    status := c.Query("status")        // ?status=published
    search := c.Query("search")        // ?search=golang
    
    // Фильтры
    filters := repository.PostFilters{
        Status: status,
        Search: search,
    }
    
    // Получаем посты
    response, err := h.postService.GetAllPosts(page, pageSize, filters)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, response)
}

// GetByID - получить пост по ID
func (h *PostHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid ID"})
        return
    }
    
    post, err := h.postService.GetPost(uint(id))
    if err != nil {
        c.JSON(404, gin.H{"error": "пост не найден"})
        return
    }
    
    c.JSON(200, post)
}

// Create - создание поста (требует auth)
func (h *PostHandler) Create(c *gin.Context) {
    // Получаем user_id из контекста (установлен AuthMiddleware)
    userID := c.MustGet("user_id").(uint)
    
    var req domain.CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    post, err := h.postService.CreatePost(userID, &req)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, post)
}

// Update - обновление поста
func (h *PostHandler) Update(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    userID := c.MustGet("user_id").(uint)
    
    var req domain.UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    post, err := h.postService.UpdatePost(uint(id), userID, &req)
    if err != nil {
        c.JSON(403, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, post)
}

// Delete - удаление поста
func (h *PostHandler) Delete(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    user := c.MustGet("user").(*domain.User)
    
    if err := h.postService.DeletePost(uint(id), user); err != nil {
        c.JSON(403, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"message": "пост удалён"})
}
```

---

## 🛣️ Routing (обновлённый)

**Обновляем:** `internal/handler/routes.go`

```go
func SetupRoutes(
    router *gin.Engine,
    authHandler *AuthHandler,
    userHandler *UserHandler,
    postHandler *PostHandler,        // НОВОЕ
    commentHandler *CommentHandler,  // НОВОЕ
    categoryHandler *CategoryHandler, // НОВОЕ
    cfg *config.Config,
) {
    // CORS
    router.Use(middleware.CORSMiddleware())
    
    // API v1
    v1 := router.Group("/api/v1")
    
    // === PUBLIC ROUTES ===
    public := v1.Group("")
    {
        // Auth
        public.POST("/auth/register", authHandler.Register)
        public.POST("/auth/login", authHandler.Login)
        
        // Posts (публичное чтение)
        public.GET("/posts", postHandler.GetAll)                    // Список постов
        public.GET("/posts/:id", postHandler.GetByID)               // Пост по ID
        public.GET("/posts/slug/:slug", postHandler.GetBySlug)      // Пост по slug
        public.GET("/posts/:id/comments", commentHandler.GetByPost) // Комментарии поста
        
        // Categories & Tags
        public.GET("/categories", categoryHandler.GetAll)
        public.GET("/tags", tagHandler.GetAll)
    }
    
    // === PROTECTED ROUTES (требуют auth) ===
    protected := v1.Group("")
    protected.Use(middleware.AuthMiddleware(cfg))
    {
        // Auth
        protected.GET("/auth/me", authHandler.Me)
        
        // Users
        protected.GET("/users", userHandler.GetAll)
        protected.GET("/users/:id", userHandler.GetByID)
        protected.PUT("/users/:id", userHandler.Update)
        
        // Posts (создание и редактирование)
        protected.POST("/posts", postHandler.Create)
        protected.PUT("/posts/:id", postHandler.Update)
        protected.DELETE("/posts/:id", postHandler.Delete)
        protected.POST("/posts/:id/publish", postHandler.Publish)
        
        // Comments
        protected.POST("/posts/:id/comments", commentHandler.Create)
        protected.PUT("/comments/:id", commentHandler.Update)
        protected.DELETE("/comments/:id", commentHandler.Delete)
        
        // Likes
        protected.POST("/posts/:id/like", likeHandler.Like)
        protected.DELETE("/posts/:id/like", likeHandler.Unlike)
    }
    
    // === MODERATOR ROUTES ===
    moderator := v1.Group("/moderator")
    moderator.Use(middleware.AuthMiddleware(cfg))
    moderator.Use(middleware.RequireModerator(userService))
    {
        moderator.POST("/comments/:id/block", commentHandler.Block)
        moderator.POST("/comments/:id/unblock", commentHandler.Unblock)
        moderator.DELETE("/posts/:id", postHandler.Delete)  // Модераторы могут удалять любые посты
    }
    
    // === ADMIN ROUTES ===
    admin := v1.Group("/admin")
    admin.Use(middleware.AuthMiddleware(cfg))
    admin.Use(middleware.RequireAdmin(userService))
    {
        admin.POST("/categories", categoryHandler.Create)
        admin.PUT("/categories/:id", categoryHandler.Update)
        admin.DELETE("/categories/:id", categoryHandler.Delete)
        
        admin.POST("/tags", tagHandler.Create)
        admin.DELETE("/tags/:id", tagHandler.Delete)
        
        admin.DELETE("/users/:id", userHandler.Delete)
        admin.PUT("/users/:id/role", userHandler.UpdateRole)  // Изменение роли пользователя
    }
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"service": "blog-api", "status": "ok"})
    })
}
```

---

## 🔐 Примеры использования с ролями

### Создание поста (любой авторизованный пользователь)

```bash
TOKEN="your-jwt-token"

curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Введение в Go",
    "content": "Go - это современный язык программирования...",
    "excerpt": "Узнайте основы Go",
    "status": "published",
    "category_id": 1,
    "tag_ids": [1, 2, 3]
  }'
```

**Response:**
```json
{
  "id": 1,
  "title": "Введение в Go",
  "slug": "vvedenie-v-go",
  "content": "...",
  "status": "published",
  "author_id": 1,
  "author": {
    "id": 1,
    "name": "Alice",
    "role": "user"
  },
  "view_count": 0,
  "created_at": "2025-10-15T10:00:00Z"
}
```

---

### Блокировка комментария (модератор)

```bash
# Только пользователи с ролью "moderator" или "admin"

curl -X POST http://localhost:8080/api/v1/moderator/comments/123/block \
  -H "Authorization: Bearer $MODERATOR_TOKEN"
```

**Response:**
```json
{
  "message": "комментарий заблокирован"
}
```

Если user:
```json
{
  "error": "insufficient permissions"
}
```

---

### Создание категории (только admin)

```bash
# Только пользователи с ролью "admin"

curl -X POST http://localhost:8080/api/v1/admin/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Технологии",
    "description": "Посты о новых технологиях"
  }'
```

---

## 🧪 Миграции

**Создаём:** `migrations/002_create_blog_tables.sql`

```sql
-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) UNIQUE NOT NULL,
    slug        VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP NULL
);

-- Tags
CREATE TABLE IF NOT EXISTS tags (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(50) UNIQUE NOT NULL,
    slug       VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Posts
CREATE TABLE IF NOT EXISTS posts (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(200) NOT NULL,
    slug         VARCHAR(250) UNIQUE NOT NULL,
    content      TEXT NOT NULL,
    excerpt      VARCHAR(500),
    cover_image  VARCHAR(500),
    status       VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    view_count   INTEGER DEFAULT 0,
    author_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id  INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW(),
    deleted_at   TIMESTAMP NULL,
    published_at TIMESTAMP NULL
);

CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_category_id ON posts(category_id);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_slug ON posts(slug);

-- Comments
CREATE TABLE IF NOT EXISTS comments (
    id         SERIAL PRIMARY KEY,
    content    TEXT NOT NULL,
    post_id    INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id  INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    is_edited  BOOLEAN DEFAULT FALSE,
    is_blocked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);

-- Likes
CREATE TABLE IF NOT EXISTS likes (
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(post_id, user_id)
);

CREATE INDEX idx_likes_post_id ON likes(post_id);

-- Post_Tags (Many-to-Many)
CREATE TABLE IF NOT EXISTS post_tags (
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);
```

**Применение миграции:**

```bash
# Через psql
psql -U postgres -d advanced_api -f migrations/002_create_blog_tables.sql

# Или через GORM Auto Migrate (в коде)
db.AutoMigrate(
    &domain.User{},
    &domain.Post{},
    &domain.Comment{},
    &domain.Category{},
    &domain.Tag{},
    &domain.Like{},
)
```

---

## 📋 Пошаговый план реализации

### Этап 1: Модели (1-2 часа)

1. ✅ Обновить `User` модель (добавить Bio, Avatar)
2. ✅ Создать `Post` модель
3. ✅ Создать `Comment` модель
4. ✅ Создать `Category` модель
5. ✅ Создать `Tag` модель
6. ✅ Создать `Like` модель
7. ✅ Создать DTOs для всех моделей

---

### Этап 2: Repository (2-3 часа)

1. ✅ `PostRepository` - CRUD + фильтры + пагинация
2. ✅ `CommentRepository` - вложенные комментарии
3. ✅ `CategoryRepository` - CRUD
4. ✅ `TagRepository` - CRUD
5. ✅ `LikeRepository` - лайк/анлайк

---

### Этап 3: Service (2-3 часа)

1. ✅ `PostService` - бизнес-логика постов
2. ✅ `CommentService` - валидация комментариев
3. ✅ `CategoryService`
4. ✅ `TagService`
5. ✅ Добавить проверки прав (CanEdit, CanDelete)

---

### Этап 4: Middleware (1 час)

1. ✅ `RequireRole()` - проверка роли
2. ✅ `RequireAdmin()` - только админы
3. ✅ `RequireModerator()` - модераторы и админы

---

### Этап 5: Handlers (2-3 часа)

1. ✅ `PostHandler` - все endpoints постов
2. ✅ `CommentHandler` - endpoints комментариев
3. ✅ `CategoryHandler`
4. ✅ `TagHandler`
5. ✅ `LikeHandler`

---

### Этап 6: Routes (1 час)

1. ✅ Настроить public routes
2. ✅ Настроить protected routes
3. ✅ Настроить moderator routes
4. ✅ Настроить admin routes

---

### Этап 7: Миграции (1 час)

1. ✅ Создать SQL миграции
2. ✅ Применить миграции
3. ✅ Добавить тестовые данные (seed)

---

### Этап 8: Тесты (2-3 часа)

1. ✅ Unit тесты для сервисов
2. ✅ Integration тесты для API
3. ✅ Тесты проверки прав

---

## 🎯 Итоговые endpoints

После реализации у вас будет:

### Public (без auth)
```
GET    /api/v1/posts                    # Список постов
GET    /api/v1/posts/:id                # Пост по ID
GET    /api/v1/posts/slug/:slug         # Пост по slug
GET    /api/v1/posts/:id/comments       # Комментарии
GET    /api/v1/categories               # Категории
GET    /api/v1/tags                     # Теги
```

### Authenticated (auth required)
```
POST   /api/v1/posts                    # Создать пост
PUT    /api/v1/posts/:id                # Обновить свой пост
DELETE /api/v1/posts/:id                # Удалить свой пост
POST   /api/v1/posts/:id/comments       # Добавить комментарий
PUT    /api/v1/comments/:id             # Редактировать свой комментарий
DELETE /api/v1/comments/:id             # Удалить свой комментарий
POST   /api/v1/posts/:id/like           # Лайкнуть пост
DELETE /api/v1/posts/:id/like           # Убрать лайк
```

### Moderator
```
POST   /api/v1/moderator/comments/:id/block    # Заблокировать комментарий
DELETE /api/v1/moderator/posts/:id             # Удалить любой пост
```

### Admin
```
POST   /api/v1/admin/categories         # Создать категорию
PUT    /api/v1/admin/categories/:id     # Обновить категорию
DELETE /api/v1/admin/categories/:id     # Удалить категорию
POST   /api/v1/admin/tags               # Создать тег
DELETE /api/v1/admin/tags/:id           # Удалить тег
PUT    /api/v1/admin/users/:id/role     # Изменить роль пользователя
DELETE /api/v1/admin/users/:id          # Удалить пользователя
```

**Итого:** ~30 endpoints (сейчас 8)

---

## 🔒 Система разрешений

### Матрица прав

| Действие | User | Moderator | Admin |
|----------|------|-----------|-------|
| Создать пост | ✅ | ✅ | ✅ |
| Редактировать свой пост | ✅ | ✅ | ✅ |
| Редактировать чужой пост | ❌ | ❌ | ✅ |
| Удалить свой пост | ✅ | ✅ | ✅ |
| Удалить чужой пост | ❌ | ✅ | ✅ |
| Создать комментарий | ✅ | ✅ | ✅ |
| Удалить свой комментарий | ✅ | ✅ | ✅ |
| Удалить чужой комментарий | ❌ | ✅ | ✅ |
| Блокировать комментарий | ❌ | ✅ | ✅ |
| Создать категорию | ❌ | ❌ | ✅ |
| Изменить роль пользователя | ❌ | ❌ | ✅ |
| Удалить пользователя | ❌ | ❌ | ✅ |

---

## 📦 Дополнительные библиотеки

### 1. Slug generation

```bash
go get github.com/gosimple/slug
```

**Использование:**
```go
import "github.com/gosimple/slug"

slug := slug.Make("Введение в Go")  // "vvedenie-v-go"
```

### 2. Pagination helper

```bash
go get github.com/biezhi/gorm-paginator/pagination
```

**Использование:**
```go
var posts []Post
var paginator pagination.Paginator

db.Scopes(pagination.Paginate(&posts, &pagination.Param{
    DB:      db,
    Page:    page,
    Limit:   pageSize,
})).Find(&posts)
```

### 3. Sanitize HTML (для безопасности)

```bash
go get github.com/microcosm-cc/bluemonday
```

**Использование:**
```go
import "github.com/microcosm-cc/bluemonday"

policy := bluemonday.UGCPolicy()
safeContent := policy.Sanitize(req.Content)
```

---

## 🧪 Тестирование ролей

**Создаём:** `tests/unit/role_test.go`

```go
func TestPostDelete_AsAuthor(t *testing.T) {
    author := &domain.User{ID: 1, Role: domain.RoleUser}
    post := &domain.Post{ID: 1, AuthorID: 1}
    
    assert.True(t, post.CanDelete(author))
}

func TestPostDelete_AsModerator(t *testing.T) {
    moderator := &domain.User{ID: 2, Role: domain.RoleModerator}
    post := &domain.Post{ID: 1, AuthorID: 1}  // Чужой пост
    
    assert.True(t, post.CanDelete(moderator))  // Модератор может
}

func TestPostDelete_AsUser(t *testing.T) {
    user := &domain.User{ID: 2, Role: domain.RoleUser}
    post := &domain.Post{ID: 1, AuthorID: 1}  // Чужой пост
    
    assert.False(t, post.CanDelete(user))  // Обычный user не может
}
```

---

## 🚀 Docker Compose обновление

Добавляем Redis для кеширования:

```yaml
# docker-compose.yml
services:
  redis:
    image: redis:7-alpine
    container_name: blog-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - app-network

volumes:
  redis_data:
```

**Использование Redis:**
- Кеш популярных постов
- Счётчики просмотров
- Session storage

---

## 📊 Пример полного flow

### Сценарий: Пользователь создаёт и публикует пост

```go
// 1. Регистрация
POST /api/v1/auth/register
{
  "email": "blogger@example.com",
  "name": "Tech Blogger",
  "password": "password123"
}

// Получаем токен
Response: {"token": "...", "user": {"id": 1, "role": "user"}}

// 2. Создание поста (draft)
POST /api/v1/posts
Authorization: Bearer <token>
{
  "title": "Мой первый пост",
  "content": "Содержание поста...",
  "status": "draft"
}

Response: {"id": 1, "status": "draft", ...}

// 3. Публикация поста
POST /api/v1/posts/1/publish
Authorization: Bearer <token>

Response: {"id": 1, "status": "published", "published_at": "..."}

// 4. Другой пользователь комментирует
POST /api/v1/posts/1/comments
Authorization: Bearer <other-token>
{
  "content": "Отличный пост!"
}

// 5. Автор отвечает на комментарий
POST /api/v1/posts/1/comments
Authorization: Bearer <token>
{
  "content": "Спасибо!",
  "parent_id": 1  // Ответ на комментарий #1
}

// 6. Модератор блокирует спам-комментарий
POST /api/v1/moderator/comments/2/block
Authorization: Bearer <moderator-token>
```

---

## 💡 Дополнительные функции

### 1. Полнотекстовый поиск (PostgreSQL)

```go
// Repository
func (r *postRepository) Search(query string) ([]domain.Post, error) {
    var posts []domain.Post
    err := r.db.
        Where("to_tsvector('russian', title || ' ' || content) @@ plainto_tsquery('russian', ?)", query).
        Find(&posts).Error
    return posts, err
}
```

**Миграция:**
```sql
-- Создание полнотекстового индекса
CREATE INDEX idx_posts_search ON posts 
USING GIN(to_tsvector('russian', title || ' ' || content));
```

### 2. View Counter с Redis

```go
// Service
func (s *postService) IncrementView(postID uint) error {
    // Увеличиваем в Redis
    key := fmt.Sprintf("post:%d:views", postID)
    s.redisClient.Incr(ctx, key)
    
    // Раз в минуту синхронизируем с БД
    // (через background worker или cron)
    return nil
}
```

### 3. Popular Posts (кеш топ-10)

```go
// Service с Redis
func (s *postService) GetPopularPosts() ([]domain.Post, error) {
    // Проверяем кеш
    cached, err := s.redisClient.Get(ctx, "popular_posts").Result()
    if err == nil {
        json.Unmarshal([]byte(cached), &posts)
        return posts, nil
    }
    
    // Если нет в кеше, берём из БД
    posts, err := s.postRepo.FindPopular(10)
    
    // Кешируем на 1 час
    json, _ := json.Marshal(posts)
    s.redisClient.Set(ctx, "popular_posts", json, time.Hour)
    
    return posts, nil
}
```

### 4. Markdown поддержка

```bash
go get github.com/gomarkdown/markdown
```

```go
import "github.com/gomarkdown/markdown"

md := []byte("# Hello\n\nThis is **bold**")
html := markdown.ToHTML(md, nil, nil)
```

### 5. Upload изображений

```go
// Handler для загрузки обложки
func (h *PostHandler) UploadCover(c *gin.Context) {
    file, _ := c.FormFile("image")
    
    // Сохраняем в /uploads
    filename := generateFilename(file.Filename)
    c.SaveUploadedFile(file, "uploads/"+filename)
    
    c.JSON(200, gin.H{
        "url": "/uploads/" + filename,
    })
}
```

---

## 📚 Полезные библиотеки для блога

```bash
# Slug generation
go get github.com/gosimple/slug

# Markdown to HTML
go get github.com/gomarkdown/markdown

# HTML sanitization
go get github.com/microcosm-cc/bluemonday

# Image processing
go get github.com/disintegration/imaging

# Pagination
go get github.com/biezhi/gorm-paginator/pagination

# Redis client
go get github.com/redis/go-redis/v9

# Full-text search (PostgreSQL)
# Встроено в PostgreSQL

# Sitemap generation
go get github.com/ikeikeikeike/go-sitemap-generator/v2/stm
```

---

## 🎓 Что вы получите

### Функциональная блог-платформа:

✅ **Посты** - создание, редактирование, публикация  
✅ **Комментарии** - с вложенностью (ответы на комментарии)  
✅ **Роли** - user, moderator, admin  
✅ **Права** - разграничение доступа  
✅ **Категории** - организация постов  
✅ **Теги** - маркировка контента  
✅ **Лайки** - социальная функция  
✅ **Поиск** - полнотекстовый поиск  
✅ **Пагинация** - для больших списков  
✅ **Markdown** - форматирование постов  
✅ **Кеширование** - Redis для производительности  

### Продвинутые возможности:

- Image upload для обложек
- Slug generation для SEO
- View counter
- Popular posts
- Draft/Published статусы
- Soft delete
- Модерация контента

---

## 📖 См. также

- [Architecture Guide](./ARCHITECTURE.md) - Как расширить архитектуру
- [GORM Documentation](./libraries/GORM.md) - Связи между моделями
- [Testing Guide](./TESTING.md) - Тестирование новых функций
- [API Documentation](./API.md) - Добавление новых endpoints

---

**Время реализации:** 10-15 часов для полной функциональности

**Сложность:** Средняя (если понимаете текущий код)

**Результат:** Production-ready blog platform с разделением ролей и полным функционалом!

