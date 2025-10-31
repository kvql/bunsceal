# Bun Scéal

Bun (base/foundation) Scéal (story) is an approach and suplimentary tooling for defining and mapping the segments of your infrastructure in an approach that is technology agnostic and follows Archictecture as Code principles

## Hirearchy

This tooling follows a 2 level heirarchy for segmenting infrastructure.

## Taxonomy Files

The SecEnvs and SecDoms are defined in individual files under the `taxonomy` directory. For full detail on any specific SecEnv or SecDom this is the source of truth.

### Example Security Environment File

[All Files can be found here: taxonomy/security-environments/](taxonomy/security-environments/)

```yaml
---
name: Production
id: production
description: |
  .......
def_sensitivity: "a"
sensitivity_reason: |
  .......
def_criticality: "1"
criticality_reason: |
  .......
def_compliance_reqs:
  - ""
```

### Example Security Domain File

[All Files can be found here: taxonomy/security-domains/](taxonomy/security-domains/)

```yaml
name: "Main Product"
id: main
description: |
 ............
env_details:
  production:
    def_sensitivity: "a"
    sensitivity_reason: |
      ......
    def_criticality: "1"
    criticality_reason: |
      ......
    def_compliance_reqs:
      - ""
  staging:
  sandbox:
```

> [!NOTE]  
> For Security Domain properties under any listed SecEnv, where the property is missing or blank, the values will be inherited from the Security Environment.
>
> e.g for sandbox in the above example, the SecDom will be list as being under the Sandbox SecEnv with the same properties. 

## Roadmap

- migrate to Json schema
- change terminology from security environment to hierarchy levels
– Allow terminology to be defined the Yaml with updated Json schema
- Allow sensitivity and criticality levels to be defined in yaml
- colour config in yaml (with default values)
- refactor rendering code to not hard code the diagrams
- web api to allow for quering the taxononmy
- web api to allow for image generation
- MCP server for LLM queries to web API