package v1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type podAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&podAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&coreV1.Pod{}),
		},
	})
}

func (thisAdapter podAdapter) tryCastObject(obj runtime.Object) (*coreV1.Pod, error) {
	casted, ok := obj.(*coreV1.Pod)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by thisAdapter instance
func (thisAdapter podAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter podAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/pod.svg")
}

// Connect creates and edge between the given node and an object on thisAdapter adapter
func (thisAdapter podAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on thisAdapter adapter with its dependencies
func (thisAdapter podAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	configMapAdapter, err := adapter.Get(reflect.TypeOf(&coreV1.ConfigMap{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	}

	secretAdapter, err := adapter.Get(reflect.TypeOf(&coreV1.Secret{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	}

	pvcAdapter, err := adapter.Get(reflect.TypeOf(&coreV1.PersistentVolumeClaim{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	}

	saAdapter, err := adapter.Get(reflect.TypeOf(&coreV1.ServiceAccount{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
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

		for _, volume := range resource.Spec.Volumes {
			if volume.ConfigMap != nil && configMapAdapter != nil {
				_, err := configMapAdapter.Connect(statefulGraph, resourceNode, volume.ConfigMap.Name)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
				}
			} else if volume.Secret != nil && secretAdapter != nil {
				_, err := secretAdapter.Connect(statefulGraph, resourceNode, volume.Secret.SecretName)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
				}
			} else if volume.PersistentVolumeClaim != nil && pvcAdapter != nil {
				_, err := pvcAdapter.Connect(statefulGraph, resourceNode, volume.PersistentVolumeClaim.ClaimName)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
				}
			} else if projectedVolume := volume.Projected; projectedVolume != nil {
				for _, projectionSource := range projectedVolume.Sources {
					if projectionSource.ConfigMap != nil && configMapAdapter != nil {
						_, err := configMapAdapter.Connect(statefulGraph, resourceNode, projectionSource.ConfigMap.Name)
						if err != nil {
							fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
						}
					} else if projectionSource.Secret != nil && secretAdapter != nil {
						_, err := secretAdapter.Connect(statefulGraph, resourceNode, projectionSource.Secret.Name)
						if err != nil {
							fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
						}
					}
				}
			}
		}

		if resource.Spec.ServiceAccountName != "" && saAdapter != nil {
			_, err := saAdapter.Connect(statefulGraph, resourceNode, resource.Spec.ServiceAccountName)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
			}
		}
	}

	return nil
}
