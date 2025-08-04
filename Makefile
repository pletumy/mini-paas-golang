.PHONY: help build run test clean docker-build docker-run docker-stop

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the backend application"
	@echo "  run           - Run the backend application"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker images"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker Compose"
	@echo "  frontend-install - Install frontend dependencies"
	@echo "  frontend-build   - Build frontend application"
	@echo "  frontend-start    - Start frontend development server"

# Backend commands
build:
	@echo "Building backend..."
	cd backend && go build -o bin/server cmd/server/main.go

run:
	@echo "Running backend..."
	cd backend && go run cmd/server/main.go

test:
	@echo "Running tests..."
	cd backend && go test ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin
	rm -rf frontend/build
	rm -rf frontend/node_modules

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

docker-stop:
	@echo "Stopping Docker Compose..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

# Frontend commands
frontend-install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

frontend-build:
	@echo "Building frontend..."
	cd frontend && npm run build

frontend-start:
	@echo "Starting frontend development server..."
	cd frontend && npm start

# Database commands
db-migrate:
	@echo "Running database migrations..."
	cd backend && go run cmd/migrate/main.go

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	$(MAKE) frontend-install
	$(MAKE) docker-build
	@echo "Development environment setup complete!"

# Production build
prod-build:
	@echo "Building for production..."
	$(MAKE) build
	$(MAKE) frontend-build
	$(MAKE) docker-build

# Kubernetes commands
k8s-apply:
	@echo "Applying Kubernetes manifests..."
	kubectl apply -f k8s/

k8s-delete:
	@echo "Deleting Kubernetes resources..."
	kubectl delete -f k8s/

# Health checks
health-check:
	@echo "Checking service health..."
	curl -f http://localhost:8080/health || echo "Backend health check failed"
	curl -f http://localhost:3000/health || echo "Frontend health check failed" 