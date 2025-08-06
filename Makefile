# USFM Parser Makefile

# Variables
BINARY_NAME=usfmp
CMD_PATH=./cmd/usfmp
PKG_PATH=./pkg/...
INTERNAL_PATH=./internal/...
BUILD_DIR=build
VERSION?=0.0.1

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"

.PHONY: all build clean test coverage lint fmt vet deps help run-sample docs serve-docs

# Default target
all: clean deps lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v $(PKG_PATH) $(INTERNAL_PATH)

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out $(PKG_PATH) $(INTERNAL_PATH)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint the code
lint:
	@echo "Running linter..."
	@if command -v $(GOLINT) > /dev/null 2>&1; then \
		$(GOLINT) run ./...; \
	else \
		echo "golangci-lint not found, skipping lint"; \
		echo "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format the code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Vet the code
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/ 2>/dev/null || cp $(BUILD_DIR)/$(BINARY_NAME) ~/go/bin/ || echo "Could not install to Go bin directory"

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux amd64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	# Linux arm64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_PATH)
	# Darwin amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	# Darwin arm64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)
	# Windows amd64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	@echo "Multi-platform build complete"

# Run the binary with sample data
run-sample: build
	@echo "Running $(BINARY_NAME) with sample data..."
	@if [ -d "bsb_usfm" ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) -f json bsb_usfm/01GENBSB.SFM | head -20; \
	else \
		echo "Sample data not found. Please ensure bsb_usfm directory exists."; \
	fi

# Run quick test with Genesis file
test-genesis: build
	@echo "Testing with Genesis file..."
	@if [ -f "bsb_usfm/01GENBSB.SFM" ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) -f txt bsb_usfm/01GENBSB.SFM | head -50; \
	else \
		echo "Genesis sample file not found"; \
	fi

# Development workflow - quick build and test
dev: fmt vet test build
	@echo "Development build complete"

# Generate and view documentation
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./pkg/usfm > docs.txt
	$(GOCMD) doc -all ./internal/formatter >> docs.txt
	@echo "Documentation written to docs.txt"

# Serve documentation locally (requires godoc tool)
serve-docs:
	@echo "Starting documentation server..."
	@echo "Visit http://localhost:6060/pkg/github.com/arenzana/usfmp/"
	@if command -v godoc > /dev/null 2>&1; then \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Help target
help:
	@echo "USFM Parser Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all         - Clean, deps, lint, test, and build (default)"
	@echo "  build       - Build the binary"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run unit tests"
	@echo "  coverage    - Run tests with coverage report"
	@echo "  lint        - Run code linter"
	@echo "  fmt         - Format source code"
	@echo "  vet         - Run go vet"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  install     - Install binary to Go bin directory"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  run-sample  - Build and run with sample data"
	@echo "  test-genesis- Test with Genesis sample file"
	@echo "  dev         - Quick development build (fmt, vet, test, build)"
	@echo "  docs        - Generate documentation to docs.txt"
	@echo "  serve-docs  - Start local documentation server"
	@echo "  help        - Show this help message"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION     - Set build version (default: $(VERSION))"
	@echo "  BINARY_NAME - Binary name (default: $(BINARY_NAME))"