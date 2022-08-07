package v1alpha3

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i istio.io/client-go/pkg/apis/networking/v1alpha3 -n networkV1alpha3 -t Gateway
//go:generate adapter -i istio.io/client-go/pkg/apis/networking/v1alpha3 -n networkV1alpha3 -t VirtualService
