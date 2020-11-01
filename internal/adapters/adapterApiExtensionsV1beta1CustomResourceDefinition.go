package adapters

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	apiExtensionsV1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterAPIExtensionsV1beta1CustomResourceDefinition struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterAPIExtensionsV1beta1CustomResourceDefinition{
		Resource{
			resourceType: reflect.TypeOf(&apiExtensionsV1beta1.CustomResourceDefinition{}),
		},
	})
}

func (adapter adapterAPIExtensionsV1beta1CustomResourceDefinition) tryCastObject(obj runtime.Object) (*apiExtensionsV1beta1.CustomResourceDefinition, error) {
	casted, ok := obj.(*apiExtensionsV1beta1.CustomResourceDefinition)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterAPIExtensionsV1beta1CustomResourceDefinition) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterAPIExtensionsV1beta1CustomResourceDefinition) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/crd.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterAPIExtensionsV1beta1CustomResourceDefinition) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterAPIExtensionsV1beta1CustomResourceDefinition) Configure(statefulGraph StatefulGraph) error {
	log.Printf("please implement a configuration for %s resources", adapter.GetType().String())
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
		log.Printf("nothing to configure for %s, node %s", resource.Name, resourceNode.Name())
	}
	return nil
}
