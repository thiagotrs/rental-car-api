package eventhandler

import (
	"errors"
	"testing"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/adapters/repository"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
)

func newStationFixture() *domain.Station {
	return &domain.Station{
		ID:         "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		Name:       "Station 1",
		Address:    "Farway Av.",
		Complement: "45, Ap. 50",
		State:      "Polar",
		City:       "Nort City",
		Cep:        "20778990",
		Capacity:   100,
		Idle:       1,
	}
}

func TestStationEventHandler_HandleCarAdded(t *testing.T) {
	dispatcher := events.NewEventDispatcher()

	stations := []domain.Station{*newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationEH := NewStationEventHandler(stationUC)

	dispatcher.Register(
		events.EventHandlerFunc(stationEH.HandleCarAdded),
		domain.CarAdded{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.CarAdded{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
			},
			errWant: nil,
		},
		{
			name: "incorrect station id input",
			eventArg: domain.CarAdded{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: application.ErrInvalidStation,
		},
		{
			name: "incorrect event input",
			eventArg: domain.CarParked{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
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

func TestStationEventHandler_HandleCarParked(t *testing.T) {
	dispatcher := events.NewEventDispatcher()

	stations := []domain.Station{*newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationEH := NewStationEventHandler(stationUC)

	dispatcher.Register(
		events.EventHandlerFunc(stationEH.HandleCarParked),
		domain.CarParked{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.CarParked{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
			},
			errWant: nil,
		},
		{
			name: "incorrect station id input",
			eventArg: domain.CarParked{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: application.ErrInvalidStation,
		},
		{
			name: "incorrect event input",
			eventArg: domain.CarAdded{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
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

func TestStationEventHandler_HandleCarUnderMaintenance(t *testing.T) {
	dispatcher := events.NewEventDispatcher()

	stations := []domain.Station{*newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationEH := NewStationEventHandler(stationUC)

	dispatcher.Register(
		events.EventHandlerFunc(stationEH.HandleCarUnderMaintenance),
		domain.CarUnderMaintenance{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.CarUnderMaintenance{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
				CarStatus: domain.Parked,
			},
			errWant: nil,
		},
		{
			name: "incorrect station id input",
			eventArg: domain.CarUnderMaintenance{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				CarStatus: domain.Parked,
			},
			errWant: application.ErrInvalidStation,
		},
		{
			name: "incorrect event input",
			eventArg: domain.CarParked{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
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

func TestStationEventHandler_HandleCarInTransfer(t *testing.T) {
	dispatcher := events.NewEventDispatcher()

	stations := []domain.Station{*newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationEH := NewStationEventHandler(stationUC)

	dispatcher.Register(
		events.EventHandlerFunc(stationEH.HandleCarInTransfer),
		domain.CarInTransfer{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.CarInTransfer{
				ID:            "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationIdFrom: stations[0].ID,
				StationIdTo:   "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: nil,
		},
		{
			name: "incorrect station id input",
			eventArg: domain.CarInTransfer{
				ID:            "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationIdFrom: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationIdTo:   "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: application.ErrInvalidStation,
		},
		{
			name: "incorrect event input",
			eventArg: domain.CarAdded{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
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

func TestStationEventHandler_HandleSyncCarParked(t *testing.T) {
	dispatcher := events.NewEventDispatcher()

	stations := []domain.Station{*newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationEH := NewStationEventHandler(stationUC)

	dispatcher.Register(
		events.EventHandlerFunc(stationEH.HandleSyncCarParked),
		domain.SyncCarParked{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.SyncCarParked{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
			},
			errWant: nil,
		},
		{
			name: "incorrect station id input",
			eventArg: domain.SyncCarParked{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: application.ErrInvalidStation,
		},
		{
			name: "incorrect event input",
			eventArg: domain.CarAdded{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
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

func TestStationEventHandler_HandleCarReserved(t *testing.T) {
	dispatcher := events.NewEventDispatcher()

	stations := []domain.Station{*newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationEH := NewStationEventHandler(stationUC)

	dispatcher.Register(
		events.EventHandlerFunc(stationEH.HandleSyncCarReserved),
		domain.SyncCarReserved{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.SyncCarReserved{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
			},
			errWant: nil,
		},
		{
			name: "incorrect station id input",
			eventArg: domain.SyncCarReserved{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: application.ErrInvalidStation,
		},
		{
			name: "incorrect event input",
			eventArg: domain.CarAdded{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: stations[0].ID,
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
