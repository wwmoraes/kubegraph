// +build wireinject

package kubegraph

import (
	"github.com/google/wire"
	"github.com/wwmoraes/dot"
	"github.com/wwmoraes/kubegraph/internal/adapters"
	"github.com/wwmoraes/kubegraph/internal/registry"
)

func InitializeKubegraph(optionsFn ...dot.GraphOptionFn) (*Kubegraph, error) {
	panic(wire.Build(NewKubegraph, dot.New, registry.Instance, adapters.GetDecoder))
}
