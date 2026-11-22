package application

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("concurrent modification or duplicate")
)
