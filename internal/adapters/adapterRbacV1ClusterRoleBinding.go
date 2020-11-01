package adapters

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterRbacV1ClusterRoleBinding struct{}

func init() {
	RegisterResourceAdapter(&adapterRbacV1ClusterRoleBinding{})
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterRbacV1ClusterRoleBinding) GetType() reflect.Type {
	return reflect.TypeOf(&rbacV1.ClusterRoleBinding{})
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterRbacV1ClusterRoleBinding) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacV1.ClusterRoleBinding)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/crb.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterRbacV1ClusterRoleBinding) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterRbacV1ClusterRoleBinding) Configure(statefulGraph StatefulGraph) error {
	clusterRoleV1beta1Adapter, err := GetAdapterFor(reflect.TypeOf(&rbacV1beta1.ClusterRole{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err))
	}
	clusterRoleV1Adapter, err := GetAdapterFor(reflect.TypeOf(&rbacV1.ClusterRole{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err))
	}

	objects, err := statefulGraph.GetObjects(adapter.GetType())
	if err != nil {
		return err
	}
	for resourceName, resourceObject := range objects {
		resource := resourceObject.(*rbacV1.ClusterRoleBinding)
		resourceNode, err := statefulGraph.GetNode(adapter.GetType(), resourceName)
		if err != nil {
			return err
		}

		if clusterRoleV1Adapter != nil {
			clusterRoleV1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
		}
		if clusterRoleV1beta1Adapter != nil {
			clusterRoleV1beta1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
		}

		for _, subject := range resource.Subjects {
			if subject.Kind == "ServiceAccount" {
				saNode, err := statefulGraph.GetNode(reflect.TypeOf(&coreV1.ServiceAccount{}), subject.Name)
				if err != nil {
					return err
				}
				adapter.Connect(statefulGraph, saNode, resourceName)
			}
		}
	}
	return nil
}
