package adapter

import (
	"reflect"

	"github.com/wwmoraes/dot"
	"k8s.io/apimachinery/pkg/runtime"
)

type NodesMap map[string]dot.Node
type TypeNodesMap map[reflect.Type]NodesMap

type ObjectsMap map[string]runtime.Object
type TypeObjectsMap map[reflect.Type]ObjectsMap

type StatefulGraph interface {
	dot.Graph
	AddStyledNode(resourceType reflect.Type, resourceObject runtime.Object, nodeName string, resourceName string, icon string) (dot.Node, error)
	LinkNode(node dot.Node, targetNodeType reflect.Type, targetNodeName string) (edge dot.Edge, err error)
	GetObjects(objectType reflect.Type) (ObjectsMap, error)
	GetNode(nodeType reflect.Type, nodeName string) (dot.Node, error)
}
