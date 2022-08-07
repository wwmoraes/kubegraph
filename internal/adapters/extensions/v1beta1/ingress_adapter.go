// Code generated by Kubegraph; DO NOT EDIT.

package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	extensionsV1beta1 "k8s.io/api/extensions/v1beta1"
)

// IngressObject alias for extensionsV1beta1.Ingress
type IngressObject = extensionsV1beta1.Ingress

// IngressAdapterObjectType reflected type of *extensionsV1beta1.Ingress
var IngressAdapterObjectType = reflect.TypeOf((*extensionsV1beta1.Ingress)(nil))

// GetIngressAdapter retrieves *IngressAdapter from the registry
func GetIngressAdapter() (*IngressAdapter, error) {
	adapter, err := registry.GetAdapter[*extensionsV1beta1.Ingress]()
	if err != nil {
		return nil, err
	}
	casted, ok := adapter.(*IngressAdapter)
	if !ok {
		return nil, fmt.Errorf("unable get adapter: %w", registry.ErrIncompatibleType)
	}
	return casted, nil
}

// IngressAdapter implements an adapter for extensionsV1beta1.Ingress
type IngressAdapter struct {
	registry.Adapter
}

// CastObject casts a runtime.Object to *extensionsV1beta1.Ingress
func (this *IngressAdapter) CastObject(obj registry.RuntimeObject) (*extensionsV1beta1.Ingress, error) {
	casted, ok := obj.(*extensionsV1beta1.Ingress)
	if !ok {
		return casted, fmt.Errorf("unable convert object: %w", registry.ErrIncompatibleType)
	}

	return casted, nil
}

// GetGraphObjects returns a map of all *extensionsV1beta1.Ingress
// nodes stored on a graph
func (this *IngressAdapter) GetGraphObjects(statefulGraph registry.StatefulGraph) (map[string]*extensionsV1beta1.Ingress, error) {
	objects, err := statefulGraph.GetObjects(IngressAdapterObjectType)
	if err != nil {
		return nil, err
	}

	castedObjects := make(map[string]*extensionsV1beta1.Ingress, len(objects))
	for key, object := range objects {
		casted, ok := object.(*extensionsV1beta1.Ingress)
		if !ok {
			return nil, fmt.Errorf("unable convert object: %w", registry.ErrIncompatibleType)
		}
		castedObjects[key] = casted
	}
	return castedObjects, nil
}

// GetGraphNode returns the node value representing a *extensionsV1beta1.Ingress
func (this *IngressAdapter) GetGraphNode(statefulGraph registry.StatefulGraph, name string) (registry.Node, error) {
	return statefulGraph.GetNode(IngressAdapterObjectType, name)
}

func (this *IngressAdapter) AddStyledNode(graph registry.StatefulGraph, obj registry.RuntimeObject) (registry.Node, error) {
	resource, err := this.CastObject(obj)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return graph.AddStyledNode(this.GetType(), obj, name, resource.Name, "icons/ing.svg")
}

// init registers IngressAdapter
func init() {
	registry.MustRegister(&IngressAdapter{
		registry.NewAdapter(
			IngressAdapterObjectType,
			"icons/ing.svg",
		),
	})
}
