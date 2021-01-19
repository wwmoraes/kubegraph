package adapter

import (
	"reflect"

	"k8s.io/apimachinery/pkg/api/meta"
)

// Registry is implemented by values that stores functions to create, connect
// and configure kubernetes resource nodes on a graphviz context
type Registry interface {
	// Register stores the adapter if none is currently registered for such type,
	// and if it is not registered with another registry instance
	Register(Resource) error
	// Deregister removes the adapter for its type
	Deregister(Resource)
	// DeregisterByType removes the adapter for the given type, if any is registered
	DeregisterByType(reflect.Type)
	// Get returns the adapter for a resource type, if it is registered
	Get(reflect.Type) (Resource, error)
	// GetMultiple returns a map of adapters with only the requested types
	GetMultiple(...reflect.Type) ResourceAdapterMap
	// GetAll returns all adapters registered
	GetAll() ResourceAdapterMap
	// GetAccessor returns the kubernetes metadata accessor shared instance
	GetAccessor() meta.MetadataAccessor
}
