package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterRbacV1beta1Role struct{}

func init() {
	RegisterResourceAdapter(&adapterRbacV1beta1Role{})
}

func (adapter adapterRbacV1beta1Role) GetType() reflect.Type {
	return reflect.TypeOf(&rbacv1beta1.Role{})
}

func (adapter adapterRbacV1beta1Role) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacv1beta1.Role)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/role.svg")
}

func (adapter adapterRbacV1beta1Role) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterRbacV1beta1Role) Configure(kgraph KubeGraph) error {
	return nil
}
