# 🌿 Git Workflow Guide

## Branching Strategy

Профессиональный workflow для командной разработки с тремя основными ветками.

### Branch Structure

```
main (production)
  ↑
  └── staging (testing)
        ↑
        └── develop (development)
              ↑
              └── feature/* (новые функции)
              └── bugfix/* (исправления багов)
              └── hotfix/* (срочные исправления)
```

---

## 🌲 Main Branches

### `main` - Production
- **Назначение**: Стабильный production код
- **Защита**: Protected branch, требует review
- **Деплой**: Автоматический в production
- **Правило**: Только merge из `staging` после тестирования

```bash
# Никогда не коммитим напрямую в main!
# Только через Pull Request из staging
```

### `staging` - Testing
- **Назначение**: Тестирование перед production
- **Деплой**: Автоматический на staging сервер
- **Правило**: Merge из `develop` после code review

### `develop` - Development
- **Назначение**: Основная ветка разработки
- **Деплой**: Автоматический на dev сервер
- **Правило**: Merge feature веток после code review

---

## 🔧 Feature Development Workflow

### 1. Создание feature ветки

```bash
# Переключаемся на develop
git checkout develop

# Обновляем develop
git pull origin develop

# Создаём новую feature ветку
git checkout -b feature/user-authentication

# Или более короткая форма
git checkout -b feature/add-pagination
```

**Именование веток:**
- `feature/` - новая функциональность
- `bugfix/` - исправление бага
- `hotfix/` - срочное исправление для production

**Примеры:**
```
feature/jwt-authentication
feature/add-pagination
bugfix/fix-login-error
hotfix/critical-security-patch
```

### 2. Работа над функцией

```bash
# Делаем изменения в коде

# Смотрим статус
git status

# Добавляем изменённые файлы
git add .
# или конкретные файлы
git add internal/handler/auth_handler.go

# Коммитим с описательным сообщением
git commit -m "feat: add JWT authentication middleware"

# Отправляем в remote
git push origin feature/user-authentication
```

### 3. Синхронизация с develop

```bash
# Регулярно синхронизируемся с develop
git checkout develop
git pull origin develop

# Возвращаемся в feature ветку
git checkout feature/user-authentication

# Вливаем изменения из develop
git merge develop

# Или используем rebase (чище история)
git rebase develop

# Разрешаем конфликты если есть
git status
# редактируем конфликтные файлы
git add .
git rebase --continue

# Отправляем обновления
git push origin feature/user-authentication --force-with-lease
```

### 4. Создание Pull Request

```bash
# После завершения работы, пушим финальную версию
git push origin feature/user-authentication

# Переходим на GitHub и создаём Pull Request:
# feature/user-authentication → develop
```

**В описании PR указываем:**
- Что реализовано
- Какие issue закрываются (`Closes #123`)
- Скриншоты (если UI)
- Как тестировать

### 5. Code Review и Merge

```bash
# После одобрения PR, merge через GitHub UI
# или локально:

git checkout develop
git merge --no-ff feature/user-authentication
git push origin develop

# Удаляем feature ветку
git branch -d feature/user-authentication
git push origin --delete feature/user-authentication
```

---

## 🐛 Bugfix Workflow

```bash
# Создаём bugfix ветку из develop
git checkout develop
git pull origin develop
git checkout -b bugfix/fix-login-validation

# Исправляем баг
# ... код ...

git add .
git commit -m "fix: correct email validation in login"
git push origin bugfix/fix-login-validation

# Pull Request: bugfix/fix-login-validation → develop
```

---

## 🚨 Hotfix Workflow (срочное исправление в production)

```bash
# Hotfix создаём из main!
git checkout main
git pull origin main
git checkout -b hotfix/critical-security-patch

# Исправляем критический баг
# ... код ...

git add .
git commit -m "hotfix: fix SQL injection vulnerability"
git push origin hotfix/critical-security-patch

# Pull Request: hotfix/critical-security-patch → main
# После merge в main, также мержим в develop и staging
```

**Важно:** Hotfix нужно влить обратно в `develop` и `staging`:

```bash
git checkout staging
git merge main
git push origin staging

git checkout develop
git merge main
git push origin develop
```

---

## 🔄 Release Workflow

### Staging → Main

```bash
# 1. Тестируем на staging
# 2. Если всё ок, создаём PR: staging → main
# 3. После merge, создаём release tag

git checkout main
git pull origin main

# Создаём тег с версией
git tag -a v1.2.0 -m "Release v1.2.0 - Add JWT authentication"
git push origin v1.2.0

# Или создаём release через GitHub UI
```

---

## 📝 Commit Message Convention

### Формат

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

```
feat:     Новая функциональность
fix:      Исправление бага
docs:     Изменения в документации
style:    Форматирование кода (без изменения логики)
refactor: Рефакторинг кода
test:     Добавление/изменение тестов
chore:    Рутинные задачи (обновление зависимостей)
perf:     Улучшение производительности
ci:       CI/CD изменения
```

### Примеры

```bash
# Хорошо ✅
git commit -m "feat: add JWT authentication middleware"
git commit -m "fix: resolve race condition in user repository"
git commit -m "docs: update API documentation"
git commit -m "test: add unit tests for auth service"

# Плохо ❌
git commit -m "fix"
git commit -m "update"
git commit -m "WIP"
git commit -m "asdfasdf"
```

### Подробный commit

```bash
git commit -m "feat: add user pagination

- Add limit and offset query parameters
- Implement pagination in repository layer
- Update API documentation
- Add integration tests

Closes #45"
```

---

## 🔍 Useful Git Commands

### Статус и история

```bash
# Текущий статус
git status

# Короткий статус
git status -s

# История коммитов
git log

# Краткая история
git log --oneline

# История с графом веток
git log --oneline --graph --all

# Последние 10 коммитов
git log -10

# Изменения в файле
git log -p internal/handler/auth_handler.go

# Кто изменял файл
git blame internal/handler/auth_handler.go
```

### Работа с изменениями

```bash
# Посмотреть изменения
git diff

# Изменения в конкретном файле
git diff internal/handler/auth_handler.go

# Изменения staged файлов
git diff --staged

# Отменить изменения в файле
git checkout -- internal/handler/auth_handler.go

# Отменить все изменения
git checkout -- .

# Убрать файл из staging
git reset HEAD internal/handler/auth_handler.go

# Добавить изменения к последнему коммиту
git add .
git commit --amend

# Изменить сообщение последнего коммита
git commit --amend -m "новое сообщение"
```

### Работа с ветками

```bash
# Список локальных веток
git branch

# Список всех веток (включая remote)
git branch -a

# Создать новую ветку
git branch feature/new-feature

# Переключиться на ветку
git checkout feature/new-feature

# Создать и переключиться (короче)
git checkout -b feature/new-feature

# Удалить локальную ветку
git branch -d feature/old-feature

# Принудительно удалить
git branch -D feature/old-feature

# Удалить remote ветку
git push origin --delete feature/old-feature

# Переименовать текущую ветку
git branch -m new-name
```

### Синхронизация с remote

```bash
# Получить изменения (без merge)
git fetch origin

# Получить и влить изменения
git pull origin develop

# Отправить изменения
git push origin feature/my-feature

# Отправить все ветки
git push --all origin

# Отправить теги
git push --tags

# Форсированный push (осторожно!)
git push origin feature/my-feature --force-with-lease
```

### Откат изменений

```bash
# Откатить последний коммит (сохранить изменения)
git reset --soft HEAD~1

# Откатить последний коммит (удалить изменения)
git reset --hard HEAD~1

# Откатить к конкретному коммиту
git reset --hard abc1234

# Создать новый коммит, отменяющий предыдущий
git revert abc1234

# Отменить merge
git revert -m 1 abc1234
```

### Stash (временное сохранение)

```bash
# Сохранить текущие изменения
git stash

# Сохранить с сообщением
git stash save "WIP: работаю над аутентификацией"

# Список stash
git stash list

# Применить последний stash
git stash apply

# Применить и удалить stash
git stash pop

# Применить конкретный stash
git stash apply stash@{2}

# Удалить stash
git stash drop stash@{0}

# Очистить все stash
git stash clear
```

### Работа с конфликтами

```bash
# После git merge или git pull с конфликтами

# Посмотреть конфликтные файлы
git status

# Открыть файлы и разрешить конфликты
# Ищем маркеры:
# <<<<<<< HEAD
# ваш код
# =======
# код из другой ветки
# >>>>>>> feature/other-branch

# После разрешения конфликтов
git add .
git commit -m "merge: resolve conflicts"

# Отменить merge
git merge --abort

# Отменить rebase
git rebase --abort
```

### Просмотр remote

```bash
# Список remote репозиториев
git remote -v

# Добавить remote
git remote add origin git@github.com:user/repo.git

# Изменить URL remote
git remote set-url origin git@github.com:user/new-repo.git

# Удалить remote
git remote remove origin
```

### Cherry-pick (взять конкретный коммит)

```bash
# Взять коммит из другой ветки
git cherry-pick abc1234

# Взять несколько коммитов
git cherry-pick abc1234 def5678

# Взять коммит без автокоммита
git cherry-pick -n abc1234
```

### Rebase vs Merge

```bash
# Merge (создаёт merge commit)
git checkout develop
git merge feature/my-feature

# Rebase (переписывает историю, чище)
git checkout feature/my-feature
git rebase develop

# Interactive rebase (изменить последние 3 коммита)
git rebase -i HEAD~3

# Продолжить rebase после разрешения конфликтов
git rebase --continue
```

---

## 🛡️ Best Practices

### 1. Коммитьте часто, но логично

```bash
# ✅ Хорошо
git commit -m "feat: add user validation"
git commit -m "feat: add password hashing"
git commit -m "test: add auth tests"

# ❌ Плохо
git commit -m "massive changes to everything"
```

### 2. Пишите описательные сообщения

```bash
# ✅ Хорошо
git commit -m "fix: resolve race condition in user repository when updating concurrent sessions"

# ❌ Плохо
git commit -m "fix bug"
```

### 3. Держите ветки актуальными

```bash
# Регулярно синхронизируйтесь с develop
git checkout develop
git pull origin develop
git checkout feature/my-feature
git merge develop
```

### 4. Не коммитьте в main/staging напрямую

```bash
# ❌ Плохо
git checkout main
git commit -m "fix"

# ✅ Хорошо
git checkout -b hotfix/critical-fix
# ... делаем исправления ...
git commit -m "hotfix: fix critical bug"
# Создаём PR
```

### 5. Используйте .gitignore

```bash
# Никогда не коммитьте:
.env              # Секреты
*.log             # Логи
node_modules/     # Зависимости
bin/              # Бинарники
coverage.html     # Временные файлы
```

### 6. Проверяйте перед push

```bash
# Убедитесь что код работает
go test ./...
go build ./...

# Проверьте что коммитите
git diff --staged

# Затем push
git push origin feature/my-feature
```

---

## 🔐 SSH Setup для GitHub

```bash
# Генерация SSH ключа
ssh-keygen -t ed25519 -C "your-email@example.com"

# Запуск ssh-agent
eval "$(ssh-agent -s)"

# Добавление ключа
ssh-add ~/.ssh/id_ed25519

# Скопировать публичный ключ
cat ~/.ssh/id_ed25519.pub

# Добавить на GitHub: Settings → SSH Keys → New SSH Key

# Тестирование
ssh -T git@github.com
```

---

## 📖 Additional Resources

- [Git Documentation](https://git-scm.com/doc)
- [GitHub Flow](https://docs.github.com/en/get-started/quickstart/github-flow)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Git Best Practices](https://sethrobertson.github.io/GitBestPractices/)

---

*Этот workflow адаптирован для командной разработки с CI/CD и code review процессом.*

