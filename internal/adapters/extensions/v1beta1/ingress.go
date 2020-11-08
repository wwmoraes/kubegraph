package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	extensionsV1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type IngressAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&IngressAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&extensionsV1beta1.Ingress{}),
		},
	})
}

func (thisAdapter IngressAdapter) tryCastObject(obj runtime.Object) (*extensionsV1beta1.Ingress, error) {
	casted, ok := obj.(*extensionsV1beta1.Ingress)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter IngressAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter IngressAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/ing.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter IngressAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter IngressAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	serviceAdapter, err := adapter.Get(reflect.TypeOf(&coreV1.Service{}))
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

		// connects default backend service
		if resource.Spec.Backend != nil && resource.Spec.Backend.ServiceName != "" {
			_, err := serviceAdapter.Connect(statefulGraph, resourceNode, resource.Spec.Backend.ServiceName)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
			}
		}

		// connects rule backends
		for _, rule := range resource.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				if path.Backend.ServiceName != "" {
					_, err := serviceAdapter.Connect(statefulGraph, resourceNode, path.Backend.ServiceName)
					if err != nil {
						fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
					}
				}
			}
		}
	}
	return nil
}
