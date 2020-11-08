package v1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type daemonSetAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&daemonSetAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&appsV1.DaemonSet{}),
		},
	})
}

func (thisAdapter *daemonSetAdapter) tryCastObject(obj runtime.Object) (*appsV1.DaemonSet, error) {
	casted, ok := obj.(*appsV1.DaemonSet)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *daemonSetAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *daemonSetAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	resourceNode, err := statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/ds.svg")
	if err != nil {
		return nil, err
	}

	podAdapter, err := adapter.Get(reflect.TypeOf(&coreV1.Pod{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	} else {
		podMetadata := resource.Spec.Template.ObjectMeta
		podMetadata.Name = resource.Name
		_, err := podAdapter.Create(statefulGraph, &coreV1.Pod{
			ObjectMeta: podMetadata,
			Spec:       resource.Spec.Template.Spec,
		})
		if err != nil {
			fmt.Println(fmt.Errorf("%s create error: %w", thisAdapter.GetType().String(), err))
		}
	}

	return resourceNode, nil
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *daemonSetAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *daemonSetAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
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

		_, err = podAdapter.Connect(statefulGraph, resourceNode, resource.Name)
		if err != nil {
			fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
		}
	}
	return nil
}
