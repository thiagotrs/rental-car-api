package domain

import "errors"

var (
	ErrInvalidEntity = errors.New("invalid entity")
	ErrInvalidModel  = errors.New("invalid model")
	ErrInvalidPolicy = errors.New("invalid policy")
)
