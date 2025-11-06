package taxonomy

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/util"
)

// GenLocalTaxonomy generates a local taxonomy file
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

func Version() string {
	prefix := "plat-sec-taxonomy"
	gitCommit := os.Getenv("GITHUB_SHA")
	if gitCommit == "" {
		if util.CheckGit() {
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
	util.Log.Println("Taxonomy version: ", file)
	return file
}
