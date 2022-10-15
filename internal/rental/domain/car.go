package domain

import (
	"fmt"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

type CarStatus uint

const (
	Maintenance CarStatus = iota + 1
	Transit
	Parked
	Reserved
)

type Car struct {
	ID        string    `json:"id" validate:"required,uuid4" db:"id"`
	Age       uint16    `json:"age" validate:"required,min=1900,max=2100" db:"age"`
	Plate     string    `json:"plate" validate:"required" db:"plate"`
	Document  string    `json:"document" validate:"required" db:"document"`
	CarModel  string    `json:"carModel" validate:"required" db:"carModel"`
	InitialKM uint64    `json:"initialKM" validate:"required" db:"initialKM"`
	FinalKM   uint64    `json:"finalKM,omitempty" db:"finalKM"`
	Status    CarStatus `json:"status" validate:"required" db:"status"`
	StationId string    `json:"stationId" validate:"uuid4" db:"stationId"`
}

func NewCar(id string, age uint16, plate, document string, carModel string, initialKM, finalKM uint64, status CarStatus, stationId string) (*Car, error) {
	car := &Car{
		ID:        id,
		Age:       age,
		Plate:     plate,
		Document:  document,
		CarModel:  carModel,
		InitialKM: initialKM,
		FinalKM:   finalKM,
		Status:    status,
		StationId: stationId,
	}

	if err := validation.ValidateEntity(car); err != nil {
		return nil, fmt.Errorf("%w\n%v", ErrInvalidEntity, err)
	}

	return car, nil
}

func (c *Car) ToTransit() error {
	if c.Status != Reserved {
		return ErrInvalidTransit
	}

	c.Status = Transit

	return nil
}

func (c *Car) Reserve() error {
	if c.Status != Parked {
		return ErrInvalidReserve
	}

	c.Status = Reserved

	return nil
}

func (c *Car) Park(finalKM uint64, stationId string) error {
	if c.Status == Parked {
		return ErrInvalidPark
	}

	if finalKM < c.InitialKM {
		return ErrInvalidPark
	}

	c.Status = Parked
	c.FinalKM = finalKM
	c.StationId = stationId

	return nil
}
