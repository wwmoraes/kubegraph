package pod

import (
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/adapter"
	coreV1 "k8s.io/api/core/v1"
)

func init() {
	if err := Register(reflect.TypeOf(&coreV1.Pod{}), "icons/pod.svg"); err != nil {
		panic(err)
	}
}

func New(resourceData adapter.Resource) adapter.Resource {
	return &adapterResource{
		resourceData,
	}
}
