package v1beta1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i sigs.k8s.io/application/api/v1beta1 -n applicationV1beta1 -t Application
