package adapter

import (
	"fmt"
	"reflect"
)

var adapters = make(map[reflect.Type]ResourceAdapter)

// Register sets up a given resource type adapter
func Register(adapter ResourceAdapter) {
	if _, exists := adapters[adapter.GetType()]; exists {
		panic(fmt.Errorf("only one adapter should be registered per type %s", adapter.GetType().String()))
	}
	adapters[adapter.GetType()] = adapter
}

// Get returns the adapter for a resource type, if it is registered
func Get(resourceType reflect.Type) (ResourceAdapter, error) {
	adapter, adapterExists := adapters[resourceType]
	if !adapterExists {
		// TODO return noop adapter to reduce ifs on adapter configure functions
		return nil, fmt.Errorf("type %s has no adapter registered", resourceType.String())
	}

	return adapter, nil
}

// GetAll returns all adapters registered
func GetAll() map[reflect.Type]ResourceAdapter {
	return adapters
}
