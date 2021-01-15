package adapters

import (
	istioScheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	apiExtensionsApiServerScheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeSchema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	aggregatorScheme "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/scheme"
	applicationV1beta1 "sigs.k8s.io/application/api/v1beta1"

	_ "github.com/wwmoraes/kubegraph/internal/adapters/admissionregistration.k8s.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/apiextensions"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/apiregistration.k8s.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/app.k8s.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/apps"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/autoscaling"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/batch"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/core"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/extensions"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/networking.istio.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/policy"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/rbac"
	// _ "github.com/wwmoraes/kubegraph/internal/adapters/dummy"
)

type DecodeFn func(data []byte, defaults *runtimeSchema.GroupVersionKind, into runtime.Object) (runtime.Object, *runtimeSchema.GroupVersionKind, error)

func init() {
	_ = aggregatorScheme.AddToScheme(scheme.Scheme)
	_ = apiExtensionsApiServerScheme.AddToScheme(scheme.Scheme)
	_ = applicationV1beta1.AddToScheme(scheme.Scheme)
	_ = istioScheme.AddToScheme(scheme.Scheme)
	// add any extra schemes here
}

// GetDecoder returns an instance of an universal deserializer decoder with all
// supported schemas already added
func GetDecoder() DecodeFn {
	return scheme.Codecs.UniversalDeserializer().Decode
}
