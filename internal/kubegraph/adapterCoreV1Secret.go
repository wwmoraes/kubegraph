package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	v1 "k8s.io/api/core/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterCoreV1Secret struct{}

func init() {
	RegisterResourceAdapter(&adapterCoreV1Secret{})
}

func (adapter adapterCoreV1Secret) GetType() reflect.Type {
	return reflect.TypeOf(&v1.Secret{})
}

func (adapter adapterCoreV1Secret) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*v1.Secret)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/secret.svg")
}

func (adapter adapterCoreV1Secret) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterCoreV1Secret) Configure(kgraph KubeGraph) error {
	return nil
}
