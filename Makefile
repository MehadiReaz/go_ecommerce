.PHONY: help build run test clean migrate seed docker-up docker-down

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh

run-api: ## Run the API server
	@echo "Starting API server..."
	@go run cmd/api/main.go

run-worker: ## Run background workers
	@echo "Starting background workers..."
	@go run cmd/worker/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -cover ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f *.log

migrate: ## Run database migrations
	@echo "Running migrations..."
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh

seed: ## Seed the database with sample data
	@echo "Seeding database..."
	@go run scripts/seed_data.go

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

tidy: ## Tidy go modules
	@echo "Tidying modules..."
	@go mod tidy

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs: ## View Docker logs
	@docker-compose logs -f

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev: docker-up migrate seed run-api ## Start development environment

.DEFAULT_GOAL := help
