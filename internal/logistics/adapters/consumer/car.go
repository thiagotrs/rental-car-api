package consumer

import (
	"encoding/json"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
)

type Topic string

const (
	OrderOpened    Topic = "order.opened"
	OrderConfirmed Topic = "order.confirmed"
	OrderClosed    Topic = "order.closed"
	OrderCanceled  Topic = "order.canceled"
)

type openedOrderMsg struct {
	ID        string `json:"id"`
	CarId     string `json:"carId"`
	StationId string `json:"stationId"`
}

type confirmedOrderMsg struct {
	ID    string `json:"id"`
	CarId string `json:"carId"`
}

type canceledOrderMsg struct {
	ID        string `json:"id"`
	CarId     string `json:"carId"`
	StationId string `json:"stationId"`
	FinalKM   uint64 `json:"finalKM"`
}

type closedOrderMsg struct {
	ID        string `json:"id"`
	CarId     string `json:"carId"`
	StationId string `json:"stationId"`
	FinalKM   uint64 `json:"finalKM"`
}

type orderConsumer struct {
	disp events.Dispatcher
}

func NewOrderConsumer(disp events.Dispatcher) *orderConsumer {
	return &orderConsumer{disp}
}

func (c *orderConsumer) ConsumeOpenedOrder(data interface{}) {
	if orderB, ok := data.([]byte); ok {
		var order openedOrderMsg
		json.Unmarshal(orderB, &order)

		c.disp.Dispatch([]events.Event{domain.SyncCarReserved{ID: order.CarId, StationId: order.StationId}})
	}
}

func (c *orderConsumer) ConsumeConfirmedOrder(data interface{}) {
	if orderB, ok := data.([]byte); ok {
		var order confirmedOrderMsg
		json.Unmarshal(orderB, &order)

		c.disp.Dispatch([]events.Event{domain.SyncCarInTransit{ID: order.CarId}})
	}
}

func (c *orderConsumer) ConsumeCanceledOrder(data interface{}) {
	if orderB, ok := data.([]byte); ok {
		var order canceledOrderMsg
		json.Unmarshal(orderB, &order)

		c.disp.Dispatch([]events.Event{domain.SyncCarParked{ID: order.CarId, StationId: order.StationId, KM: order.FinalKM}})
	}
}

func (c *orderConsumer) ConsumeClosedOrder(data interface{}) {
	if orderB, ok := data.([]byte); ok {
		var order closedOrderMsg
		json.Unmarshal(orderB, &order)

		c.disp.Dispatch([]events.Event{domain.SyncCarParked{ID: order.CarId, StationId: order.StationId, KM: order.FinalKM}})
	}
}
