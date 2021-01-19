package registry

import (
	"reflect"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
)

var once sync.Once
var registryInstance *registryData

func init() {
	once.Do(func() {
		registryInstance = &registryData{
			adapters: make(AdapterMap),
			accessor: meta.NewAccessor(),
		}
	})
}

func Register(adapter Adapter) error {
	return Instance().Register(adapter)
}

func MustRegister(adapter Adapter) {
	if err := Register(adapter); err != nil {
		panic(err)
	}
}

// AdapterMap is a collection of adapter values by their reflected type
type AdapterMap map[reflect.Type]Adapter

// Registry is implemented by values that stores functions to create, connect
// and configure kubernetes resource nodes on a graphviz context
type Registry interface {
	// Register stores the adapter if none is currently registered for such type,
	// and if it is not registered with another registry instance
	Register(Adapter) error
	// Deregister removes the adapter for its type
	Deregister(Adapter)
	// DeregisterByType removes the adapter for the given type, if any is registered
	DeregisterByType(reflect.Type)
	// Get returns the adapter for a resource type, if it is registered
	Get(reflect.Type) (Adapter, error)
	// GetMultiple returns a map of adapters with only the requested types
	GetMultiple(...reflect.Type) AdapterMap
	// GetAll returns all adapters registered
	GetAll() AdapterMap
	// GetAccessor returns the kubernetes metadata accessor shared instance
	GetAccessor() meta.MetadataAccessor
}

type registryData struct {
	adapters AdapterMap
	accessor meta.MetadataAccessor
}

// Instance returns the singleton Registry instance used across all adapters
func Instance() Registry {
	return registryInstance
}

// Register stores the adapter if none is currently registered for such type,
// and if it is not registered with another registry instance
func (thisInstance *registryData) Register(adapter Adapter) error {
	if adapter.GetRegistry() != nil {
		return ErrAdapterRegisteredElsewhere
	}
	if _, err := thisInstance.Get(adapter.GetType()); err != ErrAdapterNotFound {
		return ErrAdapterAlreadyRegistered
	}

	adapter.SetRegistry(thisInstance)
	thisInstance.adapters[adapter.GetType()] = adapter
	return nil
}

// Deregister removes the adapter for its type
func (thisInstance *registryData) Deregister(adapter Adapter) {
	thisInstance.DeregisterByType(adapter.GetType())
}

// DeregisterByType removes the adapter for the given type, if any is registered
func (thisInstance *registryData) DeregisterByType(adapterType reflect.Type) {
	delete(thisInstance.adapters, adapterType)
}

// Get returns the adapter for a resource type, if it is registered
func (thisInstance *registryData) Get(resourceType reflect.Type) (Adapter, error) {
	adapter, adapterExists := thisInstance.adapters[resourceType]
	if !adapterExists {
		return nil, ErrAdapterNotFound
	}

	return adapter, nil
}

// GetMultiple returns a map of adapters with only the requested types
func (thisInstance *registryData) GetMultiple(resourceTypes ...reflect.Type) AdapterMap {
	var resourceAdapters = make(AdapterMap, len(resourceTypes))

	for _, resourceType := range resourceTypes {
		resourceAdapters[resourceType] = thisInstance.adapters[resourceType]
	}

	return resourceAdapters
}

// GetAll returns all adapters registered
func (thisInstance *registryData) GetAll() AdapterMap {
	return thisInstance.adapters
}

// GetAccessor returns a global instance of a kubernetes metadata accessor
func (thisInstance *registryData) GetAccessor() meta.MetadataAccessor {
	return thisInstance.accessor
}
