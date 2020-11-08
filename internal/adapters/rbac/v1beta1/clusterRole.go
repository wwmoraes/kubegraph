package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type clusterRoleAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.Register(&clusterRoleAdapter{
		adapter.ResourceData{
			ResourceType: reflect.TypeOf(&rbacV1beta1.ClusterRole{}),
		},
	})
}

func (thisAdapter *clusterRoleAdapter) tryCastObject(obj runtime.Object) (*rbacV1beta1.ClusterRole, error) {
	casted, ok := obj.(*rbacV1beta1.ClusterRole)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (thisAdapter *clusterRoleAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

// Create add a graph node for the given object and stores it for further actions
func (thisAdapter *clusterRoleAdapter) Create(statefulGraph adapter.StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/c-role.svg")
}

// Connect creates and edge between the given node and an object on this adapter
func (thisAdapter *clusterRoleAdapter) Connect(statefulGraph adapter.StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (thisAdapter *clusterRoleAdapter) Configure(statefulGraph adapter.StatefulGraph) error {
	return nil
}
