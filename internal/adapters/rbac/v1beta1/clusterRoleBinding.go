package v1beta1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterRoleBindingAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&clusterRoleBindingAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&rbacV1beta1.ClusterRoleBinding{}),
			"icons/crb.svg",
		),
	})
}

func (thisAdapter *clusterRoleBindingAdapter) tryCastObject(obj runtime.Object) (*rbacV1beta1.ClusterRoleBinding, error) {
	casted, ok := obj.(*rbacV1beta1.ClusterRoleBinding)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *clusterRoleBindingAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	clusterRoleV1beta1Adapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&rbacV1beta1.ClusterRole{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	}
	clusterRoleV1Adapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&rbacV1.ClusterRole{}))
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

		if clusterRoleV1Adapter != nil {
			_, err := clusterRoleV1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
			}
		}
		if clusterRoleV1beta1Adapter != nil {
			_, err := clusterRoleV1beta1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
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
