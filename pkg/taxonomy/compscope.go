package taxonomy

import (
	"os"

	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

type CompReq struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	ReqsLink    string `yaml:"requirements_link"`
}

func LoadCompScope(filePath string) (map[string]CompReq, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		util.Log.Println("Error reading file:", err)
		return nil, err
	}
	// Parse the file into a allCompScopes struct
	var compReqs map[string]CompReq
	err = yaml.Unmarshal(data, &compReqs)
	if err != nil {
		util.Log.Println("Error parsing file:", err)
		return nil, err
	}

	return compReqs, nil
}
