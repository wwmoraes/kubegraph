package v1beta1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	networkV1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

type gatewayAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&gatewayAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&networkV1alpha3.Gateway{}),
			"icons/unknown.svg",
		),
	})
}
