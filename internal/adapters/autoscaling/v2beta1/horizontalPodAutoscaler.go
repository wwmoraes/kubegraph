package v2beta1

import (
	"fmt"
	"os"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	appsV1 "k8s.io/api/apps/v1"
	autoscalingV2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type horizontalPodAutoscalerAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&horizontalPodAutoscalerAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&autoscalingV2beta1.HorizontalPodAutoscaler{}),
		},
	})
}

func (thisAdapter horizontalPodAutoscalerAdapter) tryCastObject(obj runtime.Object) (*autoscalingV2beta1.HorizontalPodAutoscaler, error) {
	casted, ok := obj.(*autoscalingV2beta1.HorizontalPodAutoscaler)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter horizontalPodAutoscalerAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter horizontalPodAutoscalerAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/hpa.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter horizontalPodAutoscalerAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter horizontalPodAutoscalerAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	deploymentAdapter, err := adapter.Get(reflect.TypeOf(&appsV1.Deployment{}))
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
			deploymentAdapter.Connect(statefulGraph, resourceNode, resource.Spec.ScaleTargetRef.Name)
		} else {
			fmt.Fprintf(os.Stderr, "unknown scaleRef %s.%s\n", resource.Spec.ScaleTargetRef.APIVersion, resource.Spec.ScaleTargetRef.Kind)
		}
	}
	return nil
}
