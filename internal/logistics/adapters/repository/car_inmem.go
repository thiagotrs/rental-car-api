package repository

import (
	"sync"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
)

type carRepositoryInMemory struct {
	cars map[string]domain.Car
	*sync.RWMutex
}

func NewCarRepositoryInMemory(cars []domain.Car) *carRepositoryInMemory {
	carsMap := make(map[string]domain.Car)
	for _, v := range cars {
		carsMap[v.ID] = v
	}
	return &carRepositoryInMemory{carsMap, &sync.RWMutex{}}
}

func (repo carRepositoryInMemory) Find(search application.SearchCarParams) []domain.Car {
	repo.Lock()
	defer repo.Unlock()

	cars := []domain.Car{}
	for k, v := range repo.cars {
		if v.Plate == search.Plate || v.Document == search.Document || v.Model == search.Model || v.Make == search.Make || v.StationId == search.StationId || v.Age == search.Age || v.KM == search.KM || v.Status == domain.CarStatus(search.Status) {
			cars = append(cars, repo.cars[k])
		}
	}

	return cars
}

func (repo carRepositoryInMemory) FindOne(id string) (*domain.Car, error) {
	repo.Lock()
	defer repo.Unlock()

	s, exists := repo.cars[id]
	if !exists {
		return &s, application.ErrNotFoundCar
	}

	return &s, nil
}

func (repo *carRepositoryInMemory) Save(car domain.Car) error {
	repo.Lock()
	defer repo.Unlock()

	repo.cars[car.ID] = car

	return nil
}

func (repo *carRepositoryInMemory) Delete(id string) error {
	repo.Lock()
	defer repo.Unlock()

	if _, exists := repo.cars[id]; !exists {
		return application.ErrInvalidCar
	}

	delete(repo.cars, id)

	return nil
}
