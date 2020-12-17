package adapter

import (
	"reflect"

	"github.com/wwmoraes/dot"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceTransformer is implemented by values that can transform a kubernetes
// resource kind information into nodes and create edges between them and other
// kinds
type ResourceTransformer interface {
	// GetType returns the reflected type of the k8s kind managed by this instance
	GetType() reflect.Type
	// GetIconPath returns the type icon file path
	GetIconPath() string
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
