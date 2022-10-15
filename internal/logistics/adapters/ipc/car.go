package ipc

import (
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/ipc"
)

type carIPC struct {
	carUC application.CarUseCase
}

func NewCarIPC(carUC application.CarUseCase) *carIPC {
	return &carIPC{carUC}
}

func (uc carIPC) GetCar(stationId, carModel string) (*ipc.CarData, error) {
	cars := uc.carUC.SearchCars(application.SearchCarParams{
		Model:     carModel,
		StationId: stationId,
	})

	if len(cars) == 0 {
		return nil, application.ErrNotFoundCar
	}

	carData := &ipc.CarData{
		ID:        cars[0].ID,
		Age:       cars[0].Age,
		Plate:     cars[0].Plate,
		Document:  cars[0].Document,
		Model:     cars[0].Model,
		Make:      cars[0].Make,
		StationId: cars[0].StationId,
		KM:        cars[0].KM,
		Status:    uint(cars[0].Status),
	}

	return carData, nil
}
