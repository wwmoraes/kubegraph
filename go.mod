module github.com/wwmoraes/kubegraph

go 1.15

require (
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/wwmoraes/dot v0.4.1
	istio.io/client-go v1.8.1
	k8s.io/api v0.20.1
	k8s.io/apiextensions-apiserver v0.20.1
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.1
	k8s.io/kube-aggregator v0.20.1
	sigs.k8s.io/application v0.8.3
)

// replace github.com/wwmoraes/dot => github.com/wwmoraes/dot master
// replace github.com/wwmoraes/dot => ../dot
