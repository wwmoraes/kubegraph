package adapter

import (
	"reflect"
)

// ResourceData data used by kubernetes resource adapters
type ResourceData struct {
	ResourceType reflect.Type
}
