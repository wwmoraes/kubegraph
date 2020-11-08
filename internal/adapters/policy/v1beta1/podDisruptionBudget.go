package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	"github.com/wwmoraes/kubegraph/internal/utils"
	coreV1 "k8s.io/api/core/v1"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type podDisruptionBudgetAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&podDisruptionBudgetAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&policyV1beta1.PodDisruptionBudget{}),
		},
	})
}

func (thisAdapter *podDisruptionBudgetAdapter) tryCastObject(obj runtime.Object) (*policyV1beta1.PodDisruptionBudget, error) {
	casted, ok := obj.(*policyV1beta1.PodDisruptionBudget)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *podDisruptionBudgetAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *podDisruptionBudgetAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/pdb.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *podDisruptionBudgetAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *podDisruptionBudgetAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	podAdapter, err := adapter.Get(reflect.TypeOf(&coreV1.Pod{}))
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err)
	}

	objects, err := statefulGraph.GetObjects(thisAdapter.GetType())
	if err != nil {
		return err
	}

	for resourceName, resourceObject := range objects {
		resource, err := thisAdapter.tryCastObject(resourceObject)
		if err != nil {
			return err
		}
		resourceNode, err := statefulGraph.GetNode(thisAdapter.GetType(), resourceName)
		if err != nil {
			return err
		}

		objects, err := statefulGraph.GetObjects(reflect.TypeOf(&coreV1.Pod{}))
		if err != nil {
			return err
		}
		for podName, podObject := range objects {
			pod := podObject.(*coreV1.Pod)

			if utils.MatchLabels(resource.Spec.Selector.MatchLabels, pod.Labels) {
				_, err := podAdapter.Connect(statefulGraph, resourceNode, podName)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
				}
			}
		}
	}
	return nil
}
