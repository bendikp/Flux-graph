package graph

import (
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/distributed-technologies/flux-graph/pkg/helmRelease"
	"github.com/distributed-technologies/flux-graph/pkg/kustomization"
	"github.com/distributed-technologies/flux-graph/pkg/logging"
)

type kustomizationGraph struct {
	*gographviz.Graph
}

var subGraphAttr = map[string]string{
	"shape": "rectangle",
	"rank":  "same",
}
var nodeAttr = map[string]string{
	"shape": "rectangle",
}

func New(graphName string) *kustomizationGraph {
	defaultGraph := gographviz.NewGraph()
	defaultGraph.SetName(graphName)
	defaultGraph.Directed = true
	// Needed for edges to point to subgraphs (clusters) https://graphviz.org/docs/attrs/compound/
	defaultGraph.Attrs.Add("compound", "true")

	graph := &kustomizationGraph{
		Graph: defaultGraph,
	}

	return graph
}

func (g *kustomizationGraph) Generate(ks []kustomization.Kustomization) (string, error) {

	for _, ks := range ks {

		if len(ks.HRSlice) > 0 {
			sub := g.makeSubGraph(g.Name, ks.Metadata.Name)

			for _, hr := range helmRelease.HelmReleases {
				if hr.Parent != ks.Metadata.Name {
					continue
				}

				// Look up the node, if missing create one, and look it up again
				startN := g.makeMissingNode(sub.Name, hr.Name())
				logging.Debug("startN: %v", startN)

				err := g.makeDependencyNodes(sub.Name, startN, hr.GetDependencies())
				if err != nil {
					return "", err
				}

			}
		}

		// Look up the node, if missing create one, and look it up again
		startN := g.makeMissingNode(g.Name, ks.Name())
		logging.Debug("startN: %v", startN)

		err := g.makeDependencyNodes(g.Name, startN, ks.GetDependencies())
		if err != nil {
			return "", err
		}

		logging.Debug("Ks: %v", ks)

	}

	return g.String(), nil
}

func (g *kustomizationGraph) makeDependencyNodes(parentGraph string, startN string, dependencies []string) error {
	for _, dep := range dependencies {

		// Look up the node, if missing create one, and look it up again
		depN := g.makeMissingNode(parentGraph, dep)
		logging.Debug("depN: %s\n", depN)

		// Create relational arrows from the start node to the dependency nodes
		err := g.makeEdge(depN, startN)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *kustomizationGraph) makeEdge(source string, destination string) error {
	var arrowAttr = map[string]string{
		"dir": "back",
	}

	sourceSubGraph := g.getSubGraph(source)

	if sourceSubGraph != nil {
		children := g.subGraphChildren(sourceSubGraph)

		arrowAttr["ltail"] = sourceSubGraph.Name
		source = children[len(children)-1]
	}

	destinationSubGraph := g.getSubGraph(destination)

	if destinationSubGraph != nil {
		children := g.subGraphChildren(destinationSubGraph)

		arrowAttr["lhead"] = destinationSubGraph.Name
		destination = children[0]
	}

	// Create relational arrows from the start node to the dependency nodes
	err := g.AddEdge(source, destination, true, arrowAttr)
	if err != nil {
		return err
	}

	return nil
}

func (g *kustomizationGraph) makeMissingNode(parentGraph string, name string) string {
	node := g.Nodes.Lookup[name]
	subGraph := g.getSubGraph(name)

	if node == nil && subGraph == nil {
		g.addNodeToGraph(parentGraph, name)
		node = g.Nodes.Lookup[name]
		return node.Name
	}
	return name
}

func (g *kustomizationGraph) makeSubGraph(parentGraph string, subGraphName string) *gographviz.SubGraph {
	name := "cluster_" + subGraphName

	g.AddSubGraph(g.Name, name, subGraphAttr)
	sub := g.SubGraphs.SubGraphs[name]
	sub.Attrs.Add("label", subGraphName)

	return sub
}

func (g *kustomizationGraph) addNodeToGraph(parentGraph string, nodeName string) (string, error) {
	node := nodeName
	if g.IsNode(node) {
		return node, nil
	}

	err := g.AddNode(parentGraph, node, nodeAttr)
	if err != nil {
		return "", err
	}

	return node, nil
}

func (g *kustomizationGraph) getSubGraph(name string) *gographviz.SubGraph {
	subGraphName := "cluster_" + strings.ReplaceAll(name, "\"", "")
	return g.SubGraphs.SubGraphs[subGraphName]
}

func (g *kustomizationGraph) subGraphChildren(graph *gographviz.SubGraph) []string {
	subGraphChildren := g.Relations.SortedChildren(graph.Name)
	logging.Debug("subGraph children: %v", subGraphChildren)
	return subGraphChildren
}
