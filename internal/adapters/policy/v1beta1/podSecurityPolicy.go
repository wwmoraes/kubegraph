package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type podSecurityPolicyAdapter struct {
	adapter.ResourceData
}

func init() {
	adapter.MustRegister(&podSecurityPolicyAdapter{
		adapter.NewResourceData(
			reflect.TypeOf(&policyV1beta1.PodSecurityPolicy{}),
			"icons/psp.svg",
		),
	})
}

func (thisAdapter *podSecurityPolicyAdapter) tryCastObject(obj runtime.Object) (*policyV1beta1.PodSecurityPolicy, error) {
	casted, ok := obj.(*policyV1beta1.PodSecurityPolicy)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
