package v1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i k8s.io/api/rbac/v1 -n rbacV1 -t ClusterRole --icon c-role
//go:generate adapter -i k8s.io/api/rbac/v1 -n rbacV1 -t ClusterRoleBinding --icon crb
//go:generate adapter -i k8s.io/api/rbac/v1 -n rbacV1 -t Role --icon role
//go:generate adapter -i k8s.io/api/rbac/v1 -n rbacV1 -t RoleBinding --icon rb

//go:generate -command dependency go run github.com/wwmoraes/kubegraph/cmd/adapter dep
//go:generate dependency -i k8s.io/api/core/v1 -n coreV1 -t ServiceAccount
//go:generate dependency -i k8s.io/api/rbac/v1beta1 -n rbacV1beta1 -t Role --prefixed
//go:generate dependency -i k8s.io/api/rbac/v1beta1 -n rbacV1beta1 -t ClusterRole --prefixed
