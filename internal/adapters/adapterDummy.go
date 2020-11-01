package adapters

/*
 * remove the dummy struct and replace the references with a proper kubernetes API resource
 */

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type dummy struct {
	metaV1.TypeMeta
	metaV1.ObjectMeta
}

func (d dummy) GetObjectKind() schema.ObjectKind {
	return nil
}

func (d dummy) DeepCopyObject() runtime.Object {
	return dummy{}
}

type adapterCoreV1Dummy struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterCoreV1Dummy{
		Resource{
			resourceType: reflect.TypeOf(&dummy{}),
		},
	})
}

func (adapter adapterCoreV1Dummy) tryCastObject(obj runtime.Object) (*dummy, error) {
	casted, ok := obj.(*dummy)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterCoreV1Dummy) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterCoreV1Dummy) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/unknown.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterCoreV1Dummy) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterCoreV1Dummy) Configure(statefulGraph StatefulGraph) error {
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
