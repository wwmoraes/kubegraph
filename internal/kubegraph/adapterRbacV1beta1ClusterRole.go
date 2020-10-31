package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterRbacV1beta1ClusterRole struct{}

func init() {
	RegisterResourceAdapter(&adapterRbacV1beta1ClusterRole{})
}

func (adapter adapterRbacV1beta1ClusterRole) GetType() reflect.Type {
	return reflect.TypeOf(&rbacv1beta1.ClusterRole{})
}

func (adapter adapterRbacV1beta1ClusterRole) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacv1beta1.ClusterRole)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/c-role.svg")
}

func (adapter adapterRbacV1beta1ClusterRole) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterRbacV1beta1ClusterRole) Configure(kgraph KubeGraph) error {
	return nil
}
