module github.com/wwmoraes/kubegraph

go 1.15

require (
	github.com/emicklei/dot v0.15.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	istio.io/client-go v1.8.0-alpha.2
	k8s.io/api v0.19.3
	k8s.io/apiextensions-apiserver v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.19.3
	k8s.io/kube-aggregator v0.19.3
	sigs.k8s.io/application v0.8.3
)

// replace github.com/emicklei/dot => ../dot
// replace github.com/emicklei/dot => github.com/wwmoraes/dot master
replace github.com/emicklei/dot => github.com/wwmoraes/dot v0.15.1-0.20201107152930-203731997226
