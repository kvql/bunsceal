# Contributing to Bun Scéal

Thank you for your interest in contributing to Bun Scéal! This document provides guidelines and workflows for contributing to the project.

## Setup

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
   make pre-commit-install
   ```

4. Verify your setup:

   ```bash
   make ci
   ```

## Development Workflow

### 1. Create a Branch

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
- **Line Endings**: Use Linux line endings (LF)

### 3. Run Local Validation

Before committing:

```bash
make ci  # Runs all checks (same as CI pipeline)
```

If pre-commit hooks are installed, they run automatically on commit. Otherwise, manually run `make ci`.

### 4. Commit Changes

```bash
git add .
git commit -m "feat: add new feature"
```

Commit message format (conventional commits):

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

## Available Commands

Run `make help` to see all available targets. Common commands:

```bash
make fmt           # Format code
make test          # Run tests
make lint          # Run linter
make sec           # Run security scanner
make vulncheck     # Check vulnerabilities
make ci            # Run all CI checks locally
```

## Code Quality Standards

### Automated Checks

All checks are enforced by `make ci` and run in GitHub Actions:

- **Linting**: golangci-lint with strict preset - see `.golangci.yml`
- **Security**: gosec (code) + govulncheck (dependencies)
- **Testing**: Tests with race detector
- **Format**: gofmt + goimports

All checks must pass before merging.

### Test Requirements

- All new features must include tests
- Maintain or improve test coverage
- Tests run with race detector in CI

### Pre-commit Hooks

If installed, pre-commit hooks run a fast subset of checks on every commit:

- Go formatting and imports
- Fast linting subset
- Go vet
- YAML/JSON validation
- Whitespace and line ending fixes

Skip temporarily (not recommended):

```bash
git commit --no-verify
```

## Architecture Decision Records (ADRs)

Create ADRs for decisions that:

- Are one-way doors (difficult to reverse)
- Require 2+ weeks to refactor if changed
- Impact architecture, coupling, or performance

### ADR Workflow

```bash
# Create new ADR
adr new "Decision title"

# Supersede existing ADR
adr new -s <number> "New decision title"
```

Fill in required sections:

- **Context**: Why is this decision needed?
- **Decision**: What are we doing?
- **Consequences**: What are the impacts (positive and negative)?
- **Options Considered**: What alternatives were evaluated?

ADRs start as `draft` status and must be reviewed before acceptance. Location: `docs/adrs/`

See existing ADRs in `docs/adrs/` for examples.

## Troubleshooting

### Linting failures on existing code

1. Fix the issues if straightforward
2. If extensive refactoring needed, discuss in an issue first
3. Use `//nolint:linter-name // reason` sparingly with explanation

### Tool version mismatches

Ensure tool versions match CI. Update with:

```bash
make install-all
```

Check `.github/workflows/ci.yml` for required versions.

### Pre-commit is slow

Pre-commit runs a fast subset of checks. Full validation runs in CI.

Skip slow checks temporarily:

```bash
SKIP=golangci-lint git commit -m "message"
```

## Questions or Issues?

- Open an issue on GitHub
- Check existing ADRs in `docs/adrs/`
- Review project documentation in `docs/`

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
