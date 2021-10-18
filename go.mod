module github.com/wwmoraes/kubegraph

go 1.15

require (
	github.com/google/wire v0.5.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/wwmoraes/dot v0.4.1
	istio.io/client-go v1.11.4
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/kube-aggregator v0.21.1
	sigs.k8s.io/application v0.8.3
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)

// replace github.com/wwmoraes/dot => github.com/wwmoraes/dot master
// replace github.com/wwmoraes/dot => ../dot
