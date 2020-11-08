package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type secretAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&secretAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&coreV1.Secret{}),
		},
	})
}

func (thisAdapter *secretAdapter) tryCastObject(obj runtime.Object) (*coreV1.Secret, error) {
	casted, ok := obj.(*coreV1.Secret)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *secretAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *secretAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (adapter.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/secret.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *secretAdapter) Connect(statefulGraph adapter.StatefulGraph, source adapter.Node, targetName string) (adapter.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *secretAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	return nil
}
