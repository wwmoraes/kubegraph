package v1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	coreV1 "k8s.io/api/core/v1"
)

type secretAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&secretAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&coreV1.Secret{}),
			"icons/secret.svg",
		),
	})
}
