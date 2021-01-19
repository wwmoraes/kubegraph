package v1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	admissionregistrationV1beta1 "k8s.io/api/admissionregistration/v1beta1"
)

type validatingWebhookConfigurationAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&validatingWebhookConfigurationAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&admissionregistrationV1beta1.ValidatingWebhookConfiguration{}),
			"icons/unknown.svg",
		),
	})
}
