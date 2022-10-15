package application

import "github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"

type CarReadRepository interface {
	Find(search SearchCarParams) []domain.Car
	FindOne(id string) (*domain.Car, error)
}

type CarWriteRepository interface {
	Save(car domain.Car) error
	Delete(id string) error
}

type CarRepository interface {
	CarReadRepository
	CarWriteRepository
}

type StationReadRepository interface {
	FindAll() []domain.Station
	FindOne(id string) (*domain.Station, error)
}

type StationWriteRepository interface {
	Save(station domain.Station) error
	Delete(id string) error
}

type StationRepository interface {
	StationReadRepository
	StationWriteRepository
}
