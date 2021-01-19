package v1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	coreV1 "k8s.io/api/core/v1"
)

type persistentVolumeAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(NewPersistentVolumeAdapter())
}

func NewPersistentVolumeAdapter() registry.Adapter {
	return &persistentVolumeAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&coreV1.PersistentVolume{}),
			"icons/persistentVolume.svg",
		),
	}
}
