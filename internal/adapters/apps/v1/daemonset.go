package v1

import (
	"fmt"
	"log"

	"github.com/wwmoraes/kubegraph/internal/registry"
)

// Create add a graph node for the given object and stores it for further actions
func (adapter *DaemonSetAdapter) Create(graph registry.StatefulGraph, obj registry.RuntimeObject) (registry.Node, error) {
	resource, err := adapter.CastObject(obj)
	if err != nil {
		return nil, err
	}

	resourceNode, err := adapter.AddStyledNode(graph, obj)
	if err != nil {
		return nil, err
	}

	podAdapter, err := GetPodAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", adapter.GetType().String(), err))
	} else {
		podMetadata := resource.Spec.Template.ObjectMeta
		podMetadata.Name = resource.Name
		_, err := podAdapter.Create(graph, &PodObject{
			ObjectMeta: podMetadata,
			Spec:       resource.Spec.Template.Spec,
		})
		if err != nil {
			fmt.Println(fmt.Errorf("%s create error: %w", adapter.GetType().String(), err))
		}
	}

	return resourceNode, nil
}

// Configure connects the resources on this adapter with its dependencies
func (adapter *DaemonSetAdapter) Configure(graph registry.StatefulGraph) error {
	podAdapter, err := GetPodAdapter()
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %w", adapter.GetType().String(), err)
	}

	objects, err := adapter.GetGraphObjects(graph)
	if err != nil {
		return err
	}

	for name, resource := range objects {
		resourceNode, err := adapter.GetGraphNode(graph, name)
		if err != nil {
			return err
		}

		_, err = podAdapter.Connect(graph, resourceNode, resource.Name)
		if err != nil {
			fmt.Println(fmt.Errorf("%s configure error: %w", adapter.GetType().String(), err))
		}
	}
	return nil
}
