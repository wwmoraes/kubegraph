package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	v1 "k8s.io/api/core/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterCoreV1Service struct{}

func init() {
	RegisterResourceAdapter(&adapterCoreV1Service{})
}

func (adapter adapterCoreV1Service) GetType() reflect.Type {
	return reflect.TypeOf(&v1.Service{})
}

func (adapter adapterCoreV1Service) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*v1.Service)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/svc.svg")
}

func (adapter adapterCoreV1Service) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterCoreV1Service) Configure(kgraph KubeGraph) error {
	for resourceName, resourceObject := range kgraph.objects[adapter.GetType()] {
		resource := resourceObject.(*v1.Service)
		resourceNode, ok := kgraph.nodes[adapter.GetType()][resourceName]
		if !ok {
			return fmt.Errorf("node %s not found", resourceName)
		}

		podAdapter := adapters[reflect.TypeOf(&v1.Pod{})]

		for podName, podObject := range kgraph.objects[reflect.TypeOf(&v1.Pod{})] {
			pod := podObject.(*v1.Pod)

			if matchLabels(resource.Spec.Selector, pod.Labels) {
				podAdapter.Connect(kgraph, resourceNode, podName)
			}
		}

	}
	return nil
}
