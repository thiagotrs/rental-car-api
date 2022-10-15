package eventhandler

import (
	"errors"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
)

type carEventHandler struct {
	carUC application.CarUseCase
}

func NewCarEventHandler(carUC application.CarUseCase) *carEventHandler {
	return &carEventHandler{carUC}
}

func (h carEventHandler) HandleSyncCarParked(e events.Event) error {
	event, ok := e.(domain.SyncCarParked)

	if !ok {
		return errors.New("wrong event")
	}

	if err := h.carUC.SyncParkCar(event.ID, event.StationId, event.KM); err != nil {
		return err
	}

	return nil
}

func (h carEventHandler) HandleSyncCarReserved(e events.Event) error {
	event, ok := e.(domain.SyncCarReserved)

	if !ok {
		return errors.New("wrong event")
	}

	if err := h.carUC.SyncReserveCar(event.ID); err != nil {
		return err
	}

	return nil
}

func (h carEventHandler) HandleSyncCarInTransit(e events.Event) error {
	event, ok := e.(domain.SyncCarInTransit)

	if !ok {
		return errors.New("wrong event")
	}

	if err := h.carUC.SyncCarToTransit(event.ID); err != nil {
		return err
	}

	return nil
}
