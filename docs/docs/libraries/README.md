# 📚 Libraries Documentation

Подробная документация по всем библиотекам, используемым в проекте.

---

## 🍸 [Gin Web Framework](./GIN.md)

**Версия:** v1.11.0

**Описание:** Высокопроизводительный HTTP веб-фреймворк

**Что документировано:**
- 16 основных методов (gin.Default, router.Use, c.JSON, c.ShouldBindJSON, и т.д.)
- Validation tags
- Middleware pattern
- Error handling
- Best practices

**Где используем:**
- HTTP handlers
- Middleware
- Роутинг

---

## 🗄️ [GORM ORM](./GORM.md)

**Версия:** v1.25.5

**Описание:** Полнофункциональная ORM библиотека для Go

**Что документировано:**
- CRUD операции (Create, First, Find, Save, Updates, Delete)
- Auto Migration
- Struct tags
- Connection pool
- Soft delete
- Error handling
- Best practices

**Где используем:**
- Repository layer
- Database operations
- Модели данных

---

## 🔐 [JWT - JSON Web Tokens](./JWT.md)

**Версия:** v5.3.0

**Описание:** Аутентификация через JWT токены

**Что документировано:**
- Структура JWT
- Claims (данные в токене)
- Генерация токенов (NewWithClaims, SignedString)
- Валидация токенов (ParseWithClaims)
- Security best practices
- Полный auth flow

**Где используем:**
- Auth service
- Auth middleware
- Защита endpoints

---

## 🔒 [Bcrypt & Viper](./BCRYPT_VIPER.md)

### Bcrypt - Password Hashing

**Версия:** v0.18.0

**Описание:** Безопасное хеширование паролей

**Что документировано:**
- GenerateFromPassword
- CompareHashAndPassword
- Cost values
- Security best practices

**Где используем:**
- Хеширование паролей при регистрации
- Проверка паролей при входе

### Viper - Configuration

**Версия:** v1.21.0

**Описание:** Управление конфигурацией

**Что документировано:**
- Чтение .env файлов
- Environment variables
- Defaults
- Приоритет источников

**Где используем:**
- Загрузка конфигурации при старте
- Настройки БД, JWT, сервера

---

## 📋 Быстрая навигация

### По задачам:

**HTTP обработка:**
- [Gin - Handlers](./GIN.md#handlers)
- [Gin - Middleware](./GIN.md#middleware-pattern)
- [Gin - Validation](./GIN.md#validation-tags)

**База данных:**
- [GORM - CRUD](./GORM.md#crud-operations)
- [GORM - Migration](./GORM.md#auto-migration)
- [GORM - Модели](./GORM.md#struct-tags)

**Безопасность:**
- [JWT - Генерация](./JWT.md#generate---генерация-токена)
- [JWT - Валидация](./JWT.md#validate---валидация-токена)
- [Bcrypt - Хеширование](./BCRYPT_VIPER.md#hash---хеширование-пароля)
- [Bcrypt - Проверка](./BCRYPT_VIPER.md#verify---проверка-пароля)

**Конфигурация:**
- [Viper - Load config](./BCRYPT_VIPER.md#load---загрузка-конфигурации)

---

## 🔗 Связь с кодом

Каждая документация содержит:
- ✅ Сигнатуры методов
- ✅ Описание параметров
- ✅ Примеры использования
- ✅ Ссылки на файлы проекта где используется
- ✅ Best practices

---

## 📖 См. также

- [Architecture Guide](../ARCHITECTURE.md) - Где используются библиотеки в архитектуре
- [API Documentation](../API.md) - Как библиотеки работают в API
- [Testing Guide](../TESTING.md) - Тестирование с этими библиотеками

