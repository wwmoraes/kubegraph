package kubegraph

import (
	"fmt"
	"log"
	"os"
	"reflect"
	goRuntime "runtime"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/wwmoraes/kubegraph/internal/adapters"
	"github.com/wwmoraes/kubegraph/internal/utils"
	"k8s.io/apimachinery/pkg/runtime"
)

// KubeGraph graphviz wrapper that creates kubernetes resource graphs
type KubeGraph struct {
	graphviz *graphviz.Graphviz
	graph    *cgraph.Graph
	nodes    map[reflect.Type]map[string]*cgraph.Node
	objects  map[reflect.Type]map[string]runtime.Object
}

// New creates an instance of KubeGraph
func New() (KubeGraph, error) {
	gz := graphviz.New()

	graph, err := gz.Graph(graphviz.Name("kubegraph"))
	if err != nil {
		return KubeGraph{}, err
	}

	graph.
		SetNewRank(true).
		SetPad(1.0).
		SetRankDir(cgraph.TBRank).
		SetRankSeparator(0.75).
		SetNodeSeparator(0.60).
		SetMargin(0).
		SetFontSize(15).
		SetSplines("ortho").
		SetLayout("dot").
		SetStyle(cgraph.RoundedGraphStyle)

	goRuntime.SetFinalizer(graph, closeGraph)
	goRuntime.SetFinalizer(gz, closeGraphviz)

	// initialize nodes map with registered adapter types
	nodes := make(map[reflect.Type]map[string]*cgraph.Node)
	for adapterType := range adapters.GetAdapters() {
		nodes[adapterType] = make(map[string]*cgraph.Node)
	}

	// initialize object map with registered adapter types
	objects := make(map[reflect.Type]map[string]runtime.Object)
	for adapterType := range adapters.GetAdapters() {
		objects[adapterType] = make(map[string]runtime.Object)
	}

	kubegraph := KubeGraph{
		graphviz: gz,
		graph:    graph,
		nodes:    nodes,
		objects:  objects,
	}

	return kubegraph, nil
}

func (kgraph KubeGraph) createStyledNode(name string, label string, icon string) (*cgraph.Node, error) {
	node, err := kgraph.graph.CreateNode(name)
	if err != nil {
		return nil, err
	}

	// break long labels so it fits on our graph (k8s resource names can be up to
	// 253 characters long)
	labelLines := utils.StringChunks(label, 16)
	labelLinesCount := len(labelLines)
	minHeight := 1.9 + 0.4*float64(labelLinesCount)
	minWidth := 1.9
	node.SetShape(cgraph.NoneShape)
	node.SetImage(icon)
	node.SetLabelLocation(cgraph.BottomLocation)
	node.SetHeight(minHeight)
	node.SetWidth(minWidth)
	node.SetFontSize(13)
	node.SetFixedSize(true)
	node.SetImageScale(true)
	node.SetLabel(strings.Join(labelLines, "\n"))

	return node, nil
}

func (kgraph KubeGraph) getNodes(objectType reflect.Type) (map[string]*cgraph.Node, error) {
	typeNodes, typeExists := kgraph.nodes[objectType]
	if !typeExists {
		return nil, fmt.Errorf("no nodes for type %s found", objectType.String())
	}

	return typeNodes, nil
}

func (kgraph KubeGraph) addNode(nodeType reflect.Type, nodeName string, node *cgraph.Node) error {
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

func (kgraph KubeGraph) createUnknown(obj runtime.Object) (*cgraph.Node, error) {
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
	for _, adapter := range adapters.GetAdapters() {
		err := adapter.Configure(kgraph)
		if err != nil {
			log.Println(err)
		}
	}
}

// Transform creates a node on the graph for the resource
func (kgraph KubeGraph) Transform(obj runtime.Object) (*cgraph.Node, error) {
	adapter, err := adapters.GetAdapterFor(reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}

	return adapter.Create(kgraph, obj)
}

// Render generates a graph file
func (kgraph KubeGraph) Render(fileName string, format graphviz.Format) error {
	return kgraph.graphviz.RenderFilename(kgraph.graph, format, fileName)
}
