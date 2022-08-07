package v1beta1

import (
	"fmt"
	"log"

	"github.com/wwmoraes/kubegraph/internal/registry"
	"github.com/wwmoraes/kubegraph/internal/utils"
)

// Configure connects the resources on this adapter with its dependencies
func (this *RoleAdapter) Configure(statefulGraph registry.StatefulGraph) error {
	podSecurityPolicyV1beta1Adapter, err := GetPolicyV1beta1PodSecurityPolicyAdapter()
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

		for _, rule := range resource.Rules {
			if utils.ContainsString(rule.ResourceNames, "podsecuritypolicies") &&
				utils.ContainsString(rule.Verbs, "use") {
				for _, podSecurityPolicyResourceName := range rule.ResourceNames {
					if utils.ContainsString(rule.APIGroups, "policy") {
						_, err := podSecurityPolicyV1beta1Adapter.Connect(statefulGraph, resourceNode, podSecurityPolicyResourceName)
						if err != nil {
							fmt.Println(fmt.Errorf("%s configure error: %w", this.GetType().String(), err))
						}
					}
				}
			}
		}
	}
	return nil
}
