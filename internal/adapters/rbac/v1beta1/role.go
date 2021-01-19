package v1beta1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	"github.com/wwmoraes/kubegraph/internal/utils"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type roleAdapter struct {
	adapter.Resource
}

func init() {
	adapter.MustRegister(&roleAdapter{
		adapter.NewResource(
			reflect.TypeOf(&rbacV1beta1.Role{}),
			"icons/role.svg",
		),
	})
}

func (thisAdapter *roleAdapter) tryCastObject(obj runtime.Object) (*rbacV1beta1.Role, error) {
	casted, ok := obj.(*rbacV1beta1.Role)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *roleAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	podSecurityPolicyV1beta1Adapter, err := thisAdapter.GetRegistry().Get(reflect.TypeOf(&policyV1beta1.PodSecurityPolicy{}))
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

		for _, rule := range resource.Rules {
			if utils.ContainsString(rule.ResourceNames, "podsecuritypolicies") &&
				utils.ContainsString(rule.Verbs, "use") {
				for _, podSecurityPolicyResourceName := range rule.ResourceNames {
					if utils.ContainsString(rule.APIGroups, "policy") {
						_, err := podSecurityPolicyV1beta1Adapter.Connect(statefulGraph, resourceNode, podSecurityPolicyResourceName)
						if err != nil {
							fmt.Println(fmt.Errorf("%s configure error: %w", thisAdapter.GetType().String(), err))
						}
					}
				}
			}
		}
	}
	return nil
}
