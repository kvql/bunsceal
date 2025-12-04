# Code Architecture

Overview of the code architecture.

### Layered Architecture

| Layer | Components | Responsibilities | Dependencies |
|-------|------------|------------------|--------------|
| **User Layer** | `main.go`, `config.yaml` | Register plugins at startup; Load configuration; Initialize core services | → Core Taxonomy Layer |
| **Core Taxonomy Layer** | `TaxonomyService`, `PluginRegistry`, `SegL1`/`SegL2` Models, Core Validation | Orchestrate segment loading and validation; Manage plugin lifecycle; Store core fields + metadata map; Validate uniqueness and references; **No knowledge of concrete plugins** | → Plugin Interface, → Data Layer |
| **Plugin Interface** | `MetadataPlugin` interface | Define contract for all plugins; Decouple core from implementations; Enable Dependency Inversion | *(interface, no dependencies)* |
| **Plugin Implementations** | `ClassificationPlugin`, `CompliancePlugin`, Custom plugins | Implement domain-specific metadata logic; Define validation rules; Handle inheritance behavior; Provide configurable terminology | → Plugin Interface (implements) |
| **Data Layer** | `Repository`, `taxonomy.yaml` | Load segment definitions from YAML; Persist taxonomy data | *(no dependencies)* |