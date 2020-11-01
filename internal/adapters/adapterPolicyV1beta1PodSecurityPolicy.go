package adapters

/*
 * remove the dummy struct and replace the references with a proper kubernetes API resource
 */

import (
	"fmt"
	"log"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type adapterPolicyV1beta1PodSecurityPolicy struct {
	Resource
}

func init() {
	RegisterResourceAdapter(&adapterPolicyV1beta1PodSecurityPolicy{
		Resource{
			resourceType: reflect.TypeOf(&policyV1beta1.PodSecurityPolicy{}),
		},
	})
}

func (adapter adapterPolicyV1beta1PodSecurityPolicy) tryCastObject(obj runtime.Object) (*policyV1beta1.PodSecurityPolicy, error) {
	casted, ok := obj.(*policyV1beta1.PodSecurityPolicy)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), adapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (adapter adapterPolicyV1beta1PodSecurityPolicy) GetType() reflect.Type {
	return adapter.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (adapter adapterPolicyV1beta1PodSecurityPolicy) Create(statefulGraph StatefulGraph, obj runtime.Object) (*cgraph.Node, error) {
	resource, err := adapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/psp.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (adapter adapterPolicyV1beta1PodSecurityPolicy) Connect(statefulGraph StatefulGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return statefulGraph.LinkNode(source, adapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (adapter adapterPolicyV1beta1PodSecurityPolicy) Configure(statefulGraph StatefulGraph) error {
	objects, err := statefulGraph.GetObjects(adapter.GetType())
	if err != nil {
		return err
	}
	for resourceName, resourceObject := range objects {
		resource, err := adapter.tryCastObject(resourceObject)
		if err != nil {
			return err
		}
		resourceNode, err := statefulGraph.GetNode(adapter.GetType(), resourceName)
		if err != nil {
			return err
		}

		// do something with each resource
		log.Printf("%s resource %s, node %s", adapter.GetType().String(), resource.Name, resourceNode.Name())
	}
	return nil
}
