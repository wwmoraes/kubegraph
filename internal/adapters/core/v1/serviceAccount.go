package v1

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type serviceAccountAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&serviceAccountAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&coreV1.ServiceAccount{}),
		},
	})
}

func (thisAdapter *serviceAccountAdapter) tryCastObject(obj runtime.Object) (*coreV1.ServiceAccount, error) {
	casted, ok := obj.(*coreV1.ServiceAccount)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *serviceAccountAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *serviceAccountAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/sa.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *serviceAccountAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *serviceAccountAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	return nil
}
