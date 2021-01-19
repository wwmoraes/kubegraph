package v1beta1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	policyV1beta1 "k8s.io/api/policy/v1beta1"
)

type podSecurityPolicyAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&podSecurityPolicyAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&policyV1beta1.PodSecurityPolicy{}),
			"icons/psp.svg",
		),
	})
}
