package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterRbacV1ClusterRole struct{}

func init() {
	RegisterResourceAdapter(&adapterRbacV1ClusterRole{})
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterRbacV1ClusterRole) GetType() reflect.Type {
	return reflect.TypeOf(&rbacV1.ClusterRole{})
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterRbacV1ClusterRole) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacV1.ClusterRole)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/c-role.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterRbacV1ClusterRole) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterRbacV1ClusterRole) Configure(statefulGraph StatefulGraph) error {
	return nil
}
