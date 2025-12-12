# 1. Decouple data access from validation logic

Date: 2025-11-01

## Status

approved

## Context

The current implementation of taxonomy data loading (`LoadSegL1Files`, `LoadSegFiles`) violates the Single Responsibility Principle by combining multiple concerns in a single function:

- File I/O operations (reading directories and files)
- Schema validation using JSON schemas
- YAML parsing
- Business rule validation (ID and name uniqueness)
- Error collection and reporting

This monolithic structure creates several problems:

1. **Testing Challenges**: Cannot test validation logic without performing actual file I/O operations
2. **Tight Coupling**: Data source is hardcoded to filesystem, preventing alternative sources (HTTP, database, S3)
3. **Code Duplication**: Uniqueness validation logic is duplicated between SegL1 and Seg implementations
4. **Configuration Issues**: Schema path is hardcoded (`./schema`), creating testing difficult
5. **Maintainability**: Changes to validation require modifying functions that also handle I/O concerns

The architecture needs to support:

- Independent testing of business validation rules
- Reusable validation logic across different entity types
- **Multiple data sources**: S3, HTTP endpoints, and other network storage (concrete requirement)
- Dependency injection for better testability

## Decision

Implement a layered architecture that separates concerns into distinct responsibilities:

### Layer 1: Repository (Data Access)

- Define `SegL1Repository` and `SegRepository` interfaces
- Implement file-based repositories (`FileSegL1Repository`, `FileSegRepository`)
- Responsibilities: File I/O, schema validation, parsing
- Future: Easy to add HTTP, database, or other implementations

### Layer 2: Validation

- Create generic `UniquenessValidator[T]` using Go generics
- Define `ValidationError` type for structured error reporting
- Support composite validation patterns
- Responsibilities: Business rule validation only
- Reusable across SegL1, Seg, and future entity types

### Layer 3: Service (Orchestration)

- Implement `SegL1Service` and `SegService`
- Compose repository and validators via dependency injection
- Responsibilities: Coordinate loading and validation workflow
- Convert slice results to maps as needed

### Layer 4: Public API

- Refactor public functions to accept dependencies directly
- **No backwards compatibility required** - confirmed no external consumers
- Simplify API to use service layer directly

## Consequences

### Positive

- **Testability**: Validators can be tested with in-memory data, no file I/O required
- **Reusability**: Generic validators work with any entity type, eliminating duplication
- **Extensibility**: New data sources only require implementing Repository interface
- **Maintainability**: Clear separation of concerns makes code easier to understand and modify
- **Flexibility**: Schema path and other dependencies can be injected, improving testability
- **Type Safety**: Go generics provide compile-time type checking for validators

### Negative

- **Complexity**: More files and indirection adds cognitive load for developers
- **Migration Effort**: Existing code requires refactoring across 2 implementation phases
- **Breaking Changes**: Callers of `LoadSegL1Files`/`LoadSegFiles` will need updates

### Risks to Mitigate

- Document patterns clearly for team onboarding
- Validate performance impact of additional abstraction layers
- Maintain comprehensive test coverage during refactoring

## Options Considered

### Option 1: Keep Current Monolithic Structure
**Rejected**

- Pros: No migration effort, familiar to team
- Cons: Doesn't address scalability issues, testing remains difficult, duplication continues

### Option 2: Full Domain-Driven Design
**Rejected**

- What it would include:
  - **Aggregates**: `TaxonomyAggregate` as root managing all validation invariants
  - **Value Objects**: `SegmentID`, `SegmentName` as immutable types with their own validation
  - **Domain Events**: `SegL1Created`, `SegL1ValidationFailed` event streams
  - **Domain Services**: `SegmentUniquenessService`, `TaxonomyConsistencyService`
  - **Anti-Corruption Layer**: Complete isolation of domain model from infrastructure
- Pros: Maximum flexibility, comprehensive separation, rich domain model
- Cons: Overhead of entity/value object/aggregate patterns not justified for data loading and validation use case; significant implementation effort; adds abstractions (events, ACL) with no clear benefit for current requirements

### Option 3: Layered Architecture with Repository + Service Patterns

**Selected**

- Pros: Pragmatic balance of separation and simplicity, supports testing and extensibility, incremental migration path
- Cons: Some added complexity, requires team to learn new patterns
- Rationale: Provides benefits of separation without full DDD overhead

## Implementation Approach

Implementation will proceed in 2 phases (estimated 5-8 days total):

1. **Phase 1: Extract Validation** - Separate validation logic using generics, maintain current API
2. **Phase 2: Repository Pattern** - Introduce repository interfaces and services, enable S3/HTTP sources

Detailed tasks, timelines, and success criteria are documented in: [0001-implementation-plan.md](./0001-implementation-plan.md)

### Key Implementation Constraints

- **Configuration**: Schema path will be read from config file (same file used for naming conventions)
- **Testing Strategy**: Focus on critical paths and edge cases, not coverage percentages
- **Performance**: Establish baseline before Phase 1, validate after each phase
- **Breaking Changes**: No backwards compatibility shims required (no external consumers)
- **Migration**: Incremental approach allows validation in production before full repository rollout
