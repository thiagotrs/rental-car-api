package domain

import (
	"fmt"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

type Station struct {
	ID         string `json:"id" validate:"required,uuid4"`
	Name       string `json:"name" validate:"required"`
	Address    string `json:"address" validate:"required"`
	Complement string `json:"complement" validate:"required"`
	State      string `json:"state" validate:"required"`
	City       string `json:"city" validate:"required"`
	Cep        string `json:"cep" validate:"required,len=8"`
	Capacity   uint   `json:"capacity" validate:"required,gt=0"`
	Idle       uint   `json:"idle"`
}

func NewStation(name, address, complement, state, city, cep string, capacity, idle uint) (*Station, error) {
	station := &Station{
		ID:         validation.NewId(),
		Name:       name,
		Address:    address,
		Complement: complement,
		State:      state,
		City:       city,
		Cep:        cep,
		Capacity:   capacity,
		Idle:       idle,
	}

	if err := validation.ValidateEntity(station); err != nil {
		return nil, fmt.Errorf("%w\n%v", ErrInvalidEntity, err)
	}

	if idle > capacity {
		return nil, ErrInvalidCapacity
	}

	return station, nil
}

func (s *Station) SetCapacity(n uint) error {
	if n == 0 || s.Idle > n {
		return ErrInvalidCapacity
	}

	s.Capacity = n

	return nil
}

func (s *Station) SetIdle(n uint) error {
	if n > s.Capacity {
		return ErrInvalidIdle
	}

	s.Idle = n

	return nil
}
