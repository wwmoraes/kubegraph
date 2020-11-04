package v1

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type persistentVolumeClaimAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&persistentVolumeClaimAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&coreV1.PersistentVolumeClaim{}),
		},
	})
}

func (thisAdapter persistentVolumeClaimAdapter) tryCastObject(obj runtime.Object) (*coreV1.PersistentVolumeClaim, error) {
	casted, ok := obj.(*coreV1.PersistentVolumeClaim)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter persistentVolumeClaimAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter persistentVolumeClaimAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/pvc.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter persistentVolumeClaimAdapter) Connect(statefulGraph adapter.StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter persistentVolumeClaimAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	persistentVolumeAdapter, err := adapter.Get(reflect.TypeOf(&coreV1.PersistentVolume{}))
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err)
	}

	objects, err := statefulGraph.GetObjects(thisAdapter.GetType())
	if err != nil {
		return err
	}

	for resourceName, resourceObject := range objects {
		resource, err := thisAdapter.tryCastObject(resourceObject)
		if err != nil {
			return err
		}
		resourceNode, err := statefulGraph.GetNode(thisAdapter.GetType(), resourceName)
		if err != nil {
			return err
		}

		if resource.Spec.VolumeName != "" {
			persistentVolumeAdapter.Connect(statefulGraph, resourceNode, resource.Spec.VolumeName)
		}
	}

	return nil
}
