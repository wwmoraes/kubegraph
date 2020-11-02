module github.com/wwmoraes/kubegraph

go 1.15

require (
	cloud.google.com/go v0.51.0 // indirect
	github.com/Azure/go-autorest/autorest v0.9.6 // indirect
	github.com/go-bindata/go-bindata v3.1.2+incompatible // indirect
	github.com/goccy/go-graphviz v0.0.8
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	istio.io/client-go v1.8.0-alpha.2
	k8s.io/api v0.19.3
	k8s.io/apiextensions-apiserver v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.19.3
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kube-aggregator v0.19.3
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73 // indirect
	sigs.k8s.io/application v0.8.3
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)

replace github.com/goccy/go-graphviz => ../go-graphviz
