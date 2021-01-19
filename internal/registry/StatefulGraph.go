package registry

import (
	"reflect"

	"github.com/wwmoraes/dot"
	"k8s.io/apimachinery/pkg/runtime"
)

// NodesMap is a collection of Node values by their ID string
type NodesMap map[string]dot.Node

// TypeNodesMap is a collection of NodesMap's by their reflected type
type TypeNodesMap map[reflect.Type]NodesMap

// ObjectsMap is a collection of kubernetes runtime object values by their
// name (adapter-dependant, defaults to the object's metadata name value)
type ObjectsMap map[string]runtime.Object

// TypeObjectsMap is a collection of ObjectsMap's by their reflected type
type TypeObjectsMap map[reflect.Type]ObjectsMap

// StatefulGraph is implemented by dot-compatible graph values that also manages
// kubernetes resource objects, graph nodes and edges, providing methods to
// create nodes, link them (i.e. create edges) and access other object and node
// values (used during adapter configuration to find referenced objects)
type StatefulGraph interface {
	dot.Graph
	AddStyledNode(resourceType reflect.Type, resourceObject runtime.Object, nodeName string, resourceName string, icon string) (dot.Node, error)
	LinkNode(node dot.Node, targetNodeType reflect.Type, targetNodeName string) (edge dot.Edge, err error)
	GetObjects(objectType reflect.Type) (ObjectsMap, error)
	GetNode(nodeType reflect.Type, nodeName string) (dot.Node, error)
}
