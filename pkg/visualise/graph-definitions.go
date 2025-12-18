package visualise

import (
	"fmt"
	"strings"

	"github.com/awalterschulze/gographviz"
	configdomain "github.com/kvql/bunsceal/pkg/config/domain"
	"github.com/kvql/bunsceal/pkg/domain"
)

// ################################
// Function to Segment Level 1 Graph
// ################################

// See the GraphSDs function for more comments explaining how the graph is generated. That is the more complex function and therefore has more comments than GraphEnvs

func GraphL1(txy *domain.Taxonomy, cfg *configdomain.Config) (*gographviz.Graph, error) {
	// Setup the top level graph object
	title := cfg.Terminology.L1.Plural + " Overview"
	g := BaselineGraph(title, "")

	// // Add legend to the graph
	// // ------------------------
	err := AddLegend(g, 12, true)
	if err != nil {
		return nil, err
	}

	// Build rowsMap from config
	rowsLayout, err := buildRowsMap(cfg, txy)
	if err != nil {
		return nil, err
	}

	rowNodes := map[int][]string{}
	// sub graph structure:
	// - top_level_graph
	// 	- cluster_row
	// 		- env Nodes

	// setup row subgraphs
	// -------------------
	for row := 0; row < len(rowsLayout); row++ {
		orderNodes := []string{}
		rowSubGraphName := fmt.Sprintf("\"cluster_row_%d\"", row)
		rowAtt := CopyInvis()
		rowAtt["label"] = fmt.Sprintf("\"Invisible Row subgraph: %d\"", row) // help with debugging graph structure
		g.AddSubGraph("top_level_graph", rowSubGraphName, rowAtt)
		// Environment subgraphs
		envIds := rowsLayout[row]
		for _, envId := range envIds {
			label := FormatEnvLabel(txy, cfg.Terminology.L1.Singular+" - ", envId, true)
			envNodeAtt := FormatNode(label, GetClassificationValue(txy.SegL1s[envId], "sensitivity"))
			envNodeAtt["fontsize"] = "\"16\""
			envNodeName := fmt.Sprintf("\"env_node_%s\"", strings.ReplaceAll(envId, "-", "_"))
			err := g.AddNode(rowSubGraphName, envNodeName, envNodeAtt)
			if err != nil {
				return nil, err
			}
			orderNodes = append(orderNodes, envNodeName)
		}
		// Add edges to order the graph

		rowNodes[row] = orderNodes
	}
	err = AddSpacers(g, rowNodes)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// ################################
// Function to Segment Level 2 Graphs
// ################################

func GraphL2(txy *domain.Taxonomy, cfg *configdomain.Config, highlightCriticality bool, showClass bool) (*gographviz.Graph, error) {
	imageData := PrepTaxonomy(txy)
	// Setup the top level graph object
	term := cfg.Terminology
	title := term.L1.Plural + " & " + term.L2.Plural + " Layout"
	subHeading := "Overview of " + term.L2.Plural + " grouped by their respective " + term.L1.Plural
	g := BaselineGraph(title, subHeading)

	// Build rowsMap from config
	rowsLayout, err := buildRowsMap(cfg, txy)
	if err != nil {
		return nil, err
	}

	// Following code will create subgraphs for each row and add security environments as subgraphs to those rows
	// in graphviz, subgraphs are represented in their name with the prefix "cluster_"
	// sub graph structure:
	// - top_level_graph
	// 	- cluster_row
	// 		- cluster_env
	// 			- cluster_criticality
	// 				- cluster_batch

	// setup row subgraphs
	// -------------------
	rowNodes := map[int][]string{}
	for row := 0; row < len(rowsLayout); row++ {
		// edges between nodes is what the graphing algorithm will use to determine the structure of the diagram
		// To allow for different rows, the edges must be on a per-row basis
		// They must also be set at the most granular level, which is the batches for each security domain
		orderNodes := map[string]map[string][]string{}
		rowSubGraphName := fmt.Sprintf("\"cluster_row_%d\"", row)
		rowAtt := CopyInvis()
		rowAtt["label"] = fmt.Sprintf("\"Invisible Row subgraph: %d\"", row) // help with debugging graph structure
		g.AddSubGraph("top_level_graph", rowSubGraphName, rowAtt)

		// Environment subgraphs
		// --------------------
		// Now that rows have been setup, we can add the image components for each environment
		// The rowsLayout[] are slices which are ordered and therefore we can iterate over them in order,
		// tx.SegL1s is a map and therefore iterating over that would give a different order each time
		envIds := rowsLayout[row]
		for _, envId := range envIds {
			orderNodes[envId] = map[string][]string{}
			// Generate attributes object for security environment subgraph
			showEnvClass := showClass
			if highlightCriticality {
				showEnvClass = true // if highlighting criticality, we need to show the classification for the environment
			}
			label := FormatEnvLabel(txy, term.L1.Singular+" - ", envId, showEnvClass)
			envGraphAtt := FormatGraph(label, "")
			err := g.AddSubGraph(rowSubGraphName, envSubGraphName(envId), envGraphAtt)
			if err != nil {
				return nil, err
			}

			// Criticality subgraphs
			// ---------------------
			// As the colour of the nodes is set by Sensitivity, we need a way to make the criticality of nodes more visible.
			// This is done by creating a subgraph for each criticality level within the environment
			// Maps are unordered in go and therefore we need to iterate over the ordered criticality list to get a consistent image
			critGraphNames := []string{}
			for _, crit := range domain.CritOrder {
				critGraphName := focusSGName(envId, crit)
				critGraphAtt := map[string]string{}
				if highlightCriticality {
					critGraphAtt = map[string]string{
						"label":     fmt.Sprintf("\"Criticality: %s(%s)\"", crit, domain.CriticalityLevels[crit]),
						"shape":     "\"box\"",
						"color":     primaryColours.GetColour("1"),
						"fontcolor": primaryColours.GetColour("1"),
						"fontsize":  "\"14\"",
						"style":     "\"rounded,setlinewidth(1)\"",
					}
				} else {
					// Set style attribute below to \"\" if you want to see criticality subgraphs. Made invisible as it gets confusing with the environment subgraphs
					critGraphAtt = CopyInvis()
					critGraphAtt["label"] = fmt.Sprintf("\"Invisible Criticality subgraph: %s\"", crit) // help with debugging graph structure
				}
				if _, ok := imageData[envId].Criticalities[crit]; ok && (len(critGraphNames) == 0 ||
					critGraphNames[len(critGraphNames)-1] != critGraphName) {
					critGraphNames = append(critGraphNames, critGraphName)
					g.AddSubGraph(envSubGraphName(envId), focusSGName(envId, crit), critGraphAtt)
				}
			}

			// Add batch subgraphs & Segment Level 2 nodes
			// ------------------------------------------
			// Setup batch struct
			// Loop through security domains in the current environment
			batch := NewBatchVars(envId)

			for _, sdId := range imageData[envId].SortedSegs {
				sdEnvDet := imageData[envId].Segs[sdId]
				crit := GetClassificationFromOverride(sdEnvDet, "criticality")
				// Setup batch subgraphs and bump when necessary
				if batch.Count[crit] > batch.Limit || batch.Count[crit] == 0 {
					batchNodeName, err := batch.BumpBatch(crit, g)
					if err != nil {
						return nil, err
					}
					orderNodes[envId][crit] = append(orderNodes[envId][crit], batchNodeName)
				}

				// Add security domain nodes
				// -------------------------
				// Add emphasis to the label (map returns 0 if not found)
				label := FormatSdLabel(txy, "", envId, sdId, showClass, txy.SegsL2s[sdId].Prominence)
				sdNodeAtt := FormatNode(label, GetClassificationFromOverride(sdEnvDet, "sensitivity")) // attributes to format the node
				sdNodeName := fmt.Sprintf("\"sd_node_%s_%s\"", strings.ReplaceAll(envId, "-", "_"), strings.ReplaceAll(sdId, "-", "_"))
				// Add security domain node to the batch subgraph
				err := g.AddNode(batch.CurrentSGName(crit), sdNodeName, sdNodeAtt)
				if err != nil {
					return nil, err
				}
				batch.Count[crit]++

			}
		}
		// Add edges to order the graph
		// ----------------------------
		// add all nodes to a single slice before adding edges, this creates a single list of nodes to link in each row
		fullOrderNodes := []string{}
		// loop through environments and criticalities to get the order
		// To change the order of the environments update it in the config's visuals.l1_layout
		for _, env := range envIds {
			for _, c := range domain.CritOrder {
				if _, ok := imageData[env].Criticalities[c]; ok {
					fullOrderNodes = append(fullOrderNodes, orderNodes[env][c]...)
				}
			}
		}
		rowNodes[row] = fullOrderNodes
	}
	err = AddSpacers(g, rowNodes)
	if err != nil {
		return nil, err
	}

	// Add legend to the graph
	// ------------------------
	err = AddLegend(g, 12, true)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// ################################
// Function to Segment Level 2 Graphs
// ################################
// GraphCompliance  showOut is used control if out of scope domains are added to the graph
func GraphCompliance(txy *domain.Taxonomy, cfg *configdomain.Config, compName string, showOutOfScope bool) (*gographviz.Graph, error) {
	if _, ok := txy.CompReqs[compName]; !ok {
		return nil, fmt.Errorf("compliance standard %s not found in taxonomy", compName)
	}

	imageData := PrepTaxonomy(txy)
	// Setup the top level graph object
	term := cfg.Terminology
	title := term.L1.Plural + " & " + term.L2.Plural + " Layout"
	subHeading := fmt.Sprintf("Compliance Standard: %s", txy.CompReqs[compName].Name)
	g := BaselineGraph(title, subHeading)

	// Build rowsMap from config
	rowsLayout, err := buildRowsMap(cfg, txy)
	if err != nil {
		return nil, err
	}

	// setup row subgraphs
	// -------------------
	rowNodes := map[int][]string{}
	for row := 0; row < len(rowsLayout); row++ {

		orderNodes := map[string]map[string][]string{}
		rowSubGraphName := fmt.Sprintf("\"cluster_row_%d\"", row)
		rowAtt := CopyInvis()
		rowAtt["label"] = fmt.Sprintf("\"Invisible Row subgraph: %d\"", row) // help with debugging graph structure
		g.AddSubGraph("top_level_graph", rowSubGraphName, rowAtt)

		// Environment subgraphs
		// --------------------
		envIds := rowsLayout[row]
		for _, envId := range envIds {
			orderNodes[envId] = map[string][]string{}
			label := FormatEnvLabel(txy, term.L1.Singular+" - ", envId, false)
			envGraphAtt := FormatGraph(label, "")
			err := g.AddSubGraph(rowSubGraphName, envSubGraphName(envId), envGraphAtt)
			if err != nil {
				return nil, err
			}

			// Scope subgraphs
			// ---------------------
			compGraphAtt := map[string]string{
				"label":     fmt.Sprintf("\"%s - In Scope\"", compName),
				"shape":     "\"box\"",
				"color":     FontColour,
				"fontcolor": FontColour,
				"fontsize":  "\"16\"",
				"style":     "\"rounded,setlinewidth(1)\"",
			}
			name := focusSGName(envId, "in")
			g.AddSubGraph(envSubGraphName(envId), name, compGraphAtt)
			compGraphAtt["label"] = fmt.Sprintf("\"%s - Out of Scope\"", compName)
			name = focusSGName(envId, "out")
			g.AddSubGraph(envSubGraphName(envId), name, compGraphAtt)

			// Add batch subgraphs & Segment Level 2 nodes
			// ------------------------------------------
			batch := NewBatchVars(envId)

			for _, sdId := range imageData[envId].SortedSegs {
				scope := "out"
				if _, ok := imageData[envId].Segs[sdId].CompReqs[compName]; ok {
					scope = "in"
				}
				if showOutOfScope || scope == "in" {
					sdEnvDet := imageData[envId].Segs[sdId]
					//crit := sdEnvDet.Criticality
					// Setup batch subgraphs and bump when necessary
					if batch.Count[scope] > batch.Limit || batch.Count[scope] == 0 {
						batchNodeName, err := batch.BumpBatch(scope, g)
						if err != nil {
							return nil, err
						}
						orderNodes[envId][scope] = append(orderNodes[envId][scope], batchNodeName)
					}

					// Add security domain nodes
					// -------------------------
					// Add emphasis to the label (map returns 0 if not found)
					label := FormatSdLabel(txy, "", envId, sdId, false, txy.SegsL2s[sdId].Prominence)
					sdNodeAtt := FormatNode(label, GetClassificationFromOverride(sdEnvDet, "sensitivity")) // attributes to format the node
					// Make out of scope nodes less visible in the graph by removing filled style (font needs to be bright if fill removed)
					if scope == "out" {
						sdNodeAtt["fontcolor"] = sdNodeAtt["color"]
						sdNodeAtt["style"] = "\"rounded,setlinewidth(2)\""
					}
					sdNodeName := fmt.Sprintf("\"sd_node_%s_%s\"", strings.ReplaceAll(envId, "-", "_"), strings.ReplaceAll(sdId, "-", "_"))
					// Add security domain node to the batch subgraph

					err := g.AddNode(batch.CurrentSGName(scope), sdNodeName, sdNodeAtt)
					if err != nil {
						return nil, err
					}
					batch.Count[scope]++
				}
			}
		}
		// Add edges to order the graph
		// ----------------------------
		// add all nodes to a single slice before adding edges, this creates a single list of nodes to link in each row
		fullOrderNodes := []string{}
		// To change the order of the environments update it in the config's visuals.l1_layout
		for _, env := range envIds {
			fullOrderNodes = append(fullOrderNodes, orderNodes[env]["in"]...)
			if showOutOfScope {
				fullOrderNodes = append(fullOrderNodes, orderNodes[env]["out"]...)
			}
		}
		rowNodes[row] = fullOrderNodes
	}
	err = AddSpacers(g, rowNodes)
	if err != nil {
		return nil, err
	}

	// Add legend to the graph
	// ------------------------
	err = AddLegend(g, 12, true)
	if err != nil {
		return nil, err
	}

	return g, nil
}
