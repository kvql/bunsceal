.PHONY: help install-all goimports-install fmt fmt-check vet lint lint-install sec sec-install vulncheck vulncheck-install test test-race coverage coverage-html build clean ci pre-commit-install

# Default target shows help
help:
	@echo "Available targets:"
	@echo "  make install-all      - Install all required tools (golangci-lint, gosec, govulncheck, goimports)"
	@echo "  make fmt              - Format all Go files"
	@echo "  make fmt-check        - Check if files are formatted"
	@echo "  make vet              - Run go vet"
	@echo "  make lint             - Run golangci-lint"
	@echo "  make lint-install     - Install golangci-lint"
	@echo "  make sec              - Run gosec security scanner"
	@echo "  make sec-install      - Install gosec"
	@echo "  make vulncheck        - Run govulncheck for vulnerabilities"
	@echo "  make vulncheck-install - Install govulncheck"
	@echo "  make goimports-install - Install goimports"
	@echo "  make test             - Run tests"
	@echo "  make test-race        - Run tests with race detector"
	@echo "  make coverage         - Generate coverage report"
	@echo "  make coverage-html    - Generate HTML coverage report"
	@echo "  make build            - Build the application"
	@echo "  make clean            - Clean build artifacts and coverage"
	@echo "  make ci               - Run all CI checks (fmt, vet, lint, sec, vulncheck, test-race, build)"
	@echo "  make pre-commit-install - Install git pre-commit hook (uses Makefile targets, no deps required)"

# Install all required tools
install-all: goimports-install lint-install sec-install vulncheck-install
	@echo ""
	@echo "All tools installed successfully!"
	@echo "Run 'make ci' to verify your setup."

#check all
check-all: fmt-check vet lint sec vulncheck
# Formatting
fmt:
	@echo "Formatting Go files..."
	@gofmt -w -s .
	@goimports -w .

fmt-check:
	@echo "Checking Go formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Files not formatted:"; gofmt -l .; exit 1)

# Vetting
vet:
	@echo "Running go vet..."
	@go vet ./...

# Linting
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --config .golangci.yml

lint-install:
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

goimports-install:
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest

# Security scanning
sec:
	@echo "Running gosec security scanner..."
	@gosec -fmt=text -severity=medium ./...

sec-install:
	@echo "Installing gosec..."
	@go install github.com/securego/gosec/v2/cmd/gosec@latest

# Vulnerability checking
vulncheck:
	@echo "Running govulncheck..."
	@govulncheck ./...

vulncheck-install:
	@echo "Installing govulncheck..."
	@go install golang.org/x/vuln/cmd/govulncheck@latest

# Testing
test:
	@echo "Running tests..."
	@go test -v ./...

test-race:
	@echo "Running tests with race detector..."
	@go test -v -race ./...

coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

coverage-html:
	@echo "Generating HTML coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Building
build:
	@echo "Building application..."
	@go build -v -o bin/bunsceal .

# Cleanup
clean:
	@echo "Cleaning build artifacts and coverage..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# CI target - runs all checks
ci: fmt-check vet lint sec vulncheck test-race build
	@echo "All CI checks passed!"

# Git hook setup - uses built-in git hooks with Makefile targets
pre-commit-install:
	@echo "Installing git pre-commit hook..."
	@mkdir -p .git/hooks
	@printf '#!/bin/sh\n' > .git/hooks/pre-commit
	@printf '# Auto-generated git pre-commit hook\n' >> .git/hooks/pre-commit
	@printf '# Regenerate with: make pre-commit-install\n' >> .git/hooks/pre-commit
	@printf '# Skip with: git commit --no-verify\n\n' >> .git/hooks/pre-commit
	@printf 'echo "Running pre-commit checks..."\n\n' >> .git/hooks/pre-commit
	@printf '# Format code\n' >> .git/hooks/pre-commit
	@printf 'make fmt || exit 1\n\n' >> .git/hooks/pre-commit
	@printf '# Run linter\n' >> .git/hooks/pre-commit
	@printf 'make lint || {\n' >> .git/hooks/pre-commit
	@printf '  echo "✗ Linting failed. Run '\''make lint'\'' to see details"\n' >> .git/hooks/pre-commit
	@printf '  exit 1\n' >> .git/hooks/pre-commit
	@printf '}\n\n' >> .git/hooks/pre-commit
	@printf 'echo "✓ Pre-commit checks passed"\n' >> .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✓ Git pre-commit hook installed at .git/hooks/pre-commit"
	@echo ""
	@echo "The hook runs: make fmt && make lint"
	@echo ""
	@echo "To skip: git commit --no-verify"
	@echo "To customize: edit .git/hooks/pre-commit or regenerate with 'make pre-commit-install'"
