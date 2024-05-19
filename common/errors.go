package common

import "errors"

var (
	ErrNoItem = errors.New("Order must have at least one item")
)
