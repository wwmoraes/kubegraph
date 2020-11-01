package adapters

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterAppsV1StatefulSet struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterAppsV1StatefulSet{
		Resource{
			resourceType: reflect.TypeOf(&appsV1.StatefulSet{}),
		},
	})
}

func (adapter adapterAppsV1StatefulSet) tryCastObject(obj runtime.Object) (*appsV1.StatefulSet, error) {
	casted, ok := obj.(*appsV1.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterAppsV1StatefulSet) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterAppsV1StatefulSet) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	resourceNode, err := statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/sts.svg")
	if err != nil {
		return nil, err
	}

	podAdapter, err := GetAdapterFor(reflect.TypeOf(&coreV1.Pod{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err))
	} else {
		podMetadata := resource.Spec.Template.ObjectMeta
		podMetadata.Name = resource.Name
		podAdapter.Create(statefulGraph, &coreV1.Pod{
			ObjectMeta: podMetadata,
			Spec:       resource.Spec.Template.Spec,
		})
	}

	return resourceNode, nil
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterAppsV1StatefulSet) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterAppsV1StatefulSet) Configure(statefulGraph StatefulGraph) error {
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

		podAdapter.Connect(statefulGraph, resourceNode, resource.Name)
	}
	return nil
}
