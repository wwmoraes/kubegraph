package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	networkV1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
)

type gatewayAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&gatewayAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&networkV1alpha3.Gateway{}),
			"icons/unknown.svg",
		),
	})
}

func (thisAdapter *gatewayAdapter) tryCastObject(obj runtime.Object) (*networkV1alpha3.Gateway, error) {
	casted, ok := obj.(*networkV1alpha3.Gateway)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
