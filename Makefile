.PHONY: help dev dev-db docker-up docker-down docker-logs docker-restart test clean db-connect

# Default target
help:
	@echo "Ogenkidesuka Server - Available Commands:"
	@echo ""
	@echo "  make dev          - Run Go server locally (requires PostgreSQL)"
	@echo "  make dev-db       - Start PostgreSQL in Docker + Run Go locally (RECOMMENDED)"
	@echo ""
	@echo "  make docker-up    - Start all services in Docker"
	@echo "  make docker-down  - Stop all Docker services"
	@echo "  make docker-logs  - View Docker logs"
	@echo "  make docker-reset - Reset everything (removes data)"
	@echo ""
	@echo "  make db-connect   - Connect to PostgreSQL CLI"
	@echo "  make test         - Run tests (coming soon)"
	@echo "  make clean        - Clean build artifacts"
	@echo ""

# Development: PostgreSQL in Docker, Go locally (RECOMMENDED)
dev-db:
	@echo "🐘 Starting PostgreSQL in Docker..."
	docker-compose up -d postgres
	@echo "⏳ Waiting for database to be ready..."
	@sleep 5
	@echo "🔧 Setting up environment..."
	@if [ ! -f .env ]; then \
		echo "📝 Creating .env file with Docker credentials..."; \
		cp .env.example .env; \
		sed -i.bak 's/DB_USER=postgres/DB_USER=ogenkiuser/' .env; \
		sed -i.bak 's/DB_PASSWORD=your_password_here/DB_PASSWORD=ogenkipass123/' .env; \
		rm -f .env.bak; \
		echo "✅ .env file created!"; \
	fi
	@echo "🚀 Starting Go server..."
	go run cmd/server/main.go

# Development: Run Go locally (requires PostgreSQL installed)
dev:
	@echo "🚀 Starting Go server..."
	go run cmd/server/main.go

# Docker: Start all services
docker-up:
	@echo "🐳 Starting all services in Docker..."
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "📊 API: http://localhost:8080"
	@echo "💾 PostgreSQL: localhost:5432"

# Docker: Stop all services
docker-down:
	@echo "🛑 Stopping all services..."
	docker-compose down

# Docker: View logs
docker-logs:
	docker-compose logs -f

# Docker: Reset everything
docker-reset:
	@echo "⚠️  This will delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		echo "✅ Reset complete"; \
	fi

# Docker: Restart API service
docker-restart:
	docker-compose restart api

# Database: Connect to PostgreSQL
db-connect:
	@docker exec -it ogenkidesuka-db psql -U ogenkiuser -d ogenkidesuka

# Database: Connect to PostgreSQL
db-connect:
	@docker exec -it ogenkidesuka-db psql -U ogenkiuser -d ogenkidesuka

# Database: Run migrations
db-migrate:
	@echo "📊 Running database migrations..."
	@docker exec -i ogenkidesuka-db psql -U ogenkiuser -d ogenkidesuka < migrations/001_create_tables.sql
	@echo "✅ Migrations complete!"

# Database: Reset database (drop and recreate with migrations)
db-reset:
	@echo "⚠️  This will delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "🗑️  Dropping database..."; \
		docker exec -i ogenkidesuka-db psql -U ogenkiuser -d postgres -c "DROP DATABASE IF EXISTS ogenkidesuka;"; \
		echo "🔨 Creating database..."; \
		docker exec -i ogenkidesuka-db psql -U ogenkiuser -d postgres -c "CREATE DATABASE ogenkidesuka;"; \
		echo "📊 Running migrations..."; \
		docker exec -i ogenkidesuka-db psql -U ogenkiuser -d ogenkidesuka < migrations/001_create_tables.sql; \
		echo "✅ Database reset complete!"; \
	fi

# Test: Run unit tests
test:
	@echo "🧪 Running tests..."
	go test ./... -v

# Clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f cmd/server/server
	go clean
	@echo "✅ Clean complete"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"
