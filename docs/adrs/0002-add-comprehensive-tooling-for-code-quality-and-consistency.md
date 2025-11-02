# 2. Add comprehensive tooling for code quality and consistency

Date: 2025-11-01

## Status

draft

## Context

The project currently lacks automated code quality checks, security scanning, and consistent formatting enforcement. This creates several risks:

- No automated detection of security vulnerabilities in code or dependencies
- Inconsistent code formatting across the codebase
- No enforcement of Go best practices or idioms
- Quality issues only discovered during manual code review
- No standardized developer workflow
- Risk of tech debt accumulation without early detection

The project has test files and coverage tracking, but lacks comprehensive tooling integration for linting, security scanning, and automated quality gates. As a security-focused project (dealing with infrastructure taxonomy and security domains), maintaining high code quality and security standards is critical.

This is a one-way door decision requiring 2+ weeks to refactor if changed, and will impact all future development workflows and coupling decisions.

## Decision

Implement comprehensive tooling for code quality, security, and consistency using Go-native tools to avoid vendor lock-in:

**Core Tooling:**
- **golangci-lint**: Comprehensive linting with strict configuration (100+ checks)
- **gosec**: Go security scanner for OWASP vulnerability detection
- **govulncheck**: Official Go vulnerability scanner for dependency scanning
- **pre-commit framework**: Local hooks for fast feedback before commits
- **GitHub Actions**: CI pipeline for comprehensive checks on push/PR
- **Make**: Task orchestration without vendor lock-in

**Configuration:**
- Strict linting preset (fail on warnings)
- Pre-commit hooks for fast local checks (format, lint subset, tests)
- Comprehensive CI checks (full lint, security scan, tests with race detector)
- Linux line endings enforcement via .editorconfig and .gitattributes
- DRY principle: Makefile as single source of truth for commands

**Integration Points:**
- Pre-commit: gofmt, goimports, quick lint, go vet, go test
- CI: Full golangci-lint, gosec, govulncheck, comprehensive tests, build verification
- Developer workflow via Make targets: test, lint, fmt, vet, sec, vulncheck, coverage, ci

## Options Considered

**Option 1: CI-only checks (Rejected)**
- Simpler setup, no local tooling required
- Issues found only after pushing, slower feedback loop
- Wastes CI resources on formatting issues
- Rejected: Pre-commit provides faster feedback and catches trivial issues early

**Option 2: Minimal linting (Rejected)**
- Less disruptive to existing code
- Only catches bugs and security issues
- Doesn't enforce best practices or prevent tech debt
- Rejected: Strict linting aligns with project's focus on quality and prevents coupling issues

**Option 3: Vendor-specific tools (Rejected)**
- SonarQube, CodeClimate, or similar commercial tools
- Additional features but vendor lock-in
- Cost considerations for private repos
- Rejected: Goes against pragmatic programming principles and vendor lock-in avoidance

**Option 4: Pre-commit + CI with strict linting (Selected)**
- Fast local feedback via pre-commit
- Comprehensive CI checks for thorough validation
- Strict standards prevent tech debt early
- Go-native tools (no vendor lock-in)
- Aligns with pragmatic programming and DRY principles

## Consequences

**Positive:**
- Early detection of bugs, security issues, and code quality problems
- Consistent code formatting and style across team
- Reduced manual code review burden
- Vulnerability scanning for dependencies via govulncheck
- Fast feedback loop via pre-commit hooks
- Comprehensive CI validation before merge
- Standardized developer workflow via Makefile
- Prevention of tech debt accumulation

**Negative:**
- Initial setup effort for developers (install pre-commit, configure tools)
- Potential friction from strict linting on existing code
- CI runtime increases (mitigated with caching and parallel jobs)
- Learning curve for new contributors to understand tooling

**Risks to Mitigate:**
- Pre-commit hook failures may frustrate developers if too slow (use fast subset locally)
- Strict linting may require refactoring existing code (phased rollout if needed)
- CI failures need clear error messages and documentation
- Tool version drift between local and CI (pin versions in configs)
