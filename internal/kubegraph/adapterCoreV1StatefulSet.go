package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterAppsV1StatefulSet struct{}

func init() {
	RegisterResourceAdapter(&adapterAppsV1StatefulSet{})
}

func (adapter adapterAppsV1StatefulSet) GetType() reflect.Type {
	return reflect.TypeOf(&appsv1.StatefulSet{})
}

func (adapter adapterAppsV1StatefulSet) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*appsv1.StatefulSet)

	podMetadata := resource.Spec.Template.ObjectMeta
	podMetadata.Name = resource.Name
	adapters[reflect.TypeOf(&v1.Pod{})].Create(kgraph, &v1.Pod{
		ObjectMeta: podMetadata,
		Spec:       resource.Spec.Template.Spec,
	})

	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/sts.svg")
}

func (adapter adapterAppsV1StatefulSet) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterAppsV1StatefulSet) Configure(kgraph KubeGraph) error {
	for resourceName, resourceObject := range kgraph.objects[adapter.GetType()] {
		resource := resourceObject.(*appsv1.StatefulSet)
		resourceNode, ok := kgraph.nodes[adapter.GetType()][resourceName]
		if !ok {
			return fmt.Errorf("node %s not found", resourceName)
		}

		podAdapter := adapters[reflect.TypeOf(&v1.Pod{})]

		podAdapter.Connect(kgraph, resourceNode, resource.Name)
	}
	return nil
}
