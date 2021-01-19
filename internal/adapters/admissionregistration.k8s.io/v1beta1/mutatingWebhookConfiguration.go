package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	admissionregistrationV1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
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

func (thisAdapter *mutatingWebhookConfigurationAdapter) tryCastObject(obj runtime.Object) (*admissionregistrationV1beta1.MutatingWebhookConfiguration, error) {
	casted, ok := obj.(*admissionregistrationV1beta1.MutatingWebhookConfiguration)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
