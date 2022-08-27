// Code generated by Kubegraph; DO NOT EDIT.

package v1

import (
	"fmt"
	"reflect"

	"github.com/wwmoraes/kubegraph/internal/registry"
	coreV1 "k8s.io/api/core/v1"
)

// PodObject alias for coreV1.Pod
type PodObject = coreV1.Pod

// PodAdapterObjectType reflected type of *coreV1.Pod
var PodAdapterObjectType = reflect.TypeOf((*coreV1.Pod)(nil))

// GetPodAdapter retrieves the adapter that handles
// coreV1.Pod resources from the registry, if any
func GetPodAdapter() (registry.ResourceAdapter[*coreV1.Pod], error) {
	adapter, err := registry.GetAdapter[*coreV1.Pod]()
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter: %w", registry.ErrIncompatibleType)
	}

	return adapter, nil
}
