package kubegraph

import (
	"fmt"
	"log"
	"reflect"
	goRuntime "runtime"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/wwmoraes/kubegraph/internal/adapters"
	"github.com/wwmoraes/kubegraph/internal/utils"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// KubeGraph graphviz wrapper that creates kubernetes resource graphs
type KubeGraph struct {
	graphviz *graphviz.Graphviz
	graph    *cgraph.Graph
	// unknownArea *cgraph.Graph
	nodes   map[reflect.Type]map[string]*cgraph.Node
	objects map[reflect.Type]map[string]runtime.Object
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

	// unknownArea := graph.SubGraph("unknown", 1)
	// unknownArea.
	// 	SetNewRank(true).
	// 	SetPad(1.0).
	// 	SetRankDir(cgraph.TBRank).
	// 	SetRankSeparator(0.75).
	// 	SetNodeSeparator(0.60).
	// 	SetMargin(0).
	// 	SetFontSize(15).
	// 	SetSplines("ortho").
	// 	SetLayout("dot").
	// 	SetStyle(cgraph.InvisibleGraphStyle)

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
		// unknownArea: unknownArea,
		nodes:   nodes,
		objects: objects,
	}

	return kubegraph, nil
}

// AddStyledNode creates a new styled node with the given resource
func (kgraph KubeGraph) AddStyledNode(resourceType reflect.Type, resourceObject runtime.Object, nodeName string, resourceName string, icon string) (*cgraph.Node, error) {
	node, err := kgraph.createStyledNode(nodeName, resourceName, icon)
	if err != nil {
		return nil, err
	}

	kgraph.addNode(resourceType, resourceName, node)
	kgraph.addObject(resourceType, resourceName, resourceObject)

	return node, nil
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

func (kgraph KubeGraph) addNode(nodeType reflect.Type, nodeName string, node *cgraph.Node) {
	kgraph.nodes[nodeType][nodeName] = node
}

// GetNode gets a node by type/name
func (kgraph KubeGraph) GetNode(nodeType reflect.Type, nodeName string) (*cgraph.Node, error) {
	typeNodes, typeExists := kgraph.nodes[nodeType]
	if !typeExists {
		return nil, fmt.Errorf("no nodes for type %s found", nodeType.String())
	}

	node, nodeExists := typeNodes[nodeName]
	if !nodeExists {
		return nil, fmt.Errorf("node %s/%s not found", nodeType.String(), nodeName)
	}

	return node, nil
}

// GetObjects gets all objects in store
func (kgraph KubeGraph) GetObjects(objectType reflect.Type) (map[string]runtime.Object, error) {
	typeObjects, typeExists := kgraph.objects[objectType]
	if !typeExists {
		return nil, fmt.Errorf("no objects for type %s found", objectType.String())
	}

	return typeObjects, nil
}

func (kgraph KubeGraph) addObject(objectType reflect.Type, objectName string, object runtime.Object) {
	kgraph.objects[objectType][objectName] = object
}

// LinkNode links the node to the target node type/name, if it exists
func (kgraph KubeGraph) LinkNode(node *cgraph.Node, targetNodeType reflect.Type, targetNodeName string) (*cgraph.Edge, error) {
	targetNode, ok := kgraph.nodes[targetNodeType][targetNodeName]
	// TODO get or create unknown node and link here
	if !ok {
		// log.Printf("%s node %s not found, unable to link", targetNodeType, targetNodeName)
		return nil, fmt.Errorf("%s node %s not found, unable to link", targetNodeType, targetNodeName)
	}

	edgeName := fmt.Sprintf("%s-%s", node.Name(), targetNode.Name())
	return kgraph.graph.CreateEdge(edgeName, node, targetNode)
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

// CreateStatefulSetAppsV1 creates a new StatefulSet node on the graph
func (kgraph KubeGraph) CreateStatefulSetAppsV1(statefulSet *appsV1.StatefulSet) (*cgraph.Node, error) {
	name := fmt.Sprintf("%s.%s~%s", statefulSet.APIVersion, statefulSet.Kind, statefulSet.Name)
	return kgraph.createStyledNode(name, statefulSet.Name, "icons/deploy.svg")
}

// CreateSecretV1 creates a new Secret node on the graph
func (kgraph KubeGraph) CreateSecretV1(secret *v1.Secret) (*cgraph.Node, error) {
	name := fmt.Sprintf("%s.%s~%s", secret.APIVersion, secret.Kind, secret.Name)
	return kgraph.createStyledNode(name, secret.Name, "icons/secret.svg")
}

// CreateServiceV1 creates a new Service node on the graph
func (kgraph KubeGraph) CreateServiceV1(service *v1.Service) (*cgraph.Node, error) {
	name := fmt.Sprintf("%s.%s~%s", service.APIVersion, service.Kind, service.Name)
	return kgraph.createStyledNode(name, service.Name, "icons/svc.svg")
}

// CreateClusterRoleV1 creates a new ClusterRole node on the graph
func (kgraph KubeGraph) CreateClusterRoleV1(clusterRole *rbacV1.ClusterRole) (*cgraph.Node, error) {
	name := fmt.Sprintf("%s.%s~%s", clusterRole.APIVersion, clusterRole.Kind, clusterRole.Name)
	return kgraph.createStyledNode(name, clusterRole.Name, "icons/sa.svg")
}

// CreateClusterRoleV1beta1 creates a new ClusterRole node on the graph
func (kgraph KubeGraph) CreateClusterRoleV1beta1(clusterRole *rbacV1beta1.ClusterRole) (*cgraph.Node, error) {
	name := fmt.Sprintf("%s.%s~%s", clusterRole.APIVersion, clusterRole.Kind, clusterRole.Name)
	return kgraph.createStyledNode(name, clusterRole.Name, "icons/c-role.svg")
}

// CreateClusterRoleBindingV1beta1 creates a new ClusterRoleBinding node on the graph
func (kgraph KubeGraph) CreateClusterRoleBindingV1beta1(clusterRoleBinding *rbacV1beta1.ClusterRoleBinding) (*cgraph.Node, error) {
	name := fmt.Sprintf("%s.%s~%s", clusterRoleBinding.APIVersion, clusterRoleBinding.Kind, clusterRoleBinding.Name)
	return kgraph.createStyledNode(name, clusterRoleBinding.Name, "icons/crb.svg")
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
