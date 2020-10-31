package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	v1 "k8s.io/api/core/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterCoreV1PersistentVolumeClaim struct{}

func init() {
	RegisterResourceAdapter(&adapterCoreV1PersistentVolumeClaim{})
}

func (adapter adapterCoreV1PersistentVolumeClaim) GetType() reflect.Type {
	return reflect.TypeOf(&v1.PersistentVolumeClaim{})
}

func (adapter adapterCoreV1PersistentVolumeClaim) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*v1.PersistentVolumeClaim)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/pvc.svg")
}

func (adapter adapterCoreV1PersistentVolumeClaim) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterCoreV1PersistentVolumeClaim) Configure(kgraph KubeGraph) error {
	for resourceName, resourceObject := range kgraph.objects[adapter.GetType()] {
		resource := resourceObject.(*v1.PersistentVolumeClaim)
		resourceNode, ok := kgraph.nodes[adapter.GetType()][resourceName]
		if !ok {
			return fmt.Errorf("node %s not found", resourceName)
		}

		// TODO check if adapter is available
		volumeAdapter := adapters[reflect.TypeOf(&v1.Volume{})]

		if resource.Spec.VolumeName != "" {
			volumeAdapter.Connect(kgraph, resourceNode, resource.Spec.VolumeName)
		}
	}

	return nil
}
