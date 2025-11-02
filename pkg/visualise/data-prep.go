package visualise

import (
	"fmt"
	"sort"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
)

type EnvImageData struct {
	SegL2s        map[string]domain.L1Overrides
	SegL2Names    map[string]string
	SortedSegL2s  []string
	Criticalities map[string]bool
}

// TODO: validate all envs are present in the rows
func validateRows(txy *domain.Taxonomy, rowsMap map[int][]string) error {
	rows := append(rowsMap[0], rowsMap[1]...)
	rows = append(rows, rowsMap[2]...)
	rowMap := make(map[string]bool)
	for _, r := range rows {
		rowMap[r] = true
	}
	rowTotal := 0
	for _, r := range rowsMap {
		rowTotal += len(r)
	}
	if len(txy.SegL1s) != rowTotal {
		return fmt.Errorf("Number of environments in taxonomy does not match number of environments in rows")
	}
	for env := range txy.SegL1s {
		if _, ok := rowMap[env]; !ok {
			return fmt.Errorf("environment %s not found in rows", env)
		}
	}
	return nil
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
