package v1beta1

import (
	"fmt"

	"github.com/wwmoraes/kubegraph/internal/registry"
)

// Configure connects the resources on this adapter with its dependencies
func (this *IngressAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	serviceAdapter, err := GetServiceAdapter()
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

		// connects default backend service
		if resource.Spec.Backend != nil && resource.Spec.Backend.ServiceName != "" {
			_, err := serviceAdapter.Connect(statefulGraph, resourceNode, resource.Spec.Backend.ServiceName)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
			}
		}

		// connects rule backends
		for _, rule := range resource.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				if path.Backend.ServiceName != "" {
					_, err := serviceAdapter.Connect(statefulGraph, resourceNode, path.Backend.ServiceName)
					if err != nil {
						fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
					}
				}
			}
		}
	}
	return nil
}
