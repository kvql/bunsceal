# 2. Layered directory structure for taxonomy package

Date: 2025-11-01

## Status

approved

## Context

The `pkg/taxonomy` package has grown organically to ~18 Go files in a flat structure. This creates several problems:

### Discoverability Issues

Current file naming is unclear and inconsistent:
- `secenv.go` - Not obvious this handles SegL1 (security environments)
- `secdomain.go` - Not obvious this handles SegL2 (security domains)
- `compscope.go` - Not obvious this handles compliance requirements
- `business_logic.go` - Too generic, doesn't indicate what business logic
- `validator.go` - Doesn't indicate which validation (uniqueness, schema, cross-entity)

New developers must read every file to understand what it does. No clear pattern for where new functionality belongs.

### Mixed Responsibilities

Files violate Single Responsibility Principle:
- `taxonomy.go` contains: constants (lines 14-32) + inheritance logic (lines 36+) + orchestration
- `config.go` contains: data structures + file loading + defaults
- Each entity loader (`secenv.go`, `secdomain.go`) mixes: file I/O + schema validation + business rules + data transformation

### No Architectural Layers

Current structure doesn't enforce separation of concerns:
- Data access (repositories) mixed with business logic (validators) mixed with orchestration (services)
- No clear dependency direction - any file can import any other file
- Hard to test in isolation - can't mock data access without mocking file I/O

### Not Scalable

Flat structure doesn't scale:
- Already 18 files at root level
- Adding SegL3 means another "sec*" file with no clear naming pattern
- Adding S3/HTTP data sources (from ADR-0001 Phase 2) - where do they go?
- Test files growing - 6 test files at root, will be 12+ after Phase 2

### Alignment with ADR-0001

ADR-0001 introduces layered architecture (Repository, Validation, Service layers). Current flat structure doesn't support this:
- Where do `SegL1Repository` and `FileSegL1Repository` go?
- Where does `SegL1Service` orchestration logic go?
- How to organize multiple repository implementations (File, S3, HTTP)?

## Decision

**Incrementally** reorganize `pkg/taxonomy` to address immediate discoverability and separation of concerns issues, while deferring repository/service abstraction until actually needed (ADR-0001 Phase 2).

### Phase 1: Immediate Refactoring (This ADR)

Create `domain/` and `validation/` subdirectories to separate pure domain models from business logic, and rename root-level files with clear prefixes:

```text
pkg/taxonomy/
├── domain/                    # Domain models and constants (no external dependencies)
│   ├── models.go              # Taxonomy, SegL1, SegL2, CompReq, L1Overrides
│   ├── constants.go           # SensitivityLevels, CriticalityLevels, ApiVersion
│   └── config.go              # Config, TermDef, TermConfig (data structures only)
│
├── validation/                # Business rule validation (reusable validators)
│   ├── uniqueness.go          # UniquenessValidator[T], ValidationError
│   ├── cross_entity.go        # ValidateSecurityDomains, ValidateEnv, ValidateSharedServices
│   └── schema.go              # SchemaValidator, JSON schema validation
│
├── loader_segl1.go            # LoadSegL1Files (from secenv.go)
├── loader_segl2.go            # LoadSegL2Files (from secdomain.go)
├── loader_compreq.go          # LoadCompScope (from compscope.go)
├── loader_config.go           # LoadConfig, DefaultConfig (from config.go functions)
│
├── inheritance.go             # ApplyInheritance - domain operation on Taxonomy
├── publish.go                 # PublishTaxonomy workflow
├── taxonomy.go                # LoadTaxonomy, CompleteAndValidateTaxonomy
│
└── testhelpers/               # Test utilities (existing)
    ├── files.go
    └── helpers.go
```

### Phase 2: Future Refactoring (When Implementing ADR-0001 Phase 2)

When multi-source data loading is actually needed, create `repository/` and `service/` layers:

```text
pkg/taxonomy/
├── domain/                    # (unchanged from Phase 1)
├── validation/                # (unchanged from Phase 1)
│
├── repository/                # Created when needed
│   ├── segl1_repository.go    # SegL1Repository interface + FileSegL1Repository
│   ├── segl2_repository.go    # SegL2Repository interface + FileSegL2Repository
│   └── compreq_repository.go  # CompReqRepository interface + FileCompReqRepository
│
├── service/                   # Created when needed
│   ├── segl1_service.go       # SegL1Service: repository → validation → transformation
│   ├── segl2_service.go       # SegL2Service
│   └── taxonomy_service.go    # TaxonomyService: CompleteAndValidateTaxonomy
│
└── (loader_*.go renamed to match new architecture)
```

### File Migration Mapping (Phase 1)

| Current File | New Location | Rationale |
|--------------|--------------|-----------|
| `data_types.go` | `domain/models.go` | Pure domain models, no dependencies |
| `taxonomy.go` (lines 14-32) | `domain/constants.go` | Constants belong with domain |
| `config.go` (structs only) | `domain/config.go` | Data structures separate from loading logic |
| `validator.go` | `validation/uniqueness.go` | Specific validation type |
| `business_logic.go` | `validation/cross_entity.go` | Business rules validation |
| `schema_validator.go` | `validation/schema.go` | Infrastructure for validation |
| `secenv.go` | `loader_segl1.go` | Clear prefix, stays at root until Phase 2 |
| `secdomain.go` | `loader_segl2.go` | Clear prefix, stays at root until Phase 2 |
| `compscope.go` | `loader_compreq.go` | Clear prefix, stays at root until Phase 2 |
| `config.go` (functions) | `loader_config.go` | Loading logic with clear prefix |
| `taxonomy.go` (inheritance) | `inheritance.go` | Already has clear name |
| `taxonomy.go` (orchestration) | `taxonomy.go` | Keeps main orchestration at root |
| `publish.go` | `publish.go` | Already has clear name |

### Test Organization (Phase 1)

Tests follow same structure - test file next to implementation:

```text
domain/
  models_test.go                # ← unit tests for domain models

validation/
  uniqueness.go
  uniqueness_test.go            # ← from validator_test.go
  cross_entity.go
  cross_entity_test.go          # ← from business_logic_test.go
  schema.go
  schema_test.go                # ← from schema_validator_test.go

loader_segl1_test.go            # ← from loading_test.go SegL1 tests
loader_segl2_test.go            # ← from loading_test.go SegL2 tests
loader_compreq_test.go          # ← from loading_test.go CompReq tests
loader_config_test.go           # ← from config_test.go
inheritance_test.go             # ← existing
taxonomy_test.go                # ← integration tests
```

### Dependency Rules (Phase 1)

Simplified dependency direction for Phase 1:

```text
domain/         → (no dependencies)
validation/     → domain/
loader_*.go     → domain/, validation/
taxonomy.go     → domain/, validation/, loader_*.go
inheritance.go  → domain/
```

**Phase 2** will add stricter layering when repository/ and service/ are created.

## Consequences

### Positive

**Immediate Value with Reduced Risk**

- Solves discoverability problem now with `domain/` and `validation/` directories
- Clear file naming (`loader_segl1.go`, `loader_segl2.go`) makes purpose obvious
- Much smaller migration effort (~2-3 hours vs 7-9 hours)
- No premature abstraction - repository/service layers added only when needed

**Discoverability**

- Domain models clearly separated in `domain/` - new developers know where to find data structures
- Validation logic grouped in `validation/` - obvious where business rules live
- Loader files with consistent `loader_` prefix - clear pattern for data loading
- Test files colocated with implementation

**Separation of Concerns**

- Domain models isolated from external dependencies
- Validation logic reusable and testable independently
- Loading logic still separated by entity (SegL1, SegL2, CompReq)

**Incremental Migration Path**

- Phase 1 delivers value immediately (better organization)
- Phase 2 adds repository pattern only when multi-source loading is actually implemented
- Can validate Phase 1 structure before committing to full layering
- If ADR-0001 Phase 2 requirements change, haven't over-invested in wrong abstraction

**Maintainability**

- Changes localized - updating SegL1 loading only touches `loader_segl1.go`
- Domain changes isolated to `domain/` directory
- Validation changes isolated to `validation/` directory

### Negative

**Migration Effort (Reduced but Still Present)**

- Must create `domain/` and `validation/` directories and move ~8 files
- Update imports for moved files across codebase
- Update test imports
- Rename ~5 files at root level with `loader_` prefix
- Risk of missing references in comments/docs

**Import Path Changes**

- Internal imports change: `taxonomy.SegL1` → `domain.SegL1`, `taxonomy.ValidateSecurityDomains` → `validation.ValidateSecurityDomains`
- Could break external consumers if they import internal types/functions directly
- Need to validate no external code depends on internal implementation

**Potential Import Cycles**

- Subdirectories can create import cycles if not careful
- Must ensure `domain/` has no external dependencies
- Must ensure `validation/` only depends on `domain/`

**Two-Phase Migration**

- Phase 1 structure will change again in Phase 2
- Team moves files twice (once now, once during ADR-0001 Phase 2 implementation)
- Could confuse developers about "final" structure

### Risks to Mitigate

**External Consumers Breaking**

- Mitigation: Validate no external code imports internal types before starting
- Mitigation: Re-export key types from package root if needed for backwards compatibility
- Mitigation: Check internal usage before making imports public

**Import Cycles**

- Mitigation: Enforce dependency rules (`domain/` → no deps, `validation/` → `domain/` only)
- Mitigation: Review imports carefully during migration
- Mitigation: Run `go build ./...` after each file move

**Incomplete Migration**

- Mitigation: Keep old files during migration, only delete when tests pass
- Mitigation: Grep for old file names (`secenv`, `secdomain`, `compscope`) in comments/docs

**Test Failures**

- Mitigation: Run full test suite after each phase
- Mitigation: Move and verify one directory at a time (`domain/` first, then `validation/`)

**Confusion About Phase 2**

- Mitigation: Document Phase 1 as "current structure" and Phase 2 as "future when needed"
- Mitigation: Add comments in code indicating Phase 2 plans where relevant

## Options Considered

### Option 1: Keep Flat Structure, Just Rename Files

**Rejected**

- Pros: Minimal migration effort, no import changes, no learning curve
- Cons: Doesn't solve scalability, doesn't enforce layering, doesn't prepare for Phase 2
- Example: `secenv.go` → `segl1_loader.go`, `business_logic.go` → `cross_entity_validation.go`
- Rejected because: Solves naming but not architecture - still violates SRP, still can't scale to 50+ files

### Option 2: Feature-Based Packages (Bounded Contexts)

**Rejected**

```text
pkg/
  segl1/          # All SegL1 code (loading, validation, service)
  segl2/          # All SegL2 code
  compreq/        # All compliance requirement code
  taxonomy/       # Orchestration only
```

- Pros: Clear feature boundaries, easy to find all SegL1 code
- Cons: Code duplication (each has its own repository/service patterns), shared validation split across packages, taxonomy orchestration harder
- Rejected because: Over-engineered for current needs, would need aggregate root and domain events (full DDD), doesn't align with ADR-0001 layered approach

### Option 3: Hybrid - Layers Within Features

**Rejected**

```text
pkg/taxonomy/
  segl1/
    repository.go
    service.go
    validator.go
  segl2/
    repository.go
    service.go
    validator.go
  shared/
    schema_validator.go
    uniqueness.go
```

- Pros: Combines benefits of both approaches
- Cons: Duplication (every entity has repository/service/validator), shared code in awkward `shared/` package, inconsistent with Go conventions (packages by function, not feature)
- Rejected because: More complex than Option 3 (selected), harder to share validation logic, doesn't follow Go package organization best practices

### Option 4: Full Layered Directory Structure with Repository/Service

**Rejected**

- Pros: Clear separation of concerns, enforces layering, fully prepared for ADR-0001 Phase 2, scalable
- Cons: Premature abstraction (creates empty `repository/` and `service/` directories before needed), larger migration effort (7-9 hours), violates YAGNI principle
- Rejected because: Creates abstractions before they're needed, significant effort with no immediate functional value, risky if ADR-0001 Phase 2 requirements change

### Option 5: Incremental Refactoring (Selected)

**Selected** - see Decision section

- Pros: Immediate value (`domain/` and `validation/` solve current problems), reduced migration effort (~2-3 hours), no premature abstraction, validates structure before full commitment, defers repository/service until actually needed
- Cons: Two-phase migration (files move twice), doesn't fully prepare for Phase 2 upfront
- Rationale: Best balance of immediate value vs. risk, follows YAGNI, solves discoverability problems now without betting on future requirements

## Implementation Notes

### Migration Strategy (Phase 1 Only)

**Step 1: Validate External Dependencies**

1. Search for external code importing taxonomy package types/functions
2. Verify no external dependencies on internal implementation
3. Document current public API surface

**Step 2: Create domain/ Directory**

1. Create `domain/` directory
2. Copy `data_types.go` → `domain/models.go`
3. Extract constants from `taxonomy.go` → `domain/constants.go`
4. Extract structs from `config.go` → `domain/config.go`
5. Update package declarations to `package domain`
6. Keep original files temporarily

**Step 3: Create validation/ Directory**

1. Create `validation/` directory
2. Move `validator.go` → `validation/uniqueness.go`
3. Move `business_logic.go` → `validation/cross_entity.go`
4. Move `schema_validator.go` → `validation/schema.go`
5. Update package declarations to `package validation`
6. Update imports to use `domain` package
7. Keep original files temporarily

**Step 4: Rename Root-Level Files**

1. Rename `secenv.go` → `loader_segl1.go`
2. Rename `secdomain.go` → `loader_segl2.go`
3. Rename `compscope.go` → `loader_compreq.go`
4. Extract loading functions from `config.go` → `loader_config.go`
5. Update imports to use `domain` and `validation` packages

**Step 5: Update Tests**

1. Move test files to match new structure
2. Update test imports to reference `domain` and `validation`
3. Run full test suite
4. Verify no regressions

**Step 6: Cleanup**

1. Delete original files after all tests pass
2. Grep for old file names in comments/docs and update
3. Run `go build ./...` to verify build
4. Update any package documentation

### Success Criteria (Phase 1)

- All tests pass with new structure
- No import cycles between `domain/`, `validation/`, and root package
- `go build ./...` succeeds without errors
- `domain/` package has no external dependencies (only Go stdlib)
- `validation/` package only depends on `domain/` (plus Go stdlib)
- Grep for old file names (`secenv`, `secdomain`, `compscope`, `data_types`) shows no results
- No breaking changes to package exports (if external consumers exist)

### Rollback Plan (Phase 1)

If migration fails or issues discovered:

1. Delete `domain/` and `validation/` directories
2. Revert file renames (loader_* back to original names)
3. Original files kept until Step 6 - should still work
4. Revert import changes using git
5. Run tests - should work with original structure

## Related ADRs

- **ADR-0001**: Decouple data access from validation logic
  - Phase 1 of this ADR (domain/ and validation/ directories) provides immediate organizational value
  - Phase 2 of this ADR (repository/ and service/ layers) will be implemented only when ADR-0001 Phase 2 is actually needed
  - This ADR does not depend on ADR-0001 acceptance - Phase 1 stands alone
