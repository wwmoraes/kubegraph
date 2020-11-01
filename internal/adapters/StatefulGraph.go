package adapters

import (
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	"k8s.io/apimachinery/pkg/runtime"
)

// StatefulGraph graphviz-compatible struct with runtime object and node storage
type StatefulGraph interface {
	// AddStyledNode creates a new styled node with the given resource
	AddStyledNode(resourceType reflect.Type, resourceObject runtime.Object, nodeName string, resourceName string, icon string) (*cgraph.Node, error)
	// LinkNode links the node to the target node type/name, if it exists
	LinkNode(node *cgraph.Node, targetNodeType reflect.Type, targetNodeName string) (*cgraph.Edge, error)
	// GetObjects gets all objects in store
	GetObjects(objectType reflect.Type) (map[string]runtime.Object, error)
	// GetNode gets a node by type/name
	GetNode(nodeType reflect.Type, nodeName string) (*cgraph.Node, error)
}
