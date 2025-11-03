# Implementation Plan: ADR-0001 - Decouple Data Access from Validation Logic

**ADR**: [0001-decouple-data-access-from-validation-logic.md](./0001-decouple-data-access-from-validation-logic.md)

**Status**: Draft

---

## Phase 1: Extract Validation Logic (Estimated: 2-3 days)

### Phase 1 Goals

Separate validation from data access without breaking existing API

### Phase 1 Tasks

#### 1. Create `pkg/taxonomy/validator.go`

- Define `ValidationError` struct with `Field` and `Message`
- Implement `UniquenessValidator[T]` using Go 1.18+ generics
- Add `IDExtractor` and `NameExtractor` functions as parameters
- Write comprehensive unit tests with in-memory data

#### 2. Create `pkg/taxonomy/validator_test.go`

- Test uniqueness validation with various scenarios
- Test empty inputs, single items, duplicates
- Verify error messages are descriptive

#### 3. Refactor `LoadSegL1Files` and `LoadSegL2Files`

- Extract validation into separate function calls
- Use `UniquenessValidator` for ID/name checks
- Maintain identical public API signatures
- Ensure all existing tests pass

#### 4. Update existing tests

- Add specific tests for validator package
- Ensure integration tests still cover end-to-end flows

### Phase 1 Success Criteria

- All existing tests pass
- Test coverage isn't decreased for validation logic
- No breaking changes to public API
- Validation logic duplicated code removed

---

## Phase 2: Introduce Repository Pattern (Estimated: 3-5 days)

### Phase 2 Goals

Decouple data access from orchestration logic

### Phase 2 Tasks

#### 1. Create `pkg/taxonomy/repository.go`

- Define `SegL1Repository` and `SegL2Repository` interfaces
- Document interface contracts
- Define `LoadAll(source string) ([]T, error)` signature

#### 2. Implement file-based repositories

- Create `FileSegL1Repository` struct
- Create `FileSegL2Repository` struct
- Move file I/O and parsing logic from Load functions
- Accept `SchemaValidator` as dependency in constructor
- Write unit tests mocking file system where appropriate

#### 3. Create `pkg/taxonomy/service.go`

- Define `SegL1Service` and `SegL2Service` structs
- Accept repository and validator via constructor
- Implement `LoadAndValidate(source string) (map[string]T, error)`
- Handle orchestration: load → validate → convert to map

#### 4. Refactor public API functions

- `LoadSegL1Files` becomes wrapper around `SegL1Service`
- `LoadSegL2Files` becomes wrapper around `SegL2Service`
- Inject schema path as configurable parameter (with default) from config file used for naming
- Breaking changes acceptable - no external consumers

#### 5. Update all tests

- Create repository tests with mock implementations
- Create service tests with mock repositories and validators
- Update integration tests to use new structure
- Add tests for alternative repository implementations

### Phase 2 Success Criteria

- All existing tests pass
- New repository and service tests cover critical paths and edge cases
- Schema path is configurable (fixes TODO in existing code)
- Code is organized into clear, single-responsibility components
- Performance impact measured and documented (baseline vs new implementation)

---

## Phase 3: Network Data Sources (Future)

### Phase 3 Goals

Add S3 and HTTP repository implementations

### Phase 3 Prerequisites

- Phase 2 complete and stable
- Performance baseline established
- Repository interface proven with file-based implementation

### Phase 3 Tasks (TBD)

- Implement `S3SegL1Repository` and `S3SegL2Repository`
- Implement `HTTPSegL1Repository` and `HTTPSegL2Repository`
- Add retry logic and error handling for network failures
- Add integration tests with mocked network services
- Document configuration for each repository type

---

## Notes

- **Performance**: Measure baseline performance before Phase 1, validate after each phase
- **Testing Strategy**: Focus on critical paths (happy path, validation failures, I/O errors) and edge cases (empty inputs, duplicates, malformed data)
- **Configuration**: Schema path will be read from config file (same file used for naming conventions)
- **Breaking Changes**: Since there are no external consumers, we can simplify migration by removing backwards compatibility shims
