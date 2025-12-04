# 4. plugin-model

Date: 2025-12-04

## Status

Accepted

## Context

**Project stage**: Pre-alpha / greenfield design. No existing users to migrate.

The taxonomy currently supports metadata about segments (sensitivity, criticality) which isn't required for the fundamental function of mapping infrastructure segments to code. This metadata may not be useful to all users.

To generalise additional metadata and make it easier to expand, it should be moved to a plugin model with common functions like inheritance and graphing.

As a greenfield design choice, implementing the plugin architecture now avoids the need to refactor later when the API has stabilized and adoption has grown. The current validation logic for sensitivity/criticality demonstrates the type of domain-specific business logic that future metadata types will also require.

## Decision

Implement a plugin model for segment metadata.

### Plugin Responsibilities

Plugins define the following features:

- **Inheritance logic**: How metadata propagates through segment hierarchies
- **Terminology**: Domain-specific naming that matches business context
- **Schema validation**: Structural constraints on plugin-specific metadata
- **Business logic validation**: Semantic rules and relationships between metadata values
- **Validation failure strategy**: Plugin config defines behavior (strict, warn, optional)

### Plugin Design

**Static compilation approach**: Plugins are compiled into the binary at build time, not dynamically loaded at runtime.

- Uses existing config file
  - plugin object in the config
  - sets terminology and allowed values
  - sets any common functionality like inheritance
  - defines validation failure behavior (fail-fast vs warning)
- Plugin packages registered explicitly in `main.go`
  - No runtime discovery or dynamic loading
  - All plugins shipped with binary
  - Configuration determines which plugins are active

### Metadata Structure

Plugin-namespaced metadata structure (`metadata: {namespace: {key: value}}`) mimics cloud-native patterns (cloud provider tags, Kubernetes labels), prevents conflicts, and ensures extensibility. Follows proven patterns: Go's database/sql (static registration with config-driven selection) and configuration-driven behavior separation. Terminology remains configurable to avoid adoption blockers.

### Implementation Details

**Distribution**: All plugins compile into binary. Users download pre-built binary, activate plugins via YAML config - no Go toolchain or recompilation needed. Trade-off: larger binary size (acceptable for metadata plugins).

**Interface design**: Core `Taxonomy` handles plugin-independent logic. Plugins implement `MetadataPlugin` interface and register explicitly in `main()`. Benefits: core evolution without interface changes, clear plugin inclusion, simpler than `init()` pattern.

**Configuration separation**: YAML defines terminology, allowed values, inheritance rules. Plugin code implements validation logic. Users customize without forking; no recompilation for config changes.

**Validation modes**: Config-driven failure behavior: **strict** (fail on errors), **warn** (log warnings), **optional** (validate if present). Enables flexible adoption and per-plugin criticality.

**Third-party plugins**: Developers implement `MetadataPlugin` interface and submit PR with plugin code, `main.go` registration, and docs. **Option A** (recommended): merged into core repo, included in official releases, maintained by core team. **Option B** (advanced): separate repo, users fork and build custom binary. Trade-off: A reduces friction, B enables proprietary plugins.

**Key properties**: Dependency Inversion (core depends on interface), compile-time registration (`main.go`), runtime configuration (YAML-driven behavior), namespace separation (prevents conflicts), Open/Closed Principle (extend without modifying core).

### Other Options

#### 1. Continue as is

Current logic is highly coupled to the primary functionality of mapping segments. This makes the tool difficult to adopt if the specific metadata fields (sensitivity/criticality) don't align with business needs or terminology.

**Rejected because**: This limits adoption to organisations whose terminology exactly matches our assumptions.

#### 2. Configuration-based metadata

Define metadata fields via YAML/JSON config without code plugins. Could handle simple cases (enums, type checking, basic inheritance) but breaks down for complex validation: cross-field rules ("if critical, rationale required"), hierarchical constraints ("child can't be less sensitive than parent"), domain-specific logic (compliance framework relationships), and conditional inheritance. The existing codebase already has this complexity for sensitivity/criticality. Configuration alone would either violate DRY (reimplementing validation per type) or require adding a scripting language (recreating plugins poorly).

**Rejected because**: Configuration can't express the business logic complexity already present in the codebase. Code-based plugins provide a proper home for validation logic while maintaining extensibility.

#### 3. Fork-friendly architecture

Make metadata fields easily modifiable in code, encouraging users to fork and customise.

**Rejected because**: Creates maintenance burden for users who must merge upstream changes into their forks. Fragments the ecosystem and prevents sharing of domain-specific extensions (e.g., compliance frameworks) across organisations.

#### 4. Go native plugin system

**Rejected because**: platform limitations, version coupling, and poor user experience compared to static compilation."

## Consequences

### Benefits

- **Extensibility**: Add metadata types as plugins without core changes. Follows Open/Closed Principle. Plugin authors iterate independently.
- **Reduced adoption friction**: No forced metadata models. Terminology customizable via config. Pre-built binary - no Go toolchain needed.
- **Separation of concerns**: Core decoupled from domain metadata. Plugins encapsulate features. Isolated testing and maintenance.
- **Ecosystem enablement**: Third parties build compliance/industry-specific plugins. Community contributions via PR. Clear interface contract.

### Trade-offs

- **Initial complexity**: Requires plugin interfaces, registration mechanisms, and developer documentation.
- **Refactoring effort**: Current metadata tightly coupled across layers (schemas, types, validation, constants, graphing). Requires refactoring to plugin model.
- **Plugin management**: Binary size increases. New plugins need binary release. Core maintainers as gatekeepers. Coordinated releases. Custom plugins require user-maintained builds.

- **Performance testing deferred**: Pre-alpha focus on correctness over speed. Plugin overhead negligible vs I/O. Profile post-alpha if needed.
- **Interface versioning deferred**: Pre-alpha accepts breaking interface changes to find right design. Versioning added post-alpha when API stabilizes.

### Refactoring Strategy

**Pre-alpha project** - no external users, no timelines. Use test-driven refactoring: establish comprehensive tests for current sensitivity/criticality behavior (validation, inheritance, graphing), implement plugin architecture, migrate existing logic to classification plugin, ensure all tests pass, then remove old coupled code. Tests define the contract that plugin equivalents must satisfy.
