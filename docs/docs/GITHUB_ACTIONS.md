# ⚙️ GitHub Actions CI/CD Guide

## Обзор

GitHub Actions - встроенная CI/CD платформа GitHub для автоматизации тестирования, сборки и деплоя.

**Файл в проекте:** `.github/workflows/ci.yml`

**Что автоматизируем:**
- ✅ Тестирование при каждом push
- ✅ Сборка Docker образа
- ✅ Проверка кода (linting)
- ✅ Security scanning
- ✅ Code coverage reporting

---

## 📁 Структура workflow файла

```
.github/
└── workflows/
    └── ci.yml          # CI/CD pipeline
```

**Можно создать несколько workflows:**
```
.github/workflows/
├── ci.yml              # CI (тесты, линтеры)
├── deploy.yml          # Деплой
├── release.yml         # Создание релизов
└── security.yml        # Security scans
```

---

## 📄 Разбор ci.yml

### Структура файла

```yaml
name: CI/CD Pipeline      # Название workflow

on:                       # Когда запускать
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:                     # Задачи для выполнения
  test-and-build:         # Job 1
    ...
  docker-build:           # Job 2
    ...
  security-scan:          # Job 3
    ...
```

---

## 🎯 Triggers (когда запускать)

### on: push

```yaml
on:
  push:
    branches: [ main, develop ]
```

**Когда срабатывает:**
- При `git push` в ветки `main` или `develop`
- При merge Pull Request в эти ветки

**Можно также:**
```yaml
on:
  push:
    branches:
      - main
      - 'feature/**'     # Все ветки feature/*
    paths:
      - 'internal/**'    # Только при изменении internal/
      - 'cmd/**'
    tags:
      - 'v*'             # При создании тега v1.0.0
```

---

### on: pull_request

```yaml
on:
  pull_request:
    branches: [ main ]
```

**Когда срабатывает:**
- При создании PR в ветку `main`
- При обновлении PR (новые коммиты)

**Проверяет код ДО merge!**

---

### on: schedule

```yaml
on:
  schedule:
    - cron: '0 2 * * *'  # Каждый день в 2:00
```

**Полезно для:**
- Ночные тесты
- Проверка зависимостей
- Бэкапы

---

### on: workflow_dispatch

```yaml
on:
  workflow_dispatch:  # Ручной запуск
```

**Позволяет запустить вручную** через GitHub UI

---

## 🏗️ Jobs - Задачи

### Job 1: Test and Build

```yaml
test-and-build:
  name: Test & Build           # Отображаемое имя
  runs-on: ubuntu-latest       # ОС для выполнения
  
  services:                    # Дополнительные сервисы (PostgreSQL для тестов)
    postgres:
      image: postgres:15-alpine
      env:
        POSTGRES_USER: postgres
        POSTGRES_PASSWORD: postgres
        POSTGRES_DB: advanced_api_test
      ports:
        - 5432:5432
      options: >-
        --health-cmd pg_isready
        --health-interval 10s
        --health-timeout 5s
        --health-retries 5
  
  steps:                       # Шаги выполнения
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.25'
    
    # ... остальные шаги
```

---

## 📋 Steps - Шаги выполнения

### 1. Checkout Code

```yaml
- name: Checkout code
  uses: actions/checkout@v4
```

**Что делает:**
- Клонирует ваш репозиторий в runner
- `actions/checkout` - официальный action от GitHub
- `@v4` - версия action

**Параметры (опционально):**
```yaml
- uses: actions/checkout@v4
  with:
    fetch-depth: 0      # Полная история (для changelog)
    submodules: true    # Клонировать submodules
```

---

### 2. Setup Go

```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.25'
```

**Что делает:**
- Устанавливает Go указанной версии
- Настраивает PATH
- Кеширует Go модули

**Можно указать:**
```yaml
go-version: '1.25'      # Точная версия
go-version: '1.25.x'    # Любой patch
go-version: '^1.25'     # >= 1.25
```

---

### 3. Cache Go Modules

```yaml
- name: Cache Go modules
  uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

**Что делает:**
- Кеширует `~/go/pkg/mod` (скачанные модули)
- `key` - уникальный ключ кеша (если `go.sum` изменился, кеш сбрасывается)
- `restore-keys` - fallback ключи

**Результат:**
- Первый запуск: ~2 минуты (скачивает модули)
- Последующие: ~10 секунд (использует кеш) ✅

---

### 4. Download Dependencies

```yaml
- name: Download dependencies
  run: go mod download
```

**Что делает:**
- Скачивает все зависимости из `go.mod`
- Если кеш есть, берёт из кеша

---

### 5. Run Linters

```yaml
- name: Run linters
  run: |
    go fmt ./...
    go vet ./...
```

**Что проверяет:**
- `go fmt` - форматирование кода (должно совпадать с `gofmt`)
- `go vet` - статический анализ (поиск ошибок)

**Можно добавить:**
```yaml
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v3
  with:
    version: latest
```

---

### 6. Run Unit Tests

```yaml
- name: Run unit tests
  run: go test -v -race -coverprofile=coverage.txt -covermode=atomic ./tests/unit/...
```

**Флаги:**
- `-v` - verbose (подробный вывод)
- `-race` - проверка race conditions
- `-coverprofile=coverage.txt` - файл с coverage
- `-covermode=atomic` - режим coverage (для `-race`)
- `./tests/unit/...` - только unit тесты

**Что проверяет:**
- ✅ Все тесты проходят
- ✅ Нет race conditions
- ✅ Собирает coverage для отчёта

---

### 7. Run Integration Tests

```yaml
- name: Run integration tests
  env:
    DB_HOST: localhost
    DB_PORT: 5432
    DB_USER: postgres
    DB_PASSWORD: postgres
    DB_NAME: advanced_api_test
    JWT_SECRET: test-secret-for-ci
  run: go test -v ./tests/integration/...
```

**env: - переменные окружения для тестов**

**PostgreSQL доступен:**
- `localhost:5432` (из services секции)
- База `advanced_api_test` уже создана

---

### 8. Upload Coverage

```yaml
- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v3
  with:
    files: ./coverage.txt
```

**Что делает:**
- Отправляет `coverage.txt` на Codecov.io
- Создаёт красивый отчёт о покрытии кода
- Добавляет badge в README: ![codecov](https://codecov.io/gh/user/repo/branch/main/graph/badge.svg)

**Настройка:**
1. Зарегистрируйтесь на https://codecov.io/
2. Подключите репозиторий
3. Badge автоматически появится

---

### 9. Build Binary

```yaml
- name: Build binary
  run: go build -o ./bin/api cmd/api/main.go
```

**Проверяет что код компилируется**

---

## 🐳 Job 2: Docker Build

```yaml
docker-build:
  name: Docker Build
  runs-on: ubuntu-latest
  needs: test-and-build       # ⚠️ Запускается ПОСЛЕ test-and-build
  
  steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3
    
    - name: Build Docker image
      uses: docker/build-push-action@v5
      with:
        context: .
        file: ./docker/Dockerfile
        push: false              # Не пушим в registry (пока)
        tags: advanced-user-api:latest
        cache-from: type=gha     # Кеш из GitHub Actions
        cache-to: type=gha,mode=max
```

### needs: Зависимости между jobs

```yaml
needs: test-and-build
```

**Порядок выполнения:**
1. `test-and-build` (тесты)
2. Если успешно → `docker-build` (сборка Docker)
3. Если неуспешно → `docker-build` **НЕ запускается**

**Можно несколько зависимостей:**
```yaml
needs: [test, lint, security-scan]
```

---

### Docker Buildx

```yaml
- name: Set up Docker Buildx
  uses: docker/setup-buildx-action@v3
```

**Что это:**
- Расширенный builder для Docker
- Поддержка multi-platform builds
- Продвинутое кеширование

---

### Build and Push Action

```yaml
- name: Build Docker image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./docker/Dockerfile
    push: false
    tags: advanced-user-api:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**Параметры:**
- `context: .` - где искать файлы (корень репозитория)
- `file` - путь к Dockerfile
- `push: false` - не пушить в registry (для PR)
- `tags` - теги образа
- `cache-from/to` - кеш для ускорения сборки

**Для production (push в registry):**
```yaml
- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ secrets.DOCKER_USERNAME }}
    password: ${{ secrets.DOCKER_PASSWORD }}

- name: Build and push
  uses: docker/build-push-action@v5
  with:
    push: true
    tags: username/advanced-user-api:${{ github.sha }}
```

---

## 🔒 Job 3: Security Scan

```yaml
security-scan:
  name: Security Scan
  runs-on: ubuntu-latest
  
  steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        scan-type: 'fs'                    # Filesystem scan
        scan-ref: '.'
        format: 'sarif'
        output: 'trivy-results.sarif'
    
    - name: Upload Trivy results to GitHub Security
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: 'trivy-results.sarif'
```

### Trivy Scanner

**Что проверяет:**
- Уязвимости в зависимостях (`go.mod`)
- Уязвимости в системных пакетах
- Секреты в коде (API keys, пароли)
- Misconfiguration в конфигах

**Результат:**
- Отчёт в GitHub Security tab
- Можно блокировать PR при критических уязвимостях

---

## 🔐 GitHub Secrets

### Что это?

Безопасное хранение секретов (пароли, API keys, токены)

### Добавление секретов:

1. GitHub → Repository → **Settings**
2. **Secrets and variables** → **Actions**
3. **New repository secret**

**Примеры:**
- `DOCKER_USERNAME` - логин Docker Hub
- `DOCKER_PASSWORD` - пароль Docker Hub
- `DEPLOY_KEY` - SSH ключ для деплоя
- `CODECOV_TOKEN` - токен Codecov

### Использование в workflow:

```yaml
steps:
  - name: Login to Docker Hub
    uses: docker/login-action@v3
    with:
      username: ${{ secrets.DOCKER_USERNAME }}
      password: ${{ secrets.DOCKER_PASSWORD }}
  
  - name: Deploy
    env:
      DEPLOY_KEY: ${{ secrets.DEPLOY_KEY }}
    run: |
      echo "$DEPLOY_KEY" > deploy_key
      ssh -i deploy_key user@server 'cd app && git pull'
```

**Безопасность:**
- Секреты **не видны** в логах
- `***` вместо значений
- Доступны только в вашем репозитории

---

## 🔄 Context Variables

### GitHub предоставляет переменные:

```yaml
- name: Print info
  run: |
    echo "Repository: ${{ github.repository }}"
    echo "Branch: ${{ github.ref }}"
    echo "Commit SHA: ${{ github.sha }}"
    echo "Actor: ${{ github.actor }}"
    echo "Run number: ${{ github.run_number }}"
```

**Доступные:**
- `github.repository` - `AlexRodving/advanced-user-api`
- `github.ref` - `refs/heads/main`
- `github.sha` - commit hash
- `github.actor` - кто запустил workflow
- `github.event_name` - `push`, `pull_request`, etc.

**Использование для тегов:**
```yaml
tags: |
  username/app:latest
  username/app:${{ github.sha }}
  username/app:${{ github.ref_name }}
```

---

## 🎨 Матрица стратегия (Matrix Strategy)

### Тестирование на разных версиях

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.23', '1.24', '1.25']
        postgres-version: ['14', '15', '16']
    
    steps:
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      
      - name: Start PostgreSQL
        run: |
          docker run -d \
            -e POSTGRES_VERSION=${{ matrix.postgres-version }} \
            postgres:${{ matrix.postgres-version }}-alpine
```

**Результат:**
- 3 версии Go × 3 версии PostgreSQL = **9 тестовых запусков**
- Параллельно!

---

## 🔧 Services - Дополнительные сервисы

### PostgreSQL для тестов

```yaml
services:
  postgres:
    image: postgres:15-alpine
    env:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: advanced_api_test
    ports:
      - 5432:5432
    options: >-
      --health-cmd pg_isready
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5
```

**Как работает:**
1. GitHub Actions запускает PostgreSQL контейнер
2. Ждёт пока не станет `healthy`
3. Тесты подключаются к `localhost:5432`
4. После тестов контейнер удаляется

**Доступны:**
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
- MySQL: `localhost:3306`

---

## 📊 Artifacts - Артефакты

### Сохранение результатов

```yaml
- name: Build binary
  run: go build -o api cmd/api/main.go

- name: Upload artifact
  uses: actions/upload-artifact@v3
  with:
    name: api-binary
    path: api
    retention-days: 30
```

**Что делает:**
- Сохраняет файлы между jobs
- Можно скачать через GitHub UI
- Автоматически удаляется через N дней

**Использование в другом job:**
```yaml
- name: Download artifact
  uses: actions/download-artifact@v3
  with:
    name: api-binary
```

---

## 📈 Coverage Reporting

### Codecov Integration

```yaml
- name: Run tests with coverage
  run: go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v3
  with:
    files: ./coverage.txt
    flags: unittests
    name: codecov-umbrella
```

**Результат:**
- Отчёт на https://codecov.io/
- Badge для README
- Комментарии в PR с изменением coverage

**Badge:**
```markdown
[![codecov](https://codecov.io/gh/AlexRodving/advanced-user-api/branch/main/graph/badge.svg)](https://codecov.io/gh/AlexRodving/advanced-user-api)
```

---

## 🚀 Deployment

### Автоматический деплой при push в main

```yaml
deploy:
  name: Deploy to Production
  runs-on: ubuntu-latest
  needs: [test-and-build, docker-build]
  if: github.ref == 'refs/heads/main'  # Только для main ветки
  
  steps:
    - name: Deploy to server
      env:
        DEPLOY_KEY: ${{ secrets.DEPLOY_KEY }}
        SERVER_HOST: ${{ secrets.SERVER_HOST }}
      run: |
        echo "$DEPLOY_KEY" > deploy_key
        chmod 600 deploy_key
        ssh -i deploy_key -o StrictHostKeyChecking=no user@$SERVER_HOST << 'EOF'
          cd /app/advanced-user-api
          git pull origin main
          docker compose down
          docker compose up -d --build
        EOF
```

**Условие:**
- `if: github.ref == 'refs/heads/main'` - только при push в main

---

## 🎭 Conditional Steps

### Выполнение только при определённых условиях

```yaml
- name: Deploy
  if: github.event_name == 'push' && github.ref == 'refs/heads/main'
  run: ./deploy.sh

- name: Comment on PR
  if: github.event_name == 'pull_request'
  run: |
    gh pr comment ${{ github.event.pull_request.number }} \
      --body "Tests passed! ✅"
```

**Операторы:**
- `&&` - AND
- `||` - OR
- `!` - NOT
- `==` - равно

---

## 📝 Badges для README

### Добавление badge статуса CI

```markdown
[![CI/CD](https://github.com/AlexRodving/advanced-user-api/actions/workflows/ci.yml/badge.svg)](https://github.com/AlexRodving/advanced-user-api/actions/workflows/ci.yml)
```

**Показывает:**
- ✅ passing (зелёный) - тесты прошли
- ❌ failing (красный) - тесты не прошли
- ⚠️ no status (серый) - ещё не запускалось

**Другие badges:**
```markdown
[![Go Report Card](https://goreportcard.com/badge/github.com/AlexRodving/advanced-user-api)](https://goreportcard.com/report/github.com/AlexRodving/advanced-user-api)

[![codecov](https://codecov.io/gh/AlexRodving/advanced-user-api/branch/main/graph/badge.svg)](https://codecov.io/gh/AlexRodving/advanced-user-api)
```

---

## 🔔 Notifications

### Slack notification при ошибке

```yaml
- name: Notify Slack
  if: failure()
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
    text: 'Build failed! 🚨'
```

### Email notification

```yaml
- name: Send email
  if: failure()
  uses: dawidd6/action-send-mail@v3
  with:
    server_address: smtp.gmail.com
    server_port: 465
    username: ${{ secrets.EMAIL_USERNAME }}
    password: ${{ secrets.EMAIL_PASSWORD }}
    subject: CI/CD Failed
    body: Build failed for ${{ github.sha }}
    to: developer@example.com
```

---

## 🎯 Примеры реальных сценариев

### Сценарий 1: Pull Request

```
1. Developer создаёт PR: feature/new-endpoint → main
   
2. GitHub Actions автоматически запускает:
   ✅ Checkout code
   ✅ Setup Go
   ✅ Run linters (go fmt, go vet)
   ✅ Run unit tests
   ✅ Run integration tests
   ✅ Build Docker image
   ✅ Security scan
   
3. Если всё ✅ → PR можно мержить
   Если что-то ❌ → нужны исправления
```

---

### Сценарий 2: Push в main

```
1. Merge PR в main
   
2. GitHub Actions запускает:
   ✅ Все тесты
   ✅ Build Docker образа
   ✅ Security scan
   ✅ Push образа в Docker Hub (если настроено)
   ✅ Deploy в production (если настроено)
   
3. Codecov обновляет статистику coverage
```

---

### Сценарий 3: Nightly Tests

```yaml
on:
  schedule:
    - cron: '0 2 * * *'  # Каждый день в 2:00
```

**Зачем:**
- Проверка зависимостей
- Тесты на latest версиях
- Генерация отчётов

---

## 🐛 Debugging Workflows

### Просмотр логов

1. GitHub → Repository → **Actions**
2. Выбрать workflow run
3. Кликнуть на job
4. Развернуть step

**Поиск ошибок:**
```
❌ Run tests
   go test -v ./...
   FAIL: TestLogin (0.00s)
   --- FAIL: TestLogin (0.00s)
       login_test.go:25: Expected 200, got 401
```

---

### Debug mode

```yaml
- name: Setup tmate session  # Подключение к runner через SSH
  if: failure()
  uses: mxschmitt/action-tmate@v3
```

---

### Вывод переменных

```yaml
- name: Debug
  run: |
    echo "Event: ${{ github.event_name }}"
    echo "Ref: ${{ github.ref }}"
    echo "SHA: ${{ github.sha }}"
    env  # Вывести все переменные окружения
```

---

## ⏱️ Оптимизация времени выполнения

### 1. Кеширование

```yaml
# Go modules
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

# Docker layers
cache-from: type=gha
cache-to: type=gha,mode=max
```

**Результат:**
- Без кеша: ~5 минут
- С кешем: ~1 минута ✅

---

### 2. Параллельные jobs

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    # ...
  
  lint:
    runs-on: ubuntu-latest
    # ...
  
  security:
    runs-on: ubuntu-latest
    # ...
```

**Выполняются одновременно!**

---

### 3. Timeout

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    timeout-minutes: 10  # Максимум 10 минут
```

**Защищает от зависших workflows**

---

## 📊 Мониторинг и статистика

### Metrics

GitHub Actions предоставляет:
- Время выполнения каждого job
- Использование minutes (бесплатно 2,000 мин/месяц для public repos)
- История запусков

### Просмотр:

1. **Actions** tab
2. Выбрать workflow
3. Статистика справа

---

## 🎓 Best Practices

### 1. Именуйте steps понятно

```yaml
# ❌ Плохо
- run: go test

# ✅ Хорошо
- name: Run unit tests with race detector
  run: go test -v -race ./tests/unit/...
```

### 2. Используйте официальные actions

```yaml
# ✅ Официальные от GitHub
uses: actions/checkout@v4
uses: actions/setup-go@v5

# ⚠️ Community actions - проверяйте надёжность
uses: random-user/some-action@v1
```

### 3. Закрепляйте версии

```yaml
# ❌ Плохо
uses: actions/checkout@latest  # Может сломаться!

# ✅ Хорошо
uses: actions/checkout@v4      # Стабильная версия
```

### 4. Используйте matrix для тестирования

```yaml
strategy:
  matrix:
    go-version: ['1.23', '1.24', '1.25']
```

### 5. Кешируйте зависимости

```yaml
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

---

## 🔄 Полный CI/CD Pipeline

```
┌─────────────────────────────────────────┐
│     Push to main / Pull Request         │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Trigger GitHub Actions          │
└──────────────┬──────────────────────────┘
               │
               ├──────────┬────────────┬──────────┐
               ▼          ▼            ▼          ▼
         ┌─────────┐ ┌────────┐ ┌─────────┐ ┌─────────┐
         │  Test   │ │  Lint  │ │Security │ │ Build   │
         │         │ │        │ │  Scan   │ │ Docker  │
         └────┬────┘ └───┬────┘ └────┬────┘ └────┬────┘
              │          │           │          │
              └──────────┴───────────┴──────────┘
                         │
                         ▼ (все успешны)
              ┌──────────────────────┐
              │   Deploy (если main) │
              └──────────┬───────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │   Production Server  │
              └──────────────────────┘
```

---

## 📖 Примеры workflow для разных задач

### Auto-merge Dependabot PRs

```yaml
name: Auto-merge Dependabot
on:
  pull_request:

jobs:
  auto-merge:
    if: github.actor == 'dependabot[bot]'
    runs-on: ubuntu-latest
    steps:
      - name: Approve PR
        run: gh pr review --approve "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Merge PR
        run: gh pr merge --auto --merge "$PR_URL"
```

---

### Release Creation

```yaml
name: Create Release
on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build
        run: go build -o api cmd/api/main.go
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: api
          body: |
            ## Changes
            - Feature 1
            - Bug fix 2
```

---

## 💰 GitHub Actions Pricing

### Free Tier (Public repos)

- ✅ Unlimited minutes
- ✅ Unlimited storage
- ✅ 20 concurrent jobs

### Free Tier (Private repos)

- 2,000 minutes/month
- 500 MB storage
- После превышения - платно

### Self-hosted Runners

```yaml
runs-on: self-hosted  # Ваш сервер
```

**Преимущества:**
- Бесплатно
- Более мощное железо
- Доступ к локальным ресурсам

---

## 🔗 Integration с другими сервисами

### Slack

```yaml
- uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

### Telegram

```yaml
- uses: appleboy/telegram-action@master
  with:
    to: ${{ secrets.TELEGRAM_TO }}
    token: ${{ secrets.TELEGRAM_TOKEN }}
    message: Build completed!
```

### Sentry (error tracking)

```yaml
- name: Create Sentry release
  uses: getsentry/action-release@v1
  with:
    environment: production
```

---

## 📖 См. также

- [Docker Guide](./DOCKER.md) - Docker и Docker Compose
- [Deployment Guide](./DEPLOY.md) - Деплой с GitHub Actions
- [Testing Guide](./TESTING.md) - Тесты в CI/CD

---

## 📚 Дополнительные ресурсы

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Workflow Syntax](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
- [Awesome Actions](https://github.com/sdras/awesome-actions) - список полезных actions
- [GitHub Actions Marketplace](https://github.com/marketplace?type=actions)

---

**Бесплатная автоматизация для всех этапов разработки!** 🎉

