package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type secretAdapter struct {
	adapter.Resource
}

func init() {
	adapter.MustRegister(&secretAdapter{
		adapter.NewResource(
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
