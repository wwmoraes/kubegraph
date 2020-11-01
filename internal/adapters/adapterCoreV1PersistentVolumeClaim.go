package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterCoreV1PersistentVolumeClaim struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterCoreV1PersistentVolumeClaim{
		Resource{
			resourceType: reflect.TypeOf(&coreV1.PersistentVolumeClaim{}),
		},
	})
}

func (adapter adapterCoreV1PersistentVolumeClaim) tryCastObject(obj runtime.Object) (*coreV1.PersistentVolumeClaim, error) {
	casted, ok := obj.(*coreV1.PersistentVolumeClaim)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterCoreV1PersistentVolumeClaim) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterCoreV1PersistentVolumeClaim) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/pvc.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterCoreV1PersistentVolumeClaim) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterCoreV1PersistentVolumeClaim) Configure(statefulGraph StatefulGraph) error {
	volumeAdapter, err := GetAdapterFor(reflect.TypeOf(&coreV1.Volume{}))
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err)
	}

	objects, err := statefulGraph.GetObjects(adapter.GetType())
	if err != nil {
		return err
	}

	for resourceName, resourceObject := range objects {
		resource, err := adapter.tryCastObject(resourceObject)
		if err != nil {
			return err
		}
		resourceNode, err := statefulGraph.GetNode(adapter.GetType(), resourceName)
		if err != nil {
			return err
		}

		if resource.Spec.VolumeName != "" {
			volumeAdapter.Connect(statefulGraph, resourceNode, resource.Spec.VolumeName)
		}
	}

	return nil
}
