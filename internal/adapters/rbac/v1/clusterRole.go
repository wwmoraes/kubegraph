package v1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	rbacV1 "k8s.io/api/rbac/v1"
)

type clusterRoleAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&clusterRoleAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&rbacV1.ClusterRole{}),
			"icons/c-role.svg",
		),
	})
}
