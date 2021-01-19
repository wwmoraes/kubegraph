package adapter

import (
	"reflect"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
)

type ResourceMap map[reflect.Type]Resource

type registrySingleton struct {
	adapters ResourceAdapterMap
	accessor meta.MetadataAccessor
}

var once sync.Once
var registryInstance *registrySingleton

// RegistryInstance returns the singleton Registry instance
func RegistryInstance() Registry {
	once.Do(func() {
		registryInstance = &registrySingleton{
			adapters: make(ResourceAdapterMap),
			accessor: meta.NewAccessor(),
		}
	})

	return registryInstance
}

// Register stores the adapter if none is currently registered for such type,
// and if it is not registered with another registry instance
func (thisInstance *registrySingleton) Register(adapter Resource) error {
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
func (thisInstance *registrySingleton) Deregister(adapter Resource) {
	thisInstance.DeregisterByType(adapter.GetType())
}

// DeregisterByType removes the adapter for the given type, if any is registered
func (thisInstance *registrySingleton) DeregisterByType(adapterType reflect.Type) {
	delete(thisInstance.adapters, adapterType)
}

// Get returns the adapter for a resource type, if it is registered
func (thisInstance *registrySingleton) Get(resourceType reflect.Type) (Resource, error) {
	adapter, adapterExists := thisInstance.adapters[resourceType]
	if !adapterExists {
		return nil, ErrAdapterNotFound
	}

	return adapter, nil
}

// GetMultiple returns a map of adapters with only the requested types
func (thisInstance *registrySingleton) GetMultiple(resourceTypes ...reflect.Type) ResourceAdapterMap {
	var resourceAdapters = make(ResourceAdapterMap, len(resourceTypes))

	for _, resourceType := range resourceTypes {
		resourceAdapters[resourceType] = thisInstance.adapters[resourceType]
	}

	return resourceAdapters
}

// GetAll returns all adapters registered
func (thisInstance *registrySingleton) GetAll() ResourceAdapterMap {
	return thisInstance.adapters
}

// GetAccessor returns a global instance of a kubernetes metadata accessor
func (thisInstance *registrySingleton) GetAccessor() meta.MetadataAccessor {
	return thisInstance.accessor
}
