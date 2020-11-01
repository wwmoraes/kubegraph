package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterRbacV1Role struct{}

func init() {
	RegisterResourceAdapter(&adapterRbacV1Role{})
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterRbacV1Role) GetType() reflect.Type {
	return reflect.TypeOf(&rbacV1.Role{})
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterRbacV1Role) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource := obj.(*rbacV1.Role)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/role.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterRbacV1Role) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterRbacV1Role) Configure(statefulGraph StatefulGraph) error {
	return nil
}
