package registry

import (
	"reflect"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

type MetadataAccessor = meta.MetadataAccessor

type RuntimeObject = runtime.Object

// AdapterMap is a collection of adapter values by their reflected type
type AdapterMap = map[reflect.Type]Adapter

type RegistryReader interface {
	// Get returns the adapter for a resource type, if it is registered
	Get(reflect.Type) (Adapter, error)
	// GetMultiple returns a map of adapters with only the requested types
	GetMultiple(...reflect.Type) AdapterMap
	// GetAll returns all adapters registered
	GetAll() AdapterMap
	// GetAccessor returns the kubernetes metadata accessor shared instance
	GetAccessor() MetadataAccessor
}

type RegistryWriter interface {
	// Register stores the adapter if none is currently registered for such type,
	// and if it is not registered with another registry instance
	Register(Adapter) error
	// Deregister removes the adapter for its type
	Deregister(Adapter)
	// DeregisterByType removes the adapter for the given type, if any is registered
	DeregisterByType(reflect.Type)
}

// Registry is implemented by values that stores functions to create, connect
// and configure kubernetes resource nodes on a graphviz context
type Registry interface {
	RegistryReader
	RegistryWriter
}

var instance Registry

// Instance returns the singleton Registry instance used across all adapters
func Instance() Registry {
	return instance
}

func Register(adapter Adapter) error {
	return Instance().Register(adapter)
}

func MustRegister(adapter Adapter) {
	if err := Register(adapter); err != nil {
		panic(err)
	}
}
