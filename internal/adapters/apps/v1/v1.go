package v1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i k8s.io/api/apps/v1 -n appsV1 -t DaemonSet --icon ds
//go:generate adapter -i k8s.io/api/apps/v1 -n appsV1 -t Deployment --icon deploy

//go:generate -command dependency go run github.com/wwmoraes/kubegraph/cmd/adapter dep
//go:generate dependency -i k8s.io/api/core/v1 -n coreV1 -t Pod
