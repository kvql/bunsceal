package taxonomyCmd

import (
	"flag"
	"os"

	tx "github.com/kvql/bunsceal/pkg/taxonomy"
	"github.com/kvql/bunsceal/pkg/util"
	vis "github.com/kvql/bunsceal/pkg/visualise"
)

// Execute runs the taxonomy command with configured flags.
func Execute() {
	// Define command line flags
	localExport := flag.String("localExport", "", "Path for the taxonomy to be exported to a local JSON file")
	verify := flag.Bool("verify", false, "Validate the taxonomy")
	graph := flag.Bool("graph", false, "Generate diagrams to visualise the taxonomy")
	graphDir := flag.String("graphDir", ".tmp", "Directory for the graph visualisations")
	configPath := flag.String("config", "", "Path to config.yaml (default: <taxDir>/config.yaml)")

	// Parse command line flags
	flag.Parse()

	// Load configuration (with defaults if not found)
	cfg, err := tx.LoadConfig(*configPath, "schema/")
	if err != nil {
		util.Log.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Validate and Load the taxonomy and validate it
	// required for all actions
	tax, err := tx.LoadTaxonomy(cfg)

	if err != nil {
		util.Log.Println("Taxonomy content is not valid")
		os.Exit(1)
	}
	util.Log.Println("Taxonomy is valid")
	// Validate the taxonomy
	if *verify {
		if !vis.ValidateImageVersions() {
			util.Log.Println("Taxonomy images created before last update of the taxonomy, regenerate the images")
			os.Exit(1)
		}
		util.Log.Println("Images are up to date with the taxonomy")
	}
	// Generate local JSON file of the taxonomy
	if *localExport != "" {
		err := tx.GenLocalTaxonomy(tax, *localExport)
		if err != nil {
			util.Log.Println("Failed to export taxonomy to local JSON file")
			os.Exit(1)
		}
	}

	// Generate taxonomy visualisations
	if *graph {
		err := vis.RenderDiagrams(&tax, *graphDir, &cfg)
		if err != nil {
			util.Log.Print(err)
			os.Exit(1)
		}
	}
}
