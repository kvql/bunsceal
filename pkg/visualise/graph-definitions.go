package visualise

import (
	"fmt"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

// ################################
// Function to Segment Level 1 Graph
// ################################

// See the GraphSDs function for more comments explaining how the graph is generated. That is the more complex function and therefore has more comments than GraphEnvs

func GraphL1(txy domain.Taxonomy, terms domain.TermConfig, visCfg VisualsDef, allGroups []plugins.ImageGroupingData) (*gographviz.Graph, error) {
	// Setup the top level graph object
	title := terms.L1.Plural + " Overview"
	g := BaselineGraph(title, "")

	// // Add legend to the graph
	// // ------------------------
	err := AddLegend(g, allGroups, 12, true)
	if err != nil {
		return nil, err
	}

	// Build rowsMap from config
	rowsLayout, err := buildRowsMap(visCfg, txy)
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
			label := FormatEnvLabel(txy, terms.L1.Singular+" - ", envId, true)
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

func GraphL2Grouped(txy domain.Taxonomy, terms domain.TermConfig, visCfg VisualsDef, allGroups []plugins.ImageGroupingData, groupData plugins.ImageGroupingData) (*gographviz.Graph, error) {
	imageData := VisL2GroupingPrep(txy, groupData)
	// Setup the top level graph object
	title := terms.L1.Plural + " & " + terms.L2.Plural + " Layout"
	subHeading := "Overview of " + terms.L2.Plural + " grouped by their respective " + terms.L1.Plural
	g := BaselineGraph(title, subHeading)
	groupingEnabled := groupData.Namespace != "" && groupData.Key != ""

	// Build rowsMap from config
	rowsLayout, err := buildRowsMap(visCfg, txy)
	if err != nil {
		return nil, err
	}

	// Following code will create subgraphs for each row and add security environments as subgraphs to those rows
	// in graphviz, subgraphs are represented in their name with the prefix "cluster_"
	// sub graph structure:
	// - top_level_graph
	// 	- cluster_row
	// 		- cluster_env
	// 			- cluster_classification (criticality or sensitivity)
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
			label := FormatEnvLabel(txy, terms.L1.Singular+" - ", envId, true)
			envGraphAtt := FormatGraph(label, "")
			err := g.AddSubGraph(rowSubGraphName, envSubGraphName(envId), envGraphAtt)
			if err != nil {
				return nil, err
			}

			// Classification subgraphs (Criticality or Sensitivity)
			// ------------------------------------------------------
			// Create subgraphs to group nodes by classification level
			// Maps are unordered in go and therefore we need to iterate over the ordered list to get a consistent image
			if groupingEnabled {
				// Group by sensitivity
				groupGraphNames := []string{}
				for groupVal, i := range groupData.OrderMap {
					groupGraphName := focusSGName(envId, groupVal)
					groupGraphAtt := map[string]string{
						"label":     fmt.Sprintf("\"%s: %s(%s)\"", groupData.DisplayName, groupVal, groupData.ValuesMap[groupVal]),
						"shape":     "\"box\"",
						"color":     primaryColours.GetColour(i),
						"fontcolor": primaryColours.GetColour(i, true),
						"fontsize":  "\"14\"",
						"style":     "\"rounded,setlinewidth(1)\"",
					}
					if _, ok := imageData[envId].PresentGroupValues[groupVal]; ok && (len(groupGraphNames) == 0 ||
						groupGraphNames[len(groupGraphNames)-1] != groupGraphName) {
						groupGraphNames = append(groupGraphNames, groupGraphName)
						g.AddSubGraph(envSubGraphName(envId), focusSGName(envId, groupVal), groupGraphAtt)
					}
				}
			} else {
				g.AddSubGraph(envSubGraphName(envId), focusSGName(envId, unknownGroupKey), InvisAtt)
			}

			// Add batch subgraphs & Segment Level 2 nodes
			// ------------------------------------------
			// Setup batch struct
			// Loop through security domains in the current environment
			batch := NewBatchVars(envId)

			for _, segL2Id := range imageData[envId].SortedSegs {

				seg := txy.SegsL2s[segL2Id]
				var groupKey string
				// Get grouping key based on mode
				var colorIndex int
				if groupingEnabled {
					groupKey, err = seg.GetNamespacedValue(envId, groupData.Namespace, groupData.Key)
					if err != nil {
						return nil, err
					}
					colorIndex = groupData.OrderMap[groupKey]
				} else {
					groupKey = unknownGroupKey
				}
				// Setup batch subgraphs and bump when necessary
				if batch.Count[groupKey] > batch.Limit || batch.Count[groupKey] == 0 {
					batchNodeName, err := batch.BumpBatch(groupKey, g)
					if err != nil {
						return nil, err
					}
					orderNodes[envId][groupKey] = append(orderNodes[envId][groupKey], batchNodeName)
				}

				// Add L2 Seg nodes
				// -------------------------
				// Add emphasis to the label (map returns 0 if not found)
				label := FormatSdLabel(txy, "", envId, segL2Id,
					groupingEnabled, txy.SegsL2s[segL2Id].Prominence)
				l2SegNodeAtt := FormatGroupedNode(label, colorIndex)
				l2SegNodeName := fmt.Sprintf("\"l2_seg_node_%s_%s\"",
					strings.ReplaceAll(envId, "-", "_"),
					strings.ReplaceAll(segL2Id, "-", "_"))
				// Add security domain node to the batch subgraph
				err = g.AddNode(batch.CurrentSGName(groupKey), l2SegNodeName, l2SegNodeAtt)
				if err != nil {
					return nil, err
				}
				batch.Count[groupKey]++

			}
		}
		// Add edges to order the graph
		// ----------------------------
		// add all nodes to a single slice before adding edges, this creates a single list of nodes to link in each row
		fullOrderNodes := []string{}
		// loop through environments and classifications to get the order
		// To change the order of the environments update it in the config's visuals.l1_layout
		for _, env := range envIds {
			if groupingEnabled {
				for _, groupValue := range groupData.OrderedValues {
					if _, ok := imageData[env].PresentGroupValues[groupValue]; ok {
						fullOrderNodes = append(fullOrderNodes, orderNodes[env][groupValue]...)
					}
				}
			} else {
				fullOrderNodes = append(fullOrderNodes, orderNodes[env][unknownGroupKey]...)
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
	err = AddLegend(g, allGroups, 12, true)
	if err != nil {
		return nil, err
	}

	return g, nil
}
