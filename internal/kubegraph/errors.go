package kubegraph

import "errors"

var (
	ErrUnregistered = errors.New("unregistered type")
	ErrNodeNotFound = errors.New("node not found")
)
