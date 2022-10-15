package application

import (
	"time"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

type OrderUseCase interface {
	GetById(id string) (*domain.Order, error)
	Open(dateReservFrom, dateReservTo time.Time, stationFromId, stationToId, categoryId, carModel, policyId string) error
	Confirm(id string, dateFrom time.Time) error
	Close(id string, discount, tax float32, dateTo time.Time, km uint64) error
	Cancel(id string) error
}

type orderUseCase struct {
	orderRepo OrderRepository
	orderSvc  OrderService
}

func NewOrderUseCase(orderRepo OrderRepository, orderSvc OrderService) *orderUseCase {
	return &orderUseCase{
		orderRepo: orderRepo,
		orderSvc:  orderSvc,
	}
}

func (uc orderUseCase) GetById(id string) (*domain.Order, error) {
	if err := validation.ValidId(id); err != nil {
		return nil, ErrInvalidId
	}

	order, err := uc.orderRepo.FindOne(id)

	if err != nil {
		return nil, ErrNotFoundOrder
	}

	return order, nil
}

func (uc orderUseCase) Open(dateReservFrom, dateReservTo time.Time, stationFromId, stationToId, categoryId, carModel, policyId string) error {
	policy, err := uc.orderSvc.GetPolicy(categoryId, carModel, policyId)
	if err != nil {
		return ErrInvalidEntity
	}

	car, err := uc.orderSvc.GetCar(stationFromId, carModel)
	if err != nil {
		return ErrInvalidEntity
	}

	newOrder, err := domain.NewOrder(dateReservFrom, dateReservTo, *car, stationFromId, stationToId, *policy)
	if err != nil {
		return ErrInvalidEntity
	}

	if err := uc.orderRepo.Save(*newOrder); err != nil {
		return ErrInvalidOrder
	}

	return nil
}

func (uc orderUseCase) Confirm(id string, dateFrom time.Time) error {
	order, err := uc.orderRepo.FindOne(id)

	if err != nil {
		return ErrInvalidOrder
	}

	if err := order.Confirm(dateFrom); err != nil {
		return ErrInvalidOrder
	}

	if err := uc.orderRepo.Save(*order); err != nil {
		return ErrInvalidOrder
	}

	return nil
}

func (uc orderUseCase) Close(id string, discount, tax float32, dateTo time.Time, km uint64) error {
	order, err := uc.orderRepo.FindOne(id)

	if err != nil {
		return ErrInvalidOrder
	}

	if err := order.Close(discount, tax, dateTo, km); err != nil {
		return ErrInvalidOrder
	}

	if err := uc.orderRepo.Save(*order); err != nil {
		return ErrInvalidOrder
	}

	return nil
}

func (uc orderUseCase) Cancel(id string) error {
	order, err := uc.orderRepo.FindOne(id)

	if err != nil {
		return ErrInvalidOrder
	}

	if err := order.Cancel(); err != nil {
		return ErrInvalidOrder
	}

	if err := uc.orderRepo.Save(*order); err != nil {
		return ErrInvalidOrder
	}

	return nil
}
