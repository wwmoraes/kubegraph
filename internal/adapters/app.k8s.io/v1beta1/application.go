package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	"k8s.io/apimachinery/pkg/runtime"
	applicationV1beta1 "sigs.k8s.io/application/api/v1beta1"
)

// applicationAdapter a kubegraph adapter to render an specific kubernetes resource
type applicationAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&applicationAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&applicationV1beta1.Application{}),
			"icons/unknown.svg",
		),
	})
}

func (thisAdapter *applicationAdapter) tryCastObject(obj runtime.Object) (*applicationV1beta1.Application, error) {
	casted, ok := obj.(*applicationV1beta1.Application)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
