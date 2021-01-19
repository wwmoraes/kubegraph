package v2beta1

import (
	"fmt"
	"os"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	appsV1 "k8s.io/api/apps/v1"
	autoscalingV2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type horizontalPodAutoscalerAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&horizontalPodAutoscalerAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&autoscalingV2beta1.HorizontalPodAutoscaler{}),
			"icons/hpa.svg",
		),
	})
}

func (thisAdapter *horizontalPodAutoscalerAdapter) tryCastObject(obj runtime.Object) (*autoscalingV2beta1.HorizontalPodAutoscaler, error) {
	casted, ok := obj.(*autoscalingV2beta1.HorizontalPodAutoscaler)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *horizontalPodAutoscalerAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	deploymentAdapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&appsV1.Deployment{}))
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

		if resource.Spec.ScaleTargetRef.Kind == "Deployment" && resource.Spec.ScaleTargetRef.APIVersion == "apps/v1" {
			_, err := deploymentAdapter.Connect(statefulGraph, resourceNode, resource.Spec.ScaleTargetRef.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
			}
		} else {
			fmt.Fprintf(os.Stderr, "unknown scaleRef %s.%s\n", resource.Spec.ScaleTargetRef.APIVersion, resource.Spec.ScaleTargetRef.Kind)
		}
	}
	return nil
}
