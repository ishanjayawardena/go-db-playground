.PHONY: build run clean test lint db-up db-down db-logs db-shell

# Default target
all: build

# Build the application
build:
	go build -o bin/isolation-levels-lab

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Start the database
db-up:
	docker-compose up -d

# Stop the database
db-down:
	docker-compose down

# View database logs
db-logs:
	docker-compose logs -f postgres

# Connect to database shell
db-shell:
	docker exec -it isolation-levels-lab-postgres-1 psql -U isolation_user -d isolation_db

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  db-up      - Start the database"
	@echo "  db-down    - Stop the database"
	@echo "  db-logs    - View database logs"
	@echo "  db-shell   - Connect to database shell"
	@echo "  fmt        - Format code"
	@echo "  tidy       - Tidy dependencies"
	@echo "  help       - Show this help message" 