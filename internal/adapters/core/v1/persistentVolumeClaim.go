package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type persistentVolumeClaimAdapter struct {
	adapter.Resource
}

func init() {
	adapter.MustRegister(NewPersistentVolumeClaimAdapter())
}

func NewPersistentVolumeClaimAdapter() adapter.Resource {
	return &persistentVolumeClaimAdapter{
		adapter.NewResource(
			reflect.TypeOf(&coreV1.PersistentVolumeClaim{}),
			"icons/pvc.svg",
		),
	}
}

func (thisAdapter *persistentVolumeClaimAdapter) tryCastObject(obj runtime.Object) (*coreV1.PersistentVolumeClaim, error) {
	casted, ok := obj.(*coreV1.PersistentVolumeClaim)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *persistentVolumeClaimAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	persistentVolumeAdapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&coreV1.PersistentVolume{}))
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
			_, err := persistentVolumeAdapter.Connect(statefulGraph, resourceNode, resource.Spec.VolumeName)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
			}
		}
	}

	return nil
}
