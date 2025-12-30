package plugins

import (
	"fmt"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
)

var NsPrefix = "bunsceal.plugin."

// GetAllPluginSchemas returns all plugin-related schemas for validation.
// Centralises schema registration to avoid coupling in config loader.
func GetAllPluginSchemas() []schemaValidation.ExternalSchema {
	schemas := []schemaValidation.ExternalSchema{
		{
			JSON: ClassificationsConfigSchema,
			ID:   "https://github.com/kvql/bunsceal/pkg/config/schemas/plugin-classifications.json",
		},
		{
			JSON: ComplianceConfigSchema,
			ID:   "https://github.com/kvql/bunsceal/pkg/config/schemas/plugin-compliance.json",
		},
		{
			JSON: PluginsConfigSchema,
			ID:   "https://github.com/kvql/bunsceal/pkg/config/schemas/plugins.json",
		},
	}
	return schemas
}

type PluginsCommonSettings struct {
	LabelInheritance  bool `yaml:"label_inheritance"`
	RequireCompleteL1 bool `yaml:"require_complete_l1"`
}

type ConfigPlugins struct {
	Classifications *ClassificationsConfig `yaml:"classifications"`
	Compliance      *ComplianceConfig      `yaml:"compliance"`
}

type PluginValidationResult struct {
	Valid  bool
	Errors []error
}
type ImageGroupingData struct {
	DisplayName   string
	Namespace     string
	Key           string
	ValuesMap     map[string]string
	OrderedValues []string
	OrderMap      map[string]int
}

type Plugin interface {
	ValidateLabels(*domain.Seg) PluginValidationResult
	GetEnabled() bool
	GetNamespace() string
	GetImageData() []ImageGroupingData
}

type RelationalValidator interface {
	ValidateRelationship(parent, child *domain.Seg) []error
}

type Plugins map[string]Plugin

func (p Plugins) LoadPlugins(cfg ConfigPlugins) error {
	if cfg.Classifications != nil {
		p["classifications"] = NewClassificationPlugin(cfg.Classifications, NsPrefix)
	}
	if cfg.Compliance != nil {
		p["compliance"] = NewCompliancePlugin(cfg.Compliance, NsPrefix)
	}
	return nil
}

// ValidateAllSegments validates all L1 and L2 segments against all loaded plugins.
// Skips segments with no labels for efficiency.
// Returns a flat list of all validation errors across all segments and plugins.
func (p Plugins) ValidateAllSegments(l1s, l2s map[string]domain.Seg) []error {
	var allErrors []error

	// Validate L1 segments
	for id, seg := range l1s {
		// Skip segments with no labels (Option 3 optimization)
		if len(seg.Labels) == 0 {
			continue
		}

		for pluginName, plugin := range p {
			result := plugin.ValidateLabels(&seg)
			if !result.Valid {
				for _, err := range result.Errors {
					allErrors = append(allErrors, fmt.Errorf("L1 segment %s (plugin %s): %w", id, pluginName, err))
				}
			}
		}
	}

	// Validate L2 segments
	for id, seg := range l2s {
		// Skip segments with no labels (Option 3 optimization)
		if len(seg.Labels) == 0 {
			continue
		}

		for pluginName, plugin := range p {
			result := plugin.ValidateLabels(&seg)
			if !result.Valid {
				for _, err := range result.Errors {
					allErrors = append(allErrors, fmt.Errorf("L2 segment %s (plugin %s): %w", id, pluginName, err))
				}
			}
		}
	}

	return allErrors
}

func (p Plugins) ApplyPluginInheritanceAndValidate(parent domain.Seg, child *domain.Seg) []error {
	var allErrors []error

	for _, plugin := range p {
		if !plugin.GetEnabled() {
			continue
		}

		ns := plugin.GetNamespace()

		if child.LabelNamespaces[ns] == nil {
			child.LabelNamespaces[ns] = make(map[string]string)
		}

		// Inherit from parent when child is missing values (overrides are NOT inheritance sources)
		if parentLabels, pHas := parent.LabelNamespaces[ns]; pHas {
			for k, v := range parentLabels {
				if _, childHas := child.LabelNamespaces[ns][k]; !childHas {
					child.LabelNamespaces[ns][k] = v
					child.ParsedLabels[ns+"/"+k] = v
				}
			}
		}

		if validator, ok := plugin.(RelationalValidator); ok {
			allErrors = append(allErrors, validator.ValidateRelationship(&parent, child)...)
		}
	}
	return allErrors
}
