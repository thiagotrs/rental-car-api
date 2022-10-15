package eventhandler

import (
	"encoding/json"
	"errors"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/broker"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

type orderEventHandler struct {
	broker broker.Publisher
}

func NewOrderEventHandler(broker broker.Publisher) *orderEventHandler {
	return &orderEventHandler{broker}
}

func (eh orderEventHandler) HandleOpenedOrder(e events.Event) error {
	event, ok := e.(domain.OpenedOrder)

	if !ok {
		return errors.New("wrong event")
	}

	data, _ := json.Marshal(event)
	eh.broker.Publish(e.Name(), data)

	return nil
}

func (eh orderEventHandler) HandleConfirmedOrder(e events.Event) error {
	event, ok := e.(domain.ConfirmedOrder)

	if !ok {
		return errors.New("wrong event")
	}

	data, _ := json.Marshal(event)
	eh.broker.Publish(e.Name(), data)

	return nil
}

func (eh orderEventHandler) HandleClosedOrder(e events.Event) error {
	event, ok := e.(domain.ClosedOrder)

	if !ok {
		return errors.New("wrong event")
	}

	data, _ := json.Marshal(event)
	eh.broker.Publish(e.Name(), data)

	return nil
}

func (eh orderEventHandler) HandleCanceledOrder(e events.Event) error {
	event, ok := e.(domain.CanceledOrder)

	if !ok {
		return errors.New("wrong event")
	}

	data, _ := json.Marshal(event)
	eh.broker.Publish(e.Name(), data)

	return nil
}
