package infrastructure

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
)

// GenLocalTaxonomy generates a local taxonomy file
// TODO: Update JSON output to use ParsedLabels map instead of Labels array
// See: Labels field should have json:"-" and ParsedLabels should have json:"labels,omitempty"
func GenLocalTaxonomy(tx domain.Taxonomy, dir string) error {
	// Check if provided directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		err := os.Mkdir(dir, 0750)
		if err != nil {
			return err
		}
	}

	// output taxonomy as a json file
	data, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		return err
	}
	filePath := dir + "/" + Version()
	return os.WriteFile(filePath, data, 0600)
}

// CheckGit checks if git binary is available
func CheckGit() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func Version() string {
	prefix := "bunsceal-taxonomy"
	gitCommit := os.Getenv("GITHUB_SHA")
	if gitCommit == "" {
		if CheckGit() {
			cmd := exec.Command("git", "rev-parse", "HEAD")
			output, err := cmd.Output()
			if err == nil {
				gitCommit = strings.TrimSpace(string(output))
			} else {
				gitCommit = "unknown"
			}
		} else {
			gitCommit = "unknown"
		}
	}
	file := prefix + "-" + gitCommit[:7] + ".json"
	o11y.Log.Println("Taxonomy version: ", file)
	return file
}
