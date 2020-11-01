package loader

import (
	apiExtensionsApiServerScheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeSchema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	aggregatorScheme "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/scheme"
)

func getDecoder() func(data []byte, defaults *runtimeSchema.GroupVersionKind, into runtime.Object) (runtime.Object, *runtimeSchema.GroupVersionKind, error) {
	_ = aggregatorScheme.AddToScheme(scheme.Scheme)
	_ = apiExtensionsApiServerScheme.AddToScheme(scheme.Scheme)
	return scheme.Codecs.UniversalDeserializer().Decode
}