package v1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	admissionregistrationV1beta1 "k8s.io/api/admissionregistration/v1beta1"
)

type mutatingWebhookConfigurationAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&mutatingWebhookConfigurationAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&admissionregistrationV1beta1.MutatingWebhookConfiguration{}),
			"icons/unknown.svg",
		),
	})
}
