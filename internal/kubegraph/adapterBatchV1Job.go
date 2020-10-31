package kubegraph

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-graphviz/cgraph"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type adapterBatchV1Job struct{}

func init() {
	RegisterResourceAdapter(&adapterBatchV1Job{})
}

func (adapter adapterBatchV1Job) GetType() reflect.Type {
	return reflect.TypeOf(&batchv1.Job{})
}

func (adapter adapterBatchV1Job) Create(kgraph KubeGraph, obj k8sruntime.Object) (*cgraph.Node, error) {
	resource := obj.(*batchv1.Job)
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)

	podMetadata := resource.Spec.Template.ObjectMeta
	podMetadata.Name = resource.Name
	adapters[reflect.TypeOf(&v1.Pod{})].Create(kgraph, &v1.Pod{
		ObjectMeta: podMetadata,
		Spec:       resource.Spec.Template.Spec,
	})

	return kgraph.addStyledNode(adapter.GetType(), obj, name, resource.Name, "icons/job.svg")
}

func (adapter adapterBatchV1Job) Connect(kgraph KubeGraph, source *cgraph.Node, targetName string) (*cgraph.Edge, error) {
	return kgraph.linkNode(source, adapter.GetType(), targetName)
}

func (adapter adapterBatchV1Job) Configure(kgraph KubeGraph) error {
	for resourceName, resourceObject := range kgraph.objects[adapter.GetType()] {
		resource := resourceObject.(*batchv1.Job)
		resourceNode, ok := kgraph.nodes[adapter.GetType()][resourceName]
		if !ok {
			return fmt.Errorf("node %s not found", resourceName)
		}

		podAdapter := adapters[reflect.TypeOf(&v1.Pod{})]

		podAdapter.Connect(kgraph, resourceNode, resource.Name)
	}
	return nil
}
