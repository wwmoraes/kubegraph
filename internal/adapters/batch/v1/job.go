package v1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/wwmoraes/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type jobAdapter struct {
	adapter.Resource
}

func init() {
	adapter.MustRegister(&jobAdapter{
		adapter.NewResource(
			reflect.TypeOf(&batchV1.Job{}),
			"icons/job.svg",
		),
	})
}

func (thisAdapter *jobAdapter) tryCastObject(obj runtime.Object) (*batchV1.Job, error) {
	casted, ok := obj.(*batchV1.Job)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *jobAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	resourceNode, err := statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/job.svg")
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
func (thisAdapter *jobAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
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
