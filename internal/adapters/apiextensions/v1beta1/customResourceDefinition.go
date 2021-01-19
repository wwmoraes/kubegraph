package v1beta1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	apiExtensionsV1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

type customResourceDefinitionAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&customResourceDefinitionAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&apiExtensionsV1beta1.CustomResourceDefinition{}),
			"icons/crd.svg",
		),
	})
}
