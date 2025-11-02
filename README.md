# Bun Scéal

Bun (base/foundation) Scéal (story) is an approach and suplimentary tooling for defining and mapping the segments of your infrastructure in an approach that is technology agnostic and follows Archictecture as Code principles.

This solves the problem of mapping internal mental models of how your infra is split up into visual diagram and machine readable source of truth for cross referencing in policy as code

## Problem

The idea behind this project came from a number of differnet problems but which had a single underlying root cause. These are roughly, compliance scoping, control rollouts, network design and architecting cloud and network for isolation.

The underlying problem is that without having a source of truth on how you segment your infrastructure, all the other solutions built on the assumption of segmentation are build on quicksand.

One of the most common security requirements is the need for environment isolation, however most companies never actually define what an environment means to them or have a source of truth for this. This leads to having conflicting understandings, unaligned labeling (e.g Production vs prod). 

These problems get worse as you look at internal segments within an environment. It has been many years since it was good practice to have a flat open internal network. 


## Principles

- Naming is hard, therefore tooling should not be opinoated on naming, this should be driven by user.
- 
  
User guides:

- terms should be single purpose within the context of you whole audience and defined, don't use terms which are already established in your org with a different understanding. 

## Hirearchy

This tooling follows a 2 level heirarchy for segmenting infrastructure.

Segmentation Level 1: top level segment which contains tier 2 segments.
Segmentation Level 2: Second level segment which breaks up tier 1 segments.

## Metadata Inheritance

Interitance flows down the levels, therefore l2 defined under L1 will inherit the metadata of it's L1.

Metadata Precedence:

Lower level taxes precedence of upder levels, eg. L2 metadata takes precedence over inherited L1 metadata.

## Naming and plugins

L1 and L2 names are configurable

terminology.json
## Taxonomy Files

The SegL1s and SegL2s are defined in individual files under the `taxonomy` directory. For full detail on any specific SegL1or SegL2 this is the source of truth.

### Example Segment Level 1 File

[All Files can be found here: taxonomy/security-environments/](taxonomy/security-environments/)

```yaml
---
name: Production
id: production
description: |
  .......
def_sensitivity: "a"
sensitivity_reason: |
  .......
def_criticality: "1"
criticality_reason: |
  .......
def_compliance_reqs:
  - ""
```

### Example Segment Level 2 File

[All Files can be found here: taxonomy/security-domains/](taxonomy/security-domains/)

```yaml
name: "Main Product"
id: main
description: |
 ............
env_details:
  production:
    def_sensitivity: "a"
    sensitivity_reason: |
      ......
    def_criticality: "1"
    criticality_reason: |
      ......
    def_compliance_reqs:
      - ""
  staging:
  sandbox:
```

> [!NOTE]  
> For Segment Level 2 properties under any listed SegL1, where the property is missing or blank, the values will be inherited from the Segment Level 1.
>
> e.g for sandbox in the above example, the SegL2 will be list as being under the Sandbox SegL1with the same properties. 

## Development

### Prerequisites

- Go 1.22 or later
- [golangci-lint](https://golangci-lint.run/usage/install/) - Install with `make lint-install`
- [gosec](https://github.com/securego/gosec) - Install with `make sec-install`
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) - Install with `make vulncheck-install`
- [pre-commit](https://pre-commit.com/#install) (optional but recommended) - Install with `brew install pre-commit` or `pip install pre-commit`

### Getting Started

1. Clone the repository
2. Install development tools:

   ```bash
   make install-all
   ```

3. (Optional) Install pre-commit hooks for local validation:

   ```bash
   make pre-commit-install
   ```

### Available Make Targets

Run `make help` to see all available targets. Common commands:

```bash
make test          # Run tests
make test-race     # Run tests with race detector
make lint          # Run linter
make fmt           # Format code
make sec           # Run security scanner
make vulncheck     # Check for vulnerabilities
make coverage      # Generate coverage report
make ci            # Run all CI checks locally
```

### Code Quality

This project uses strict linting and security scanning:

- **golangci-lint**: Comprehensive linting with 100+ checks
- **gosec**: Security vulnerability scanning
- **govulncheck**: Dependency vulnerability scanning
- **Pre-commit hooks**: Fast local validation before commits
- **GitHub Actions**: Comprehensive CI checks on push/PR

All checks must pass before merging. Run `make ci` locally to validate your changes.

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development workflow and guidelines.

## Roadmap

### Completed

- ✅ Migrate to JSON Schema (Draft 2020-12) for validation
- ✅ Comprehensive test coverage for business logic validation (100%)
- ✅ Schema validation with santhosh-tekuri/jsonschema/v6
- ✅ Format-agnostic validation (supports YAML and JSON)
- ✅ Comprehensive tooling for code quality and security
- ✅ change terminology from security environment to hierarchy levels

### In Progress

- Make business logic configurable, implement rules in go, default status for some on/some off, override in config file

### Planned

- Refactor file loading functions to accept configurable paths for better testability
- Allow sensitivity and criticality levels to be defined in yaml
- evaluate metadata extension model for sen, crit, etc.
- colour config in yaml (with default values)
- refactor rendering code to not hard code the diagrams
- web api to allow for querying the taxonomy
- web api to allow for image generation
- MCP server for LLM queries to web API
