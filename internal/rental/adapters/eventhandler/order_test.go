package eventhandler

import (
	"errors"
	"testing"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/broker"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

func Run(t *testing.T, channel <-chan interface{}) {
	t.Helper()
	for data := range channel {
		go func(data interface{}) {
			orderB, ok := data.([]byte)
			if !ok {
				t.Error("invalid data type")
			}
			t.Log(string(orderB))
		}(data)
	}
}

func TestOrderEventHandler_HandleOpenedOrder(t *testing.T) {
	pubsub := broker.NewPubSub()
	dispatcher := events.NewEventDispatcher()
	orderEH := NewOrderEventHandler(pubsub)

	channel := pubsub.Subscribe(domain.OpenedOrder{}.Name())
	go Run(t, channel)
	defer pubsub.Close()

	dispatcher.Register(
		events.EventHandlerFunc(orderEH.HandleOpenedOrder),
		domain.OpenedOrder{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.OpenedOrder{
				ID:    "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				CarId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			errWant: nil,
		},
		{
			name: "incorrect event input",
			eventArg: domain.ConfirmedOrder{
				ID:    "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				CarId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			errWant: events.ErrNoneHandler,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := dispatcher.Dispatch([]events.Event{tc.eventArg})

			if !errors.Is(err, tc.errWant) {
				t.Error("wrong err", err, tc.errWant)
			}
		})
	}
}

func TestOrderEventHandler_HandleConfirmedOrder(t *testing.T) {
	pubsub := broker.NewPubSub()
	dispatcher := events.NewEventDispatcher()
	orderEH := NewOrderEventHandler(pubsub)

	channel := pubsub.Subscribe(domain.ConfirmedOrder{}.Name())
	go Run(t, channel)
	defer pubsub.Close()

	dispatcher.Register(
		events.EventHandlerFunc(orderEH.HandleConfirmedOrder),
		domain.ConfirmedOrder{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.ConfirmedOrder{
				ID:    "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				CarId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			errWant: nil,
		},
		{
			name: "incorrect event input",
			eventArg: domain.OpenedOrder{
				ID:    "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				CarId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			errWant: events.ErrNoneHandler,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := dispatcher.Dispatch([]events.Event{tc.eventArg})

			if !errors.Is(err, tc.errWant) {
				t.Error("wrong err", err, tc.errWant)
			}
		})
	}
}

func TestOrderEventHandler_HandleClosedOrder(t *testing.T) {
	pubsub := broker.NewPubSub()
	dispatcher := events.NewEventDispatcher()
	orderEH := NewOrderEventHandler(pubsub)

	channel := pubsub.Subscribe(domain.ClosedOrder{}.Name())
	go Run(t, channel)
	defer pubsub.Close()

	dispatcher.Register(
		events.EventHandlerFunc(orderEH.HandleClosedOrder),
		domain.ClosedOrder{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.ClosedOrder{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				CarId:     "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
				StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
				FinalKM:   13000,
			},
			errWant: nil,
		},
		{
			name: "incorrect event input",
			eventArg: domain.OpenedOrder{
				ID:    "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				CarId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			errWant: events.ErrNoneHandler,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := dispatcher.Dispatch([]events.Event{tc.eventArg})

			if !errors.Is(err, tc.errWant) {
				t.Error("wrong err", err, tc.errWant)
			}
		})
	}
}

func TestOrderEventHandler_HandleCanceledOrder(t *testing.T) {
	pubsub := broker.NewPubSub()
	dispatcher := events.NewEventDispatcher()
	orderEH := NewOrderEventHandler(pubsub)

	channel := pubsub.Subscribe(domain.CanceledOrder{}.Name())

	go Run(t, channel)
	defer pubsub.Close()

	dispatcher.Register(
		events.EventHandlerFunc(orderEH.HandleCanceledOrder),
		domain.CanceledOrder{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.CanceledOrder{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				CarId:     "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
				StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
				FinalKM:   13000,
			},
			errWant: nil,
		},
		{
			name: "incorrect event input",
			eventArg: domain.OpenedOrder{
				ID:    "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				CarId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			errWant: events.ErrNoneHandler,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := dispatcher.Dispatch([]events.Event{tc.eventArg})

			if !errors.Is(err, tc.errWant) {
				t.Error("wrong err", err, tc.errWant)
			}
		})
	}
}
