// +build wireinject

package pod

import (
	"reflect"

	"github.com/google/wire"
	"github.com/wwmoraes/kubegraph/internal/adapter"
)

func Register(resourceType reflect.Type, iconPath string) error {
	panic(wire.Build(adapter.RegisterSet, New))
}
