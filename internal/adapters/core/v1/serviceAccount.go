package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type serviceAccountAdapter struct {
	adapter.Resource
}

func init() {
	adapter.MustRegister(&serviceAccountAdapter{
		adapter.NewResource(
			reflect.TypeOf(&coreV1.ServiceAccount{}),
			"icons/sa.svg",
		),
	})
}

func (thisAdapter *serviceAccountAdapter) tryCastObject(obj runtime.Object) (*coreV1.ServiceAccount, error) {
	casted, ok := obj.(*coreV1.ServiceAccount)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
