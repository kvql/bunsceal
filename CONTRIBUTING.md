# Contributing to Bun Scéal

Thank you for your interest in contributing to Bun Scéal! This document provides guidelines and workflows for contributing to the project.

## Development Environment Setup

### Prerequisites

- Go 1.22 or later
- Git
- Make
- (Optional) [pre-commit](https://pre-commit.com/#install)

### Initial Setup

1. Fork and clone the repository:

   ```bash
   git clone git@github.com:YOUR_USERNAME/bunsceal.git
   cd bunsceal
   ```

2. Install development tools:

   ```bash
   make install-all
   ```

3. (Optional but recommended) Install pre-commit hooks:

   ```bash
   # Install pre-commit if you don't have it
   # macOS: brew install pre-commit
   # Other: pip install pre-commit

   make pre-commit-install
   ```

4. Verify your setup:

   ```bash
   make ci
   ```

## Development Workflow

### 1. Create a Branch

Create a feature or fix branch from `main`:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

### 2. Make Changes

Follow these guidelines:

- **Code Style**: Run `make fmt` to format your code
- **Testing**: Add tests for new functionality
- **DRY Principle**: Avoid code duplication
- **Decoupling**: Minimize coupling where possible
- **Line Endings**: Use Linux line endings (LF) - enforced automatically

### 3. Run Local Validation

Before committing, validate your changes:

```bash
# Format code
make fmt

# Run all checks (same as CI)
make ci

# Or run individual checks:
make test          # Run tests
make lint          # Run linter
make sec           # Run security scanner
make vulncheck     # Check vulnerabilities
```

### 4. Commit Changes

If you installed pre-commit hooks, they will run automatically on commit. Otherwise, manually run `make ci` before committing.

```bash
git add .
git commit -m "feat: add new feature"
```

Commit message format (loosely following conventional commits):

- `feat:` New feature
- `fix:` Bug fix
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `docs:` Documentation changes
- `chore:` Maintenance tasks

### 5. Push and Create Pull Request

```bash
git push origin your-branch-name
```

Create a pull request on GitHub with:

- Clear description of changes
- Reference to any related issues
- Test results if applicable

## Code Quality Standards

This project enforces strict quality standards:

### Linting

- **golangci-lint** with strict preset (100+ checks)
- Covers style, bugs, performance, and best practices
- Configuration: `.golangci.yml`

### Security

- **gosec**: Detects security vulnerabilities in code
- **govulncheck**: Scans dependencies for known vulnerabilities
- Both run in CI and optionally via pre-commit

### Testing

- All new features must include tests
- Maintain or improve test coverage
- Tests run with race detector in CI

### Pre-commit Hooks

If installed, pre-commit hooks run on every commit:

- Go formatting (gofmt, goimports)
- Fast linting subset
- Go vet
- Quick test suite
- YAML/JSON validation
- Whitespace and line ending fixes

To skip hooks temporarily (not recommended):

```bash
git commit --no-verify
```

## CI Pipeline

GitHub Actions runs comprehensive checks on all PRs:

1. **Lint**: Full golangci-lint with strict config
2. **Security**: gosec + govulncheck
3. **Test**: Tests with race detector across Go 1.22 and 1.23
4. **Build**: Verify clean build
5. **Format Check**: Ensure code is properly formatted

All checks must pass before merging.

## Architecture Decision Records (ADRs)

Significant architectural decisions should be documented as ADRs:

### When to Create an ADR

Create an ADR for decisions that:

- Are one-way doors (difficult to reverse)
- Require 2+ weeks to refactor if changed
- Impact architecture, coupling, or performance
- Affect development workflow

### ADR Workflow

1. Create a new ADR:

   ```bash
   adr new "Decision title"
   ```

2. Fill in the ADR sections:
   - **Context**: Why is this decision needed?
   - **Decision**: What are we doing?
   - **Consequences**: What are the impacts (positive and negative)?
   - **Options Considered**: What alternatives were evaluated?

3. ADRs start as `draft` status
4. Never self-approve - ADRs must be reviewed
5. Location: `docs/adrs/`

For more details, see existing ADRs in `docs/adrs/`.

## Common Make Targets

```bash
make help          # Show all available targets
make install-all   # Install all required tools (golangci-lint, gosec, govulncheck, goimports)
make fmt           # Format all Go files
make vet           # Run go vet
make lint          # Run golangci-lint
make sec           # Run gosec
make vulncheck     # Run govulncheck
make test          # Run tests
make test-race     # Run tests with race detector
make coverage      # Generate coverage report
make coverage-html # Generate HTML coverage report
make build         # Build the application
make ci            # Run all CI checks locally
make clean         # Clean build artifacts
```

## Troubleshooting

### Pre-commit is slow

Pre-commit runs a fast subset of checks. For full validation, CI will run comprehensive checks.

To skip slow checks temporarily:

```bash
SKIP=golangci-lint git commit -m "message"
```

### Linting failures on existing code

If strict linting flags issues in existing code:

1. Fix the issues if straightforward
2. If extensive refactoring needed, discuss in an issue first
3. Use `//nolint:linter-name // reason` sparingly with explanation

### Tool version mismatches

Ensure tool versions match CI:

- golangci-lint: v1.55.2
- Go: 1.22+

Update tools:

```bash
make install-all
```

## Questions or Issues?

- Open an issue on GitHub
- Check existing ADRs in `docs/adrs/`
- Review project documentation in `docs/`

## License

By contributing, you agree that your contributions will be licensed under the project's license.
