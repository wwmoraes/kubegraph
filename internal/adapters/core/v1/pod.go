package v1

import (
	"fmt"
	"log"

	"github.com/wwmoraes/kubegraph/internal/registry"
)

// Configure connects the resources on thisAdapter adapter with its dependencies
func (this *PodAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	configMapAdapter, err := GetConfigMapAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err))
	}

	secretAdapter, err := GetSecretAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err))
	}

	pvcAdapter, err := GetPersistentVolumeClaimAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err))
	}

	saAdapter, err := GetServiceAccountAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err))
	}

	objects, err := this.GetGraphObjects(statefulGraph)
	if err != nil {
		return err
	}

	for name, pod := range objects {
		resourceNode, err := this.GetGraphNode(statefulGraph, name)
		if err != nil {
			return err
		}

		for _, volume := range pod.Spec.Volumes {
			if volume.ConfigMap != nil && configMapAdapter != nil {
				_, err := configMapAdapter.Connect(statefulGraph, resourceNode, volume.ConfigMap.Name)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
				}
			} else if volume.Secret != nil && secretAdapter != nil {
				_, err := secretAdapter.Connect(statefulGraph, resourceNode, volume.Secret.SecretName)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
				}
			} else if volume.PersistentVolumeClaim != nil && pvcAdapter != nil {
				_, err := pvcAdapter.Connect(statefulGraph, resourceNode, volume.PersistentVolumeClaim.ClaimName)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
				}
			} else if projectedVolume := volume.Projected; projectedVolume != nil {
				for _, projectionSource := range projectedVolume.Sources {
					if projectionSource.ConfigMap != nil && configMapAdapter != nil {
						_, err := configMapAdapter.Connect(statefulGraph, resourceNode, projectionSource.ConfigMap.Name)
						if err != nil {
							fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
						}
					} else if projectionSource.Secret != nil && secretAdapter != nil {
						_, err := secretAdapter.Connect(statefulGraph, resourceNode, projectionSource.Secret.Name)
						if err != nil {
							fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
						}
					}
				}
			}
		}

		if pod.Spec.ServiceAccountName != "" && saAdapter != nil {
			_, err := saAdapter.Connect(statefulGraph, resourceNode, pod.Spec.ServiceAccountName)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
			}
		}
	}

	return nil
}
