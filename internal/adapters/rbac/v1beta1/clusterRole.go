package v1beta1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
)

type clusterRoleAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&clusterRoleAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&rbacV1beta1.ClusterRole{}),
			"icons/c-role.svg",
		),
	})
}
