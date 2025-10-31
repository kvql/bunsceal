package visualise

import (
	"fmt"
	"sort"

	tx "github.com/kvql/bunsceal/pkg/taxonomy"
)

type EnvImageData struct {
	SecDoms       map[string]tx.EnvDetails
	SecDomNames   map[string]string
	SortedSecDoms []string
	Criticalities map[string]bool
}

// TODO: validate all envs are present in the rows
func validateRows(txy *tx.Taxonomy, rowsMap map[int][]string) error {
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
	if len(txy.SecEnvironments) != rowTotal {
		return fmt.Errorf("Number of environments in taxonomy does not match number of environments in rows")
	}
	for env, _ := range txy.SecEnvironments {
		if _, ok := rowMap[env]; !ok {
			return fmt.Errorf("environment %s not found in rows", env)
		}
	}
	return nil
}

func PrepTaxonomy(txy *tx.Taxonomy) map[string]EnvImageData {
	data := make(map[string]EnvImageData)
	for envId, _ := range txy.SecEnvironments {
		data[envId] = EnvImageData{
			SecDoms:       make(map[string]tx.EnvDetails),
			SecDomNames:   make(map[string]string),
			Criticalities: make(map[string]bool),
			SortedSecDoms: make([]string, 0),
		}
	}

	for _, sd := range txy.SecDomains {
		for envId, det := range sd.EnvDetails {
			envData := data[envId]
			envData.SecDoms[sd.ID] = det
			envData.SecDomNames[sd.ID] = sd.Name
			envData.Criticalities[det.DefCriticality] = true
			envData.SortedSecDoms = append(data[envId].SortedSecDoms, sd.ID)
			data[envId] = envData
		}
	}
	for _, envData := range data {
		sort.Strings(envData.SortedSecDoms)
	}
	return data
}
