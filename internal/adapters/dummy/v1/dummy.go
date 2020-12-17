package v1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

/*
 * remove the dummyResource struct and replace the references with a proper kubernetes API resource
 */
type dummyResource struct {
	metaV1.TypeMeta
	metaV1.ObjectMeta
}

func (d dummyResource) GetObjectKind() schema.ObjectKind {
	return nil
}

func (d dummyResource) DeepCopyObject() runtime.Object {
	return dummyResource{}
}

// dummyAdapter a kubegraph adapter to render an specific kubernetes resource
type dummyAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.MustRegister(&dummyAdapter{
		adapter.NewResourceData(
			reflect.TypeOf(&dummyResource{}),
			"icons/unknown.svg",
		),
	})
}

func (thisAdapter *dummyAdapter) tryCastObject(obj runtime.Object) (*dummyResource, error) {
	casted, ok := obj.(*dummyResource)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *dummyAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	log.Printf("please implement a configuration for %s resources", thisAdapter.GetType().String())
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

		// do something with each resource
		log.Printf("nothing to configure for %s, node %s", resource.Name, resourceNode.GetAttributeString("label"))
	}
	return nil
}
