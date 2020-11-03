package v1beta1

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type podSecurityPolicyAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&podSecurityPolicyAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&policyV1beta1.PodSecurityPolicy{}),
		},
	})
}

func (thisAdapter podSecurityPolicyAdapter) tryCastObject(obj runtime.Object) (*policyV1beta1.PodSecurityPolicy, error) {
	casted, ok := obj.(*policyV1beta1.PodSecurityPolicy)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter podSecurityPolicyAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter podSecurityPolicyAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/psp.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter podSecurityPolicyAdapter) Connect(statefulGraph adapter.StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter podSecurityPolicyAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
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
		log.Printf("%s resource %s, node %s", thisAdapter.GetType().String(), resource.Name, resourceNode.Name())
	}
	return nil
}
