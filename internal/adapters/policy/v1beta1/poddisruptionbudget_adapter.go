// Code generated by Kubegraph; DO NOT EDIT.

package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
)

// PodDisruptionBudgetObject alias for policyV1beta1.PodDisruptionBudget
type PodDisruptionBudgetObject = policyV1beta1.PodDisruptionBudget

// PodDisruptionBudgetAdapterObjectType reflected type of *policyV1beta1.PodDisruptionBudget
var PodDisruptionBudgetAdapterObjectType = reflect.TypeOf((*policyV1beta1.PodDisruptionBudget)(nil))

// GetPodDisruptionBudgetAdapter retrieves *PodDisruptionBudgetAdapter from the registry
func GetPodDisruptionBudgetAdapter() (*PodDisruptionBudgetAdapter, error) {
	adapter, err := registry.GetAdapter[*policyV1beta1.PodDisruptionBudget]()
	if err != nil {
		return nil, err
	}
	casted, ok := adapter.(*PodDisruptionBudgetAdapter)
	if !ok {
		return nil, fmt.Errorf("unable get adapter: %w", registry.ErrIncompatibleType)
	}
	return casted, nil
}

// PodDisruptionBudgetAdapter implements an adapter for policyV1beta1.PodDisruptionBudget
type PodDisruptionBudgetAdapter struct {
	registry.Adapter
}

// CastObject casts a runtime.Object to *policyV1beta1.PodDisruptionBudget
func (this *PodDisruptionBudgetAdapter) CastObject(obj registry.RuntimeObject) (*policyV1beta1.PodDisruptionBudget, error) {
	casted, ok := obj.(*policyV1beta1.PodDisruptionBudget)
	if !ok {
		return casted, fmt.Errorf("unable convert object: %w", registry.ErrIncompatibleType)
	}

	return casted, nil
}

// GetGraphObjects returns a map of all *policyV1beta1.PodDisruptionBudget
// nodes stored on a graph
func (this *PodDisruptionBudgetAdapter) GetGraphObjects(statefulGraph registry.StatefulGraph) (map[string]*policyV1beta1.PodDisruptionBudget, error) {
	objects, err := statefulGraph.GetObjects(PodDisruptionBudgetAdapterObjectType)
	if err != nil {
		return nil, err
	}

	castedObjects := make(map[string]*policyV1beta1.PodDisruptionBudget, len(objects))
	for key, object := range objects {
		casted, ok := object.(*policyV1beta1.PodDisruptionBudget)
		if !ok {
			return nil, fmt.Errorf("unable convert object: %w", registry.ErrIncompatibleType)
		}
		castedObjects[key] = casted
	}
	return castedObjects, nil
}

// GetGraphNode returns the node value representing a *policyV1beta1.PodDisruptionBudget
func (this *PodDisruptionBudgetAdapter) GetGraphNode(statefulGraph registry.StatefulGraph, name string) (registry.Node, error) {
	return statefulGraph.GetNode(PodDisruptionBudgetAdapterObjectType, name)
}

func (this *PodDisruptionBudgetAdapter) AddStyledNode(graph registry.StatefulGraph, obj registry.RuntimeObject) (registry.Node, error) {
	resource, err := this.CastObject(obj)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return graph.AddStyledNode(this.GetType(), obj, name, resource.Name, "icons/pdb.svg")
}

// init registers PodDisruptionBudgetAdapter
func init() {
	registry.MustRegister(&PodDisruptionBudgetAdapter{
		registry.NewAdapter(
			PodDisruptionBudgetAdapterObjectType,
			"icons/pdb.svg",
		),
	})
}
