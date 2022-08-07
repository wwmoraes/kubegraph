package v2beta1

import (
	"fmt"
	"os"

	"github.com/wwmoraes/kubegraph/internal/registry"
)

// Configure connects the resources on this adapter with its dependencies
func (this *HorizontalPodAutoscalerAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	deploymentAdapter, err := GetDeploymentAdapter()
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err)
	}

	objects, err := this.GetGraphObjects(statefulGraph)
	if err != nil {
		return err
	}
	for name, resource := range objects {
		resourceNode, err := this.GetGraphNode(statefulGraph, name)
		if err != nil {
			return err
		}

		if resource.Spec.ScaleTargetRef.Kind == "Deployment" && resource.Spec.ScaleTargetRef.APIVersion == "apps/v1" {
			_, err := deploymentAdapter.Connect(statefulGraph, resourceNode, resource.Spec.ScaleTargetRef.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
			}
		} else {
			fmt.Fprintf(os.Stderr, "unknown scaleRef %s.%s\n", resource.Spec.ScaleTargetRef.APIVersion, resource.Spec.ScaleTargetRef.Kind)
		}
	}
	return nil
}
