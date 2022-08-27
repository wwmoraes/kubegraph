// Code generated by Kubegraph; DO NOT EDIT.

package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	admissionregistrationV1beta1 "k8s.io/api/admissionregistration/v1beta1"
)

// MutatingWebhookConfigurationObject alias for admissionregistrationV1beta1.MutatingWebhookConfiguration
type MutatingWebhookConfigurationObject = admissionregistrationV1beta1.MutatingWebhookConfiguration

// MutatingWebhookConfigurationAdapterObjectType reflected type of *admissionregistrationV1beta1.MutatingWebhookConfiguration
var MutatingWebhookConfigurationAdapterObjectType = reflect.TypeOf((*admissionregistrationV1beta1.MutatingWebhookConfiguration)(nil))

// GetMutatingWebhookConfigurationAdapter retrieves *MutatingWebhookConfigurationAdapter from the registry
func GetMutatingWebhookConfigurationAdapter() (*MutatingWebhookConfigurationAdapter, error) {
	adapter, err := registry.GetAdapter[*admissionregistrationV1beta1.MutatingWebhookConfiguration]()
	if err != nil {
		return nil, err
	}
	casted, ok := adapter.(*MutatingWebhookConfigurationAdapter)
	if !ok {
		return nil, fmt.Errorf("unable get adapter: %w", registry.ErrIncompatibleType)
	}
	return casted, nil
}

// MutatingWebhookConfigurationAdapter implements an adapter for admissionregistrationV1beta1.MutatingWebhookConfiguration
type MutatingWebhookConfigurationAdapter struct {
	registry.Adapter
}

// CastObject casts a runtime.Object to *admissionregistrationV1beta1.MutatingWebhookConfiguration
func (this *MutatingWebhookConfigurationAdapter) CastObject(obj registry.RuntimeObject) (*admissionregistrationV1beta1.MutatingWebhookConfiguration, error) {
	casted, ok := obj.(*admissionregistrationV1beta1.MutatingWebhookConfiguration)
	if !ok {
		return casted, fmt.Errorf("unable convert object: %w", registry.ErrIncompatibleType)
	}

	return casted, nil
}

// GetGraphObjects returns a map of all *admissionregistrationV1beta1.MutatingWebhookConfiguration
// nodes stored on a graph
func (this *MutatingWebhookConfigurationAdapter) GetGraphObjects(statefulGraph registry.StatefulGraph) (map[string]*admissionregistrationV1beta1.MutatingWebhookConfiguration, error) {
	objects, err := statefulGraph.GetObjects(MutatingWebhookConfigurationAdapterObjectType)
	if err != nil {
		return nil, err
	}

	castedObjects := make(map[string]*admissionregistrationV1beta1.MutatingWebhookConfiguration, len(objects))
	for key, object := range objects {
		casted, ok := object.(*admissionregistrationV1beta1.MutatingWebhookConfiguration)
		if !ok {
			return nil, fmt.Errorf("unable convert object: %w", registry.ErrIncompatibleType)
		}
		castedObjects[key] = casted
	}
	return castedObjects, nil
}

// GetGraphNode returns the node value representing a *admissionregistrationV1beta1.MutatingWebhookConfiguration
func (this *MutatingWebhookConfigurationAdapter) GetGraphNode(statefulGraph registry.StatefulGraph, name string) (registry.Node, error) {
	return statefulGraph.GetNode(MutatingWebhookConfigurationAdapterObjectType, name)
}

func (this *MutatingWebhookConfigurationAdapter) AddStyledNode(graph registry.StatefulGraph, obj registry.RuntimeObject) (registry.Node, error) {
	resource, err := this.CastObject(obj)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return graph.AddStyledNode(this.GetType(), obj, name, resource.Name, "icons/unknown.svg")
}

// init registers MutatingWebhookConfigurationAdapter
func init() {
	registry.MustRegister(&MutatingWebhookConfigurationAdapter{
		registry.NewAdapter(
			MutatingWebhookConfigurationAdapterObjectType,
			"icons/unknown.svg",
		),
	})
}
