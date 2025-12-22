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
	EnforceOrder    bool              `yaml:"enforce_order"`
	Values          map[string]string `yaml:"values"`
	Order           []string          `yaml:"order"`
}

type ClassificationsPlugin struct {
	Config     *ClassificationsConfig
	Namespace  string
	OrderIndex map[string]map[string]int // defKey -> value -> index
}

func NewClassificationPlugin(config *ClassificationsConfig, prefix string) *ClassificationsPlugin {
	orderIndex := make(map[string]map[string]int)
	for defKey, def := range config.Definitions {
		orderIndex[defKey] = make(map[string]int)
		for i, v := range def.Order {
			orderIndex[defKey][v] = i
		}
	}
	return &ClassificationsPlugin{
		Config:     config,
		Namespace:  prefix + "classifications",
		OrderIndex: orderIndex,
	}
}

func (p ClassificationsPlugin) validateNamespaceLabels(labels map[string]string, ctx string, errs *[]error) int {
	foundKeys := 0
	for defKey, def := range p.Config.Definitions {
		classification, hasClass := labels[defKey]
		rationale, hasRat := labels[defKey+"_rationale"]

		if hasClass && !hasRat {
			foundKeys++
			*errs = append(*errs, fmt.Errorf("%s has %s but missing %s_rationale", ctx, defKey, defKey))
		}
		if hasRat && !hasClass {
			foundKeys++
			*errs = append(*errs, fmt.Errorf("%s has %s_rationale but missing %s", ctx, defKey, defKey))
		}
		if hasClass && hasRat {
			foundKeys++
			foundKeys++
			if _, exists := def.Values[classification]; !exists {
				*errs = append(*errs, fmt.Errorf("%s invalid value %s for %s", ctx, classification, defKey))
			}
			if len(rationale) < p.Config.RationaleLength {
				*errs = append(*errs, fmt.Errorf("%s %s_rationale too short (min %d chars)", ctx, defKey, p.Config.RationaleLength))
			}
		}
	}
	return foundKeys
}

func (p ClassificationsPlugin) ValidateLabels(seg *domain.Seg) PluginValidationResult {
	result := PluginValidationResult{Valid: false, Errors: []error{}}

	segLabels := seg.LabelNamespaces[p.Namespace]
	foundKeys := p.validateNamespaceLabels(segLabels, "segment "+seg.ID, &result.Errors)

	// L1 segments must have all classification definitions (if RequireCompleteL1 is enabled)
	if seg.Level == "1" && p.Config.Common.RequireCompleteL1 {
		numKeysExpected := len(p.Config.Definitions) * 2
		if foundKeys != numKeysExpected {
			result.Errors = append(result.Errors, fmt.Errorf("segment %s missing classification labels. Expected %d, found %d", seg.ID, numKeysExpected, foundKeys))
		}
		if len(segLabels) != numKeysExpected {
			result.Errors = append(result.Errors, fmt.Errorf("segment %s has extra labels. Expected %d, got %d", seg.ID, numKeysExpected, len(segLabels)))
		}
	}

	for parentID, override := range seg.L1Overrides {
		if len(override.LabelNamespaces[p.Namespace]) > 0 {
			p.validateNamespaceLabels(override.LabelNamespaces[p.Namespace], fmt.Sprintf("segment %s l1_override[%s]", seg.ID, parentID), &result.Errors)
		}
	}

	result.Valid = len(result.Errors) == 0
	return result
}

// ValidateRelationship checks parent >= child in severity order for all definitions
func (p ClassificationsPlugin) ValidateRelationship(parent, child *domain.Seg) []error {
	var errs []error
	override, hasOverride := child.L1Overrides[parent.ID]

	for defKey, def := range p.Config.Definitions {
		if !def.EnforceOrder {
			continue
		}

		parentValue := parent.LabelNamespaces[p.Namespace][defKey]
		childValue := child.LabelNamespaces[p.Namespace][defKey]
		if hasOverride && len(override.LabelNamespaces[p.Namespace]) > 0 {
			childValue = override.LabelNamespaces[p.Namespace][defKey]
		}

		parentIdx, pOk := p.OrderIndex[defKey][parentValue]
		childIdx, cOk := p.OrderIndex[defKey][childValue]

		if pOk && cOk && childIdx < parentIdx {
			errs = append(errs, fmt.Errorf("child %s has higher %s (%s) than parent %s (%s)",
				child.ID, defKey, childValue, parent.ID, parentValue))
		}
	}
	return errs
}

func (p ClassificationsPlugin) GetEnabled() bool {
	return p.Config.Common.LabelInheritance
}

func (p ClassificationsPlugin) GetNamespace() string {
	return p.Namespace
}

func (p ClassificationsPlugin) GetImageData() []ImageGroupingData {
	dataList := []ImageGroupingData{}
	for key, def := range p.Config.Definitions {
		data := ImageGroupingData{
			Namespace:     p.Namespace,
			DisplayName:   def.DescriptiveName,
			OrderedValues: def.Order,
			OrderMap:      p.OrderIndex[key],
			Key:           key,
			ValuesMap:     def.Values,
		}
		dataList = append(dataList, data)
	}
	return dataList
}
