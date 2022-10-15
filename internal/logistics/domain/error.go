package domain

import "errors"

var (
	ErrInvalidEntity   = errors.New("invalid entity")
	ErrInvalidCapacity = errors.New("invalid capacity")
	ErrInvalidCurrCars = errors.New("invalid current cars number")
	ErrInvalidIdle     = errors.New("invalid idle cars number")

	ErrInvalidMaintenance = errors.New("invalid maintenance")
	ErrInvalidTransit     = errors.New("invalid transit")
	ErrInvalidTransfer    = errors.New("invalid transfer")
	ErrInvalidReserve     = errors.New("invalid reserve")
	ErrInvalidReservation = errors.New("invalid reservation")
	ErrInvalidPark        = errors.New("invalid park")
)
