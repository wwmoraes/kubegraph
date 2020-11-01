package adapters

/*
 * remove the dummy struct and replace the references with a proper kubernetes API resource
 */

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	autoscalingV2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterAutoscalingV2beta1HorizontalPodAutoscaler struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterAutoscalingV2beta1HorizontalPodAutoscaler{
		Resource{
			resourceType: reflect.TypeOf(&autoscalingV2beta1.HorizontalPodAutoscaler{}),
		},
	})
}

func (adapter adapterAutoscalingV2beta1HorizontalPodAutoscaler) tryCastObject(obj runtime.Object) (*autoscalingV2beta1.HorizontalPodAutoscaler, error) {
	casted, ok := obj.(*autoscalingV2beta1.HorizontalPodAutoscaler)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterAutoscalingV2beta1HorizontalPodAutoscaler) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterAutoscalingV2beta1HorizontalPodAutoscaler) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/hpa.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterAutoscalingV2beta1HorizontalPodAutoscaler) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterAutoscalingV2beta1HorizontalPodAutoscaler) Configure(statefulGraph StatefulGraph) error {
	objects, err := statefulGraph.GetObjects(adapter.GetType())
	if err != nil {
		return err
	}
	for resourceName, resourceObject := range objects {
		resource, err := adapter.tryCastObject(resourceObject)
		if err != nil {
			return err
		}
		resourceNode, err := statefulGraph.GetNode(adapter.GetType(), resourceName)
		if err != nil {
			return err
		}

		// do something with each resource
		log.Printf("%s resource %s, node %s", adapter.GetType().String(), resource.Name, resourceNode.Name())
	}
	return nil
}
