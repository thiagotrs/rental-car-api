package domain

import (
	"fmt"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

type CarStatus uint

const (
	Maintenance CarStatus = iota + 1
	Transit
	Parked
	Reserved
	Transfer
)

type Car struct {
	ID        string         `json:"id" validate:"required,uuid4"`
	Age       uint16         `json:"age" validate:"required,min=1900,max=2100"`
	Plate     string         `json:"plate" validate:"required"`
	Document  string         `json:"document" validate:"required"`
	Model     string         `json:"model" validate:"required"`
	Make      string         `json:"make" validate:"required"`
	StationId string         `json:"stationId" validate:"uuid4" db:"stationId"`
	KM        uint64         `json:"km" validate:"required"`
	Status    CarStatus      `json:"status" validate:"required"`
	Events    []events.Event `json:"-" bson:"-"`
}

func NewCar(age uint16, km uint64, plate, document, stationId, model, make string) (*Car, error) {
	newCar := &Car{
		ID:        validation.NewId(),
		Age:       age,
		Plate:     plate,
		Document:  document,
		Model:     model,
		Make:      make,
		StationId: stationId,
		KM:        km,
		Status:    Parked,
	}

	if err := validation.ValidateEntity(newCar); err != nil {
		return nil, fmt.Errorf("%w\n%v", ErrInvalidEntity, err)
	}

	newCar.Events = append(newCar.Events, CarAdded{
		ID:        newCar.ID,
		StationId: newCar.StationId,
	})

	return newCar, nil
}

func (c *Car) ToMaintenance(stationId string, km uint64) error {
	if c.Status != Transfer && c.Status != Parked {
		return ErrInvalidMaintenance
	}

	if km < c.KM {
		return ErrInvalidMaintenance
	}

	c.Status = Maintenance
	c.KM = km

	c.Events = append(c.Events, CarUnderMaintenance{
		ID:        c.ID,
		StationId: c.StationId,
		CarStatus: Parked,
	})

	c.StationId = stationId

	return nil
}

func (c *Car) Transfer(stationId string) error {
	if c.Status != Parked {
		return ErrInvalidTransfer
	}

	c.Status = Transfer

	c.Events = append(c.Events, CarInTransfer{
		ID:            c.ID,
		StationIdFrom: c.StationId,
		StationIdTo:   stationId,
	})

	c.StationId = stationId

	return nil
}

func (c *Car) Park(stationId string, km uint64) error {
	if c.Status != Maintenance && c.Status != Transit && c.Status != Reserved && c.Status != Transfer {
		return ErrInvalidPark
	}

	if km < c.KM {
		return ErrInvalidPark
	}

	c.Status = Parked
	c.StationId = stationId
	c.KM = km

	c.Events = append(c.Events, CarParked{
		ID:        c.ID,
		StationId: c.StationId,
		KM:        c.KM,
	})

	return nil
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
