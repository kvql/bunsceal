# 3. configurable-business-logic-rules

Date: 2025-11-03

## Status

Accepted

## Context

Outside of the schema validation there is a need for further validation logic to ensure the schema aligns with policy on how the segments should be defined and relate to each other.

Generally this isn't something that adds value to be opinionated on and therefore must be configurable.

### Options

#### Hard coded logic

Keeping logic hard coded into the package

Pros:

- no need to refactor

Cons:

- User would need to maintain a fork to customise the logic

#### Configurable logic

Add abstraction on the validation logic to allow configuration of the rules.

Pros:

- Validation logic can be enabled/disabled as needed
- No forking needed unless adding new logic
- Easier to add and rollout new logic in future, e.g. defaults off for few versions with warnings then on by default.

Cons:

- Initial refactor cost

## Decision

Adding configurable verification logic.

## Consequences

Positive:

- Easier to test validation logic through discrete, isolated rules
- Simpler to add new validation logic in future without code changes
- No breaking changes for users as existing validation logic will be migrated
- Rules can be selectively enabled/disabled via YAML configuration

Negative:

- Initial refactor cost to migrate existing validation logic
- Additional complexity in configuration management

## Implementation

### Architecture Decisions

**Rule Execution Model:**

- Rules are independent from each other and unordered
- Error aggregation: collect all errors and then fail, this is to enable faster troubleshooting.
- Rules are built-in and hard-coded, but take configuration input from YAML file

**Configuration Lifecycle:**

- Configuration loaded at initialization time via `NewRuleSet(config)`
- No runtime reload capability required
- Only enabled rules are instantiated and stored in RuleSet

**Deferred Decisions:**

- Performance optimization not considered at this point
- DSL-based rules (e.g., Rego) can be added later as a rule type taking query as part of configuration, priority is making current logic configurable first

### Package Structure

```text
pkg/taxonomy/
├── domain/
    |──config
|── loader_config.go         # Configuration loading and parsing
└── logic_rules.go     # Rule interface and Individual rule definitions
```

### Schema and Config

Utilise existing configuration code and schema, coupling is acceptable for now.

Schema validation in pkg/domain/schemas/config.json must be updated to include the `rules` property, otherwise LoadConfig will reject any config file with rules defined due to `additionalProperties: false`.

### Core Types

pkg/taxonomy/logic_rules.go

```go
type LogicRule interface {
    Validate(taxonomy Taxonomy) []error
}

type ValidationResult struct {
    RuleName string
    Errors   []error
}

type LogicRuleSet struct {
    LogicRules map[string]LogicRule
}

func NewLogicRuleSet(config Config) *LogicRuleSet {
    rs := &LogicRuleSet{LogicRules: make(map[string]LogicRule)}

    if config.LogicRules.SharedService.Enabled {
        rs.LogicRules["SharedService"] = NewLogicRuleSharedService(config.LogicRules.SharedService)
    }

    return rs
}

func (rs *LogicRuleSet) ValidateAll(taxonomy Taxonomy) []ValidationResult {
    var results []ValidationResult

    for name, logicRule := range rs.LogicRules {
        if errs := logicRule.Validate(taxonomy); len(errs) > 0 {
            results = append(results, ValidationResult{
                RuleName: name,
                Errors:   errs,
            })
        }
    }

    return results
}

type LogicRuleSharedService struct {
    config GeneralBooleanConfig
}

func NewSharedServiceRule(config GeneralBooleanConfig) *LogicRuleSharedService {
    return &LogicRuleSharedService{config: config}
}

func (r *LogicRuleSharedService) Validate(taxonomy Taxonomy) []error {
    // implementation
}
```

pkg/taxonomy/domain/config.go

```go

type Config struct {
    //... existing struct
    Rules LogicRulesConfig
}

type LogicRulesConfig struct {
    SharedService GeneralBooleanConfig `yaml:"shared_service"`
}

type GeneralBooleanConfig struct {
    Enabled bool `yaml:"enabled"`
}
```

### Integration point

Replace current `CompleteAndValidateTaxonomy` func with a calling func `ApplyInheritance` and new validation directly in the `LoadTaxonomy` func. This solves the breaking change in the new response from the verification functions

### Migrations:

`ValidateL2Definition` & `ValidateL1Definitions` will be moved to `ApplyInheritance` as they aren't business logic but validating that schema is actually cross-referencing existing objects

`ValidateSharedServices` will be moved to the new rule system
Code is pretty much the same only changing response slightly

`Uniqueness` check for IDs to remain in the loader files but create a new rule for uniqueness on other fields which aren't map keys
