# Getting Started with Bunsceal

This guide will help you understand Bunsceal and create your first infrastructure taxonomy.

## What is Bunsceal?

Bunsceal (Irish: "Bun" = foundation, "Scéal" = story) is a technology-agnostic tool for defining and visualizing infrastructure segmentation. It transforms internal mental models of infrastructure boundaries into:

- Visual diagrams showing segment relationships
- Machine-readable JSON for policy-as-code integration
- Validated taxonomy with automated schema and business logic validation

## Installation

```bash
git clone https://github.com/kvql/bunsceal.git
cd bunsceal
make install
# Or: make build (binary at ./bin/bunsceal)
```

Verify installation:

```bash
bunsceal -config example/config.yaml
```

## Core Concepts

### Segment Hierarchy

Bunsceal uses a two-level hierarchy:

**Level 1 (L1)**: Top-level segments that are completely isolated. Examples: Production/Staging/Dev environments, trust boundaries, separate networks. Use when segments should have zero connectivity by default.

**Level 2 (L2)**: Internal subdivisions within L1 segments. Examples: application tiers (web/app/database), security zones (PCI/general), functional areas (compute/security tooling). Use when segments are related but need different security/compliance requirements.

L1 segments contain L2 segments. L2 segments do not span multiple L1 segments.

### Metadata and Inheritance

Each segment can be classified with:

- **Criticality**: Impact of availability loss (see [criticality.md](criticality.md))
- **Sensitivity**: Impact of confidentiality/integrity loss (see [sensitivity.md](sensitivity.md))
- **Compliance requirements**: Regulatory frameworks that apply

L2 segments inherit metadata from their parent L1 segment. L2 can override inherited values by explicitly setting them. Lower levels take precedence.

### Configurable Terminology

L1 and L2 names are configurable in `config.yaml` to match your organization (e.g., "Security Domain" instead of "Segment").

## Your First Taxonomy

### Step 1: Copy the Example

```bash
cp -r example/ my-taxonomy
cd my-taxonomy
```

The example includes complete working taxonomy with config, compliance requirements, three L1 segments (dev/staging/production), and one L2 segment (security-tooling).

### Step 2: Create an L1 Segment

Create a file under `taxonomy/environments/`:

```yaml
---
name: Production
id: production
description: |
  Production environment containing all live customer-facing systems.
def_sensitivity: "a"
sensitivity_reason: |
  Contains customer data and PII requiring highest protection level.
def_criticality: "1"
criticality_reason: |
  Mission-critical systems with 24/7 availability requirement.
def_compliance_reqs:
  - "SOC2"
  - "GDPR"
```

Required fields: `name`, `id`, `description`, `def_sensitivity`, `sensitivity_reason`, `def_criticality`, `criticality_reason`, `def_compliance_reqs`

### Step 3: Create an L2 Segment

Create a file under `taxonomy/segments/`:

```yaml
---
name: "Security Tooling"
id: sec-tooling
description: |
  Security monitoring and detection infrastructure.
env_details:
  production:
    def_sensitivity: "a"
    sensitivity_reason: |
      Has visibility into all production systems and logs.
    def_criticality: "1"
    criticality_reason: |
      Required for threat detection and incident response.
    def_compliance_reqs:
      - "SOC2"
  staging:
    # Inherits all metadata from staging L1
  dev:
    # Inherits all metadata from dev L1
```

The `env_details` section maps this L2 to L1 segments. Leave L1 entries empty to inherit all metadata. Override specific fields as needed.

### Step 4: Validate

```bash
bunsceal -config config.yaml
```

If validation passes:

```text
Taxonomy is valid
```

If validation fails, the tool reports specific schema violations or business logic rule failures.

### Step 5: Generate Visualization

```bash
bunsceal -config config.yaml -graph -graphDir ./diagrams
```

Generates GraphViz diagrams showing L1 overview, L2 segments within each L1, and metadata inheritance.

**Requires**: [GraphViz](https://graphviz.org/download/) installed.

### Step 6: Export for Policy-as-Code

```bash
bunsceal -config config.yaml -localExport taxonomy.json
```

Creates JSON file for integration with policy-as-code tools (OPA, Sentinel, cloud policy engines, etc.).

### Step 7: Verify Diagram Freshness

After taxonomy changes:

```bash
bunsceal -config config.yaml -verify
```

Validates that taxonomy schema is valid and diagrams reflect current state.

## Workflow

### Creating New Segments

1. Identify whether you need L1 (new isolated boundary) or L2 (subdivision of existing L1)
2. Create YAML file in appropriate directory (`taxonomy/environments/` or `taxonomy/segments/`)
3. Copy format from existing file or refer to schema in `schema/` directory
4. Define all required metadata fields
5. **Provide rationale**: All criticality and sensitivity classifications must include reasoning

### Updating Segments

1. Modify the YAML file
2. Update rationale if classification levels change
3. Validate with `bunsceal -config config.yaml`
4. Regenerate diagrams: `bunsceal -config config.yaml -graph -graphDir ./diagrams`
5. Export updated taxonomy if using policy-as-code integration

### Isolation Principles

**L1 segments**:

- Should not connect with each other (except through designated shared services)
- Represent complete isolation boundaries
- Network connectivity denied by default

**L2 segments**:

- Represent subsets of an L1 and may communicate with other L2 in the same L1
- Connections should be explicit and denied by default
- Exist to apply different compliance/security requirements within a single L1
