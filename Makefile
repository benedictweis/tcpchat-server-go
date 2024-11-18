# Main verify target that runs all checks
.PHONY: verify
verify: clean build format-check lint-check

# Clean target to remove build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	go clean

# Build target to compile the project
.PHONY: build
build:
	@echo "Building the project..."
	go build ./...

# Target to check code formatting
.PHONY: format-check
format-check:
	@echo "Checking formatting..."
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then \
		echo "The following files need formatting:"; \
		echo "$$out"; \
		exit 1; \
	else \
		echo "All files are properly formatted."; \
	fi

# Target to format code automatically
.PHONY: format
format:
	@echo "Formatting files..."
	gofmt -w .

# Lint-check target to analyze code with linter
.PHONY: lint-check
lint-check:
	@echo "Running linter checks..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint is not installed. Please install it first."; \
		exit 1; \
	}
	golangci-lint run

# Lint-fix target to fix linting issues automatically
.PHONY: lint-fix
lint-fix:
	@echo "Fixing linting issues..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint is not installed. Please install it first."; \
		exit 1; \
	}
	golangci-lint run --fix

# Test target to run unit tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -v
