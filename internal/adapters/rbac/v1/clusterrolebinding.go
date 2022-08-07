package v1

import (
	"fmt"
	"log"

	"github.com/wwmoraes/kubegraph/internal/registry"
)

// Configure connects the resources on this adapter with its dependencies
func (this *ClusterRoleBindingAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	clusterRoleV1beta1Adapter, err := GetRbacV1beta1ClusterRoleAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err))
	}
	clusterRoleV1Adapter, err := GetClusterRoleAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err))
	}

	saAdapter, err := GetServiceAccountAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err))
	}

	objects, err := this.GetGraphObjects(statefulGraph)
	if err != nil {
		return err
	}
	for resourceName, resource := range objects {
		resourceNode, err := this.GetGraphNode(statefulGraph, resourceName)
		if err != nil {
			return err
		}

		if clusterRoleV1Adapter != nil {
			_, err := clusterRoleV1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
			}
		}
		if clusterRoleV1beta1Adapter != nil {
			_, err := clusterRoleV1beta1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
			}
		}

		for _, subject := range resource.Subjects {
			if subject.Kind == "ServiceAccount" {
				saNode, err := saAdapter.GetGraphNode(statefulGraph, subject.Name)
				if err != nil {
					return err
				}
				_, err = this.Connect(statefulGraph, saNode, resourceName)
				if err != nil {
					fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
				}
			}
		}
	}
	return nil
}
