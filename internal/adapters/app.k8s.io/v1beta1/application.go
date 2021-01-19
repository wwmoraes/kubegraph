package v1beta1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	applicationV1beta1 "sigs.k8s.io/application/api/v1beta1"
)

// applicationAdapter a kubegraph adapter to render an specific kubernetes resource
type applicationAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&applicationAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&applicationV1beta1.Application{}),
			"icons/unknown.svg",
		),
	})
}
