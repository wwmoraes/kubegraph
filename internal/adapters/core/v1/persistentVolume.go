package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type persistentVolumeAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.MustRegister(NewPersistentVolumeAdapter())
}

func NewPersistentVolumeAdapter() adapter.ResourceTransformer {
	return &persistentVolumeAdapter{
		adapter.NewResourceData(
			reflect.TypeOf(&coreV1.PersistentVolume{}),
			"icons/persistentVolume.svg",
		),
	}
}

func (thisAdapter *persistentVolumeAdapter) tryCastObject(obj runtime.Object) (*coreV1.PersistentVolume, error) {
	casted, ok := obj.(*coreV1.PersistentVolume)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
