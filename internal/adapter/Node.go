package adapter

import (
	"fmt"

	"github.com/emicklei/dot"
)

type Node interface {
	ID() string
	Attrs(labelvalues ...interface{})
	Value(label string) interface{}
}

func TryGetDotNode(node Node) (*dot.Node, error) {
	dotNode, ok := node.(*dot.Node)
	if !ok {
		return nil, fmt.Errorf("source node %s is not compatible with *dot.Node", node.ID())
	}

	return dotNode, nil
}
