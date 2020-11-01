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

type adapterRbacV1RoleBinding struct{}

func init() {
	RegisterResourceAdapter(&adapterRbacV1RoleBinding{})
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterRbacV1RoleBinding) GetType() reflect.Type {
	return reflect.TypeOf(&rbacV1.RoleBinding{})
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterRbacV1RoleBinding) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacV1.RoleBinding)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/rb.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterRbacV1RoleBinding) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterRbacV1RoleBinding) Configure(statefulGraph StatefulGraph) error {
	roleV1Adapter, err := GetAdapterFor(reflect.TypeOf(&rbacV1.Role{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err))
	}
	roleV1beta1Adapter, err := GetAdapterFor(reflect.TypeOf(&rbacV1beta1.Role{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err))
	}

	objects, err := statefulGraph.GetObjects(adapter.GetType())
	if err != nil {
		return err
	}
	for resourceName, resourceObject := range objects {
		resource := resourceObject.(*rbacV1.RoleBinding)
		resourceNode, err := statefulGraph.GetNode(adapter.GetType(), resourceName)
		if err != nil {
			return err
		}

		if roleV1Adapter != nil {
			roleV1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
		}
		if roleV1beta1Adapter != nil {
			roleV1beta1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
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
