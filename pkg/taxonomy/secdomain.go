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

type EnvDetails struct {
	DefSensitivity    string   `yaml:"def_sensitivity"`
	SensitivityReason string   `yaml:"sensitivity_reason"`
	DefCriticality    string   `yaml:"def_criticality"`
	CriticalityReason string   `yaml:"criticality_reason"`
	DefCompReqs       []string `yaml:"def_compliance_reqs"`
	CompReqs          map[string]CompReq
}
type version struct {
	Version string `yaml:"version"`
}

type SecDomain struct {
	Name        string                `yaml:"name"`
	ID          string                `yaml:"id"`
	Description string                `yaml:"description"`
	EnvDetails  map[string]EnvDetails `yaml:"env_details"`
}

const sdidPattern = "^[a-z0-9-]{1,15}$"

var rexSdId = regexp.MustCompile(sdidPattern)

func (sd SecDomain) Validate() (bool, []string) {
	var tests []string
	outcome := true
	if sd.Name == "" {
		tests = append(tests, "Name is empty")
		outcome = false
	}
	if sd.Description == "" {
		tests = append(tests, "Description is empty")
		outcome = false
	}
	if sd.EnvDetails == nil {
		tests = append(tests, "Environment details are empty")
		outcome = false
	}
	if !rexSdId.MatchString(sd.ID) {
		tests = append(tests, "ID does not meet requirement "+sdidPattern)
		outcome = false
	}
	// Loop through each environment the SD is defined in and validate details
	for envID, env := range sd.EnvDetails {
		// if Sensitivity or reason is provided, validate values
		if env.DefSensitivity != "" || env.SensitivityReason != "" {
			if _, ok := SensitivityLevels[env.DefSensitivity]; !ok {
				outcome = false
				tests = append(tests, fmt.Sprintf("Invalid Sensitivity level (%s) SecEnv (%s) defined in SD (%s)", env.DefSensitivity, envID, sd.ID))
			}
			if !descRex.MatchString(env.SensitivityReason) {
				outcome = false
				tests = append(tests, fmt.Sprintf("Sensitivity reason does not meet requirement %s", descPattern))
			}
		}
		// If Criticality or reason is provided, validate values
		if env.DefCriticality != "" || env.CriticalityReason != "" {
			if _, ok := CriticalityLevels[env.DefCriticality]; !ok {
				outcome = false
				tests = append(tests, fmt.Sprintf("Invalid Criticality level (%s) SecEnv (%s) defined in SD (%s)", env.DefSensitivity, envID, sd.ID))
			}
			if !descRex.MatchString(env.CriticalityReason) {
				outcome = false
				tests = append(tests, fmt.Sprintf("Criticality reason does not meet requirement %s", descPattern))
			}
		}
	}

	return outcome, tests
}

// LoadSDFiles loads all security domain files from the given directory
func LoadSDFiles(secDomainDir string) (map[string]SecDomain, error) {
	secDomains := make(map[string]SecDomain)
	err := filepath.WalkDir(secDomainDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// Load the file and parse it into a SecDomain struct
			secDomain, err := parseSDFile(path)
			if err != nil {
				return err
			}
			secDomains[secDomain.ID] = secDomain
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// validate id is unique across all domains. print error with non unique domains and exit if not
	validations := make([]string, 0)
	idMap := make(map[string]bool) //used to validate that the security domain id is unique
	outcome := true
	for _, sd := range secDomains {
		if _, ok := idMap[sd.ID]; ok {
			validations = append(validations, "ID for "+sd.Name+"is not unique: "+sd.ID)
			outcome = false
		} else {
			idMap[sd.ID] = true
		}
	}
	if !outcome {
		for _, result := range validations {
			util.Log.Println(result)
		}
		return nil, errors.New("Security domain validation failed, directory: " + secDomainDir)
	}

	return secDomains, nil
}

func parseSDFile(filePath string) (SecDomain, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return SecDomain{}, errors.New("Failed to parse file" + filePath + err.Error())
	}
	var fileVersion version
	err = yaml.Unmarshal(data, &fileVersion)
	if err != nil {
		return SecDomain{}, errors.New("Failed to parse file" + filePath + err.Error())
	}
	switch fileVersion.Version {
	case "1.0":
		// Unmarshal the YAML data into a SecDomain struct
		var secDomain SecDomain
		err = yaml.Unmarshal(data, &secDomain)
		if err != nil {

			return SecDomain{}, errors.New("Failed to parse file" + filePath + err.Error())
		}
		// Validate the Security domain file content
		pass, results := secDomain.Validate()
		if !pass {
			for line := range results {
				util.Log.Println(results[line])
			}
			return SecDomain{}, errors.New("Security domain validation failed, file: " + filePath)
		}
		return secDomain, nil
	default:
		return SecDomain{}, errors.New("Unsupported security domain file version: " + filePath + fileVersion.Version)
	}
}
