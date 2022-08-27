package v1

import (
	"fmt"
	"log"

	"github.com/wwmoraes/kubegraph/internal/registry"
)

// Create add a graph node for the given object and stores it for further actions
func (this *StatefulSetAdapter) Create(graph registry.StatefulGraph, obj registry.RuntimeObject) (registry.Node, error) {
	resource, err := this.CastObject(obj)
	if err != nil {
		return nil, err
	}

	resourceNode, err := this.AddStyledNode(graph, obj)
	if err != nil {
		return nil, err
	}

	podAdapter, err := GetPodAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err))
	} else {
		podMetadata := resource.Spec.Template.ObjectMeta
		podMetadata.Name = resource.Name
		_, err := podAdapter.Create(graph, &PodObject{
			ObjectMeta: podMetadata,
			Spec:       resource.Spec.Template.Spec,
		})
		if err != nil {
			fmt.Println(fmt.Errorf("%s create error: %w", this.GetType().String(), err))
		}
	}

	return resourceNode, nil
}

// Configure connects the resources on this adapter with its dependencies
func (this *StatefulSetAdapter) Configure(graph registry.StatefulGraph) error {
	podAdapter, err := GetPodAdapter()
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err)
	}

	objects, err := graph.GetObjects(this.GetType())
	if err != nil {
		return err
	}

	for resourceName, resourceObject := range objects {
		resource, err := this.CastObject(resourceObject)
		if err != nil {
			return err
		}
		resourceNode, err := graph.GetNode(this.GetType(), resourceName)
		if err != nil {
			return err
		}

		_, err = podAdapter.Connect(graph, resourceNode, resource.Name)
		if err != nil {
			fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
		}
	}
	return nil
}
