//go:build wireinject
// +build wireinject

package kubegraph

import (
	"github.com/google/wire"
	"github.com/wwmoraes/dot"
	"github.com/wwmoraes/kubegraph/internal/adapters"
)

func InitializeKubegraph(optionsFn ...dot.GraphOptionFn) (*Kubegraph, error) {
	panic(wire.Build(NewKubegraph, dot.New, adapters.GetDecoder))
}
