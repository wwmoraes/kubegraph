package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterRbacV1beta1RoleBinding struct{}

func init() {
	RegisterResourceAdapter(&adapterRbacV1beta1RoleBinding{})
}

func (adapter adapterRbacV1beta1RoleBinding) GetType() reflect.Type {
	return reflect.TypeOf(&rbacv1beta1.RoleBinding{})
}

func (adapter adapterRbacV1beta1RoleBinding) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacv1beta1.RoleBinding)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/rb.svg")
}

func (adapter adapterRbacV1beta1RoleBinding) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterRbacV1beta1RoleBinding) Configure(kgraph KubeGraph) error {
	roleV1Adapter := adapters[reflect.TypeOf(&rbacv1.Role{})]

	for resourceName, resourceObject := range kgraph.objects[adapter.GetType()] {
		resource := resourceObject.(*rbacv1beta1.RoleBinding)
		resourceNode, ok := kgraph.nodes[adapter.GetType()][resourceName]
		if !ok {
			return fmt.Errorf("node %s not found", resourceName)
		}

		roleV1Adapter.Connect(kgraph, resourceNode, resource.RoleRef.Name)
		// TODO log warning about CRB trying to use Role

		for _, subject := range resource.Subjects {
			if subject.Kind == "ServiceAccount" {
				saNode := kgraph.nodes[reflect.TypeOf(&v1.ServiceAccount{})][subject.Name]
				adapter.Connect(kgraph, saNode, resourceName)
			}
		}
	}
	return nil
}
