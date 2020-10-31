package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	v1 "k8s.io/api/core/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterCoreV1Pod struct{}

func init() {
	RegisterResourceAdapter(&adapterCoreV1Pod{})
}

func (adapter adapterCoreV1Pod) GetType() reflect.Type {
	return reflect.TypeOf(&v1.Pod{})
}

func (adapter adapterCoreV1Pod) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*v1.Pod)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/pod.svg")
}

func (adapter adapterCoreV1Pod) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterCoreV1Pod) Configure(kgraph KubeGraph) error {
	for podName, podObject := range kgraph.objects[adapter.GetType()] {
		pod := podObject.(*v1.Pod)
		podNode, ok := kgraph.nodes[adapter.GetType()][podName]
		if !ok {
			return fmt.Errorf("node %s not found", podName)
		}

		// TODO check if adapter is available
		configMapAdapter := adapters[reflect.TypeOf(&v1.ConfigMap{})]
		secretAdapter := adapters[reflect.TypeOf(&v1.Secret{})]
		pvcAdapter := adapters[reflect.TypeOf(&v1.PersistentVolumeClaim{})]

		for _, volume := range pod.Spec.Volumes {
			if volume.ConfigMap != nil && configMapAdapter != nil {
				configMapAdapter.Connect(kgraph, podNode, volume.ConfigMap.Name)
			} else if volume.Secret != nil && secretAdapter != nil {
				secretAdapter.Connect(kgraph, podNode, volume.Secret.SecretName)
			} else if volume.PersistentVolumeClaim != nil && pvcAdapter != nil {
				pvcAdapter.Connect(kgraph, podNode, volume.PersistentVolumeClaim.ClaimName)
			} else if projectedVolume := volume.Projected; projectedVolume != nil {
				for _, projectionSource := range projectedVolume.Sources {
					if projectionSource.ConfigMap != nil && configMapAdapter != nil {
						configMapAdapter.Connect(kgraph, podNode, projectionSource.ConfigMap.Name)
					} else if projectionSource.Secret != nil && secretAdapter != nil {
						secretAdapter.Connect(kgraph, podNode, projectionSource.Secret.Name)
					}
				}
			}
		}

		if pod.Spec.ServiceAccountName != "" {
			saAdapter := adapters[reflect.TypeOf(&v1.ServiceAccount{})]
			saAdapter.Connect(kgraph, podNode, pod.Spec.ServiceAccountName)
		}
	}

	return nil
}
