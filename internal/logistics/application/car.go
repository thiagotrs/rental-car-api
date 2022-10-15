package application

import (
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

type CarUseCase interface {
	SearchCars(search SearchCarParams) []domain.Car
	GetCarById(id string) (*domain.Car, error)
	AddCar(age uint16, km uint64, plate, document, stationId, model, make string) error
	DeleteCar(id string) error
	MoveCarToMaintenance(id, stationId string, km uint64) error
	ParkCar(id, stationId string, km uint64) error
	TransferCar(id, stationId string) error
	SyncParkCar(id, stationId string, km uint64) error
	SyncCarToTransit(id string) error
	SyncReserveCar(id string) error
}

type carUseCase struct {
	carRepo     CarRepository
	stationRepo StationReadRepository
}

func NewCarUseCase(carRepo CarRepository, stationRepo StationReadRepository) *carUseCase {
	return &carUseCase{
		carRepo:     carRepo,
		stationRepo: stationRepo,
	}
}

type SearchCarParams struct {
	Plate     string `json:"plate" db:"plate"`
	Document  string `json:"document" db:"document"`
	Model     string `json:"model" db:"model"`
	Make      string `json:"make" db:"make"`
	StationId string `json:"stationId" db:"stationId"`
	Age       uint16 `json:"age" db:"age"`
	KM        uint64 `json:"km" db:"km"`
	Status    uint   `json:"status" db:"status"`
	// Limit     uint
	// Offset    uint
}

func (uc carUseCase) SearchCars(params SearchCarParams) []domain.Car {
	return uc.carRepo.Find(params)
}

func (uc carUseCase) GetCarById(id string) (*domain.Car, error) {
	if err := validation.ValidId(id); err != nil {
		return nil, ErrInvalidId
	}

	car, err := uc.carRepo.FindOne(id)

	if err != nil {
		return nil, ErrNotFoundCar
	}

	return car, nil
}

func (uc carUseCase) AddCar(age uint16, km uint64, plate, document, stationId, model, make string) error {
	s, err := uc.stationRepo.FindOne(stationId)
	if err != nil {
		return ErrInvalidEntity
	}

	if (s.Capacity - s.Idle) == 0 {
		return ErrStationMaxCapacity
	}

	newCar, err := domain.NewCar(age, km, plate, document, stationId, model, make)
	if err != nil {
		return ErrInvalidEntity
	}

	if err := uc.carRepo.Save(*newCar); err != nil {
		return ErrInvalidCar
	}

	return nil
}

func (uc carUseCase) DeleteCar(id string) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	car, err := uc.carRepo.FindOne(id)
	if err != nil {
		return ErrInvalidCar
	}

	if car.Status != domain.Maintenance {
		return ErrCarNotInMaintenance
	}

	if err := uc.carRepo.Delete(id); err != nil {
		return ErrInvalidCar
	}

	return nil
}

func (uc carUseCase) MoveCarToMaintenance(id, stationId string, km uint64) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	car, err := uc.carRepo.FindOne(id)
	if err != nil {
		return ErrInvalidCar
	}

	if err := car.ToMaintenance(stationId, km); err != nil {
		return ErrInvalidMaintenance
	}

	if err := uc.carRepo.Save(*car); err != nil {
		return ErrInvalidCar
	}

	return nil
}

func (uc carUseCase) ParkCar(id, stationId string, km uint64) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	s, err := uc.stationRepo.FindOne(stationId)
	if err != nil {
		return ErrInvalidEntity
	}

	if (s.Capacity - s.Idle) == 0 {
		return ErrStationMaxCapacity
	}

	car, err := uc.carRepo.FindOne(id)
	if err != nil {
		return ErrInvalidCar
	}

	if err := car.Park(stationId, km); err != nil {
		return ErrInvalidPark
	}

	if err := uc.carRepo.Save(*car); err != nil {
		return err
		// return ErrInvalidCar
	}

	return nil
}

func (uc carUseCase) TransferCar(id, stationId string) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	s, err := uc.stationRepo.FindOne(stationId)
	if err != nil {
		return ErrInvalidEntity
	}

	if (s.Capacity - s.Idle) == 0 {
		return ErrStationMaxCapacity
	}

	car, err := uc.carRepo.FindOne(id)
	if err != nil {
		return ErrInvalidCar
	}

	if err := car.Transfer(stationId); err != nil {
		return ErrInvalidTransfer
	}

	if err := uc.carRepo.Save(*car); err != nil {
		return ErrInvalidCar
	}

	return nil
}

func (uc carUseCase) SyncParkCar(id, stationId string, km uint64) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	s, err := uc.stationRepo.FindOne(stationId)
	if err != nil {
		return ErrInvalidEntity
	}

	if (s.Capacity - s.Idle) == 0 {
		return ErrStationMaxCapacity
	}

	car, err := uc.carRepo.FindOne(id)
	if err != nil {
		return ErrInvalidCar
	}

	if err := car.Park(stationId, km); err != nil {
		return ErrInvalidPark
	}

	car.Events = []events.Event{}

	if err := uc.carRepo.Save(*car); err != nil {
		return ErrInvalidCar
	}

	return nil
}

func (uc carUseCase) SyncCarToTransit(id string) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	car, err := uc.carRepo.FindOne(id)
	if err != nil {
		return ErrInvalidCar
	}

	if err := car.ToTransit(); err != nil {
		return ErrInvalidTransit
	}

	if err := uc.carRepo.Save(*car); err != nil {
		return ErrInvalidCar
	}

	return nil
}

func (uc carUseCase) SyncReserveCar(id string) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	car, err := uc.carRepo.FindOne(id)
	if err != nil {
		return ErrInvalidCar
	}

	if err := car.Reserve(); err != nil {
		return ErrInvalidTransfer
	}

	if err := uc.carRepo.Save(*car); err != nil {
		return ErrInvalidCar
	}

	return nil
}
