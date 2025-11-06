package visualise

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/awalterschulze/gographviz"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/util"
)

type ImageConfig struct {
	graphFunc func() (*gographviz.Graph, error)
	filename  string
}

// RenderGraph generates a graph image from a gographviz.Graph object
func RenderGraph(g *gographviz.Graph, dir string, name string) error {
	output := g.String()
	// function argument dir is only for where to save the image
	tmpDir := ".tmp/"
	// validate paths inputs and set defaults
	if dir == "" {
		dir = tmpDir
	}
	if dir[len(dir)-1:] != "/" {
		dir = dir + "/"
	}

	// Making name mandatory
	if name == "" {
		return errors.New("no name provided for the diagram")
	}
	graphFile := name + ".graph"

	// Create .tmp directory for temporary graph files
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		// Create .tmp/ directory if it doesn't exist
		err := os.Mkdir(tmpDir, 0750)
		if err != nil {
			return err
		}
	}

	// Delete the file if it already exists
	if _, err := os.Stat(tmpDir + graphFile); err == nil {
		err := os.Remove(tmpDir + graphFile)
		if err != nil {
			return err
		}
	}

	// #nosec G304 -- tmpDir is a controlled constant, graphFile is sanitised
	tmp, err := os.Create(tmpDir + graphFile)
	if err != nil {
		return err
	}
	defer tmp.Close()

	// Write the graph to a temporary file
	tmp.WriteString(output)

	// Sanitise output path to prevent path traversal
	outputPath := filepath.Clean(filepath.Join(dir, name))

	// #nosec G204 -- Using exec.Command with separate args (not shell), paths are sanitised
	cmd := exec.Command("dot", "-Tpng", tmp.Name(), "-o", outputPath)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		util.Log.Println("Failed to generate image:", err)
		util.Log.Println("Standard output:", string(cmdOutput))
		return err
	}

	util.Log.Println("Generated image at:", dir+name)
	return nil
}

// RenderDiagrams generates all the diagrams for the taxonomy
func RenderDiagrams(tax *domain.Taxonomy, dir string, cfg *domain.Config) error {
	// Generate the security domain graph
	graphConfigs := []ImageConfig{
		{func() (*gographviz.Graph, error) { return GraphL2(tax, cfg, false, false) }, "l2_Segments_overview.png"},
		{func() (*gographviz.Graph, error) { return GraphL1(tax, cfg) }, "l1_segments_overview.png"},
		{func() (*gographviz.Graph, error) { return GraphL2(tax, cfg, true, false) }, "criticality_overview_all.png"},
		{func() (*gographviz.Graph, error) { return GraphL2(tax, cfg, false, true) }, "sensitivity_overview_all.png"},
		{func() (*gographviz.Graph, error) { return GraphCompliance(tax, cfg, "pci-dss", true) }, "compliance_overview_pci.png"},
	}

	for _, config := range graphConfigs {
		g, err := config.graphFunc()
		if err != nil {
			return err
		}
		err = RenderGraph(g, dir, config.filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateImageVersions() bool {
	images := []string{
		"./docs/images/l2_Segments_overview.png",
		"./docs/images/l1_Segments_overview.png",
		"./docs/images/criticality_overview_all.png",
		"./docs/images/sensitivity_overview_all.png",
		"./docs/images/compliance_overview_pci.png",
	}
	for _, img := range images {

		if ok, err := ValidateImageVersion("./taxonomy", img); !ok && err == nil {
			return false
		} else if err != nil {
			util.Log.Println("error validating image version:", err)
			return false
		}
	}
	return true
}

// ValidateImageVersions checks if the image is up to date with the taxonomy based on latest commit times
func ValidateImageVersion(txDir string, imagePath string) (bool, error) {
	if !util.CheckGit() {
		util.Log.Println("Git binary not found")
		util.Log.Println("PATH environment variable:", os.Getenv("PATH"))
		p, _ := os.Getwd()
		util.Log.Println("Execution directory path:", p)
		return false, errors.New("git binary not found")
	}
	txTime, err := util.GetLatestCommitTime(txDir)
	if err != nil {
		tmp := fmt.Sprintf("Error getting latest commit time: %s", err)
		return false, errors.New(tmp)
	}
	imgTime, err := util.GetLatestCommitTime(imagePath)
	if err != nil {
		tmp := fmt.Sprintf("Error getting latest commit time: %s", err)
		return false, errors.New(tmp)
	}
	if txTime.After(imgTime) {
		util.Log.Println("Image is out of date")
		return false, nil
	}
	return true, nil
}
