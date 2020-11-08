package v1

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	"github.com/wwmoraes/kubegraph/internal/utils"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type serviceAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&serviceAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&coreV1.Service{}),
		},
	})
}

func (thisAdapter *serviceAdapter) tryCastObject(obj runtime.Object) (*coreV1.Service, error) {
	casted, ok := obj.(*coreV1.Service)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *serviceAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *serviceAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/svc.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *serviceAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *serviceAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	podAdapter, err := adapter.Get(reflect.TypeOf(&coreV1.Pod{}))
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

		objects, err := statefulGraph.GetObjects(reflect.TypeOf(&coreV1.Pod{}))
		if err != nil {
			return err
		}
		for podName, podObject := range objects {
			pod := podObject.(*coreV1.Pod)

			if utils.MatchLabels(resource.Spec.Selector, pod.Labels) {
				_, err := podAdapter.Connect(statefulGraph, resourceNode, podName)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
				}
			}
		}
	}
	return nil
}
