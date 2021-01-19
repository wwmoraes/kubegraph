package v1

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	rbacV1 "k8s.io/api/rbac/v1"
)

type roleAdapter struct {
	registry.Adapter
}

func init() {
	registry.MustRegister(&roleAdapter{
		registry.NewAdapter(
			reflect.TypeOf(&rbacV1.Role{}),
			"icons/role.svg",
		),
	})
}
