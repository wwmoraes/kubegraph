package v1beta1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i k8s.io/api/policy/v1beta1 -n policyV1beta1 -t PodDisruptionBudget --icon pdb
//go:generate adapter -i k8s.io/api/policy/v1beta1 -n policyV1beta1 -t PodSecurityPolicy --icon psp

//go:generate -command dependency go run github.com/wwmoraes/kubegraph/cmd/adapter dep
//go:generate dependency -i k8s.io/api/core/v1 -n coreV1 -t Pod
