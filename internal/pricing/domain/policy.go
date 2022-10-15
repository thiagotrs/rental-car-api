package domain

import (
	"fmt"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

type Unit uint

const (
	PerKM Unit = iota + 1
	PerDay
	PerWeek
)

type Policy struct {
	ID      string  `json:"id" db:"id"`
	Name    string  `json:"name" validate:"required" db:"name"`
	Price   float32 `json:"price" validate:"required,gte=0" db:"price"`
	Unit    Unit    `json:"unit" validate:"required,gt=0" db:"unit"`
	MinUnit uint    `json:"minUnit" validate:"required,gt=0" db:"minUnit"`
}

func NewPolicy(name string, price float32, unit Unit, minUnit uint) (*Policy, error) {
	policy := &Policy{
		ID:      validation.NewId(),
		Name:    name,
		Price:   price,
		Unit:    unit,
		MinUnit: minUnit,
	}

	if err := validation.ValidateEntity(policy); err != nil {
		return nil, fmt.Errorf("%w\n%v", ErrInvalidEntity, err)
	}

	return policy, nil
}
