package v1beta1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type roleBindingAdapter struct {
	adapter.Resource
}

func init() {
	adapter.MustRegister(&roleBindingAdapter{
		adapter.NewResource(
			reflect.TypeOf(&rbacV1beta1.RoleBinding{}),
			"icons/rb.svg",
		),
	})
}

func (thisAdapter *roleBindingAdapter) tryCastObject(obj runtime.Object) (*rbacV1beta1.RoleBinding, error) {
	casted, ok := obj.(*rbacV1beta1.RoleBinding)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *roleBindingAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	roleV1Adapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&rbacV1.Role{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	}
	roleV1beta1Adapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&rbacV1beta1.Role{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
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

		if roleV1Adapter != nil {
			_, err := roleV1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
			}
		}
		if roleV1beta1Adapter != nil {
			_, err := roleV1beta1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
			}
		}

		for _, subject := range resource.Subjects {
			if subject.Kind == "ServiceAccount" {
				saNode, err := statefulGraph.GetNode(reflect.TypeOf(&coreV1.ServiceAccount{}), subject.Name)
				if err != nil {
					return err
				}
				_, err = thisAdapter.Connect(statefulGraph, saNode, resourceName)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
				}
			}
		}
	}
	return nil
}
