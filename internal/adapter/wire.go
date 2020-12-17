// +build wireinject

package adapter

import (
	"github.com/google/wire"
)

var RegisterSet = wire.NewSet(Register, NewResourceData)
