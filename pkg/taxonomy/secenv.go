package taxonomy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

type SecEnv struct {
	Name              string   `yaml:"name"`
	ID                string   `yaml:"id"`
	Description       string   `yaml:"description"`
	DefSensitivity    string   `yaml:"def_sensitivity"`
	SensitivityReason string   `yaml:"sensitivity_reason"`
	DefCriticality    string   `yaml:"def_criticality"`
	CriticalityReason string   `yaml:"criticality_reason"`
	DefCompReqs       []string `yaml:"def_compliance_reqs"`
}

var envPattern = "^[a-z-]{1,15}$"

var envRex = regexp.MustCompile(envPattern)

// Validate the security environment struct to ensure required fields are present and meet any expected formats
func (env SecEnv) Validate() (bool, []string) {
	var tests []string
	outcome := true
	if env.Name == "" {
		tests = append(tests, "Name is empty")
		outcome = false
	}
	if !envRex.MatchString(env.ID) {
		tests = append(tests, fmt.Sprintf("ID (%s) does not match (%s)", env.ID, envPattern))
		outcome = false
	}
	if !descRex.MatchString(env.Description) {
		tests = append(tests, "Description is empty")
		outcome = false
	}
	if _, ok := SensitivityLevels[env.DefSensitivity]; !ok {
		tests = append(tests, fmt.Sprintf("Default sensitivity level (%s) is not a valid level", env.DefSensitivity))
		outcome = false
	}
	if !descRex.MatchString(env.SensitivityReason) {
		tests = append(tests, fmt.Sprintf("Sensitivity reason does not meet requirement %s", descPattern))
		outcome = false
	}
	if _, ok := CriticalityLevels[env.DefCriticality]; !ok {
		tests = append(tests, fmt.Sprintf("Default criticality level (%s) is not a valid level", env.DefCriticality))
		outcome = false
	}
	if !descRex.MatchString(env.CriticalityReason) {
		tests = append(tests, fmt.Sprintf("Criticality reason does not meet requirement %s", descPattern))
		outcome = false
	}
	return outcome, tests
}

// LoadSecEnvFiles Parse all security environment files from the provided directory,
// validate and return a map of SecEnv structs
func LoadSecEnvFiles(secEnvDir string) (map[string]SecEnv, error) {
	files, err := os.ReadDir(secEnvDir)
	if err != nil {
		return nil, err
	}
	valid := true
	secEnvs := make(map[string]SecEnv)
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(secEnvDir, file.Name())
			// Load the file and parse it into a SecEnv struct
			secEnv, err := parseSecEnvFile(filePath)
			if err != nil {
				util.Log.Printf("Error parsing file: %s\n", filePath)
				valid = false
			}
			secEnvs[secEnv.ID] = secEnv
		}
	}
	if valid {
		return secEnvs, nil
	} else {
		return nil, errors.New("loading security environments failed")
	}
}

func parseSecEnvFile(filePath string) (SecEnv, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return SecEnv{}, err
	}

	// Unmarshal the YAML data into a SecEnv struct
	var secEnv SecEnv
	err = yaml.Unmarshal(data, &secEnv)
	if err != nil {
		return SecEnv{}, err
	}

	pass, results := secEnv.Validate()
	if !pass {
		for _, result := range results {
			util.Log.Println(result)
		}
		return SecEnv{}, errors.New("Security environment validation failed, file: " + filePath)
	}

	return secEnv, nil
}
