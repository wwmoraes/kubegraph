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

type adapterRbacV1ClusterRoleBinding struct{}

func init() {
	RegisterResourceAdapter(&adapterRbacV1ClusterRoleBinding{})
}

func (adapter adapterRbacV1ClusterRoleBinding) GetType() reflect.Type {
	return reflect.TypeOf(&rbacv1.ClusterRoleBinding{})
}

func (adapter adapterRbacV1ClusterRoleBinding) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacv1.ClusterRoleBinding)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/crb.svg")
}

func (adapter adapterRbacV1ClusterRoleBinding) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterRbacV1ClusterRoleBinding) Configure(kgraph KubeGraph) error {
	clusterRoleV1Adapter := adapters[reflect.TypeOf(&rbacv1.ClusterRole{})]
	clusterRoleV1beta1Adapter := adapters[reflect.TypeOf(&rbacv1beta1.ClusterRole{})]

	for resourceName, resourceObject := range kgraph.objects[adapter.GetType()] {
		resource := resourceObject.(*rbacv1.ClusterRoleBinding)
		resourceNode, ok := kgraph.nodes[adapter.GetType()][resourceName]
		if !ok {
			return fmt.Errorf("node %s not found", resourceName)
		}

		clusterRoleV1Adapter.Connect(kgraph, resourceNode, resource.RoleRef.Name)
		clusterRoleV1beta1Adapter.Connect(kgraph, resourceNode, resource.RoleRef.Name)
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
