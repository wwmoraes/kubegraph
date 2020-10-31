package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

var adapters = make(map[reflect.Type]ResourceAdapter)

// ResourceAdapter instructions on how to deal with a given kubernetes kind
type ResourceAdapter interface {
	// GetType returns the reflected type of the k8s kind managed by this instance
	GetType() reflect.Type
	// Create add a graph node for the given object and stores it for further actions
	Create(KubeGraph, k8sruntime.Object) (*cgraph.Node, error)
	// Connect creates and edge between the given node and an object on this adapter
	Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error)
	// Configure connects the resources on this adapter with its dependencies
	Configure(kgraph KubeGraph) error
}

// RegisterResourceAdapter sets up a given resource type adapter
func RegisterResourceAdapter(adapter ResourceAdapter) {
	if _, exists := adapters[adapter.GetType()]; exists {
		panic(fmt.Errorf("only one adapter should be registered per type %s", adapter.GetType().String()))
	}
	adapters[adapter.GetType()] = adapter
}
