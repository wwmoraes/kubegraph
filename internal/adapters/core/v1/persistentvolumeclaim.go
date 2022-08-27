package v1

import (
	"fmt"

	"github.com/wwmoraes/kubegraph/internal/registry"
)

// Configure connects the resources on this adapter with its dependencies
func (this *PersistentVolumeClaimAdapter) Configure(graph registry.StatefulGraph) error {
	persistentVolumeAdapter, err := GetPersistentVolumeAdapter()
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err)
	}

	objects, err := this.GetGraphObjects(graph)
	if err != nil {
		return err
	}

	for resourceName, resource := range objects {
		resourceNode, err := this.GetGraphNode(graph, resourceName)
		if err != nil {
			return err
		}

		if resource.Spec.VolumeName != "" {
			_, err := persistentVolumeAdapter.Connect(graph, resourceNode, resource.Spec.VolumeName)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
			}
		}
	}

	return nil
}
