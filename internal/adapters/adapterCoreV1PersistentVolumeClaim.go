package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterCoreV1PersistentVolumeClaim struct{}

func init() {
	RegisterResourceAdapter(&adapterCoreV1PersistentVolumeClaim{})
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterCoreV1PersistentVolumeClaim) GetType() reflect.Type {
	return reflect.TypeOf(&coreV1.PersistentVolumeClaim{})
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterCoreV1PersistentVolumeClaim) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource := obj.(*coreV1.PersistentVolumeClaim)
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
		resource := resourceObject.(*coreV1.PersistentVolumeClaim)
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
