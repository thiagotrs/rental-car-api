package application

import (
	"errors"
	"fmt"

	"github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

var (
	ErrInvalidEntity = fmt.Errorf("%w", domain.ErrInvalidEntity)

	ErrInvalidOrder  = errors.New("invalid order")
	ErrNotFoundOrder = errors.New("not found order")
	ErrInvalidId     = errors.New("invalid order id")

	ErrInvalidPolicy = errors.New("invalid policy")
	ErrInvalidCar    = errors.New("invalid car")
)
