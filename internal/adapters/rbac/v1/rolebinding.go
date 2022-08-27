package v1

import (
	"fmt"
	"log"

	"github.com/wwmoraes/kubegraph/internal/registry"
)

// Configure connects the resources on this adapter with its dependencies
func (this *RoleBindingAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	roleV1Adapter, err := GetRoleAdapter()
	if err != nil {
		log.Println(fmt.Errorf("warning[%s configure]: %w", this.GetType().String(), err))
	}
	roleV1beta1Adapter, err := GetRbacV1beta1RoleAdapter()
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

		if roleV1Adapter != nil {
			_, err := roleV1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
			}
		}
		if roleV1beta1Adapter != nil {
			_, err := roleV1beta1Adapter.Connect(statefulGraph, resourceNode, resource.RoleRef.Name)
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
