.PHONY: all build test clean run-nplus1 help

# Default target
all: build

# Build all labs
build:
	@echo "Building all labs..."
	@cd gorm-nplus1-lab && go build

# Run tests
test:
	@echo "Running tests..."
	@cd gorm-nplus1-lab && go test -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f gorm-nplus1-lab/test.db
	@rm -f gorm-nplus1-lab/gorm-nplus1-lab

# Run the N+1 problem lab
run-nplus1:
	@echo "Running N+1 problem lab..."
	@cd gorm-nplus1-lab && go run main.go

# Show help
help:
	@echo "Available targets:"
	@echo "  make all        - Build all labs (default)"
	@echo "  make build      - Build all labs"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make run-nplus1 - Run the N+1 problem lab"
	@echo "  make help       - Show this help message"

# Set default target
.DEFAULT_GOAL := help 