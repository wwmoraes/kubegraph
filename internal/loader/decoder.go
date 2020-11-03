package loader

import (
	istioScheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	apiExtensionsApiServerScheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeSchema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	aggregatorScheme "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/scheme"
	applicationV1beta1 "sigs.k8s.io/application/api/v1beta1"
)

func init() {
	_ = aggregatorScheme.AddToScheme(scheme.Scheme)
	_ = apiExtensionsApiServerScheme.AddToScheme(scheme.Scheme)
	_ = applicationV1beta1.AddToScheme(scheme.Scheme)
	_ = istioScheme.AddToScheme(scheme.Scheme)
	// add any extra schemes here
}

func getDecoder() func(data []byte, defaults *runtimeSchema.GroupVersionKind, into runtime.Object) (runtime.Object, *runtimeSchema.GroupVersionKind, error) {
	return scheme.Codecs.UniversalDeserializer().Decode
}
