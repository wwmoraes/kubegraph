package v2beta1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	"github.com/wwmoraes/kubegraph/internal/adapter"
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
func (thisAdapter horizontalPodAutoscalerAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/hpa.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter horizontalPodAutoscalerAdapter) Connect(statefulGraph adapter.StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter horizontalPodAutoscalerAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
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

		// do something with each resource
		log.Printf("%s resource %s, node %s", thisAdapter.GetType().String(), resource.Name, resourceNode.Name())
	}
	return nil
}
