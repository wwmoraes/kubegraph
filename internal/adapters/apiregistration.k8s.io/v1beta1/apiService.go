package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	"k8s.io/apimachinery/pkg/runtime"
	apiregistrationV1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
)

type apiServiceAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&apiServiceAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&apiregistrationV1beta1.APIService{}),
			"icons/unknown.svg",
		),
	})
}

func (thisAdapter *apiServiceAdapter) tryCastObject(obj runtime.Object) (*apiregistrationV1beta1.APIService, error) {
	casted, ok := obj.(*apiregistrationV1beta1.APIService)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
