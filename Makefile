# Load environment variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Database configuration with fallback to defaults
DB_SERVER ?= localhost
DB_PORT ?= 5432
DB_USER ?= root
DB_PASS ?= 
DB_NAME ?= public
DB_SSL_MODE ?= false
DB_TIME_ZONE ?= Asia/Jakarta

# Construct database URL
DATABASE_URL = postgres://$(DB_USER):$(DB_PASS)@$(DB_SERVER):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)&timezone=$(DB_TIME_ZONE)

# Run the application
run:
	@echo "Starting payroll application..."
	go run main.go

# Build the application
build:
	@echo "Building payroll application..."
	go build -o bin/payrolls main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Run tests with verbose output
test-v:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...

# Generate wire code
wire:
	@echo "Generating wire code..."
	go generate ./infrastructure/http

# Run migration create
# Migration commands usage:
#   make migrate-create name=migration_name
#   Example: make migrate-create name=create_users_table
#   Example: make migrate-create name=add_email_to_users
migrate-create:
	@echo "Running migrations create..."
	migrate create -ext sql -dir database/migrations $(name)

# Run migration up
migrate-up:
	@echo "Running migrations up..."
	@echo "Database URL: $(DATABASE_URL)"
	migrate -path database/migrations -database "$(DATABASE_URL)" up

# Run migration down
migrate-down:
	@echo "Running migrations down..."
	@echo "Database URL: $(DATABASE_URL)"
	migrate -path database/migrations -database "$(DATABASE_URL)" down 1

# Run migration down
migrate-down-all:
	@echo "Running migrations down all..."
	@echo "Database URL: $(DATABASE_URL)"
	migrate -path database/migrations -database "$(DATABASE_URL)" down

# Show migration version
migrate-version:
	@echo "Current migration version:"
	@echo "Database URL: $(DATABASE_URL)"
	migrate -path database/migrations -database "$(DATABASE_URL)" version

# Run database seeder
seed:
	@echo "Running database seeder..."
	go run cmd/seeder/main.go -type=all

# Run specific user seeder
seed-users:
	@echo "Running user seeder..."
	go run cmd/seeder/main.go -type=users

# Show database configuration
db-config:
	@echo "Database Configuration:"
	@echo "  Server: $(DB_SERVER)"
	@echo "  Port: $(DB_PORT)"
	@echo "  User: $(DB_USER)"
	@echo "  Database: $(DB_NAME)"
	@echo "  SSL Mode: $(DB_SSL_MODE)"
	@echo "  Timezone: $(DB_TIME_ZONE)"
	@echo "  Full URL: $(DATABASE_URL)"
	@if exist .env (echo "  .env file: Found") else (echo "  .env file: Not found (using defaults)")

