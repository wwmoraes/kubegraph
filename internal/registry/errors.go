package registry

import (
	"errors"
)

// ErrAdapterNotFound means the requested adapter was not found on the registry
var ErrAdapterNotFound = errors.New("adapter not found for such type")

// ErrAdapterAlreadyRegistered means there's already an adapter registered for
// a given reflected type
var ErrAdapterAlreadyRegistered = errors.New("only one adapter should be registered per type")

// ErrAdapterRegisteredElsewhere means the adapter has been registered within
// another registry instance, and thus cannot be registered on another one
var ErrAdapterRegisteredElsewhere = errors.New("adapter is already registered on another registry instance")
