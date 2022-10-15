package eventhandler

import (
	"errors"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
)

type stationEventHandler struct {
	stationUC application.StationUseCase
}

func NewStationEventHandler(stationUC application.StationUseCase) *stationEventHandler {
	return &stationEventHandler{stationUC}
}

func (h stationEventHandler) HandleCarAdded(e events.Event) error {
	event, ok := e.(domain.CarAdded)

	if !ok {
		return errors.New("wrong event")
	}

	if err := h.stationUC.ChangeStationIdle(event.StationId, "add"); err != nil {
		return err
	}

	return nil
}

func (h stationEventHandler) HandleCarParked(e events.Event) error {
	event, ok := e.(domain.CarParked)

	if !ok {
		return errors.New("wrong event")
	}

	if err := h.stationUC.ChangeStationIdle(event.StationId, "add"); err != nil {
		return err
	}

	return nil
}

func (h stationEventHandler) HandleCarUnderMaintenance(e events.Event) error {
	event, ok := e.(domain.CarUnderMaintenance)

	if !ok {
		return errors.New("wrong event")
	}

	if event.CarStatus != domain.Parked {
		return nil
	}

	if err := h.stationUC.ChangeStationIdle(event.StationId, "del"); err != nil {
		return err
	}

	return nil
}

func (h stationEventHandler) HandleCarInTransfer(e events.Event) error {
	event, ok := e.(domain.CarInTransfer)

	if !ok {
		return errors.New("wrong event")
	}

	if err := h.stationUC.ChangeStationIdle(event.StationIdFrom, "del"); err != nil {
		return err
	}

	return nil
}

func (h stationEventHandler) HandleSyncCarParked(e events.Event) error {
	event, ok := e.(domain.SyncCarParked)

	if !ok {
		return errors.New("wrong event")
	}

	if err := h.stationUC.ChangeStationIdle(event.StationId, "add"); err != nil {
		return err
	}

	return nil
}

func (h stationEventHandler) HandleSyncCarReserved(e events.Event) error {
	event, ok := e.(domain.SyncCarReserved)

	if !ok {
		return errors.New("wrong event")
	}

	if err := h.stationUC.ChangeStationIdle(event.StationId, "del"); err != nil {
		return err
	}

	return nil
}
