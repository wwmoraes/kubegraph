package v1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	coreV1 "k8s.io/api/core/v1"
)

type serviceAccountAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&serviceAccountAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&coreV1.ServiceAccount{}),
			"icons/sa.svg",
		),
	})
}
