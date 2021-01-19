package v1beta1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	networkV1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type virtualServiceAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&virtualServiceAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&networkV1alpha3.VirtualService{}),
			"icons/unknown.svg",
		),
	})
}
