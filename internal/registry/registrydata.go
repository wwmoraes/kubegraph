package registry

import (
	"reflect"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
)

var once sync.Once

func init() {
	once.Do(func() {
		instance = &registryData{
			adapters: make(AdapterMap),
			accessor: meta.NewAccessor(),
		}
	})
}

// registryData implements Registry to store adapters
type registryData struct {
	adapters AdapterMap
	accessor meta.MetadataAccessor
}

// Register stores the adapter if none is currently registered for such type,
// and if it is not registered with another registry instance
func (registry *registryData) Register(adapter Adapter) error {
	if _, err := registry.Get(adapter.GetType()); err == nil {
		return ErrAdapterAlreadyRegistered
	}

	registry.adapters[adapter.GetType()] = adapter
	return nil
}

// Deregister removes the adapter for its type
func (registry *registryData) Deregister(adapter Adapter) {
	registry.DeregisterByType(adapter.GetType())
}

// DeregisterByType removes the adapter for the given type, if any is registered
func (registry *registryData) DeregisterByType(adapterType reflect.Type) {
	delete(registry.adapters, adapterType)
}

// Get returns the adapter for a resource type, if it is registered
func (registry *registryData) Get(resourceType reflect.Type) (Adapter, error) {
	adapter, adapterExists := registry.adapters[resourceType]
	if !adapterExists {
		return nil, ErrAdapterNotFound
	}

	return adapter, nil
}

// GetMultiple returns a map of adapters with only the requested types
func (registry *registryData) GetMultiple(resourceTypes ...reflect.Type) AdapterMap {
	var resourceAdapters = make(AdapterMap, len(resourceTypes))

	for _, resourceType := range resourceTypes {
		resourceAdapters[resourceType] = registry.adapters[resourceType]
	}

	return resourceAdapters
}

// GetAll returns all adapters registered
func (registry *registryData) GetAll() AdapterMap {
	return registry.adapters
}

// GetAccessor returns a global instance of a kubernetes metadata accessor
func (registry *registryData) GetAccessor() MetadataAccessor {
	return registry.accessor
}
