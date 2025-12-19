package visualise

import (
	"fmt"
	"strings"
	"time"

	"github.com/awalterschulze/gographviz"
	"github.com/kvql/bunsceal/pkg/domain"
)

type ColorFont struct {
	Color string
	Font  string
}

// #######################
// Global Variables for configuring the graph
// #######################

type PrimaryColours map[string]string

var primaryColours = PrimaryColours{
	"1": "#58CDCA",
	"2": "#8F58CD",
	"3": "#CD585B",
	"4": "#96CD58",
}

var BgColour = "\"#1E6566\""
var FontColour = "\"#BFECEC\""
var DebugColour = "\"#FF1A00\""
var TsnColour = "\"#1E6580\""

func (pc PrimaryColours) GetColour(k string) string {
	return "\"" + pc[k] + "\""
}

var SenseColourMap = map[string]ColorFont{
	"A": {primaryColours.GetColour("4"), primaryColours.GetColour("2")},
	"B": {primaryColours.GetColour("3"), primaryColours.GetColour("1")},
	"C": {primaryColours.GetColour("2"), primaryColours.GetColour("4")},
	"D": {primaryColours.GetColour("1"), primaryColours.GetColour("3")},
}
var CritColourMap = map[string]ColorFont{
	"1": {primaryColours.GetColour("4"), primaryColours.GetColour("2")},
	"2": {primaryColours.GetColour("3"), primaryColours.GetColour("1")},
	"3": {primaryColours.GetColour("2"), primaryColours.GetColour("4")},
	"4": {primaryColours.GetColour("1"), primaryColours.GetColour("3")},
	"5": {"\"#A0E1E1\"", "\"#320707\""},
}

// Classification order and labels (hardcoded for now, to be migrated to plugin config)
var SenseOrder = []string{"A", "B", "C", "D"}
var CritOrder = []string{"1", "2", "3", "4", "5"}
var SensitivityLabels = map[string]string{
	"A": "High",
	"B": "Medium",
	"C": "Low",
	"D": "N/A",
}
var CriticalityLabels = map[string]string{
	"1": "Critical",
	"2": "High",
	"3": "Medium",
	"4": "Low",
	"5": "N/A",
}

var visibility = "\"invis\"" //"\"\"" for visible, "\"invis\"" for invisible
// Formatting attributes for invisible nodes and edges
var InvisAtt = map[string]string{
	"style":     visibility,
	"label":     "\"\"",
	"color":     DebugColour,
	"fontcolor": DebugColour,
}

// #################################
// Formatting variables
// #################################

// default formatting for graph nodes
var NodeFormat = map[string]string{
	"shape":     "\"box\"",
	"color":     primaryColours.GetColour("1"),
	"fontcolor": primaryColours.GetColour("2"),
	"fontname":  "\"Arial Bold\"",
	"fontsize":  "\"14\"",
	"width":     "\"2.5\"",
	"style":     "\"rounded,filled,setlinewidth(0)\"",
}

// default formatting for graph nodes
var GraphFormat = map[string]string{
	"color":     primaryColours.GetColour("4"),
	"fontcolor": primaryColours.GetColour("3"),
	"fontname":  "\"Arial Bold\"",
	"fontsize":  "\"18\"",
	"width":     "\"2.5\"",
	"style":     "\"rounded,setlinewidth(2)\"",
}

func CopyInvis() map[string]string {
	inv := make(map[string]string)
	for k, v := range InvisAtt {
		inv[k] = v
	}
	return inv
}

var Directed = true
var BatchSize = 10

// #######################
// Add Legend to the Graph
// #######################

var LegendGraphAtt = map[string]string{
	"shape":     "\"box\"",
	"color":     primaryColours.GetColour("2"),
	"fontcolor": primaryColours.GetColour("4"),
	"fontname":  "\"Arial Bold\"",
	"fontsize":  "\"12\"",
	"width":     "\"\"",
	"style":     "\"rounded,setlinewidth(1)\"",
	"nodesep":   "\"2\"",
	"label":     "\"\nLegend:\nClassification: Sensitivity+Criticality\nColours: Based on Sensitivity\"",
}

// AddLegend adds a legend to the graph. Set stack to true to stack the legend nodes vertically
func AddLegend(g *gographviz.Graph, font int, stack bool) error {
	legSGName := "\"cluster_a_legend\""
	LegendGraphAtt["fontsize"] = fmt.Sprintf("\"%d\"", font-2)
	g.AddSubGraph("top_level_graph", legSGName, LegendGraphAtt)
	nodes := make([]string, 0)
	for _, s := range SenseOrder {
		label := fmt.Sprintf("\"Sensitivity: %s (%s)\"", s, SensitivityLabels[s])
		nodeAtt := FormatNode(label, s)
		nodeAtt["fontsize"] = fmt.Sprintf("\"%d\"", font-2)
		nodeAtt["width"] = "\"\""
		nodeName := fmt.Sprintf("\"legend_%s\"", s)
		nodes = append(nodes, nodeName)
		err := g.AddNode(legSGName, nodeName, nodeAtt)
		if err != nil {
			return err
		}
	}
	if !stack {
		err := AddEdges(g, nodes)
		if err != nil {
			return err
		}
	}
	return nil
}

// #######################
// Functions for structuring the graph
// #######################
// Attempts to abstract some of the graphviz complexity away when defining the different graphs
// A lot of hidden features are added the graph to control the graphing algorithm. To see these and understand them set the visibility variable to "" and run the code

// BaselineGraph creates a new graph with default settings
func BaselineGraph(title string, subHeading string) *gographviz.Graph {

	title = "\"" + title + "\\n" + strings.Repeat("_", len(title)) + "\n" + subHeading + "\""
	// Setup the top level graph object
	g := gographviz.NewGraph()
	if err := g.SetName("top_level_graph"); err != nil {
		panic(err)
	}
	if err := g.SetDir(Directed); err != nil {
		panic(err)
	}

	g.AddAttr("top_level_graph", "label", title)
	g.AddAttr("top_level_graph", "rankdir", "\"LR\"") // Left to right graph
	g.AddAttr("top_level_graph", "splines", "\"line\"")
	g.AddAttr("top_level_graph", "center", "\"true\"")
	g.AddAttr("top_level_graph", "bgcolor", BgColour)
	g.AddAttr("top_level_graph", "color", "\"white\"")
	g.AddAttr("top_level_graph", "fontcolor", FontColour)
	g.AddAttr("top_level_graph", "fontsize", "\"24\"")
	g.AddAttr("top_level_graph", "fontname", "\"Arial Bold\"")
	g.AddAttr("top_level_graph", "nodesep", "\"0.1\"") // Increase space between nodes
	g.AddAttr("top_level_graph", "labelloc", "\"t\"")  // Moves title to top of graph

	// Adding mostly invisible timestamp to the graph. This ensures that every graph has a unique hash and gets committed after running the code. Without this the images wouldn't be
	// updated for every taxonomy change and therefore the CI validation would fail.
	tsnFormat := map[string]string{
		"color":     TsnColour,
		"label":     fmt.Sprintf("\"%s\"", time.Now().Format("2006-01-02 15:04:05")),
		"fontcolor": TsnColour,
		"fontsize":  "\"5\"",
	}
	// making visible for debugging
	if visibility == "\"\"" {
		tsnFormat["fontsize"] = "\"12\""
		tsnFormat["fontcolor"] = DebugColour
	}
	g.AddNode("top_level_graph", "time", tsnFormat)
	return g
}

func AddEdges(g *gographviz.Graph, edges []string) error {
	edgeAtt := CopyInvis()
	edgeAtt["minlen"] = "\"1\"" // int value

	for n := 0; n < len(edges)-1; n++ {
		err := g.AddEdge(edges[n], edges[n+1], Directed, edgeAtt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Graph and Node naming functions
// -------------------------------
// These functions are used to generate the names of the subgraphs and nodes in the graph.
// Having functions for this means less copy and pasting of the fmt.Sprintf() format and removes the need to
// track the names in a variable.

// envSubGraphName returns the name of the environment subgraph
func envSubGraphName(envID string) string {
	return fmt.Sprintf("\"cluster_%s\"", strings.ReplaceAll(envID, "-", "_"))
}

// FocusSGName returns the name of the subgraph for grouping Segment Level 2s within an environment by a particular focus. e.g criticality, compliance scope.
// cluster name doesn't include what the focus is as it just needs to be predicable and unique not human readable. graphviz displays the label not the graph name
func focusSGName(envID string, focusValue string) string {
	return fmt.Sprintf("\"cluster_focus_%s_%s\"", strings.ReplaceAll(envID, "-", "_"), focusValue)
}

// Struct and functions for managing batch subgraphs
// -------------------------------------------------
// -------------------------------------------------
type BatchVars struct {
	Limit        int
	Count        map[string]int // var to keep track of the number of nodes in each batch
	Current      map[string]int // var to keep track of the current batch
	EnvFormatted string
	Focus        string
}

func NewBatchVars(env string) *BatchVars {
	return &BatchVars{
		Limit:        BatchSize,
		Count:        map[string]int{},
		Current:      map[string]int{},
		EnvFormatted: strings.ReplaceAll(env, "-", "_"),
	}
}
func (b *BatchVars) BumpBatch(key string, g *gographviz.Graph) (string, error) {
	b.Count[key] = 0
	b.Current[key]++
	batchAtt := CopyInvis()
	batchAtt["label"] = fmt.Sprintf("\"Batch %d\"", b.Current[key])
	// Add subgraph for batch to the criticality subgraph
	err := g.AddSubGraph(focusSGName(b.EnvFormatted, key), b.CurrentSGName(key), batchAtt)
	if err != nil {
		return "", err
	}
	batchNodeName := fmt.Sprintf("\"batch_nodes_%s_%s_%d\"", b.EnvFormatted, key, b.Current[key])
	batchNodeAtt := CopyInvis()
	batchNodeAtt["label"] = batchNodeName
	err = g.AddNode(b.CurrentSGName(key), batchNodeName, batchNodeAtt)
	if err != nil {
		return "", err
	}
	return batchNodeName, nil
}

// SGName returns the name of the subgraph for a batch of nodes
func (b *BatchVars) CurrentSGName(key string) string {
	return fmt.Sprintf("\"cluster_batch_%s_%s_%d\"", b.EnvFormatted, key, b.Current[key])
}

// Spacer function
// -------------------------------------------------
// -------------------------------------------------
// function to add spacing nodes

func AddSpacerNodes(g *gographviz.Graph, row int, spacer int) string {
	spacerAtt := CopyInvis()
	nodeName := fmt.Sprintf("\"spacer_node_%d_%d\"", row, spacer)
	spacerAtt["label"] = nodeName
	g.AddNode(fmt.Sprintf("\"cluster_row_%d\"", row), nodeName, spacerAtt)
	return fmt.Sprintf("\"spacer_node_%d_%d\"", row, spacer)
}

func AddSpacers(g *gographviz.Graph, rowNodes map[int][]string) error {
	largestRow := 0
	for _, v := range rowNodes {
		if len(v) > largestRow {
			largestRow = len(v)
		}
	}

	// if largestRow % 2 != 0 {
	// 	largestRow++
	// }
	// add spacer nodes either side of row nodes so each row is the same length
	for row := 0; row < len(rowNodes); row++ {
		nodes := rowNodes[row]
		if len(nodes) < largestRow {
			spacers := int(float64(largestRow - len(nodes)))
			if spacers%2 != 0 {
				spacers++
			}
			// add spacer nodes to the start of the row
			for i := 0; i < spacers/2; i++ {
				nodes = append([]string{AddSpacerNodes(g, row, i)}, nodes...)
			}
			// add spacer nodes to the end of the row
			for i := spacers / 2; i < spacers; i++ {
				nodes = append(nodes, AddSpacerNodes(g, row, i))
			}
			rowNodes[row] = nodes
		}
		err := AddEdges(g, rowNodes[row])
		if err != nil {
			return err
		}
	}
	return nil
}

// #################################
// Formatting functions
// #################################

// FormatNode returns a map of attributes for a node in graphviz format
func FormatNode(label string, colourLookup string) map[string]string {
	node := make(map[string]string)
	// make a copy of the default node format
	for k, v := range NodeFormat {
		node[k] = v
	}
	node["label"] = label
	if _, ok := SenseColourMap[colourLookup]; ok {
		node["color"] = SenseColourMap[colourLookup].Color
		node["fontcolor"] = SenseColourMap[colourLookup].Font
	} else if _, ok := CritColourMap[colourLookup]; ok {
		node["color"] = CritColourMap[colourLookup].Color
		node["fontcolor"] = CritColourMap[colourLookup].Font
	}
	return node
}

// GetClassificationValue reads classification from plugin labels
func GetClassificationValue(seg domain.Seg, classType string) string {
	ns := "bunsceal.plugin.classifications"
	if val, exists := seg.LabelNamespaces[ns][classType]; exists {
		return val
	}
	return ""
}

// GetClassificationFromOverride reads from L1Override labels with fallback to child segment labels.
// Note: L2 segments can have different classifications per parent. This correctly
// reads from the specific parent's override first, then falls back to child labels if override is empty.
// This handles the case where inheritance populates child labels but not overrides.
func GetClassificationFromOverride(override domain.L1Overrides, child domain.Seg, classType string) string {
	ns := "bunsceal.plugin.classifications"

	// Check override first
	if val, exists := override.LabelNamespaces[ns][classType]; exists && val != "" {
		return val
	}

	// Fall back to child segment labels
	if val, exists := child.LabelNamespaces[ns][classType]; exists {
		return val
	}

	return ""
}

// FormatNode returns a map of attributes for a node in graphviz format
func FormatGraph(label string, colourLookup string) map[string]string {
	graph := make(map[string]string)
	for k, v := range GraphFormat {
		graph[k] = v
	}
	graph["label"] = label
	if _, ok := SenseColourMap[colourLookup]; ok {
		graph["color"] = SenseColourMap[colourLookup].Color
		graph["fontcolor"] = SenseColourMap[colourLookup].Font
	} else if _, ok := CritColourMap[colourLookup]; ok {
		graph["color"] = CritColourMap[colourLookup].Color
		graph["fontcolor"] = CritColourMap[colourLookup].Font
	} else {
		graph["color"] = "\"#A0E1E1\""
		graph["fontcolor"] = "\"#A0E1E1\""
	}

	return graph
}

// FormatLabel returns a formatted label for a node with full details
func FormatLabel(name string, sense string, crit string) string {

	// label := fmt.Sprintf("\"%s\\n%s\\nSensitivity: %s\\nCriticality: %s\\nCompliance Reqs: \\n%s\"",
	// 	name, strings.Repeat("_", len(name)),
	// 	strings.ToUpper(sense),
	// 	strings.ToUpper(crit),
	// 	strings.Join(compReqs, ", ")
	// )

	// graph without Compliance reqs until process around them is finalised
	label := fmt.Sprintf("\"%s\\n%s\\nClassification: %s%s\"",
		name, strings.Repeat("_", len(name)),
		strings.ToUpper(sense),
		strings.ToUpper(crit),
	)
	return label
}

// FormatLabel returns a formatted label for a node with full detailsname string, sense string, crit string) string {
func FormatEnvLabel(txy domain.Taxonomy, prefix string, envID string, showClass bool) string {

	// graph without Compliance reqs until process around them is finalised
	label := prefix + txy.SegL1s[envID].Name

	if showClass {
		label = fmt.Sprintf("\"%s\\n%s\\nClassification: %s%s\"",
			label, strings.Repeat("_", len(label)),
			strings.ToUpper(GetClassificationValue(txy.SegL1s[envID], "sensitivity")),
			strings.ToUpper(GetClassificationValue(txy.SegL1s[envID], "criticality")),
		)
	} else {
		label = fmt.Sprintf("\"%s \n%s\n(ID:%s)\"", label,
			strings.Repeat("_", len(label)),
			txy.SegL1s[envID].ID)
	}
	return label
}

// FormatLabel returns a formatted label for a node with full detailsname string, sense string, crit string) string {
func FormatSdLabel(txy domain.Taxonomy, prefix string, envID string, sdID string, showClass bool, emphasis int) string {

	// graph without Compliance reqs until process around them is finalised
	label := fmt.Sprintf("%s%s", prefix, txy.SegsL2s[sdID].Name)

	if showClass {
		label = fmt.Sprintf("\"%s%s\\n%s\\nClassification: %s%s%s\"",
			strings.Repeat("\n", emphasis),
			label, strings.Repeat("_", len(label)),
			strings.ToUpper(GetClassificationFromOverride(txy.SegsL2s[sdID].L1Overrides[envID], txy.SegsL2s[sdID], "sensitivity")),
			strings.ToUpper(GetClassificationFromOverride(txy.SegsL2s[sdID].L1Overrides[envID], txy.SegsL2s[sdID], "criticality")),
			strings.Repeat("\n", emphasis),
		)
	} else {
		label = fmt.Sprintf("\"%s%s (%s)%s\"", strings.Repeat("\n", emphasis), label, txy.SegsL2s[sdID].ID, strings.Repeat("\n", emphasis))
	}
	return label
}
