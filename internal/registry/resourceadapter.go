package registry

import (
	"fmt"
	"reflect"
)

// ResourceAdapter implements specialized methods for retrieving kubernetes
// resources from an adapter with the proper runtime type
type ResourceAdapter[T RuntimeObject] interface {
	Adapter

	CastObject(RuntimeObject) (T, error)
	GetGraphObjects(StatefulGraph) (map[string]T, error)
}

func GetAdapter[T RuntimeObject]() (ResourceAdapter[T], error) {
	// call to .Elem() is needed as T is a pointer to a pointer
	adapter, err := Instance().Get(reflect.TypeOf((*T)(nil)).Elem())
	if err != nil {
		return nil, err
	}

	casted, ok := adapter.(ResourceAdapter[T])
	if !ok {
		return nil, fmt.Errorf("failed to get adapter: %w", ErrIncompatibleType)
	}
	return casted, nil
}

func GetGraphObjects[T RuntimeObject](statefulGraph StatefulGraph) (map[string]T, error) {
	reflected := reflect.TypeOf((*T)(nil))
	// call to .Elem() is needed as T is a pointer to a pointer
	objects, err := statefulGraph.GetObjects(reflected.Elem())
	if err != nil {
		return nil, err
	}

	castedObjects := make(map[string]T, len(objects))
	for key, object := range objects {
		casted, ok := object.(T)
		if !ok {
			return nil, fmt.Errorf("failed to get objects: %w", ErrIncompatibleType)
		}
		castedObjects[key] = casted
	}
	return castedObjects, nil
}
