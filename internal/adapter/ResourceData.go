package adapter

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/dot"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceData data used by kubernetes resource adapters
type ResourceData struct {
	resourceType reflect.Type
	iconPath     string
	registry     Registry
}

func NewResourceData(resourceType reflect.Type, iconPath string) ResourceData {
	return ResourceData{
		resourceType: resourceType,
		iconPath:     iconPath,
	}
}

// GetIconPath returns the type icon file path
func (data *ResourceData) GetIconPath() string {
	return data.iconPath
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (data *ResourceData) GetType() reflect.Type {
	return data.resourceType
}

// GetRegistry returns registry this adapter is registered to
func (data *ResourceData) GetRegistry() Registry {
	return data.registry
}

// GetAccessor returns a global instance of a kubernetes metadata accessor
func (data *ResourceData) GetAccessor() meta.MetadataAccessor {
	return data.registry.GetAccessor()
}

func (data *ResourceData) SetRegistry(registry Registry) {
	data.registry = registry
}

// Create add a graph node for the given object and stores it for further actions
func (data *ResourceData) Create(graph StatefulGraph, obj runtime.Object) (dot.Node, error) {
	accessor := data.GetAccessor()
	apiVersion, _ := accessor.APIVersion(obj)
	kind, _ := accessor.Kind(obj)
	name, _ := accessor.Name(obj)

	nodeName := fmt.Sprintf("%s.%s~%s", apiVersion, kind, name)
	resourceNode, err := graph.AddStyledNode(data.GetType(), obj, nodeName, name, data.GetIconPath())
	if err != nil {
		return nil, err
	}

	return resourceNode, nil
}

// Connect creates and edge between the given node and an object on this adapter
func (data *ResourceData) Connect(graph StatefulGraph, source dot.Node, targetName string) (dot.Edge, error) {
	return graph.LinkNode(source, data.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (data *ResourceData) Configure(graph StatefulGraph) error {
	return nil
}
