package visualise

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/awalterschulze/gographviz"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

type ImageConfig struct {
	graphFunc func() (*gographviz.Graph, error)
	filename  string
}

// renderGraph generates a graph image from a gographviz.Graph object
func renderGraph(g *gographviz.Graph, dir string, name string) error {
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
		o11y.Log.Println("Failed to generate image:", err)
		o11y.Log.Println("Standard output:", string(cmdOutput))
		return err
	}

	o11y.Log.Println("Generated image at:", dir+name)
	return nil
}

// RenderDiagrams generates all the diagrams for the taxonomy
func RenderDiagrams(tax domain.Taxonomy, dir string, terms domain.TermConfig, visCfg VisualsDef, pluginMap plugins.Plugins) error {
	// Collect image data from all plugins
	var groupData []plugins.ImageGroupingData
	for _, plugin := range pluginMap {
		if plugin != nil {
			groupData = append(groupData, plugin.GetImageData()...)
		}
	}

	graphConfigs := []ImageConfig{
		{func() (*gographviz.Graph, error) {
			return GraphL2Grouped(tax, terms, visCfg, groupData, plugins.ImageGroupingData{})
		}, "l2_Segments_overview.png"},
		{func() (*gographviz.Graph, error) { return GraphL1(tax, terms, visCfg, groupData) }, "l1_segments_overview.png"},
	}

	for _, group := range groupData {
		graphConfigs = append(graphConfigs, ImageConfig{func() (*gographviz.Graph, error) { return GraphL2Grouped(tax, terms, visCfg, groupData, group) }, group.Key + "_overview.png"})
	}

	for _, config := range graphConfigs {
		g, err := config.graphFunc()
		if err != nil {
			o11y.Log.Printf("error generating graph for: %s", config.filename)
			return err
		}
		err = renderGraph(g, dir, config.filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateImageVersions(pluginMap plugins.Plugins) bool {

	var groupData []plugins.ImageGroupingData
	for _, plugin := range pluginMap {
		if plugin != nil {
			groupData = append(groupData, plugin.GetImageData()...)
		}
	}

	images := []string{
		"./docs/images/l2_Segments_overview.png",
		"./docs/images/l1_Segments_overview.png",
	}
	for _, group := range groupData {
		images = append(images, ".docs/images/"+group.Key+"_overview.png")
	}
	for _, img := range images {

		if ok, err := ValidateImageVersion("./taxonomy", img); !ok && err == nil {
			return false
		} else if err != nil {
			o11y.Log.Println("error validating image version:", err)
			return false
		}
	}
	return true
}

// ValidateImageVersions checks if the image is up to date with the taxonomy based on latest commit times
func ValidateImageVersion(txDir string, imagePath string) (bool, error) {
	if !CheckGit() {
		o11y.Log.Println("Git binary not found")
		o11y.Log.Println("PATH environment variable:", os.Getenv("PATH"))
		p, _ := os.Getwd()
		o11y.Log.Println("Execution directory path:", p)
		return false, errors.New("git binary not found")
	}
	if hasHistory, err := HasGitHistory(); !hasHistory {
		if err != nil {
			o11y.Log.Println("Error checking git history:", err)
			return false, errors.New("error checking git history")
		}
		o11y.Log.Println("Repository is a shallow clone, commit times may be inaccurate")
		return false, errors.New("repository is a shallow clone, full git history required")
	}
	txTime, err := GetLatestCommitTime(txDir)
	if err != nil {
		tmp := fmt.Sprintf("Error getting latest commit time: %s", err)
		return false, errors.New(tmp)
	}
	imgTime, err := GetLatestCommitTime(imagePath)
	if err != nil {
		tmp := fmt.Sprintf("Error getting latest commit time: %s", err)
		return false, errors.New(tmp)
	}
	if txTime.After(imgTime) {
		o11y.Log.Println("Image is out of date")
		return false, nil
	}
	o11y.Log.Printf("directory time: %s", txTime)
	o11y.Log.Printf("image time: %s", imgTime)
	return true, nil
}
