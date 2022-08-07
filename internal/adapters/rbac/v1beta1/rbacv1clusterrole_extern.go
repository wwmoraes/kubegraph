// Code generated by Kubegraph; DO NOT EDIT.

package v1beta1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	rbacV1 "k8s.io/api/rbac/v1"
)

// RbacV1ClusterRoleObject alias for rbacV1.ClusterRole
type RbacV1ClusterRoleObject = rbacV1.ClusterRole

// RbacV1ClusterRoleAdapterObjectType reflected type of *rbacV1.ClusterRole
var RbacV1ClusterRoleAdapterObjectType = reflect.TypeOf((*rbacV1.ClusterRole)(nil))

// GetRbacV1ClusterRoleAdapter retrieves the adapter that handles
// rbacV1.ClusterRole resources from the registry, if any
func GetRbacV1ClusterRoleAdapter() (registry.ResourceAdapter[*rbacV1.ClusterRole], error) {
	adapter, err := registry.GetAdapter[*rbacV1.ClusterRole]()
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter: %w", registry.ErrIncompatibleType)
	}

	return adapter, nil
}
