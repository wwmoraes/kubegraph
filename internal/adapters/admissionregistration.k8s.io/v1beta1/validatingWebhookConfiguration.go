package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	admissionregistrationV1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
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

func (thisAdapter *validatingWebhookConfigurationAdapter) tryCastObject(obj runtime.Object) (*admissionregistrationV1beta1.ValidatingWebhookConfiguration, error) {
	casted, ok := obj.(*admissionregistrationV1beta1.ValidatingWebhookConfiguration)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
