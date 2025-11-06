package visualise

import (
	"fmt"
	"sort"

	"github.com/kvql/bunsceal/pkg/domain"
)

type EnvImageData struct {
	SegL2s        map[string]domain.L1Overrides
	SegL2Names    map[string]string
	SortedSegL2s  []string
	Criticalities map[string]bool
}

// buildRowsMap creates a rowsMap from the config's L1Layout.
// Returns map[int][]string to maintain ordering when iterating sequentially.
// If no L1Layout is configured, defaults to all L1s on row 0.
// If an L1 ID exists in the taxonomy but is not in the config layout,
// it will be added to the last row to ensure all L1s are included in visualisations.
// Returns an error if the config references L1 IDs that don't exist in the taxonomy.
func buildRowsMap(cfg *domain.Config, txy *domain.Taxonomy) (map[int][]string, error) {
	result := make(map[int][]string)
	seenL1s := make(map[string]bool)

	// If no L1Layout is configured, default to all L1s on a single row
	if len(cfg.Visuals.L1Layout) == 0 {
		allL1s := []string{}
		for l1Id := range txy.SegL1s {
			allL1s = append(allL1s, l1Id)
		}
		return map[int][]string{0: allL1s}, nil
	}

	// Populate result map from config and track max row number
	maxRowNum := -1
	invalidL1s := []string{}
	for rowStr, l1Ids := range cfg.Visuals.L1Layout {
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

func PrepTaxonomy(txy *domain.Taxonomy) map[string]EnvImageData {
	data := make(map[string]EnvImageData)
	for envId := range txy.SegL1s {
		data[envId] = EnvImageData{
			SegL2s:        make(map[string]domain.L1Overrides),
			SegL2Names:    make(map[string]string),
			Criticalities: make(map[string]bool),
			SortedSegL2s:  make([]string, 0),
		}
	}

	for _, sd := range txy.SegL2s {
		for envId, det := range sd.L1Overrides {
			envData := data[envId]
			envData.SegL2s[sd.ID] = det
			envData.SegL2Names[sd.ID] = sd.Name
			envData.Criticalities[det.Criticality] = true
			envData.SortedSegL2s = append(data[envId].SortedSegL2s, sd.ID)
			data[envId] = envData
		}
	}
	for _, envData := range data {
		sort.Strings(envData.SortedSegL2s)
	}
	return data
}
