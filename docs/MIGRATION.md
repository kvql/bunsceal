# L2 Parent Selection Migration Guide

## Overview

**BREAKING CHANGE:** L2 segments now require an explicit `l1_parents` field to define parent relationships. The old pattern of using `l1_overrides` map keys for parent selection is no longer supported.

**What changed:**
- `l1_parents` field is now **required** in all SegL2 YAML files
- `l1_overrides` map keys are no longer used to determine parent relationships
- The `l1_overrides` field remains for storing override data only

## Required Changes

All L2 segment YAML files **must** be updated to include the `l1_parents` field. Files without this field will fail schema validation.

## Format Comparison

### Old Format (No Longer Supported)
```yaml
version: "1.0"
name: Security
id: sec
description: Security domain for security tooling
l1_overrides:
  staging:
    sensitivity: "D"
    sensitivity_rationale: "Staging has relaxed security controls"
    criticality: "5"
    criticality_rationale: "Non-production environment"
  shared-service:
    sensitivity: "A"
    sensitivity_rationale: "Shared services are highest risk"
    criticality: "2"
    criticality_rationale: "Critical shared infrastructure"
```

### New Format (Required)
```yaml
version: "1.0"
name: Security
id: sec
description: Security domain for security tooling
l1_parents:
  - staging
  - shared-service
l1_overrides:
  staging:
    sensitivity: "D"
    sensitivity_rationale: "Staging has relaxed security controls"
    criticality: "5"
    criticality_rationale: "Non-production environment"
  shared-service:
    sensitivity: "A"
    sensitivity_rationale: "Shared services are highest risk"
    criticality: "2"
    criticality_rationale: "Critical shared infrastructure"
```

### Full Inheritance (New Capability)
```yaml
version: "1.0"
name: Monitoring
id: mon
description: Monitoring domain with full inheritance
l1_parents:
  - staging
  - shared-service
l1_overrides: {}  # Inherits everything from parents
```

## Benefits

### 1. Explicit Parent Relationships
Parent-child relationships are now clearly defined in a dedicated field, making the taxonomy structure immediately visible.

### 2. Full Inheritance Support
You can now define parents without needing override data, enabling pure inheritance patterns. The system will populate all values at runtime.

### 3. Better Validation
Parent existence is checked separately from override data validation, providing clearer error messages.

### 4. Separation of Concerns
- **l1_parents**: Defines the relationship structure
- **l1_overrides**: Contains only the override data

This separation makes the intent clearer and reduces coupling between relationship management and data overrides.

## Migration Steps

### Step 1: Review Current Configuration
Check your current L2 segment YAML files to understand which L1 parents each segment has:

```bash
# Look at l1_overrides keys to identify parents
grep -A 20 "l1_overrides:" example/taxonomy/segments/*.yaml
```

### Step 2: Add l1_parents Field
For each L2 segment file, add the `l1_parents` field listing the parent L1 IDs:

**Before:**
```yaml
l1_overrides:
  prod:
    sensitivity: "A"
  staging:
    sensitivity: "D"
```

**After:**
```yaml
l1_parents:
  - prod
  - staging
l1_overrides:
  prod:
    sensitivity: "A"
  staging:
    sensitivity: "D"
```

### Step 3: Test
Run your taxonomy loading to ensure no errors:

```bash
go test ./...
```

Files without `l1_parents` will fail schema validation with a clear error message.

### Step 4: Simplify (Optional)
If a segment fully inherits from its parents, you can now simplify:

```yaml
l1_parents:
  - prod
  - staging
l1_overrides: {}  # Will inherit all values
```

## Schema Changes

The JSON schema now enforces:
- `l1_parents` is **required** (must be present in all SegL2 files)
- `l1_parents` must have at least one entry (minItems: 1)
- `l1_parents` must have unique entries (no duplicates)
- `l1_overrides` can be empty (for full inheritance scenarios)

## Troubleshooting

### Error: "missing properties: 'l1_parents'"
**Cause:** SegL2 YAML file is missing the required `l1_parents` field

**Fix:** Add the `l1_parents` field with at least one L1 parent ID

```yaml
# Wrong:
version: "1.0"
name: Security
id: sec
l1_overrides:
  prod: {...}

# Correct:
version: "1.0"
name: Security
id: sec
l1_parents:
  - prod
l1_overrides:
  prod: {...}
```

### Error: "SegL2 'X' has parent 'Y' but no override data after inheritance"
**Cause:** Internal error - inheritance didn't populate override data

**Fix:** Check that the parent L1 segment exists in your taxonomy

## Questions & Support

**Q: Do I need to migrate immediately?**
A: Yes. This is a breaking change. All SegL2 YAML files must have the `l1_parents` field.

**Q: Can I have both l1_parents and l1_overrides?**
A: Yes! `l1_parents` defines the parent relationships (required), and `l1_overrides` provides override data (optional - can be empty for full inheritance).

**Q: What if I want to remove a parent?**
A: Remove the L1 from `l1_parents`. You can also remove it from `l1_overrides` if you had override data.

**Q: Can parents have different overrides?**
A: Yes! That's the whole point of the override system. Each parent can have different sensitivity, criticality, and compliance requirements.

## Examples

See `example/taxonomy/segments/` for complete examples:
- `security.yaml` - Using new format with overrides
- `monitoring.yaml` - Using new format with full inheritance
