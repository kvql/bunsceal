package visualise

import (
	"fmt"
	"sort"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

type EnvImageData struct {
	Segs          map[string]domain.L1Overrides
	SegNames      map[string]string
	SortedSegs    []string
	Criticalities map[string]bool
	Sensitivities map[string]bool
}

// buildRowsMap creates a rowsMap from the config's L1Layout.
// Returns map[int][]string to maintain ordering when iterating sequentially.
// If no L1Layout is configured, defaults to all L1s on row 0.
// If an L1 ID exists in the taxonomy but is not in the config layout,
// it will be added to the last row to ensure all L1s are included in visualisations.
// Returns an error if the config references L1 IDs that don't exist in the taxonomy.
func buildRowsMap(cfg VisualsDef, txy domain.Taxonomy) (map[int][]string, error) {
	result := make(map[int][]string)
	seenL1s := make(map[string]bool)

	// If no L1Layout is configured, default to all L1s on a single row
	if len(cfg.L1Layout) == 0 {
		allL1s := []string{}
		for l1Id := range txy.SegL1s {
			allL1s = append(allL1s, l1Id)
		}
		return map[int][]string{0: allL1s}, nil
	}

	// Populate result map from config and track max row number
	maxRowNum := -1
	invalidL1s := []string{}
	for rowStr, l1Ids := range cfg.L1Layout {
		// Convert string key to int
		var rowNum int
		fmt.Sscanf(rowStr, "%d", &rowNum)

		validL1s := []string{}
		for _, l1Id := range l1Ids {
			// Validate that L1 ID exists in taxonomy
			if _, exists := txy.SegL1s[l1Id]; !exists {
				invalidL1s = append(invalidL1s, l1Id)
			} else {
				validL1s = append(validL1s, l1Id)
				seenL1s[l1Id] = true
			}
		}
		result[rowNum] = validL1s

		if rowNum > maxRowNum {
			maxRowNum = rowNum
		}
	}

	// Return error if config has invalid L1 IDs
	if len(invalidL1s) > 0 {
		return nil, fmt.Errorf("config references L1 IDs that don't exist in taxonomy: %v", invalidL1s)
	}

	// Find any L1s that are in the taxonomy but not in the layout
	missingL1s := []string{}
	for l1Id := range txy.SegL1s {
		if !seenL1s[l1Id] {
			missingL1s = append(missingL1s, l1Id)
		}
	}

	// Add missing L1s to the last row
	if len(missingL1s) > 0 {
		result[maxRowNum+1] = missingL1s
	}

	return result, nil
}

type L2GroupingData struct {
	SortedSegs         []string
	PresentGroupValues map[string]bool
}

func VisL2GroupingPrep(txy domain.Taxonomy, groupData *plugins.ImageGroupingData) map[string]L2GroupingData {

	data := make(map[string]L2GroupingData)
	for envId := range txy.SegL1s {
		data[envId] = L2GroupingData{
			PresentGroupValues: make(map[string]bool),
			SortedSegs:         make([]string, 0),
		}
	}

	for _, segL2 := range txy.SegsL2s {
		// REFACTORED: Iterate over L1Parents instead of L1Overrides keys
		for _, l1ID := range segL2.L1Parents {
			envData := data[l1ID]
			crit, err := segL2.GetNamespacedValue(l1ID, groupData.Namespace, groupData.Key)
			if err != nil {
				return make(map[string]L2GroupingData)
			}
			envData.PresentGroupValues[crit] = true
			envData.SortedSegs = append(data[l1ID].SortedSegs, segL2.ID)
			data[l1ID] = envData
		}
	}
	for _, envData := range data {
		sort.Strings(envData.SortedSegs)
	}
	return data
}

func PrepTaxonomy(txy domain.Taxonomy) map[string]EnvImageData {
	data := make(map[string]EnvImageData)
	for envId := range txy.SegL1s {
		data[envId] = EnvImageData{
			Segs:          make(map[string]domain.L1Overrides),
			SegNames:      make(map[string]string),
			Criticalities: make(map[string]bool),
			Sensitivities: make(map[string]bool),
			SortedSegs:    make([]string, 0),
		}
	}

	for _, segL2 := range txy.SegsL2s {
		// REFACTORED: Iterate over L1Parents instead of L1Overrides keys
		for _, l1ID := range segL2.L1Parents {
			// Lookup override (should exist after inheritance)
			det, exists := segL2.L1Overrides[l1ID]
			if !exists {
				// This should not happen if inheritance ran correctly
				o11y.Log.Printf("WARNING: Seg '%s' has parent '%s' but no override data after inheritance - skipping\n", segL2.ID, l1ID)
				continue
			}

			envData := data[l1ID]
			envData.Segs[segL2.ID] = det
			envData.SegNames[segL2.ID] = segL2.Name
			crit := GetClassificationFromOverride(det, segL2, "criticality")
			envData.Criticalities[crit] = true
			sens := GetClassificationFromOverride(det, segL2, "sensitivity")
			envData.Sensitivities[sens] = true
			envData.SortedSegs = append(data[l1ID].SortedSegs, segL2.ID)
			data[l1ID] = envData
		}
	}
	for _, envData := range data {
		sort.Strings(envData.SortedSegs)
	}
	return data
}
