package kubegraph

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/emicklei/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"

	// self-register adapters
	_ "github.com/wwmoraes/kubegraph/internal/adapters"
	"github.com/wwmoraes/kubegraph/internal/utils"
	"k8s.io/apimachinery/pkg/runtime"
)

// KubeGraph graphviz wrapper that creates kubernetes resource graphs
type KubeGraph struct {
	graph   *dot.Graph
	nodes   map[reflect.Type]map[string]*dot.Node
	objects map[reflect.Type]map[string]runtime.Object
}

// New creates an instance of KubeGraph
func New() (kubegraph KubeGraph, err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			kubegraph = KubeGraph{}
			err = fmt.Errorf("%++v", recoverErr)
		}
	}()

	graph := dot.NewGraph(dot.Directed)
	graph.ID("kubegraph")

	graph.Attrs(
		"rankdir", "TB",
		"ranksep", "0.75",
		"newrank", "true",
		"nodesep", "0.6",
		"pad", "1.0",
		"fontsize", "15",
		"layout", "dot",
		"margin", "0",
		"splines", "ortho",
		"style", "rounded",
	)

	// initialize nodes and objects maps with registered adapter types
	nodes := make(map[reflect.Type]map[string]*dot.Node)
	objects := make(map[reflect.Type]map[string]runtime.Object)
	for adapterType := range adapter.GetAll() {
		nodes[adapterType] = make(map[string]*dot.Node)
		objects[adapterType] = make(map[string]runtime.Object)
	}

	kubegraph = KubeGraph{
		graph:   graph,
		nodes:   nodes,
		objects: objects,
	}

	return kubegraph, nil
}

func (kgraph KubeGraph) createStyledNode(name string, label string, icon string) (*dot.Node, error) {
	node := kgraph.graph.Node(name)

	// break long labels so it fits on our graph (k8s resource names can be up to
	// 253 characters long)
	labelLines := utils.StringChunks(label, 16)
	labelLinesCount := len(labelLines)
	minHeight := 1.9 + 0.4*float64(labelLinesCount)
	minWidth := 1.9
	node.Attrs(
		"shape", "none",
		"image", icon,
		"labelloc", "b",
		"height", fmt.Sprintf("%f", minHeight),
		"width", fmt.Sprintf("%f", minWidth),
		"fontsize", "13",
		"fixedsize", "true",
		"imagescale", "true",
		"label", strings.Join(labelLines, "\n"),
	)

	return node, nil
}

func (kgraph KubeGraph) getNodes(objectType reflect.Type) (map[string]*dot.Node, error) {
	typeNodes, typeExists := kgraph.nodes[objectType]
	if !typeExists {
		return nil, fmt.Errorf("no nodes for type %s found", objectType.String())
	}

	return typeNodes, nil
}

func (kgraph KubeGraph) addNode(nodeType reflect.Type, nodeName string, node *dot.Node) error {
	nodes, err := kgraph.getNodes(nodeType)
	if err != nil {
		return err
	}

	nodes[nodeName] = node
	return nil
}

func (kgraph KubeGraph) addObject(objectType reflect.Type, objectName string, object runtime.Object) error {
	objects, err := kgraph.GetObjects(objectType)
	if err != nil {
		return err
	}

	objects[objectName] = object
	return nil
}

// nolint:unused // future implementation of not found nodes
func (kgraph KubeGraph) createUnknown(obj runtime.Object) (*dot.Node, error) {
	obj.GetObjectKind()
	metadata, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	name := fmt.Sprintf(
		"%s.%s~%s",
		metadata["apiVersion"].(string),
		metadata["kind"].(string),
		metadata["metadata"].(map[string]interface{})["name"].(string),
	)

	label := fmt.Sprintf(
		"%s.%s\n%s",
		metadata["apiVersion"].(string),
		metadata["kind"].(string),
		metadata["metadata"].(map[string]interface{})["name"].(string),
	)

	// node, err := kgraph.unknownArea.CreateNode(name)
	// node.SetLabel(label)
	// if err != nil {
	// 	return nil, err
	// }

	node, err := kgraph.createStyledNode(name, label, "icons/unknown.svg")
	if err != nil {
		return nil, err
	}

	return node, err
}

// ConnectNodes creates edges between the nodes
func (kgraph KubeGraph) ConnectNodes() {
	for _, adapter := range adapter.GetAll() {
		err := adapter.Configure(kgraph)
		if err != nil {
			log.Println(err)
		}
	}
}

// Transform creates a node on the graph for the resource
func (kgraph KubeGraph) Transform(obj runtime.Object) (*dot.Node, error) {
	objectAdapter, err := adapter.Get(reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}

	return objectAdapter.Create(kgraph, obj)
}

// Write write the graph contents to a writer using simple TAB indentation
func (kgraph KubeGraph) Write(target io.Writer) {
	kgraph.graph.Write(target)
}
