package registry

import (
	"errors"
)

var (
	// ErrAdapterNotFound means the requested adapter was not found on the registry
	ErrAdapterNotFound = errors.New("adapter not found")

	// ErrAdapterAlreadyRegistered means there's already an adapter registered for
	// a given reflected type
	ErrAdapterAlreadyRegistered = errors.New("adapter already registered")

	// ErrUnimplemented means an interface method has no overload implemented
	ErrUnimplemented = errors.New("Unimplemented")

	ErrIncompatibleType = errors.New("incompatible type")
)
