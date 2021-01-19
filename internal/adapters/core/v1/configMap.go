package v1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	coreV1 "k8s.io/api/core/v1"
)

type configMapAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(NewConfigMapAdapter())
}

func NewConfigMapAdapter() registry.Adapter {
	return &configMapAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&coreV1.ConfigMap{}),
			"icons/cm.svg",
		),
	}
}
