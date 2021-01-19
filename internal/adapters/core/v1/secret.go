package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

func (thisAdapter *secretAdapter) tryCastObject(obj runtime.Object) (*coreV1.Secret, error) {
	casted, ok := obj.(*coreV1.Secret)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
