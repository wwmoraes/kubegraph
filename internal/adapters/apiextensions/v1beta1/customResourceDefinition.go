package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	apiExtensionsV1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type customResourceDefinitionAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&customResourceDefinitionAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&apiExtensionsV1beta1.CustomResourceDefinition{}),
		},
	})
}

func (thisAdapter *customResourceDefinitionAdapter) tryCastObject(obj runtime.Object) (*apiExtensionsV1beta1.CustomResourceDefinition, error) {
	casted, ok := obj.(*apiExtensionsV1beta1.CustomResourceDefinition)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *customResourceDefinitionAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *customResourceDefinitionAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (adapter.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/crd.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *customResourceDefinitionAdapter) Connect(statefulGraph adapter.StatefulGraph, source adapter.Node, targetName string) (adapter.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *customResourceDefinitionAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	return nil
}
