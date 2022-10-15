package domain

import (
	"fmt"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

type Unit uint

const (
	PerKm = iota + 1
	PerDay
	PerWeek
)

type Policy struct {
	ID         string  `json:"id" validate:"required,uuid4" db:"id"`
	Name       string  `json:"name" validate:"required" db:"name"`
	Price      float32 `json:"price" validate:"required,gte=0" db:"price"`
	Unit       Unit    `json:"unit" validate:"required,gt=0" db:"unit"`
	MinUnit    uint    `json:"minUnit" validate:"required" db:"minUnit"`
	CarModel   string  `json:"carModel" validate:"required" db:"carModel"`
	CategoryId string  `json:"categoryId" validate:"required,uuid4" db:"categoryId"`
}

func NewPolicy(id, name string, price float32, unit Unit, minUnit uint, carModel, categoryId string) (*Policy, error) {
	policy := &Policy{
		ID:         id,
		Name:       name,
		Price:      price,
		Unit:       unit,
		MinUnit:    minUnit,
		CarModel:   carModel,
		CategoryId: categoryId,
	}

	if err := validation.ValidateEntity(policy); err != nil {
		return nil, fmt.Errorf("%w\n%v", ErrInvalidEntity, err)
	}

	return policy, nil
}
