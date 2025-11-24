.PHONY: build test test-coverage clean run help

BINARY_NAME=stellargen
BUILD_DIR=bin
OUTPUT_DIR=output

help: ## Display this help message
	@echo "Stellargen Makefile Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compile the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./tests/...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./tests/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Remove build artifacts
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(OUTPUT_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

run: build ## Build and run with default settings
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

run-csv: build ## Generate CSV files (1000 stars)
	@echo "Generating CSV files..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --num-stars=1000 --output-format=csv --output-dir=$(OUTPUT_DIR)

run-json: build ## Generate JSON files (500 stars)
	@echo "Generating JSON files..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --num-stars=500 --output-format=json --output-dir=$(OUTPUT_DIR)

run-parquet: build ## Generate Parquet files (10000 stars)
	@echo "Generating Parquet files..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --num-stars=10000 --output-format=parquet --output-dir=$(OUTPUT_DIR)

run-cassandra: build ## Insert data into Cassandra
	@echo "Inserting data into Cassandra..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --num-stars=1000 --output-format=cassandra --config=examples/config.yaml

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run ./...

all: clean fmt vet test build ## Clean, format, test, and build
