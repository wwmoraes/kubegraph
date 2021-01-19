package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	rbacV1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
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

func (thisAdapter *clusterRoleAdapter) tryCastObject(obj runtime.Object) (*rbacV1beta1.ClusterRole, error) {
	casted, ok := obj.(*rbacV1beta1.ClusterRole)
	if !ok {
		return nil, fmt.Errorf("unable to cast object %s to %s", reflect.TypeOf(obj), thisAdapter.GetType().String())
	}

	return casted, nil
}
