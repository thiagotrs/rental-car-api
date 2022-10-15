package repository

import (
	"sync"

	"github.com/thiagotrs/rentalcar-ddd/internal/rental/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

type orderRepositoryInMemory struct {
	orders map[string]domain.Order
	*sync.RWMutex
}

func NewOrderRepositoryInMemory(orders []domain.Order) *orderRepositoryInMemory {
	ordersMap := make(map[string]domain.Order)
	for _, v := range orders {
		ordersMap[v.ID] = v
	}
	return &orderRepositoryInMemory{ordersMap, &sync.RWMutex{}}
}

func (repo orderRepositoryInMemory) FindOne(id string) (*domain.Order, error) {
	repo.Lock()
	defer repo.Unlock()

	s, exists := repo.orders[id]
	if !exists {
		return &s, application.ErrNotFoundOrder
	}

	return &s, nil
}

func (repo *orderRepositoryInMemory) Save(order domain.Order) error {
	repo.Lock()
	defer repo.Unlock()

	repo.orders[order.ID] = order

	return nil
}
