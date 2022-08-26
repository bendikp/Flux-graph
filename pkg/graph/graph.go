package graph

import (
	"github.com/awalterschulze/gographviz"
	"github.com/distributed-technologies/flux-graph/pkg/kustomization"
	"github.com/distributed-technologies/flux-graph/pkg/logging"
)

type kustomizationGraph struct {
	*gographviz.Graph
}

var nodeAttr = map[string]string{
	"shape": "rectangle",
}
var ArrowAttr = map[string]string{
	"dir": "back",
}

func New(graphName string) *kustomizationGraph {
	defaultGraph := gographviz.NewGraph()
	defaultGraph.SetName(graphName)
	defaultGraph.Directed = true

	graph := &kustomizationGraph{
		Graph: defaultGraph,
	}

	return graph
}

func (g *kustomizationGraph) Generate(ks []kustomization.Kustomization) (string, error) {
	for _, ks := range ks {
		// Look up the node, if missing create one, and look it up again
		startN := g.Nodes.Lookup[ks.Name()]
		if startN == nil {
			g.addNodeToGraph(ks.Name())
			startN = g.Nodes.Lookup[ks.Name()]
		}
		logging.Debug("startN: %v", startN.Name)

		for _, dep := range ks.GetDepndencies() {
			// Look up the node, if missing create one, and look it up again
			depN := g.Nodes.Lookup[dep]
			if depN == nil {
				_, err := g.addNodeToGraph(dep)
				if err != nil {
					return "", err
				}
				depN = g.Nodes.Lookup[dep]
			}
			logging.Debug("depN: %s\n", depN.Name)

			// Create relational arrows from the start node to the dependency nodes
			err := g.AddEdge(depN.Name, startN.Name, true, ArrowAttr)
			if err != nil {
				return "", err
			}
		}

	}

	return g.String(), nil
}

func (g *kustomizationGraph) addNodeToGraph(nodeName string) (string, error) {
	node := nodeName
	logging.Debug("node: %v", node)
	if g.IsNode(node) {
		return node, nil
	}

	err := g.AddNode(g.Name, node, nodeAttr)
	if err != nil {
		return "", err
	}

	return node, nil
}
