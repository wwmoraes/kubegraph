// Code generated by Kubegraph; DO NOT EDIT.

package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	rbacV1 "k8s.io/api/rbac/v1"
)

// RoleBindingObject alias for rbacV1.RoleBinding
type RoleBindingObject = rbacV1.RoleBinding

// RoleBindingAdapterObjectType reflected type of *rbacV1.RoleBinding
var RoleBindingAdapterObjectType = reflect.TypeOf((*rbacV1.RoleBinding)(nil))

// GetRoleBindingAdapter retrieves *RoleBindingAdapter from the registry
func GetRoleBindingAdapter() (*RoleBindingAdapter, error) {
	adapter, err := registry.GetAdapter[*rbacV1.RoleBinding]()
	if err != nil {
		return nil, err
	}
	casted, ok := adapter.(*RoleBindingAdapter)
	if !ok {
		return nil, fmt.Errorf("unable get adapter: %w", registry.ErrIncompatibleType)
	}
	return casted, nil
}

// RoleBindingAdapter implements an adapter for rbacV1.RoleBinding
type RoleBindingAdapter struct {
	registry.Adapter
}

// CastObject casts a runtime.Object to *rbacV1.RoleBinding
func (this *RoleBindingAdapter) CastObject(obj registry.RuntimeObject) (*rbacV1.RoleBinding, error) {
	casted, ok := obj.(*rbacV1.RoleBinding)
	if !ok {
		return casted, fmt.Errorf("unable convert object: %w", registry.ErrIncompatibleType)
	}

	return casted, nil
}

// GetGraphObjects returns a map of all *rbacV1.RoleBinding
// nodes stored on a graph
func (this *RoleBindingAdapter) GetGraphObjects(statefulGraph registry.StatefulGraph) (map[string]*rbacV1.RoleBinding, error) {
	objects, err := statefulGraph.GetObjects(RoleBindingAdapterObjectType)
	if err != nil {
		return nil, err
	}

	castedObjects := make(map[string]*rbacV1.RoleBinding, len(objects))
	for key, object := range objects {
		casted, ok := object.(*rbacV1.RoleBinding)
		if !ok {
			return nil, fmt.Errorf("unable convert object: %w", registry.ErrIncompatibleType)
		}
		castedObjects[key] = casted
	}
	return castedObjects, nil
}

// GetGraphNode returns the node value representing a *rbacV1.RoleBinding
func (this *RoleBindingAdapter) GetGraphNode(statefulGraph registry.StatefulGraph, name string) (registry.Node, error) {
	return statefulGraph.GetNode(RoleBindingAdapterObjectType, name)
}

func (this *RoleBindingAdapter) AddStyledNode(graph registry.StatefulGraph, obj registry.RuntimeObject) (registry.Node, error) {
	resource, err := this.CastObject(obj)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return graph.AddStyledNode(this.GetType(), obj, name, resource.Name, "icons/rb.svg")
}

// init registers RoleBindingAdapter
func init() {
	registry.MustRegister(&RoleBindingAdapter{
		registry.NewAdapter(
			RoleBindingAdapterObjectType,
			"icons/rb.svg",
		),
	})
}
