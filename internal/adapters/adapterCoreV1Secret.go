package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterCoreV1Secret struct{}

func init() {
	RegisterResourceAdapter(&adapterCoreV1Secret{})
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterCoreV1Secret) GetType() reflect.Type {
	return reflect.TypeOf(&coreV1.Secret{})
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterCoreV1Secret) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource := obj.(*coreV1.Secret)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/secret.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterCoreV1Secret) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterCoreV1Secret) Configure(statefulGraph StatefulGraph) error {
	return nil
}
