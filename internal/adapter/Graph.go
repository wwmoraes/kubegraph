package adapter

import (
	"github.com/wwmoraes/dot"
	"github.com/wwmoraes/dot/attributes"
)

type Graph interface {
	dot.Graph
}

func NewGraph() Graph {
	dotGraph := dot.NewGraph(&dot.GraphOptions{
		ID:   "kubegraph",
		Type: attributes.GraphTypeDirected,
	})

	return dotGraph
}
