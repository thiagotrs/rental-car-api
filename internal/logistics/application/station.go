package application

import (
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

type StationUseCase interface {
	GetStations() []domain.Station
	GetStationById(id string) (*domain.Station, error)
	AddStation(name, address, complement, state, city, cep string, capacity, idle uint) error
	DeleteStation(id string) error
	ChangeStationCapacity(id string, capacity uint) error
	ChangeStationIdle(id string, op string) error
}

type stationUseCase struct {
	stationRepo StationRepository
}

func NewStationUseCase(stationRepo StationRepository) *stationUseCase {
	return &stationUseCase{stationRepo}
}

func (uc stationUseCase) GetStations() []domain.Station {
	return uc.stationRepo.FindAll()
}

func (uc stationUseCase) GetStationById(id string) (*domain.Station, error) {
	if err := validation.ValidId(id); err != nil {
		return nil, ErrInvalidId
	}

	station, err := uc.stationRepo.FindOne(id)

	if err != nil {
		return nil, ErrNotFoundStation
	}

	return station, nil
}

func (uc stationUseCase) AddStation(name, address, complement, state, city, cep string, capacity, idle uint) error {
	newStation, err := domain.NewStation(name, address, complement, state, city, cep, capacity, idle)

	if err != nil {
		return ErrInvalidEntity
	}

	if err := uc.stationRepo.Save(*newStation); err != nil {
		return ErrInvalidStation
	}

	return nil
}

func (uc stationUseCase) DeleteStation(id string) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	station, err := uc.stationRepo.FindOne(id)

	if err != nil {
		return ErrInvalidStation
	}

	if station.Idle > 0 {
		return ErrStationHasCars
	}

	if err := uc.stationRepo.Delete(id); err != nil {
		return ErrInvalidStation
	}

	return nil
}

func (uc stationUseCase) ChangeStationCapacity(id string, capacity uint) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	station, err := uc.stationRepo.FindOne(id)

	if err != nil {
		return ErrInvalidStation
	}

	if err := station.SetCapacity(capacity); err != nil {
		return ErrInvalidCapacity
	}

	if err := uc.stationRepo.Save(*station); err != nil {
		return ErrInvalidStation
	}

	return nil
}

func (uc stationUseCase) ChangeStationIdle(id string, op string) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	m, err := uc.stationRepo.FindOne(id)
	if err != nil {
		return ErrInvalidStation
	}

	switch op {
	case "add":
		if err := m.SetIdle(m.Idle + 1); err != nil {
			return err
		}
	case "del":
		if err := m.SetIdle(m.Idle - 1); err != nil {
			return err
		}
	default:
		return ErrInvalidIdle
	}

	if err := uc.stationRepo.Save(*m); err != nil {
		return ErrInvalidStation
	}

	return nil
}
