package v1beta1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	"github.com/wwmoraes/kubegraph/internal/utils"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type roleAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&roleAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&rbacV1beta1.Role{}),
		},
	})
}

func (thisAdapter roleAdapter) tryCastObject(obj runtime.Object) (*rbacV1beta1.Role, error) {
	casted, ok := obj.(*rbacV1beta1.Role)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter roleAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter roleAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/role.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter roleAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter roleAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	podSecurityPolicyV1beta1Adapter, err := adapter.Get(reflect.TypeOf(&policyV1beta1.PodSecurityPolicy{}))
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
						podSecurityPolicyV1beta1Adapter.Connect(statefulGraph, resourceNode, podSecurityPolicyResourceName)
					}
				}
			}
		}
	}
	return nil
}
