package adapters

import (
	"k8s.io/apimachinery/pkg/runtime"
	runtimeSchema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

type DecodeFn func(data []byte, defaults *runtimeSchema.GroupVersionKind, into runtime.Object) (runtime.Object, *runtimeSchema.GroupVersionKind, error)

// GetDecoder returns an instance of an universal deserializer decoder with all
// supported schemas already added
func GetDecoder() DecodeFn {
	return scheme.Codecs.UniversalDeserializer().Decode
}
