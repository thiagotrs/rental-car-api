package application

import (
	"errors"
	"fmt"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
)

var (
	ErrInvalidId      = errors.New("invalid id")
	ErrStationHasCars = errors.New("station has already cars")

	ErrInvalidEntity   = fmt.Errorf("%w", domain.ErrInvalidEntity)
	ErrInvalidCapacity = fmt.Errorf("%w", domain.ErrInvalidCapacity)
	ErrInvalidCurrCars = fmt.Errorf("%w", domain.ErrInvalidCurrCars)
	ErrInvalidIdle     = fmt.Errorf("%w", domain.ErrInvalidIdle)

	ErrNotFoundStation = errors.New("station not found")
	ErrInvalidStation  = errors.New("invalid station")

	ErrStationMaxCapacity  = errors.New("station with max capacity")
	ErrCarNotInMaintenance = errors.New("car is not in maintenance")
	ErrNotFoundCar         = errors.New("not found car")
	ErrInvalidCar          = errors.New("invalid car")
	ErrInvalidModel        = errors.New("invalid model")
	ErrNotFoundModel       = errors.New("model not found")
	ErrModelHasCars        = errors.New("model has already cars")

	ErrInvalidMaintenance = fmt.Errorf("%w", domain.ErrInvalidMaintenance)
	ErrInvalidTransit     = fmt.Errorf("%w", domain.ErrInvalidTransit)
	ErrInvalidReserve     = fmt.Errorf("%w", domain.ErrInvalidReserve)
	ErrInvalidPark        = fmt.Errorf("%w", domain.ErrInvalidPark)
	ErrInvalidTransfer    = fmt.Errorf("%w", domain.ErrInvalidTransfer)
)
