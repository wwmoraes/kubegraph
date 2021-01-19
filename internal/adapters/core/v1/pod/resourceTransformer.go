package pod

import (
	"fmt"
	"log"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterResource struct {
	adapter.Resource
}

func (thisAdapter *adapterResource) tryCastObject(obj runtime.Object) (*coreV1.Pod, error) {
	casted, ok := obj.(*coreV1.Pod)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Configure connects the resources on thisAdapter adapter with its dependencies
func (thisAdapter *adapterResource) Configure(statefulGraph adapter.StatefulGraph) error {
	configMapAdapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&coreV1.ConfigMap{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	}

	secretAdapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&coreV1.Secret{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	}

	pvcAdapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&coreV1.PersistentVolumeClaim{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	}

	saAdapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&coreV1.ServiceAccount{}))
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
