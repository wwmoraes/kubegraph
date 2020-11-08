package v1

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	"k8s.io/apimachinery/pkg/runtime"
	apiregistrationV1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
)

type apiServiceAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&apiServiceAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&apiregistrationV1beta1.APIService{}),
		},
	})
}

func (thisAdapter *apiServiceAdapter) tryCastObject(obj runtime.Object) (*apiregistrationV1beta1.APIService, error) {
	casted, ok := obj.(*apiregistrationV1beta1.APIService)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *apiServiceAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *apiServiceAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/unknown.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *apiServiceAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *apiServiceAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	return nil
}
