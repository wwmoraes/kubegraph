package adapters

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterCoreV1Pod struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterCoreV1Pod{
		Resource{
			resourceType: reflect.TypeOf(&coreV1.Pod{}),
		},
	})
}

func (adapter adapterCoreV1Pod) tryCastObject(obj runtime.Object) (*coreV1.Pod, error) {
	casted, ok := obj.(*coreV1.Pod)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterCoreV1Pod) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterCoreV1Pod) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/pod.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterCoreV1Pod) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterCoreV1Pod) Configure(statefulGraph StatefulGraph) error {
	configMapAdapter, err := GetAdapterFor(reflect.TypeOf(&coreV1.ConfigMap{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err))
	}

	secretAdapter, err := GetAdapterFor(reflect.TypeOf(&coreV1.Secret{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err))
	}

	pvcAdapter, err := GetAdapterFor(reflect.TypeOf(&coreV1.PersistentVolumeClaim{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err))
	}

	saAdapter, err := GetAdapterFor(reflect.TypeOf(&coreV1.ServiceAccount{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err))
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

		for _, volume := range resource.Spec.Volumes {
			if volume.ConfigMap != nil && configMapAdapter != nil {
				configMapAdapter.Connect(statefulGraph, resourceNode, volume.ConfigMap.Name)
			} else if volume.Secret != nil && secretAdapter != nil {
				secretAdapter.Connect(statefulGraph, resourceNode, volume.Secret.SecretName)
			} else if volume.PersistentVolumeClaim != nil && pvcAdapter != nil {
				pvcAdapter.Connect(statefulGraph, resourceNode, volume.PersistentVolumeClaim.ClaimName)
			} else if projectedVolume := volume.Projected; projectedVolume != nil {
				for _, projectionSource := range projectedVolume.Sources {
					if projectionSource.ConfigMap != nil && configMapAdapter != nil {
						configMapAdapter.Connect(statefulGraph, resourceNode, projectionSource.ConfigMap.Name)
					} else if projectionSource.Secret != nil && secretAdapter != nil {
						secretAdapter.Connect(statefulGraph, resourceNode, projectionSource.Secret.Name)
					}
				}
			}
		}

		if resource.Spec.ServiceAccountName != "" && saAdapter != nil {
			saAdapter.Connect(statefulGraph, resourceNode, resource.Spec.ServiceAccountName)
		}
	}

	return nil
}
