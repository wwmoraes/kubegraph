package v1

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	admissionregistrationV1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type validatingWebhookConfigurationAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&validatingWebhookConfigurationAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&admissionregistrationV1beta1.ValidatingWebhookConfiguration{}),
		},
	})
}

func (thisAdapter validatingWebhookConfigurationAdapter) tryCastObject(obj runtime.Object) (*admissionregistrationV1beta1.ValidatingWebhookConfiguration, error) {
	casted, ok := obj.(*admissionregistrationV1beta1.ValidatingWebhookConfiguration)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter validatingWebhookConfigurationAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter validatingWebhookConfigurationAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/unknown.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter validatingWebhookConfigurationAdapter) Connect(statefulGraph adapter.StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter validatingWebhookConfigurationAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	return nil
}
