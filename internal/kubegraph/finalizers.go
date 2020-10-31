package kubegraph

import (
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func closeGraph(graph *cgraph.Graph) error {
	return graph.Close()
}

func closeGraphviz(graphviz *graphviz.Graphviz) error {
	return graphviz.Close()
}
