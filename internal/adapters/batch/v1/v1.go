package v1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i k8s.io/api/batch/v1 -n batchV1 -t Job --icon job

//go:generate -command dependency go run github.com/wwmoraes/kubegraph/cmd/adapter dep
//go:generate dependency -i k8s.io/api/core/v1 -n coreV1 -t Pod
