package v1

import (
	"fmt"

	"github.com/wwmoraes/kubegraph/internal/registry"
	"github.com/wwmoraes/kubegraph/internal/utils"
)

// Configure connects the resources on this adapter with its dependencies
func (this *ServiceAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	podAdapter, err := GetPodAdapter()
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err)
	}

	objects, err := statefulGraph.GetObjects(this.GetType())
	if err != nil {
		return err
	}

	for resourceName, resourceObject := range objects {
		resource, err := this.CastObject(resourceObject)
		if err != nil {
			return err
		}
		resourceNode, err := statefulGraph.GetNode(this.GetType(), resourceName)
		if err != nil {
			return err
		}

		objects, err := podAdapter.GetGraphObjects(statefulGraph)
		if err != nil {
			return err
		}
		for name, pod := range objects {
			if utils.MatchLabels(resource.Spec.Selector, pod.Labels) {
				_, err := podAdapter.Connect(statefulGraph, resourceNode, name)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
				}
			}
		}
	}
	return nil
}
