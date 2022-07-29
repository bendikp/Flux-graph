package graphvizWrapper

import (
	"bytes"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type GraphWrap struct {
	G  *cgraph.Graph
	Gv *graphviz.Graphviz
}

func (gw GraphWrap) MakeNode(name string) *cgraph.Node {
	n, err := gw.G.CreateNode(name)
	if err != nil {
		log.Fatal(err)
	}
	n.SetStyle(cgraph.RoundedNodeStyle)
	n.SetShape(cgraph.RectangleShape)
	n.SetFixedSize(false)
	return n
}

func (gw GraphWrap) RenderSVG(output string) {
	var buf bytes.Buffer
	if err := gw.Gv.Render(gw.G, "dot", &buf); err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())

	if err := gw.Gv.RenderFilename(gw.G, graphviz.SVG, fmt.Sprintf("%s/%s", output, "graph.svg")); err != nil {
		log.Fatal(err)
	}
}
