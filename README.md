# Bun Scéal

Bun (base/foundation) Scéal (story) is an approach and supplementary tooling for defining and mapping the segments of your infrastructure in an approach that is technology agnostic and follows Architecture as Code principles.

This solves the problem of mapping internal mental models of how your infra is split up into visual diagram and machine readable source of truth for cross referencing in policy as code.

## Quick Start

### Installation

```bash
# Install from source
git clone https://github.com/kvql/bunsceal.git
cd bunsceal
make install

# Or build locally
make build
```

### Basic Usage

The `example/` directory contains a working taxonomy to get you started. It includes a complete configuration with three L1 segments (dev, staging, production environments) and one L2 segment demonstrating inheritance.

```bash
# Validate your taxonomy (runs automatically with all commands)
bunsceal -config example/config.yaml

# Verify taxonomy and check that diagrams are up to date
bunsceal -config example/config.yaml -verify

# Generate visualization diagrams
bunsceal -config example/config.yaml -graph -graphDir ./output

# Export to JSON for policy-as-code integration
bunsceal -config example/config.yaml -localExport taxonomy.json
```

Copy the `example/` directory as a starting point for your own infrastructure taxonomy.

**Next**: See [Getting Started Guide](docs/getting-started.md) for detailed concepts, tutorials, and workflow.

## Problem

The idea behind this project came from a number of different problems but which had a single underlying root cause. These are roughly, compliance scoping, control rollouts, network design and architecting cloud and network for isolation.

The underlying problem is that without having a source of truth on how you segment your infrastructure, all the other solutions built on the assumption of segmentation are built on quicksand.

One of the most common security requirements is the need for environment isolation, however most companies never actually define what an environment means to them or have a source of truth for this. This leads to having conflicting understandings, unaligned labeling (e.g Production vs prod).

These problems get worse as you look at internal segments within an environment. It has been many years since it was good practice to have a flat open internal network.

## Principles

- Naming is hard, therefore tooling should not be opinionated on naming, this should be driven by user.

User guides:

- terms should be single purpose within the context of you whole audience and defined, don't use terms which are already established in your org with a different understanding.

## Key Features

- **Two-level hierarchy**: L1 segments (isolated boundaries) contain L2 segments (internal subdivisions)
- **Metadata inheritance**: L2 segments inherit criticality, sensitivity, and compliance requirements from L1
- **Schema validation**: Automated validation against JSON schemas and configurable business logic rules
- **GraphViz visualization**: Generate diagrams showing segment relationships and metadata
- **JSON export**: Export taxonomy for integration with policy-as-code tools (OPA, Sentinel, cloud policies)
- **Configurable terminology**: Use your organization's terminology (e.g., "Environment" instead of "Segment")

See [Getting Started Guide](docs/getting-started.md) for detailed explanations of hierarchy, metadata inheritance, and usage workflows.

## Roadmap

- Plugin model and metadata object for additional data and business logic
- colour config in yaml (with default values)
- refactor rendering code to not hard code the diagram list
- web api to allow for querying the taxonomy
- web api to allow for image generation

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for setup, workflow, and guidelines.

Run `make help` to see all available commands.

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.
