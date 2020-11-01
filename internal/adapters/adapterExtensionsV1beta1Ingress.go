package adapters

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	coreV1 "k8s.io/api/core/v1"
	extensionsV1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterExtensionsV1beta1Ingress struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterExtensionsV1beta1Ingress{
		Resource{
			resourceType: reflect.TypeOf(&extensionsV1beta1.Ingress{}),
		},
	})
}

func (adapter adapterExtensionsV1beta1Ingress) tryCastObject(obj runtime.Object) (*extensionsV1beta1.Ingress, error) {
	casted, ok := obj.(*extensionsV1beta1.Ingress)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterExtensionsV1beta1Ingress) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterExtensionsV1beta1Ingress) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/ing.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterExtensionsV1beta1Ingress) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterExtensionsV1beta1Ingress) Configure(statefulGraph StatefulGraph) error {
	serviceAdapter, err := GetAdapterFor(reflect.TypeOf(&coreV1.Service{}))
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

		// connects default backend service
		if resource.Spec.Backend != nil && resource.Spec.Backend.ServiceName != "" {
			serviceAdapter.Connect(statefulGraph, resourceNode, resource.Spec.Backend.ServiceName)
		}

		// connects rule backends
		for _, rule := range resource.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				if path.Backend.ServiceName != "" {
					serviceAdapter.Connect(statefulGraph, resourceNode, path.Backend.ServiceName)
				}
			}
		}
	}
	return nil
}
