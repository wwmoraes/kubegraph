package adapter

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/emicklei/dot"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type dummyResource struct {
	metaV1.TypeMeta
	metaV1.ObjectMeta
}

func (d dummyResource) GetObjectKind() schema.ObjectKind {
	return nil
}

func (d dummyResource) DeepCopyObject() runtime.Object {
	return dummyResource{}
}

type dummyAdapter struct {
	ResourceData
}

func (thisAdapter dummyAdapter) tryCastObject(obj runtime.Object) (*dummyResource, error) {
	casted, ok := obj.(*dummyResource)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}

func (thisAdapter dummyAdapter) GetType() reflect.Type {
	return thisAdapter.ResourceType
}

func (thisAdapter dummyAdapter) Create(statefulGraph StatefulGraph, obj runtime.Object) (*dot.Node, error) {
	resource, err := thisAdapter.tryCastObject(obj)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s~%s", resource.APIVersion, resource.Kind, resource.Name)
	return statefulGraph.AddStyledNode(thisAdapter.GetType(), obj, name, resource.Name, "icons/unknown.svg")
}

func (thisAdapter dummyAdapter) Connect(statefulGraph StatefulGraph, source *dot.Node, targetName string) (*dot.Edge, error) {
	return statefulGraph.LinkNode(source, thisAdapter.GetType(), targetName)
}

func (thisAdapter dummyAdapter) Configure(statefulGraph StatefulGraph) error {
	return nil
}

var dummyAdapterInstance = &dummyAdapter{
	ResourceData{
		ResourceType: reflect.TypeOf(&dummyResource{}),
	},
}

func TestRegister(t *testing.T) {
	t.Run("successfully register valid adapter", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The code did panic")
			}
		}()
		Register(dummyAdapterInstance)
	})
	t.Run("panic on register existing adapter", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		Register(dummyAdapterInstance)
	})
}

func TestGet(t *testing.T) {
	t.Run("get existing adapter", func(t *testing.T) {
		got, err := Get(dummyAdapterInstance.ResourceType)
		if err != nil {
			t.Errorf("Get() error = %v, wantErr %v", err, nil)
			return
		}
		if !reflect.DeepEqual(got, dummyAdapterInstance) {
			t.Errorf("Get() = %v, want %v", got, dummyAdapterInstance)
		}
	})
	t.Run("try get un-existant adapter", func(t *testing.T) {
		resourceType := reflect.TypeOf(struct{}{})
		wantErr := fmt.Errorf("type %s has no adapter registered", resourceType.String())
		got, err := Get(resourceType)
		if err == nil {
			t.Errorf("Get() error = %v, wantErr %v", err, wantErr)
			return
		}
		if got != nil {
			t.Errorf("Get() = %v, want %v", got, nil)
		}
	})
}

func TestGetAll(t *testing.T) {
	t.Run("get all adapters", func(t *testing.T) {
		want := map[reflect.Type]ResourceAdapter{
			dummyAdapterInstance.ResourceType: dummyAdapterInstance,
		}

		if got := GetAll(); !reflect.DeepEqual(got, want) {
			t.Errorf("GetAll() = %v, want %v", got, want)
		}
	})
}
