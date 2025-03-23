.PHONY: help build test clean run-isolation run-nplus1

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-20s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

build: ## Build all labs
	@echo "Building isolation levels lab..."
	@cd isolation-levels-lab && go build
	@echo "Building GORM N+1 lab..."
	@cd gorm-nplus1-lab && go build

test: ## Run tests for all labs
	@echo "Running tests for isolation levels lab..."
	@cd isolation-levels-lab && go test -v
	@echo "Running tests for GORM N+1 lab..."
	@cd gorm-nplus1-lab && go test -v

clean: ## Clean build artifacts
	@echo "Cleaning isolation levels lab..."
	@cd isolation-levels-lab && go clean
	@echo "Cleaning GORM N+1 lab..."
	@cd gorm-nplus1-lab && go clean

run-isolation: ## Run the isolation levels lab
	@echo "Starting isolation levels lab..."
	@cd isolation-levels-lab && go run main.go

run-nplus1: ## Run the GORM N+1 query lab
	@echo "Starting GORM N+1 lab..."
	@cd gorm-nplus1-lab && go run main.go

lint: ## Run linter on all labs
	@echo "Linting isolation levels lab..."
	@cd isolation-levels-lab && golangci-lint run
	@echo "Linting GORM N+1 lab..."
	@cd gorm-nplus1-lab && golangci-lint run

tidy: ## Run go mod tidy on all labs
	@echo "Running go mod tidy on isolation levels lab..."
	@cd isolation-levels-lab && go mod tidy
	@echo "Running go mod tidy on GORM N+1 lab..."
	@cd gorm-nplus1-lab && go mod tidy 