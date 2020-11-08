package adapter

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceAdapter instructions on how to deal with a given kubernetes kind
type ResourceAdapter interface {
	// GetType returns the reflected type of the k8s kind managed by this instance
	GetType() reflect.Type
	// Create add a graph node for the given object and stores it for further actions
	Create(sgraph StatefulGraph, obj runtime.Object) (Node, error)
	// Connect creates and edge between the given node and an object on this adapter
	Connect(sgraph StatefulGraph, source Node, targetName string) (Edge, error)
	// Configure connects the resources on this adapter with its dependencies
	Configure(sgraph StatefulGraph) error
}
