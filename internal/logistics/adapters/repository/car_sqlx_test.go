package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	// _ "github.com/lib/pq"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
)

func newCarFixture() *domain.Car {
	c, _ := domain.NewCar(2020, 12000, "KST-9016", "abc.123.op-x", "83369771-f9a4-48b7-b87b-463f19f7b187", "Uno", "FIAT")
	return c
}

func newCarInvalidFixture() *domain.Car {
	return &domain.Car{
		ID:        "",
		Age:       2020,
		Plate:     "KST-9016",
		Document:  "abc.123.op-x",
		Model:     "Uno",
		Make:      "FIAT",
		StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
		KM:        12000,
		Status:    domain.Parked,
	}
}

func GetCarDBConn(t *testing.T) *sqlx.DB {
	t.Helper()
	// connStr := "postgres://postgres:admin1234@localhost/?sslmode=disable"
	// db, err := sqlx.Open("postgres", connStr)
	db, err := sqlx.Open("sqlite3", "data/logistics_test.db")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func InitCarDB(t *testing.T, db *sqlx.DB, cars []domain.Car) {
	t.Helper()
	const saveCars = "INSERT INTO cars VALUES (:id, :age, :plate, :document, :model, :make, :stationId, :km, :status)"

	for _, s := range cars {
		if _, err := db.NamedExec(saveCars, s); err != nil {
			t.Fatal(err)
		}
	}
}

func ClearCarDB(t *testing.T, db *sqlx.DB) {
	t.Helper()
	const deleteAllCars = "DELETE FROM cars"

	if _, err := db.Exec(deleteAllCars); err != nil {
		t.Fatal(err)
	}
}

func SlicesEqual(a, b []domain.Car) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if reflect.DeepEqual(v, b[i]) {
			return false
		}
	}
	return true
}

type dispatcherMock struct {
	expectedDispatchErr error
	calls               map[string]uint
}

func (d *dispatcherMock) Dispatch(events []events.Event) error {
	d.calls["Dispatch"] = d.calls["Dispatch"] + 1
	return d.expectedDispatchErr
}
func (d *dispatcherMock) Register(h events.EventHandler, eventName string) {
	d.calls["Register"] = d.calls["Register"] + 1
}

func TestCarRepositorySqlx_Find(t *testing.T) {
	db := GetCarDBConn(t)
	defer db.Close()

	cars := []domain.Car{*newCarFixture(), *newCarFixture()}
	InitCarDB(t, db, cars)

	defer ClearCarDB(t, db)

	repo := NewCarRepositorySqlx(context.Background(), db, events.NewEventDispatcher())

	testCases := []struct {
		name      string
		searchArg application.SearchCarParams
		wantCars  []domain.Car
	}{
		{
			name:      "correct input",
			searchArg: application.SearchCarParams{Age: 2020},
			wantCars:  cars,
		},
		{
			name:      "correct input all",
			searchArg: application.SearchCarParams{},
			wantCars:  cars,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := repo.Find(tc.searchArg)

			if !SlicesEqual(c, tc.wantCars) {
				t.Error("unexpected result") //c, tc.wantCars
			}
		})
	}
}

func TestCarRepositorySqlx_FindOne(t *testing.T) {
	db := GetCarDBConn(t)
	defer db.Close()

	cars := []domain.Car{*newCarFixture()}
	InitCarDB(t, db, cars)

	defer ClearCarDB(t, db)

	repo := NewCarRepositorySqlx(context.Background(), db, events.NewEventDispatcher())

	testCases := []struct {
		name      string
		idArg     string
		wantIsCar bool
		wantError error
	}{
		{
			name:      "correct input",
			idArg:     cars[0].ID,
			wantIsCar: true,
			wantError: nil,
		},
		{
			name:      "incorrect id input",
			idArg:     "invalid-id",
			wantIsCar: false,
			wantError: application.ErrNotFoundCar,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := repo.FindOne(tc.idArg)

			if reflect.ValueOf(s).IsNil() == tc.wantIsCar {
				t.Error("unexpected result", s)
			}

			if !errors.Is(err, tc.wantError) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCarRepositorySqlx_Save(t *testing.T) {
	db := GetCarDBConn(t)
	defer db.Close()

	cars := []domain.Car{*newCarFixture()}
	InitCarDB(t, db, cars)

	defer ClearCarDB(t, db)

	dispatcher := &dispatcherMock{calls: make(map[string]uint)}
	repo := NewCarRepositorySqlx(context.Background(), db, dispatcher)

	testCases := []struct {
		name          string
		carArg        domain.Car
		wantIsCar     bool
		wantError     error
		wantDispErr   error
		wantDispCalls uint
	}{
		{
			name:          "correct input",
			carArg:        *newCarFixture(),
			wantIsCar:     true,
			wantError:     nil,
			wantDispErr:   nil,
			wantDispCalls: 1,
		},
		{
			name:          "correct update input",
			carArg:        cars[0],
			wantError:     nil,
			wantDispErr:   nil,
			wantDispCalls: 1,
		},
		{
			name:          "incorrect car input",
			carArg:        *newCarInvalidFixture(),
			wantError:     application.ErrInvalidCar,
			wantDispErr:   nil,
			wantDispCalls: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dispatcher.expectedDispatchErr = tc.wantDispErr
			dispatcher.calls["Dispatch"] = 0

			err := repo.Save(tc.carArg)

			if !errors.Is(err, tc.wantError) {
				t.Error(err)
			}

			if dispatcher.calls["Dispatch"] != tc.wantDispCalls {
				t.Error("invalid dispatcher call", dispatcher.calls["Dispatch"])
			}
		})
	}
}

func TestCarRepositorySqlx_Delete(t *testing.T) {
	db := GetCarDBConn(t)
	defer db.Close()

	cars := []domain.Car{*newCarFixture()}
	InitCarDB(t, db, cars)

	defer ClearCarDB(t, db)

	repo := NewCarRepositorySqlx(context.Background(), db, events.NewEventDispatcher())

	testCases := []struct {
		name      string
		idArg     string
		wantError error
	}{
		{
			name:      "correct input",
			idArg:     cars[0].ID,
			wantError: nil,
		},
		{
			name:      "incorrect id input",
			idArg:     "invalid-id",
			wantError: application.ErrInvalidCar,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Delete(tc.idArg)

			if !errors.Is(err, tc.wantError) {
				t.Error(err)
			}
		})
	}
}
