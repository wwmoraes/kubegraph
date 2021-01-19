package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

func (thisAdapter *configMapAdapter) tryCastObject(obj runtime.Object) (*coreV1.ConfigMap, error) {
	casted, ok := obj.(*coreV1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
