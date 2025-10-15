.PHONY: help run test test-coverage docker-up docker-down docker-build migrate-up migrate-down swagger lint fmt

help: ## Показать помощь
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Запустить приложение локально
	go run cmd/api/main.go

build: ## Собрать бинарник
	go build -o bin/api cmd/api/main.go

test: ## Запустить все тесты
	go test -v ./...

test-coverage: ## Запустить тесты с покрытием
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Открыть coverage.html в браузере"

test-unit: ## Только unit тесты
	go test -v ./tests/unit/...

test-integration: ## Только integration тесты
	go test -v ./tests/integration/...

docker-up: ## Запустить через Docker Compose
	docker-compose up -d

docker-down: ## Остановить Docker Compose
	docker-compose down

docker-logs: ## Логи Docker Compose
	docker-compose logs -f

docker-build: ## Собрать Docker образ
	docker build -t advanced-api:latest -f docker/Dockerfile .

migrate-up: ## Применить миграции
	docker exec -i advanced-api-postgres psql -U postgres -d advanced_api < migrations/001_create_users.sql

migrate-down: ## Откатить миграции
	docker exec -i advanced-api-postgres psql -U postgres -d advanced_api -c "DROP TABLE IF EXISTS users;"

swagger: ## Генерация Swagger документации
	swag init -g cmd/api/main.go

lint: ## Проверка кода
	golangci-lint run

fmt: ## Форматирование кода
	go fmt ./...
	goimports -w .

deps: ## Установить зависимости
	go mod download
	go mod tidy

clean: ## Очистка
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker-compose down -v

