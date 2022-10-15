package repository

import (
	"sync"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
)

type stationRepositoryInMemory struct {
	stations map[string]domain.Station
	*sync.RWMutex
}

func NewStationRepositoryInMemory(stations []domain.Station) *stationRepositoryInMemory {
	stationsMap := make(map[string]domain.Station)
	for _, v := range stations {
		stationsMap[v.ID] = v
	}
	return &stationRepositoryInMemory{stationsMap, &sync.RWMutex{}}
}

func (repo stationRepositoryInMemory) FindAll() []domain.Station {
	repo.Lock()
	defer repo.Unlock()

	stations := []domain.Station{}
	for k := range repo.stations {
		stations = append(stations, repo.stations[k])
	}

	return stations
}

func (repo stationRepositoryInMemory) FindOne(id string) (*domain.Station, error) {
	repo.Lock()
	defer repo.Unlock()

	s, exists := repo.stations[id]
	if !exists {
		return &s, application.ErrNotFoundStation
	}

	return &s, nil
}

func (repo *stationRepositoryInMemory) Save(station domain.Station) error {
	repo.Lock()
	defer repo.Unlock()

	repo.stations[station.ID] = station

	return nil
}

func (repo *stationRepositoryInMemory) Delete(id string) error {
	repo.Lock()
	defer repo.Unlock()

	if _, exists := repo.stations[id]; !exists {
		return application.ErrInvalidStation
	}

	delete(repo.stations, id)

	return nil
}
