package plugins

import (
	"fmt"

	"github.com/kvql/bunsceal/pkg/domain"
)

type ClassificationsConfig struct {
	Common          PluginsCommonSettings               `yaml:"common_settings"`
	RationaleLength int                                 `yaml:"rationale_length"`
	Definitions     map[string]ClassificationDefinition `yaml:"definitions"`
}

type ClassificationDefinition struct {
	DescriptiveName string            `yaml:"name"`
	Description     string            `yaml:"description"`
	Values          map[string]string `yaml:"values"`
	Order           []string          `yaml:"order"`
}

type ClassificationsPlugin struct {
	Config    *ClassificationsConfig
	Namespace string
}

func NewClassificationPlugin(config *ClassificationsConfig, prefix string) *ClassificationsPlugin {
	return &ClassificationsPlugin{
		Config:    config,
		Namespace: prefix + "classifications",
	}
}

func (p ClassificationsPlugin) ValidateLabels(seg *domain.Seg) PluginValidationResult {
	result := PluginValidationResult{
		Valid:  false,
		Errors: []error{},
	}

	numKeysExpected := len(p.Config.Definitions) * 2
	numKeys := len(seg.LabelNamespaces[p.Namespace])
	if numKeys != numKeysExpected {
		result.Errors = append(result.Errors, fmt.Errorf("segment %s, incorrect number of classification labels, missing rationale or classification. Expected %d, got %d", seg.ID, numKeysExpected, numKeys))
	}
	for k := range p.Config.Definitions {
		// Check if required labels are present and rationale is required
		if _, exists := seg.LabelNamespaces[p.Namespace][k]; !exists {
			result.Errors = append(result.Errors, fmt.Errorf("segment %s, missing label %s", seg.ID, p.Namespace+"/"+k))
		}
		// check values
		v := seg.LabelNamespaces[p.Namespace][k]
		if _, exists := p.Config.Definitions[k].Values[v]; !exists {
			result.Errors = append(result.Errors, fmt.Errorf("segment %s, invalid value %s for label %s", seg.ID, v, p.Namespace+"/"+k))
		}
		// check rationale is present
		if _, exists := seg.LabelNamespaces[p.Namespace][k+"_rationale"]; !exists {
			result.Errors = append(result.Errors, fmt.Errorf("segment %s, missing label %s", seg.ID, p.Namespace+"/"+k+"_rationale"))
		}
		// check rationale length
		if len(seg.LabelNamespaces[p.Namespace][k+"_rationale"]) < p.Config.RationaleLength {
			result.Errors = append(result.Errors, fmt.Errorf("segment %s, rationale too short for label %s", seg.ID, p.Namespace+"/"+k))
		}
	}
	result.Valid = true
	return result
}
