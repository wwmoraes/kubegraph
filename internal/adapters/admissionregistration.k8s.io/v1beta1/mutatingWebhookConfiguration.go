package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	admissionregistrationV1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type mutatingWebhookConfigurationAdapter struct {
	adapter.Resource
}

func init() {
	adapter.MustRegister(&mutatingWebhookConfigurationAdapter{
		adapter.NewResource(
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
