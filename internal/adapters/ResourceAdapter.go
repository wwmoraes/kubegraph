package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	"k8s.io/apimachinery/pkg/runtime"
)

var resourceAdapters = make(map[reflect.Type]ResourceAdapter)

// ResourceAdapter instructions on how to deal with a given kubernetes kind
type ResourceAdapter interface {
	// GetType returns the reflected type of the k8s kind managed by this instance
	GetType() reflect.Type
	// Create add a graph node for the given object and stores it for further actions
	Create(sgraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error)
	// Connect creates and edge between the given node and an object on this adapter
	Connect(sgraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error)
	// Configure connects the resources on this adapter with its dependencies
	Configure(sgraph StatefulGraph) error
}

// RegisterResourceAdapter sets up a given resource type adapter
func RegisterResourceAdapter(adapter ResourceAdapter) {
	if _, exists := resourceAdapters[adapter.GetType()]; exists {
		panic(fmt.Errorf("only one adapter should be registered per type %s", adapter.GetType().String()))
	}
	resourceAdapters[adapter.GetType()] = adapter
}

// GetAdapterFor returns the adapter for a resource type, if it is registered
func GetAdapterFor(resourceType reflect.Type) (ResourceAdapter, error) {
	adapter, adapterExists := resourceAdapters[resourceType]
	if !adapterExists {
		// TODO return noop adapter to reduce ifs on adapter configure functions
		return nil, fmt.Errorf("type %s has no adapter registered", resourceType.String())
	}

	return adapter, nil
}

// GetAdapters returns all adapters registered
func GetAdapters() map[reflect.Type]ResourceAdapter {
	return resourceAdapters
}
