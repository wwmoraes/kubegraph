package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	apiExtensionsV1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type customResourceDefinitionAdapter struct {
	adapter.Resource
}

func init() {
	adapter.MustRegister(&customResourceDefinitionAdapter{
		adapter.NewResource(
			reflect.TypeOf(&apiExtensionsV1beta1.CustomResourceDefinition{}),
			"icons/crd.svg",
		),
	})
}

func (thisAdapter *customResourceDefinitionAdapter) tryCastObject(obj runtime.Object) (*apiExtensionsV1beta1.CustomResourceDefinition, error) {
	casted, ok := obj.(*apiExtensionsV1beta1.CustomResourceDefinition)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
