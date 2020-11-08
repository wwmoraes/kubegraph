package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	networkV1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
)

type virtualServiceAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&virtualServiceAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&networkV1alpha3.VirtualService{}),
		},
	})
}

func (thisAdapter *virtualServiceAdapter) tryCastObject(obj runtime.Object) (*networkV1alpha3.VirtualService, error) {
	casted, ok := obj.(*networkV1alpha3.VirtualService)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *virtualServiceAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *virtualServiceAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (adapter.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/unknown.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *virtualServiceAdapter) Connect(statefulGraph adapter.StatefulGraph, source adapter.Node, targetName string) (adapter.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *virtualServiceAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	return nil
}
