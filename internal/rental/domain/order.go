package domain

import (
	"fmt"
	"time"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

type OrderStatus uint

const (
	Opened OrderStatus = iota + 1
	Confirmed
	Closed
	Canceled
)

type Order struct {
	ID             string         `json:"id" validate:"required,uuid4" db:"id"`
	DateFrom       *time.Time     `json:"dateFrom,omitempty" db:"dateFrom"`
	DateTo         *time.Time     `json:"dateTo,omitempty" db:"dateTo"`
	DateReservFrom time.Time      `json:"dateReservFrom" validate:"required" db:"dateReservFrom"`
	DateReservTo   time.Time      `json:"dateReservTo" validate:"required" db:"dateReservTo"`
	Status         OrderStatus    `json:"status" validate:"required" db:"status"`
	Car            Car            `json:"car" validate:"required"`
	StationFromId  string         `json:"stationFromId" validate:"required,uuid4" db:"stationFromId"`
	StationToId    string         `json:"stationToId" validate:"required,uuid4" db:"stationToId"`
	Policy         Policy         `json:"policy" validate:"required"`
	Discount       float32        `json:"discount,omitempty" db:"discount"`
	Tax            float32        `json:"tax,omitempty" db:"tax"`
	Events         []events.Event `json:"-" bson:"-"`
}

func NewOrder(
	dateReservFrom time.Time,
	dateReservTo time.Time,
	car Car,
	stationFromId string,
	stationToId string,
	policy Policy,
) (*Order, error) {
	if dateReservFrom.After(dateReservTo) {
		return nil, ErrInvalidReservedDate
	}

	if car.StationId != stationFromId {
		return nil, ErrInvalidCarStation
	}

	if car.CarModel != policy.CarModel {
		return nil, ErrInvalidCarStation
	}

	if err := car.Reserve(); err != nil {
		return nil, err
	}

	newOrder := &Order{
		ID:             validation.NewId(),
		DateReservFrom: dateReservFrom,
		DateReservTo:   dateReservTo,
		Status:         Opened,
		Car:            car,
		StationFromId:  stationFromId,
		StationToId:    stationToId,
		Policy:         policy,
	}

	if err := validation.ValidateEntity(newOrder); err != nil {
		return nil, fmt.Errorf("%w\n%v", ErrInvalidEntity, err)
	}

	newOrder.Events = append(newOrder.Events, OpenedOrder{
		ID:        newOrder.ID,
		CarId:     car.ID,
		StationId: stationFromId,
	})

	return newOrder, nil
}

func (r *Order) Confirm(dateFrom time.Time) error {
	if r.Status != Opened {
		return ErrClose
	}

	if r.DateReservFrom.After(dateFrom) || dateFrom.After(r.DateReservTo) {
		return ErrIvalidConfirmDate
	}

	if err := r.Car.ToTransit(); err != nil {
		return err
	}

	r.Status = Confirmed
	r.DateFrom = &dateFrom

	r.Events = append(r.Events, ConfirmedOrder{
		ID:    r.ID,
		CarId: r.Car.ID,
	})

	return nil
}

func (r *Order) Close(discount, tax float32, dateTo time.Time, finalKM uint64) error {
	if r.Status != Confirmed {
		return ErrClose
	}

	if r.DateFrom != nil {
		if dateTo.Before(*r.DateFrom) {
			return ErrIvalidCloseDate
		}
	}

	if tax < 0 {
		return ErrIvalidCloseTax
	}

	if discount < 0 {
		return ErrIvalidCloseDiscount
	}

	if err := r.Car.Park(finalKM, r.StationToId); err != nil {
		return err
	}

	r.Status = Closed
	r.DateTo = &dateTo
	r.Discount = discount
	r.Tax = tax

	r.Events = append(r.Events, ClosedOrder{
		ID:        r.ID,
		CarId:     r.Car.ID,
		StationId: r.StationToId,
		FinalKM:   finalKM,
	})

	return nil
}

func (r *Order) Cancel() error {
	if r.Status != Opened {
		return ErrClose
	}

	if err := r.Car.Park(r.Car.InitialKM, r.Car.StationId); err != nil {
		return err
	}

	r.Status = Canceled

	r.Events = append(r.Events, CanceledOrder{
		ID:        r.ID,
		CarId:     r.Car.ID,
		StationId: r.Car.StationId,
		FinalKM:   r.Car.FinalKM,
	})

	return nil
}
