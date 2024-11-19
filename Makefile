# Main verify target that runs all checks
.PHONY: verify
verify: clean build license-check format-check test lint-clean lint-check

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

# License-check target to check that all files contain a license header
.PHONY: license-check
license-check:
	@echo "Checking license information..."
	sh ./hack/license.sh ./license-code.txt "*.go"

# License-add target to add a license header to each file if not already present
.PHONY: license-add
license-add:
	@echo "Adding license information..."
	sh ./hack/license.sh ./license-code.txt "*.go" true


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
	gofmt -l -w .

# Test target to run unit tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -v
# Lint-clean cleans the linting cache
.PHONY: lint-clean
lint-clean: installed-linter
	@echo "Running linter checks..."
	golangci-lint cache clean

# Lint-check target to analyze code with linter
.PHONY: lint-check
lint-check: installed-linter
	@echo "Running linter checks..."
	golangci-lint run --max-issues-per-linter 0 --max-same-issues 0

# Lint-fix target to fix linting issues automatically
.PHONY: lint-fix
lint-fix: installed-linter
	@echo "Fixing linting issues..."
	golangci-lint run --fix

# Installed-linter check that the linter is installed
.PHONY: installed-linter
installed-linter:
	@sh ./hack/installed.sh "golangci-lint"