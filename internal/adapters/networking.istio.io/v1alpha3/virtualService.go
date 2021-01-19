package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	networkV1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
)

type virtualServiceAdapter struct {
	adapter.Resource
}

func init() {
	adapter.MustRegister(&virtualServiceAdapter{
		adapter.NewResource(
			reflect.TypeOf(&networkV1alpha3.VirtualService{}),
			"icons/unknown.svg",
		),
	})
}

func (thisAdapter *virtualServiceAdapter) tryCastObject(obj runtime.Object) (*networkV1alpha3.VirtualService, error) {
	casted, ok := obj.(*networkV1alpha3.VirtualService)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
