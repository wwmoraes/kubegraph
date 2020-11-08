package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	"k8s.io/apimachinery/pkg/runtime"
	applicationV1beta1 "sigs.k8s.io/application/api/v1beta1"
)

// applicationAdapter a kubegraph adapter to render an specific kubernetes resource
type applicationAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&applicationAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&applicationV1beta1.Application{}),
		},
	})
}

func (thisAdapter *applicationAdapter) tryCastObject(obj runtime.Object) (*applicationV1beta1.Application, error) {
	casted, ok := obj.(*applicationV1beta1.Application)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *applicationAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *applicationAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/unknown.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *applicationAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *applicationAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	return nil
}
