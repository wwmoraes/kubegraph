package adapter

import (
	"io"

	"github.com/emicklei/dot"
)

type Graph interface {
	Node(id string) *dot.Node
	Edge(fromNode, toNode *dot.Node, labels ...string) *dot.Edge
	Attrs(labelvalues ...interface{})
	Write(w io.Writer)
}

func NewGraph() Graph {
	dotGraph := dot.NewGraph(dot.Directed)

	dotGraph.ID("kubegraph")

	return dotGraph
}
