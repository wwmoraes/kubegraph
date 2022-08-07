package v1beta1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1 -n apiregistrationV1beta1 -t APIService
