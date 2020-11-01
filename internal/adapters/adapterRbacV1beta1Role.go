package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterRbacV1beta1Role struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterRbacV1beta1Role{
		Resource{
			resourceType: reflect.TypeOf(&rbacV1beta1.Role{}),
		},
	})
}

func (adapter adapterRbacV1beta1Role) tryCastObject(obj runtime.Object) (*rbacV1beta1.Role, error) {
	casted, ok := obj.(*rbacV1beta1.Role)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterRbacV1beta1Role) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterRbacV1beta1Role) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/role.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterRbacV1beta1Role) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterRbacV1beta1Role) Configure(statefulGraph StatefulGraph) error {
	return nil
}
