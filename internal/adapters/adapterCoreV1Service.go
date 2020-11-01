package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	"github.com/wwmoraes/kubegraph/internal/utils"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterCoreV1Service struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterCoreV1Service{
		Resource{
			resourceType: reflect.TypeOf(&coreV1.Service{}),
		},
	})
}

func (adapter adapterCoreV1Service) tryCastObject(obj runtime.Object) (*coreV1.Service, error) {
	casted, ok := obj.(*coreV1.Service)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterCoreV1Service) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterCoreV1Service) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/svc.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterCoreV1Service) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterCoreV1Service) Configure(statefulGraph StatefulGraph) error {
	podAdapter, err := GetAdapterFor(reflect.TypeOf(&coreV1.Pod{}))
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err)
	}

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

		objects, err := statefulGraph.GetObjects(reflect.TypeOf(&coreV1.Pod{}))
		if err != nil {
			return err
		}
		for podName, podObject := range objects {
			pod := podObject.(*coreV1.Pod)

			if utils.MatchLabels(resource.Spec.Selector, pod.Labels) {
				podAdapter.Connect(statefulGraph, resourceNode, podName)
			}
		}
	}
	return nil
}
