package kubegraph

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/wwmoraes/dot"
	"github.com/wwmoraes/dot/attributes"
	"github.com/wwmoraes/dot/constants"
	"github.com/wwmoraes/kubegraph/internal/adapters"
	"github.com/wwmoraes/kubegraph/internal/registry"
	"github.com/wwmoraes/kubegraph/internal/utils"
	"k8s.io/apimachinery/pkg/runtime"

	// self-register adapters
	_ "github.com/wwmoraes/kubegraph/internal/adapters"
)

// New creates an instance of Kubegraph
func New() (*Kubegraph, error) {
	return InitializeKubegraph(
		dot.WithID("kubegraph"),
		dot.WithType(dot.GraphTypeDirected),
	)
}

// Kubegraph graphviz wrapper that creates kubernetes resource graphs
type Kubegraph struct {
	dot.Graph
	k8sNodes   registry.TypeNodesMap
	k8sObjects registry.TypeObjectsMap
	registry   registry.Registry
	decode     adapters.DecodeFn
}

// NewKubegraph creates an instance of Kubegraph with the provided dot Graph
// and Registry instance
func NewKubegraph(graph dot.Graph, registryInstance registry.Registry, decode adapters.DecodeFn) *Kubegraph {
	graph.SetAttributes(attributes.Map{
		constants.KeyRankDir:  attributes.NewString("TB"),
		constants.KeyRankSep:  attributes.NewString("0.75"),
		constants.KeyNewRank:  attributes.NewString("true"),
		constants.KeyNodeSep:  attributes.NewString("0.6"),
		constants.KeyPad:      attributes.NewString("1.0"),
		constants.KeyFontSize: attributes.NewString("15"),
		constants.KeyLayout:   attributes.NewString("dot"),
		constants.KeyMargin:   attributes.NewString("0"),
		constants.KeySplines:  attributes.NewString("ortho"),
		constants.KeyStyle:    attributes.NewString("rounded"),
	})

	nodes := make(registry.TypeNodesMap)
	objects := make(registry.TypeObjectsMap)

	for adapterType := range registryInstance.GetAll() {
		nodes[adapterType] = make(registry.NodesMap)
		objects[adapterType] = make(registry.ObjectsMap)
	}

	return &Kubegraph{
		graph,
		nodes,
		objects,
		registryInstance,
		decode,
	}
}

// ConnectNodes creates edges between the nodes
func (kgraph *Kubegraph) ConnectNodes() {
	for _, registryAdapter := range registry.Instance().GetAll() {
		err := registryAdapter.Configure(kgraph)
		if err != nil {
			log.Println(err)
		}
	}
}

// Transform creates a node on the graph for the resource
func (kgraph *Kubegraph) Transform(obj runtime.Object) (dot.Node, error) {
	objectAdapter, err := kgraph.registry.Get(reflect.TypeOf(obj))
	if err != nil {
		return nil, err
	}

	return objectAdapter.Create(kgraph, obj)
}

func (graph *Kubegraph) createStyledNode(name string, label string, icon string) (dot.Node, error) {
	node := graph.Node(name)

	// break long labels so it fits on our graph (k8s resource names can be up to
	// 253 characters long)
	labelLines := utils.StringChunks(label, 16)
	labelLinesCount := len(labelLines)
	minHeight := 1.9 + 0.4*float64(labelLinesCount)
	minWidth := 1.9
	node.SetAttributes(attributes.Map{
		constants.KeyShape:      attributes.NewString("none"),
		constants.KeyImage:      attributes.NewString(icon),
		constants.KeyLabelLoc:   attributes.NewString("b"),
		constants.KeyHeight:     attributes.NewString(fmt.Sprintf("%f", minHeight)),
		constants.KeyWidth:      attributes.NewString(fmt.Sprintf("%f", minWidth)),
		constants.KeyFontSize:   attributes.NewString("13"),
		constants.KeyFixedSize:  attributes.NewString("true"),
		constants.KeyImageScale: attributes.NewString("true"),
		constants.KeyLabel:      attributes.NewString(strings.Join(labelLines, "\n")),
	})

	return node, nil
}

func (graph *Kubegraph) addNode(nodeType reflect.Type, nodeName string, node dot.Node) error {
	nodes, err := graph.getNodes(nodeType)
	if err != nil {
		return err
	}

	nodes[nodeName] = node
	return nil
}

func (graph *Kubegraph) getNodes(objectType reflect.Type) (registry.NodesMap, error) {
	typeNodes, typeExists := graph.k8sNodes[objectType]
	if !typeExists {
		return nil, fmt.Errorf("no nodes for type %s found", objectType.String())
	}

	return typeNodes, nil
}

func (graph *Kubegraph) addObject(objectType reflect.Type, objectName string, object runtime.Object) error {
	objects, err := graph.GetObjects(objectType)
	if err != nil {
		return err
	}

	objects[objectName] = object
	return nil
}

// AddStyledNode creates a new styled node with the given resource
func (graph *Kubegraph) AddStyledNode(resourceType reflect.Type, resourceObject runtime.Object, nodeName string, resourceName string, icon string) (dot.Node, error) {

	node, err := graph.createStyledNode(nodeName, resourceName, icon)
	if err != nil {
		return nil, err
	}

	if err := graph.addNode(resourceType, resourceName, node); err != nil {
		return nil, err
	}
	if err := graph.addObject(resourceType, resourceName, resourceObject); err != nil {
		// TODO remove node added previously
		return nil, err
	}

	return node, nil
}

// LinkNode links the node to the target node type/name, if it exists
func (graph *Kubegraph) LinkNode(node dot.Node, targetNodeType reflect.Type, targetNodeName string) (edge dot.Edge, err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			edge = nil
			err = fmt.Errorf("%++v", recoverErr)
		}
	}()

	targetNode, ok := graph.k8sNodes[targetNodeType][targetNodeName]
	// TODO get or create unknown node and link here
	if !ok {
		// log.Printf("%s node %s not found, unable to link", targetNodeType, targetNodeName)
		return nil, fmt.Errorf("%s node %s not found, unable to link", targetNodeType, targetNodeName)
	}

	edge = graph.Edge(node, targetNode)
	edge.SetAttributeString(constants.KeyLabel, "")
	return edge, nil
}

// GetObjects gets all objects in store
func (graph *Kubegraph) GetObjects(objectType reflect.Type) (registry.ObjectsMap, error) {
	typeObjects, typeExists := graph.k8sObjects[objectType]
	if !typeExists {
		return nil, fmt.Errorf("no objects for type %s found", objectType.String())
	}

	return typeObjects, nil
}

// GetNode gets a node by type/name
func (graph *Kubegraph) GetNode(nodeType reflect.Type, nodeName string) (dot.Node, error) {
	typeNodes, typeExists := graph.k8sNodes[nodeType]
	if !typeExists {
		return nil, fmt.Errorf("no nodes for type %s found", nodeType.String())
	}

	node, nodeExists := typeNodes[nodeName]
	if !nodeExists {
		return nil, fmt.Errorf("node %s/%s not found", nodeType.String(), nodeName)
	}

	return node, nil
}

// LoadFromData normalizes input data, decodes resources and transform them
func (instance *Kubegraph) LoadFromData(data io.Reader) error {
	log.Println("reading all data...")
	fileBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	// normalize line breaks
	log.Println("normalizing linebreaks...")
	fileString := string(fileBytes[:])
	fileString = strings.ReplaceAll(fileString, "\r\n", "\n")
	fileString = strings.ReplaceAll(fileString, "\r", "\n")

	// removes all comments from yaml and json
	log.Println("removing comments and empty lines...")
	commentLineMatcher, err := regexp.Compile("^[ ]*((#|//).*)?$")
	if err != nil {
		return err
	}
	fileStringLines := strings.Split(fileString, "\n")
	var cleanFileString strings.Builder
	for _, line := range fileStringLines {
		if commentLineMatcher.MatchString(line) {
			continue
		}
		if line == "\n" || line == "" {
			continue
		}

		_, err := cleanFileString.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			return err
		}
	}
	fileString = cleanFileString.String()

	log.Println("splitting documents...")
	documents := strings.Split(fileString, "---")

	for _, document := range documents {
		if document == "\n" || document == "" {
			continue
		}

		obj, _, err := instance.decode([]byte(document), nil, nil)
		if err != nil {
			log.Printf("unable to decode document: %++v\n", err)
			continue
		}

		_, err = instance.Transform(obj)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("connecting nodes...")
	instance.ConnectNodes()

	return nil
}
