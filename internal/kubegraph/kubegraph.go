package kubegraph

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	apiMachineryRuntime "k8s.io/apimachinery/pkg/runtime"
)

// KubeGraph graphviz helper to create a kubernetes resource graph
type KubeGraph struct {
	graphviz *graphviz.Graphviz
	graph    *cgraph.Graph
	// unknownArea *cgraph.Graph
	nodes   map[reflect.Type]map[string]*cgraph.Node
	objects map[reflect.Type]map[string]apiMachineryRuntime.Object
}

// New creates an instance of KubernetesGraph
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

	runtime.SetFinalizer(graph, closeGraph)
	runtime.SetFinalizer(gz, closeGraphviz)

	// initialize nodes map with registered adapter types
	nodes := make(map[reflect.Type]map[string]*cgraph.Node)
	for adapterType := range adapters {
		nodes[adapterType] = make(map[string]*cgraph.Node)
	}

	// initialize object map with registered adapter types
	objects := make(map[reflect.Type]map[string]apiMachineryRuntime.Object)
	for adapterType := range adapters {
		objects[adapterType] = make(map[string]apiMachineryRuntime.Object)
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

func (kgraph KubeGraph) addStyledNode(resourceType reflect.Type, resourceObject apiMachineryRuntime.Object, nodeName string, resourceName string, icon string) (*cgraph.Node, error) {
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
	labelLines := stringChunks(label, 16)
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

func (kgraph KubeGraph) addObject(objectType reflect.Type, objectName string, object apiMachineryRuntime.Object) {
	kgraph.objects[objectType][objectName] = object
}

func (kgraph KubeGraph) linkNode(node *cgraph.Node, targetNodeType reflect.Type, targetNodeName string) (*cgraph.Edge, error) {
	targetNode, ok := kgraph.nodes[targetNodeType][targetNodeName]
	if !ok {
		// log.Printf("%s node %s not found, unable to link", targetNodeType, targetNodeName)
		return nil, fmt.Errorf("%s node %s not found, unable to link", targetNodeType, targetNodeName)
	}

	return kgraph.graph.CreateEdge("", node, targetNode)
}

func (kgraph KubeGraph) createUnknown(obj apiMachineryRuntime.Object) (*cgraph.Node, error) {
	obj.GetObjectKind()
	metadata, _ := apiMachineryRuntime.DefaultUnstructuredConverter.ToUnstructured(obj)
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
func (kgraph KubeGraph) ConnectNodes() error {
	for _, adapter := range adapters {
		err := adapter.Configure(kgraph)
		if err != nil {
			return err
		}
	}

	return nil
}

// Transform creates a node on the graph for the resource
func (kgraph KubeGraph) Transform(obj apiMachineryRuntime.Object) (*cgraph.Node, error) {
	adapter, ok := adapters[reflect.TypeOf(obj)]
	if !ok {
		log.Printf(
			"Unsupported resource %s.%s\n",
			obj.GetObjectKind().GroupVersionKind().Version,
			obj.GetObjectKind().GroupVersionKind().Kind,
		)
		// return kgraph.createUnknown(obj), nil
		return nil, nil
	}

	return adapter.Create(kgraph, obj)
}

// Render generates a graph file
func (kgraph KubeGraph) Render(fileName string, format graphviz.Format) error {
	return kgraph.graphviz.RenderFilename(kgraph.graph, format, fileName)
}
