package application

import (
	"errors"
	"fmt"

	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/domain"
)

var (
	ErrInvalidEntity = fmt.Errorf("%w", domain.ErrInvalidEntity)
	ErrInvalidModel  = fmt.Errorf("%w", domain.ErrInvalidModel)
	ErrInvalidPolicy = fmt.Errorf("%w", domain.ErrInvalidPolicy)

	ErrInvalidCategory  = errors.New("invalid category")
	ErrNotFoundCategory = errors.New("not found category")
	ErrNotFoundPolicy   = errors.New("not found policy")
	ErrInvalidId        = errors.New("invalid category id")
)
