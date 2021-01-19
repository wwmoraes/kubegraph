package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	"github.com/wwmoraes/kubegraph/internal/utils"
	coreV1 "k8s.io/api/core/v1"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type podDisruptionBudgetAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&podDisruptionBudgetAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&policyV1beta1.PodDisruptionBudget{}),
			"icons/pdb.svg",
		),
	})
}

func (thisAdapter *podDisruptionBudgetAdapter) tryCastObject(obj runtime.Object) (*policyV1beta1.PodDisruptionBudget, error) {
	casted, ok := obj.(*policyV1beta1.PodDisruptionBudget)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *podDisruptionBudgetAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	podAdapter, err := registry.Instance().Get(reflect.TypeOf(&coreV1.Pod{}))
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
