package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	rbacv1 "k8s.io/api/rbac/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterRbacV1ClusterRole struct{}

func init() {
	RegisterResourceAdapter(&adapterRbacV1ClusterRole{})
}

func (adapter adapterRbacV1ClusterRole) GetType() reflect.Type {
	return reflect.TypeOf(&rbacv1.ClusterRole{})
}

func (adapter adapterRbacV1ClusterRole) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacv1.ClusterRole)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/c-role.svg")
}

func (adapter adapterRbacV1ClusterRole) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterRbacV1ClusterRole) Configure(kgraph KubeGraph) error {
	return nil
}
