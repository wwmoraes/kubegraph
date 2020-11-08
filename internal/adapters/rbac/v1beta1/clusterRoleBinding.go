package v1beta1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterRoleBindingAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&clusterRoleBindingAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&rbacV1beta1.ClusterRoleBinding{}),
		},
	})
}

func (thisAdapter clusterRoleBindingAdapter) tryCastObject(obj runtime.Object) (*rbacV1beta1.ClusterRoleBinding, error) {
	casted, ok := obj.(*rbacV1beta1.ClusterRoleBinding)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter clusterRoleBindingAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter clusterRoleBindingAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/crb.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter clusterRoleBindingAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter clusterRoleBindingAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	clusterRoleV1beta1Adapter, err := adapter.Get(reflect.TypeOf(&rbacV1beta1.ClusterRole{}))
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %v", thisAdapter.GetType().String(), err))
	}
	clusterRoleV1Adapter, err := adapter.Get(reflect.TypeOf(&rbacV1.ClusterRole{}))
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
