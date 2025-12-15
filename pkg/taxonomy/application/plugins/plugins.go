package plugins

import (
	"fmt"

	"github.com/kvql/bunsceal/pkg/domain"
)

type PluginsCommonSettings struct {
	LabelInheritance bool `yaml:"label_inheritance"`
}

type ConfigPlugins struct {
	Classifications *ClassificationsConfig `yaml:"classifications"`
}

type PluginValidationResult struct {
	Valid  bool
	Errors []error
}

type Plugin interface {
	ValidateLabels(*domain.Seg) PluginValidationResult
	//ApplyInheritance(seg *domain.Seg) error
}

type Plugins struct {
	Plugins map[string]Plugin
}

func (p *Plugins) LoadPlugins(cfg ConfigPlugins) error {
	if cfg.Classifications != nil {
		p.Plugins["classifications"] = NewClassificationPlugin(cfg.Classifications, NsPrefix)
	}
	return nil
}

var NsPrefix = "bunsceal.plugin."

func ApplyPluginInheritance(parent domain.Seg, child *domain.Seg, p ClassificationsPlugin) error {
	if !p.Config.Common.LabelInheritance {
		return nil
	}
	// return error if child has labels which the parent doens't have
	if childLabels, cExists := child.LabelNamespaces[p.Namespace]; cExists {
		if _, pExists := parent.LabelNamespaces[p.Namespace]; !pExists {
			return fmt.Errorf("child %s has labels for %s, but parent %s does not", child.ID, p.Namespace, parent.ID)
		}
		for pk, pv := range parent.LabelNamespaces[p.Namespace] {
			// set if child doesn't have the label
			if _, exists := childLabels[pk]; !exists {
				child.LabelNamespaces[p.Namespace][pk] = pv
				child.ParsedLabels[p.Namespace+"/"+pk] = pv
			}
		}
	}
	return nil
}
