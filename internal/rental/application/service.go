package application

import "github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"

type PolicyService interface {
	GetPolicy(categoryId, modelId, policyId string) (*domain.Policy, error)
}

type CarService interface {
	GetCar(stationId, modelId string) (*domain.Car, error)
}

type OrderService interface {
	PolicyService
	CarService
}
