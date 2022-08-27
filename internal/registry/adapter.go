package registry

import (
	"reflect"
)

// Adapter is implemented by values that can transform a kubernetes object
// kind information into nodes and create edges between them
type Adapter interface {
	// IconPath returns the type icon file path
	IconPath() string
	// GetType returns the reflected type of the k8s kind managed by this instance
	GetType() reflect.Type
	// Create add a graph node for the given object and stores it for further actions
	Create(StatefulGraph, RuntimeObject) (Node, error)
	// Connect creates and edge between the given node and an object on this adapter
	Connect(graph StatefulGraph, source Node, targetName string) (Edge, error)
	// Configure connects the resources on this adapter with its dependencies
	Configure(graph StatefulGraph) error
	// GetGraphNode returns the node registered under name
	GetGraphNode(StatefulGraph, string) (Node, error)
}
