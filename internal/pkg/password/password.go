package password

import (
	"golang.org/x/crypto/bcrypt" // Bcrypt для хеширования паролей
)

// ================================================================
// PASSWORD UTILITIES - Хеширование и проверка паролей
// ================================================================

// Hash хеширует пароль с помощью bcrypt
// Параметры:
//   - password: пароль в открытом виде (например, "secret123")
// Возвращает:
//   - string: хеш пароля (например, "$2a$10$N9qo8uLOickgx2ZMRZoMy...")
//   - error: ошибка хеширования
//
// Bcrypt - это алгоритм хеширования паролей:
// - Односторонний (нельзя восстановить пароль из хеша)
// - Медленный (защита от brute-force атак)
// - С солью (каждый хеш уникален, даже для одинаковых паролей)
func Hash(password string) (string, error) {
	// bcrypt.GenerateFromPassword() - создаёт хеш пароля
	// Параметры:
	//   - []byte(password): пароль в байтах
	//   - bcrypt.DefaultCost: сложность хеширования (10)
	//     Чем выше cost, тем медленнее хеширование (и безопаснее)
	//     DefaultCost = 10 - хороший баланс между скоростью и безопасностью
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Конвертируем []byte в string и возвращаем
	return string(hashedBytes), nil
}

// Verify проверяет, соответствует ли пароль хешу
// Параметры:
//   - hashedPassword: хеш из БД (например, "$2a$10$N9qo8uLO...")
//   - password: пароль от пользователя для проверки (например, "secret123")
// Возвращает:
//   - bool: true если пароль правильный, false если неправильный
//
// Использование:
//   if password.Verify(user.Password, inputPassword) {
//       // Пароль правильный - пользователь аутентифицирован
//   } else {
//       // Пароль неправильный - отклоняем вход
//   }
func Verify(hashedPassword, password string) bool {
	// bcrypt.CompareHashAndPassword() - сравнивает хеш и пароль
	// Параметры:
	//   - []byte(hashedPassword): хеш из БД
	//   - []byte(password): пароль для проверки
	// Возвращает:
	//   - nil: если пароль правильный
	//   - error: если пароль неправильный
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	
	// Если err == nil, значит пароль правильный
	return err == nil
}

// ================================================================
// ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ
// ================================================================

// Пример 1: Регистрация пользователя
// func Register(email, password string) {
//     // Хешируем пароль перед сохранением в БД
//     hashedPassword, _ := password.Hash(password)
//     
//     user := &User{
//         Email: email,
//         Password: hashedPassword, // Сохраняем ХЕШ, не сам пароль!
//     }
//     
//     db.Create(user)
// }

// Пример 2: Вход пользователя
// func Login(email, password string) bool {
//     // Находим пользователя по email
//     user := findUserByEmail(email)
//     
//     // Проверяем пароль
//     if password.Verify(user.Password, password) {
//         // Пароль правильный - генерируем JWT токен
//         return true
//     }
//     
//     // Пароль неправильный
//     return false
// }

