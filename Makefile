# TemplateStore & Blog Platform - Makefile
# Development commands for easy project management

.PHONY: help start stop restart status clean build test db-setup db-seed logs dev quick

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN := \033[32m
BLUE := \033[34m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m

# Configuration
BACKEND_PORT := 8080
FRONTEND_PORT := 3000
DB_NAME := template_store

help: ## Show this help message
	@echo "$(BLUE)TemplateStore & Blog Platform - Development Commands$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Quick Start:$(NC)"
	@echo "  make start    # Start all services"
	@echo "  make logs     # View logs"
	@echo "  make stop     # Stop all services"

##@ Development Environment

start: ## Start all services (database, backend, frontend)
	@echo "$(BLUE)🚀 Starting TemplateStore development environment...$(NC)"
	@./scripts/start-dev.sh

stop: ## Stop all services
	@echo "$(BLUE)🛑 Stopping TemplateStore development environment...$(NC)"
	@./scripts/stop-dev.sh

stop-all: ## Stop all services including PostgreSQL
	@echo "$(BLUE)🛑 Stopping all services including PostgreSQL...$(NC)"
	@./scripts/stop-dev.sh --with-postgres

restart: stop start ## Restart all services

quick: ## Quick start (minimal error checking)
	@echo "$(BLUE)⚡ Quick starting TemplateStore...$(NC)"
	@./scripts/quick-start.sh

status: ## Show service status
	@echo "$(BLUE)📊 Service Status:$(NC)"
	@echo -n "  Database (PostgreSQL): "
	@if pg_isready -h localhost -p 5432 >/dev/null 2>&1; then echo "$(GREEN)✅ Running$(NC)"; else echo "$(RED)❌ Not Running$(NC)"; fi
	@echo -n "  Backend API: "
	@if curl -s http://localhost:$(BACKEND_PORT)/health >/dev/null 2>&1; then echo "$(GREEN)✅ Running (http://localhost:$(BACKEND_PORT))$(NC)"; else echo "$(RED)❌ Not Running$(NC)"; fi
	@echo -n "  Frontend Web: "
	@if curl -s http://localhost:$(FRONTEND_PORT) >/dev/null 2>&1; then echo "$(GREEN)✅ Running (http://localhost:$(FRONTEND_PORT))$(NC)"; else echo "$(RED)❌ Not Running$(NC)"; fi

##@ Database Operations

db-setup: ## Create database and run migrations
	@echo "$(BLUE)🗄️  Setting up database...$(NC)"
	@createdb $(DB_NAME) 2>/dev/null || echo "Database already exists"
	@echo "$(GREEN)✅ Database setup complete$(NC)"

db-seed: ## Seed database with initial data
	@echo "$(BLUE)🌱 Seeding database...$(NC)"
	@curl -s -X POST http://localhost:$(BACKEND_PORT)/api/v1/categories/seed | jq . || echo "Categories seeded"
	@curl -s -X POST http://localhost:$(BACKEND_PORT)/api/v1/users/seed | jq . || echo "Users seeded"
	@echo "$(GREEN)✅ Database seeded$(NC)"

db-reset: ## Drop and recreate database
	@echo "$(BLUE)🔄 Resetting database...$(NC)"
	@dropdb $(DB_NAME) 2>/dev/null || echo "Database didn't exist"
	@createdb $(DB_NAME)
	@echo "$(GREEN)✅ Database reset complete$(NC)"

db-connect: ## Connect to database with psql
	@echo "$(BLUE)🔗 Connecting to database...$(NC)"
	@psql -h localhost -d $(DB_NAME)

##@ Building & Testing

build: ## Build the application
	@echo "$(BLUE)🔨 Building application...$(NC)"
	@go build -o bin/server cmd/server/main.go
	@go build -o bin/web cmd/web/main.go
	@echo "$(GREEN)✅ Build complete$(NC)"

build-all: ## Build for multiple platforms
	@echo "$(BLUE)🔨 Building for multiple platforms...$(NC)"
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -o bin/server-linux-amd64 cmd/server/main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/server-darwin-amd64 cmd/server/main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/server-windows-amd64.exe cmd/server/main.go
	@GOOS=linux GOARCH=amd64 go build -o bin/web-linux-amd64 cmd/web/main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/web-darwin-amd64 cmd/web/main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/web-windows-amd64.exe cmd/web/main.go
	@echo "$(GREEN)✅ Multi-platform build complete$(NC)"

test: ## Run all tests
	@echo "$(BLUE)🧪 Running tests...$(NC)"
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)🧪 Running tests with coverage...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

lint: ## Run linter
	@echo "$(BLUE)🔍 Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not installed, using go vet$(NC)"; \
		go vet ./...; \
	fi

##@ Development Tools

deps: ## Install/update dependencies
	@echo "$(BLUE)📦 Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies updated$(NC)"

logs: ## Show logs from all services
	@echo "$(BLUE)📋 Recent logs:$(NC)"
	@echo ""
	@echo "$(YELLOW)Backend logs:$(NC)"
	@if [ -f tmp/logs/backend.log ]; then tail -n 10 tmp/logs/backend.log; else echo "No backend logs found"; fi
	@echo ""
	@echo "$(YELLOW)Frontend logs:$(NC)"
	@if [ -f tmp/logs/frontend.log ]; then tail -n 10 tmp/logs/frontend.log; else echo "No frontend logs found"; fi

logs-follow: ## Follow logs in real-time
	@echo "$(BLUE)📋 Following logs (Ctrl+C to stop)...$(NC)"
	@if [ -f tmp/logs/backend.log ] && [ -f tmp/logs/frontend.log ]; then \
		tail -f tmp/logs/backend.log tmp/logs/frontend.log; \
	elif [ -f tmp/logs/backend.log ]; then \
		tail -f tmp/logs/backend.log; \
	elif [ -f tmp/logs/frontend.log ]; then \
		tail -f tmp/logs/frontend.log; \
	else \
		echo "$(RED)No log files found$(NC)"; \
	fi

open: ## Open application in browser
	@echo "$(BLUE)🌐 Opening application in browser...$(NC)"
	@if command -v open >/dev/null 2>&1; then \
		open http://localhost:$(FRONTEND_PORT); \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open http://localhost:$(FRONTEND_PORT); \
	elif command -v wslview >/dev/null 2>&1; then \
		wslview http://localhost:$(FRONTEND_PORT); \
	else \
		echo "Please open http://localhost:$(FRONTEND_PORT) in your browser"; \
	fi

##@ Utilities & Cleanup

clean: ## Clean build artifacts and temporary files
	@echo "$(BLUE)🧹 Cleaning up...$(NC)"
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f coverage.out coverage.html
	@go clean -cache
	@echo "$(GREEN)✅ Cleanup complete$(NC)"

clean-all: clean ## Clean everything including Go module cache
	@echo "$(BLUE)🧹 Deep cleaning...$(NC)"
	@go clean -modcache
	@echo "$(GREEN)✅ Deep cleanup complete$(NC)"

docker-build: ## Build Docker image
	@echo "$(BLUE)🐳 Building Docker image...$(NC)"
	@docker build -t templatestore:latest .
	@echo "$(GREEN)✅ Docker image built$(NC)"

docker-run: docker-build ## Run application in Docker
	@echo "$(BLUE)🐳 Running in Docker...$(NC)"
	@docker-compose up

##@ API Testing

api-health: ## Test API health endpoint
	@echo "$(BLUE)🩺 Testing API health...$(NC)"
	@curl -s http://localhost:$(BACKEND_PORT)/health | jq . || echo "API not responding"

api-test: ## Run basic API tests
	@echo "$(BLUE)🧪 Testing API endpoints...$(NC)"
	@echo "Health check:"
	@curl -s http://localhost:$(BACKEND_PORT)/health | jq .
	@echo "\nAPI info:"
	@curl -s http://localhost:$(BACKEND_PORT)/api/v1/ | jq .
	@echo "\nCategories:"
	@curl -s http://localhost:$(BACKEND_PORT)/api/v1/categories | jq .
	@echo "\nTemplates:"
	@curl -s http://localhost:$(BACKEND_PORT)/api/v1/templates | jq .
	@echo "\nBlog posts:"
	@curl -s http://localhost:$(BACKEND_PORT)/api/v1/blog | jq .

##@ Information

info: ## Show project information
	@echo "$(BLUE)📋 Project Information:$(NC)"
	@echo ""
	@echo "  Project: TemplateStore & Blog Platform"
	@echo "  Backend Port: $(BACKEND_PORT)"
	@echo "  Frontend Port: $(FRONTEND_PORT)"
	@echo "  Database: $(DB_NAME)"
	@echo ""
	@echo "$(BLUE)🌐 URLs:$(NC)"
	@echo "  • Main Application: http://localhost:$(FRONTEND_PORT)"
	@echo "  • Test Page: http://localhost:$(FRONTEND_PORT)/test.html"
	@echo "  • Backend API: http://localhost:$(BACKEND_PORT)/api/v1/"
	@echo "  • Health Check: http://localhost:$(BACKEND_PORT)/health"
	@echo ""
	@echo "$(BLUE)📂 Important Files:$(NC)"
	@echo "  • Backend: cmd/server/main.go"
	@echo "  • Frontend: cmd/web/main.go"
	@echo "  • Frontend Assets: web/"
	@echo "  • Scripts: scripts/"
	@echo ""

env: ## Show environment information
	@echo "$(BLUE)🔧 Environment Information:$(NC)"
	@echo ""
	@echo "Go version: $$(go version)"
	@echo "PostgreSQL: $$(pg_config --version 2>/dev/null || echo 'Not found')"
	@echo "Current user: $$(whoami)"
	@echo "Current directory: $$(pwd)"
	@echo ""
	@echo "$(BLUE)🔌 Port Status:$(NC)"
	@echo "Port $(BACKEND_PORT): $$(lsof -Pi :$(BACKEND_PORT) -sTCP:LISTEN -t >/dev/null 2>&1 && echo 'In Use' || echo 'Available')"
	@echo "Port $(FRONTEND_PORT): $$(lsof -Pi :$(FRONTEND_PORT) -sTCP:LISTEN -t >/dev/null 2>&1 && echo 'In Use' || echo 'Available')"
	@echo "Port 5432: $$(lsof -Pi :5432 -sTCP:LISTEN -t >/dev/null 2>&1 && echo 'In Use' || echo 'Available')"

##@ Development Workflow

dev: start db-seed open ## Full development setup (start, seed, open browser)
	@echo "$(GREEN)🎉 Development environment ready!$(NC)"

fresh: stop clean start db-seed ## Fresh start (stop, clean, start, seed)
	@echo "$(GREEN)🎉 Fresh development environment ready!$(NC)"
