package v1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/wwmoraes/dot"
	"github.com/wwmoraes/kubegraph/internal/registry"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type statefulSetAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&statefulSetAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&appsV1.StatefulSet{}),
			"icons/sts.svg",
		),
	})
}

func (thisAdapter *statefulSetAdapter) tryCastObject(obj runtime.Object) (*appsV1.StatefulSet, error) {
	casted, ok := obj.(*appsV1.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *statefulSetAdapter) Create(statefulGraph registry.StatefulGraph, obj runtime.Object) (dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	resourceNode, err := statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/sts.svg")
	if err != nil {
		return nil, err
	}

	podAdapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&coreV1.Pod{}))
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

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *statefulSetAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	podAdapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&coreV1.Pod{}))
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
