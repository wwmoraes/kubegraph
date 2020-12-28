package adapter

import (
	"fmt"

	"github.com/wwmoraes/dot"
)

// Graph implements the external package-specific Graph interface
type Graph interface {
	dot.Graph
}

// NewGraph creates a new instance of a Graph, and returns an error if any
// happens on the underlying dot constructor
func NewGraph() (Graph, error) {
	graph, err := dot.New(
		dot.WithID("kubegraph"),
		dot.WithType(dot.GraphTypeDirected),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to initialize graph: %w", err)
	}

	return graph, nil
}
