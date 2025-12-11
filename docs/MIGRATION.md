# L2 Parent Selection Migration Guide

## Overview

As of this release, L2 segments now use an explicit `l1_parents` field to define parent relationships instead of relying on `l1_overrides` map keys.

**What's changing:** The practice of using `l1_overrides` map keys to determine parent relationships. The `l1_overrides` field itself remains for storing override data.

## Migration Timeline

### Phase 1: Backward Compatibility (Current)
- Both old and new formats are fully supported
- Files with only `l1_overrides` automatically migrate at load time
- No immediate action required
- Migration logs show which files are being auto-migrated

### Phase 2: Dual Format (Recommended - Adopt Now)
- Add `l1_parents` field to your L2 segment YAML files
- Keep `l1_overrides` for override data
- See examples below

### Phase 3: New Format Only (Future)
- `l1_overrides` keys will be ignored for parent relationships
- Only `l1_parents` will determine parent-child relationships
- `l1_overrides` becomes purely override data

## Format Comparison

### Old Format (Still Supported)
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

### New Format (Recommended)
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

Check logs for migration messages - once you've added `l1_parents`, you should no longer see migration logs for those files.

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
- At least one of `l1_parents` or `l1_overrides` must be present
- If both are present, `l1_overrides` keys must be a subset of `l1_parents`
- `l1_overrides` can now be empty (previously required minProperties: 1)

## Troubleshooting

### Error: "L1Override key 'X' not found in L1Parents"
**Cause:** You have an override for an L1 that isn't listed in `l1_parents`

**Fix:** Either add the L1 to `l1_parents` or remove the override

```yaml
# Wrong:
l1_parents:
  - prod
l1_overrides:
  prod: {...}
  staging: {...}  # ERROR: staging not in l1_parents

# Correct:
l1_parents:
  - prod
  - staging
l1_overrides:
  prod: {...}
  staging: {...}
```

### Error: "SegL2 'X' has parent 'Y' but no override data after inheritance"
**Cause:** Internal error - inheritance didn't populate override data

**Fix:** Check that the parent L1 segment exists in your taxonomy

### Migration Logs
During Phase 1, you may see logs like:
```
Migrated L1Parents for SegL2 'sec' from L1Overrides keys in file security.yaml
```

These are informational and indicate automatic migration is working. To stop seeing these logs, add the `l1_parents` field to your YAML files.

## Questions & Support

**Q: Do I need to migrate immediately?**
A: No. The system supports both formats indefinitely. However, we recommend migrating to take advantage of the new capabilities.

**Q: Can I have both l1_parents and l1_overrides?**
A: Yes! This is the recommended approach. Use `l1_parents` to define relationships and `l1_overrides` for override data.

**Q: What if I want to remove a parent?**
A: Remove the L1 from both `l1_parents` and `l1_overrides`. The system will validate consistency.

**Q: Can parents have different overrides?**
A: Yes! That's the whole point of the override system. Each parent can have different sensitivity, criticality, and compliance requirements.

## Examples

See `example/taxonomy/segments/` for complete examples:
- `security.yaml` - Using new format with overrides
- `monitoring.yaml` - Using new format with full inheritance
