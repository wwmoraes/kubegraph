package adapters

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterBatchV1Job struct{}

func init() {
	RegisterResourceAdapter(&adapterBatchV1Job{})
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterBatchV1Job) GetType() reflect.Type {
	return reflect.TypeOf(&batchV1.Job{})
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterBatchV1Job) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource := obj.(*batchV1.Job)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	resourceNode, err := statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/job.svg")
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
func (adapter adapterBatchV1Job) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterBatchV1Job) Configure(statefulGraph StatefulGraph) error {
	podAdapter, err := GetAdapterFor(reflect.TypeOf(&coreV1.Pod{}))
	if err != nil {
		return fmt.Errorf("warning[%s configure]: %v", adapter.GetType().String(), err)
	}

	objects, err := statefulGraph.GetObjects(adapter.GetType())
	if err != nil {
		return err
	}

	for resourceName, resourceObject := range objects {
		resource := resourceObject.(*batchV1.Job)
		resourceNode, err := statefulGraph.GetNode(adapter.GetType(), resourceName)
		if err != nil {
			return err
		}

		podAdapter.Connect(statefulGraph, resourceNode, resource.Name)
	}
	return nil
}
