package registry

import (
	"fmt"
	"reflect"
)

// adapterData data used by kubernetes adapters
type adapterData struct {
	resourceType reflect.Type
	iconPath     string
}

func NewAdapter(resourceType reflect.Type, iconPath string) Adapter {
	return &adapterData{
		resourceType: resourceType,
		iconPath:     iconPath,
	}
}

// IconPath returns the type icon file path
func (data *adapterData) IconPath() string {
	return data.iconPath
}

// GetType returns the reflected type of the k8s kind managed by this instance
func (data *adapterData) GetType() reflect.Type {
	return data.resourceType
}

// Create add a graph node for the given object and stores it for further actions
func (data *adapterData) Create(graph StatefulGraph, obj RuntimeObject) (Node, error) {
	accessor := Instance().GetAccessor()
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
func (data *adapterData) Connect(graph StatefulGraph, source Node, targetName string) (Edge, error) {
	return graph.LinkNode(source, data.GetType(), targetName)
}

// Configure connects the resources on this adapter with its dependencies
func (data *adapterData) Configure(graph StatefulGraph) error {
	// returns nil as there's nothing to configure
	return nil
}

// GetGraphNode returns the node registered under name
func (data *adapterData) GetGraphNode(StatefulGraph, string) (Node, error) {
	return nil, ErrUnimplemented
}
