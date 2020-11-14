package adapter

import (
	"fmt"

	"github.com/wwmoraes/dot"
)

type Graph interface {
	dot.Graph
}

func NewGraph() (Graph, error) {
	graph, err := dot.NewGraph(
		dot.WithID("kubegraph"),
		dot.WithType(dot.GraphTypeDirected),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to initialize graph: %w", err)
	}

	return graph, nil
}
