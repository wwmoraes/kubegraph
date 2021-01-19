package v1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
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
