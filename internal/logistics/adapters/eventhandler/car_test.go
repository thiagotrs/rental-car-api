package eventhandler

import (
	"errors"
	"testing"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/adapters/repository"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
)

func newCarFixture() *domain.Car {
	return &domain.Car{
		ID:        "83369771-f9a4-48b7-b87b-463f19f7b187",
		Age:       2020,
		Plate:     "KST-9016",
		Document:  "abc.123.op-x",
		Model:     "Uno",
		Make:      "FIAT",
		StationId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		KM:        12000,
		Status:    domain.Parked,
	}
}

func TestCarEventHandler_HandleSyncCarParked(t *testing.T) {
	dispatcher := events.NewEventDispatcher()

	cars := []domain.Car{*newCarFixture()}
	cars[0].Status = domain.Maintenance
	stations := []domain.Station{*newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	carRepo := repository.NewCarRepositoryInMemory(cars)
	carUC := application.NewCarUseCase(carRepo, stationRepo)
	carEH := NewCarEventHandler(carUC)

	dispatcher.Register(
		events.EventHandlerFunc(carEH.HandleSyncCarParked),
		domain.SyncCarParked{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.SyncCarParked{
				ID:        cars[0].ID,
				StationId: cars[0].StationId,
				KM:        cars[0].KM + 50,
			},
			errWant: nil,
		},
		{
			name: "incorrect car id input",
			eventArg: domain.SyncCarParked{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				KM:        cars[0].KM + 50,
			},
			errWant: application.ErrInvalidEntity,
		},
		{
			name: "incorrect event input",
			eventArg: domain.CarAdded{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: cars[0].ID,
			},
			errWant: events.ErrNoneHandler,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cars[0].StationId = "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc"
			cars[0].Status = domain.Maintenance
			cars[0].KM = 12000

			err := dispatcher.Dispatch([]events.Event{tc.eventArg})

			if !errors.Is(err, tc.errWant) {
				t.Error("wrong err", err, tc.errWant)
			}
		})
	}
}

func TestCarEventHandler_HandleSyncCarReserved(t *testing.T) {
	dispatcher := events.NewEventDispatcher()

	cars := []domain.Car{*newCarFixture()}
	stations := []domain.Station{*newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	carRepo := repository.NewCarRepositoryInMemory(cars)
	carUC := application.NewCarUseCase(carRepo, stationRepo)
	carEH := NewCarEventHandler(carUC)

	dispatcher.Register(
		events.EventHandlerFunc(carEH.HandleSyncCarReserved),
		domain.SyncCarReserved{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.SyncCarReserved{
				ID:        cars[0].ID,
				StationId: cars[0].StationId,
			},
			errWant: nil,
		},
		{
			name: "incorrect car id input",
			eventArg: domain.SyncCarReserved{
				ID:        "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
				StationId: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: application.ErrInvalidCar,
		},
		{
			name: "incorrect event input",
			eventArg: domain.CarAdded{
				ID:        cars[0].ID,
				StationId: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: events.ErrNoneHandler,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cars[0].Status = domain.Parked

			err := dispatcher.Dispatch([]events.Event{tc.eventArg})

			if !errors.Is(err, tc.errWant) {
				t.Error("wrong err", err, tc.errWant)
			}
		})
	}
}

func TestCarEventHandler_HandleSyncCarInTransit(t *testing.T) {
	dispatcher := events.NewEventDispatcher()

	cars := []domain.Car{*newCarFixture()}
	cars[0].Status = domain.Reserved
	stations := []domain.Station{*newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	carRepo := repository.NewCarRepositoryInMemory(cars)
	carUC := application.NewCarUseCase(carRepo, stationRepo)
	carEH := NewCarEventHandler(carUC)

	dispatcher.Register(
		events.EventHandlerFunc(carEH.HandleSyncCarInTransit),
		domain.SyncCarInTransit{}.Name())

	testCases := []struct {
		name     string
		eventArg events.Event
		errWant  error
	}{
		{
			name: "correct input",
			eventArg: domain.SyncCarInTransit{
				ID: cars[0].ID,
			},
			errWant: nil,
		},
		{
			name: "incorrect car id input",
			eventArg: domain.SyncCarInTransit{
				ID: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: application.ErrInvalidCar,
		},
		{
			name: "incorrect event input",
			eventArg: domain.CarAdded{
				ID:        cars[0].ID,
				StationId: "5ce5a1a1-f324-4c8b-8c92-d7e820cbb238",
			},
			errWant: events.ErrNoneHandler,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cars[0].Status = domain.Reserved

			err := dispatcher.Dispatch([]events.Event{tc.eventArg})

			if !errors.Is(err, tc.errWant) {
				t.Error("wrong err", err, tc.errWant)
			}
		})
	}
}
