package service

import (
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/ipc"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

type orderServiceIPC struct {
	logistics ipc.LogisticsIPC
	pricing   ipc.PricingIPC
}

func NewOrderServiceIPC(logistics ipc.LogisticsIPC, pricing ipc.PricingIPC) *orderServiceIPC {
	return &orderServiceIPC{logistics, pricing}
}

func (svc orderServiceIPC) GetPolicy(categoryId, carModel, policyId string) (*domain.Policy, error) {
	policy, err := svc.pricing.GetPolicy(categoryId, carModel, policyId)
	if err != nil {
		return nil, application.ErrInvalidPolicy
	}

	return &domain.Policy{
		ID:         policy.ID,
		Name:       policy.Name,
		Price:      policy.Price,
		Unit:       domain.Unit(policy.Unit),
		MinUnit:    policy.MinUnit,
		CarModel:   carModel,
		CategoryId: categoryId,
	}, nil
}

func (svc orderServiceIPC) GetCar(stationId, modelId string) (*domain.Car, error) {
	car, err := svc.logistics.GetCar(stationId, modelId)
	if err != nil {
		return nil, application.ErrInvalidCar
	}

	return &domain.Car{
		ID:        car.ID,
		Age:       car.Age,
		Plate:     car.Plate,
		Document:  car.Document,
		CarModel:  car.Model,
		InitialKM: car.KM,
		Status:    domain.CarStatus(car.Status),
		StationId: car.StationId,
	}, nil
}
