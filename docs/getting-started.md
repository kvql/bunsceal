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
labels:
  - "bunsceal.plugin.classifications/sensitivity:A"
  - "bunsceal.plugin.classifications/sensitivity_rationale:Contains customer data and PII requiring highest protection level."
  - "bunsceal.plugin.classifications/criticality:1"
  - "bunsceal.plugin.classifications/criticality_rationale:Mission-critical systems with 24/7 availability requirement."
  - "bunsceal.plugin.compliance/soc2:in-scope"
  - "bunsceal.plugin.compliance/soc2_rationale:Production systems must comply with SOC2 requirements."
  - "bunsceal.plugin.compliance/gdpr:in-scope"
  - "bunsceal.plugin.compliance/gdpr_rationale:Processes customer PII subject to GDPR."
```

Required fields: `name`, `id`, `description`, `labels`

Labels use plugin namespaces:
- Classifications: `bunsceal.plugin.classifications/{classification}:{value}` and `bunsceal.plugin.classifications/{classification}_rationale:{text}`
- Compliance: `bunsceal.plugin.compliance/{requirement}:{in-scope|out-of-scope}` and `bunsceal.plugin.compliance/{requirement}_rationale:{text}`

### Step 3: Create an L2 Segment

Create a file under `taxonomy/segments/`:

```yaml
---
name: "Security Tooling"
id: sec-tooling
description: |
  Security monitoring and detection infrastructure.
l1_parents:
  - production
  - staging
  - dev
l1_overrides:
  production:
    labels:
      - "bunsceal.plugin.classifications/sensitivity:A"
      - "bunsceal.plugin.classifications/sensitivity_rationale:Has visibility into all production systems and logs."
      - "bunsceal.plugin.classifications/criticality:1"
      - "bunsceal.plugin.classifications/criticality_rationale:Required for threat detection and incident response."
      - "bunsceal.plugin.compliance/soc2:in-scope"
      - "bunsceal.plugin.compliance/soc2_rationale:Security tooling must comply with SOC2 requirements."
  staging:
    # Inherits all metadata from staging L1
  dev:
    # Inherits all metadata from dev L1
```

The `l1_parents` array lists which L1 segments contain this L2. The `l1_overrides` section allows overriding inherited labels for specific L1 parents. Leave L1 entries empty (or omit entirely) to inherit all metadata from the parent L1.

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
3. Copy format from existing file or refer to schema in `pkg/domain/schemas/` directory
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
