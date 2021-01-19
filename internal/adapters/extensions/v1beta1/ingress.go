package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	coreV1 "k8s.io/api/core/v1"
	extensionsV1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type IngressAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&IngressAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&extensionsV1beta1.Ingress{}),
			"icons/ing.svg",
		),
	})
}

func (thisAdapter *IngressAdapter) tryCastObject(obj runtime.Object) (*extensionsV1beta1.Ingress, error) {
	casted, ok := obj.(*extensionsV1beta1.Ingress)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *IngressAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	serviceAdapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&coreV1.Service{}))
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
