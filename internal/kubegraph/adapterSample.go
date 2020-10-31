package kubegraph

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type dummy struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

func (d dummy) GetObjectKind() schema.ObjectKind {
	return nil
}

func (d dummy) DeepCopyObject() k8sruntime.Object {
	return dummy{}
}

type adapterCoreV1Dummy struct{}

func init() {
	RegisterResourceAdapter(&adapterCoreV1Dummy{})
}

func (adapter adapterCoreV1Dummy) GetType() reflect.Type {
	return reflect.TypeOf(&dummy{})
}

func (adapter adapterCoreV1Dummy) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*dummy)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/unknown.svg")
}

func (adapter adapterCoreV1Dummy) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterCoreV1Dummy) Configure(kgraph KubeGraph) error {
	for resourceName, resourceObject := range kgraph.objects[adapter.GetType()] {
		resource := resourceObject.(*dummy)
		resourceNode, ok := kgraph.nodes[adapter.GetType()][resourceName]
		if !ok {
			return fmt.Errorf("node %s not found", resourceName)
		}

		// do something with each resource
		log.Printf("dummy resource %s, node %s", resource.Name, resourceNode.Name())
	}
	return nil
}
