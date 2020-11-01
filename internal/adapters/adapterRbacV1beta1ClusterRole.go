package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterRbacV1beta1ClusterRole struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterRbacV1beta1ClusterRole{
		Resource{
			resourceType: reflect.TypeOf(&rbacV1beta1.ClusterRole{}),
		},
	})
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterRbacV1beta1ClusterRole) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterRbacV1beta1ClusterRole) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacV1beta1.ClusterRole)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/c-role.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterRbacV1beta1ClusterRole) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterRbacV1beta1ClusterRole) Configure(statefulGraph StatefulGraph) error {
	return nil
}
