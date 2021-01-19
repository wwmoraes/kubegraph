package adapter

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/dot"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

// Resource is implemented by values that can transform a kubernetes object
// kind information into nodes and create edges between them
type Resource interface {
	// IconPath returns the type icon file path
	IconPath() string
	// GetType returns the reflected type of the k8s kind managed by this instance
	GetType() reflect.Type
	// GetRegistry returns this adapter parent registry where it is registered at
	GetRegistry() Registry
	// Create add a graph node for the given object and stores it for further actions
	Create(graph StatefulGraph, obj runtime.Object) (dot.Node, error)
	// Connect creates and edge between the given node and an object on this adapter
	Connect(graph StatefulGraph, source dot.Node, targetName string) (dot.Edge, error)
	// Configure connects the resources on this adapter with its dependencies
	Configure(graph StatefulGraph) error
	// setRegistry stores a pointer to the registry where this adapter is registered at
	SetRegistry(Registry)
	// GetAccessor returns a global instance of a kubernetes metadata accessor
	GetAccessor() meta.MetadataAccessor
}

// resourceData data used by kubernetes adapters
type resourceData struct {
	resourceType reflect.Type
	iconPath     string
	registry     Registry
}

func NewResource(resourceType reflect.Type, iconPath string) Resource {
	return &resourceData{
		resourceType: resourceType,
		iconPath:     iconPath,
	}
}

// IconPath returns the type icon file path
func (data *resourceData) IconPath() string {
	return data.iconPath
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (data *resourceData) GetType() reflect.Type {
	return data.resourceType
}

// GetRegistry returns registry this adapter is registered to
func (data *resourceData) GetRegistry() Registry {
	return data.registry
}

// GetAccessor returns a global instance of a kubernetes metadata accessor
func (data *resourceData) GetAccessor() meta.MetadataAccessor {
	return data.registry.GetAccessor()
}

func (data *resourceData) SetRegistry(registry Registry) {
	data.registry = registry
}

// Create add a graph node for the given object and stores it for further actions
func (data *resourceData) Create(graph StatefulGraph, obj runtime.Object) (dot.Node, error) {
	accessor := data.GetAccessor()
	apiVersion, _ := accessor.APIVersion(obj)
	kind, _ := accessor.Kind(obj)
	name, _ := accessor.Name(obj)

	nodeName := fmt.Sprintf("%s.%s~%s", apiVersion, kind, name)
	resourceNode, err := graph.AddStyledNode(data.GetType(), obj, nodeName, name, data.IconPath())
	if err != nil {
		return nil, err
	}

	return resourceNode, nil
}

// Connect creates and edge between the given node and an object on this adapter
func (data *resourceData) Connect(graph StatefulGraph, source dot.Node, targetName string) (dot.Edge, error) {
	return graph.LinkNode(source, data.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (data *resourceData) Configure(graph StatefulGraph) error {
	return nil
}
