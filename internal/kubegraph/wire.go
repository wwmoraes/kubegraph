// +build wireinject

package kubegraph

import (
	"github.com/google/wire"
	"github.com/wwmoraes/dot"
	"github.com/wwmoraes/kubegraph/internal/adapter"
)

func InitializeKubegraph(optionsFn ...dot.GraphOptionFn) (*KubeGraph, error) {
	panic(wire.Build(NewKubegraph, dot.New, adapter.RegistryInstance))
}
