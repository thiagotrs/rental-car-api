package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	// _ "github.com/lib/pq"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

func newCarFixture() *domain.Car {
	return &domain.Car{
		ID:        "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		Age:       2020,
		Plate:     "KST-9016",
		Document:  "abc.123.op-x",
		CarModel:  "UNO",
		InitialKM: 12000,
		Status:    domain.Parked,
		StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
	}
}

func newPolicyFixture() *domain.Policy {
	return &domain.Policy{
		ID:         "5ecf09ce-8c41-4faa-a4e5-824af9c80892",
		Name:       "Promo default",
		Price:      30.5,
		Unit:       domain.PerDay,
		MinUnit:    5,
		CarModel:   "UNO",
		CategoryId: "479ab9e7-ad16-4864-8e49-29b15e4b390e",
	}
}

func newOrderFixture() *domain.Order {
	o, _ := domain.NewOrder(
		time.Now(),
		time.Now().Add(time.Hour*24*5),
		*newCarFixture(),
		"83369771-f9a4-48b7-b87b-463f19f7b187",
		"2520aade-a397-4e3c-a589-39c6ae5c2eff",
		*newPolicyFixture())
	return o
}

func GetDBConn(t *testing.T) *sqlx.DB {
	t.Helper()
	// connStr := "postgres://postgres:admin1234@localhost/?sslmode=disable"
	// db, err := sqlx.Open("postgres", connStr)
	db, err := sqlx.Open("sqlite3", "data/rental_test.db")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func InitDB(t *testing.T, db *sqlx.DB, orders []domain.Order) {
	t.Helper()
	const (
		saveOrder = `
		INSERT INTO orders (id, "dateFrom", "dateTo", "dateReservFrom", "dateReservTo", status, "stationFromId", "stationToId", discount, tax)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		saveCarOrder = `
		INSERT INTO ocars (id, "orderId", age, plate, document, "carModel", "initialKM", "finalKM", status, "stationId") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		savePolicyOrder = `
		INSERT INTO opolicies (id, "orderId", name, price, unit, "minUnit", "carModel", "categoryId") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	)

	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}

	for _, order := range orders {
		if _, err := db.Exec(saveOrder, order.ID, order.DateFrom, order.DateTo, order.DateReservFrom, order.DateReservTo, order.Status, order.StationFromId, order.StationToId, order.Discount, order.Tax); err != nil {
			t.Fatal("ORDER", err)
		}

		if _, err := tx.Exec(saveCarOrder, order.Car.ID, order.ID, order.Car.Age, order.Car.Plate, order.Car.Document, order.Car.CarModel, order.Car.InitialKM, order.Car.FinalKM, order.Car.Status, order.Car.StationId); err != nil {
			t.Fatal("CAR", err)
		}

		if _, err := tx.Exec(savePolicyOrder, order.Policy.ID, order.ID, order.Policy.Name, order.Policy.Price, order.Policy.Unit, order.Policy.MinUnit, order.Policy.CarModel, order.Policy.CategoryId); err != nil {
			t.Fatal("POLICY", err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func ClearDB(t *testing.T, db *sqlx.DB) {
	t.Helper()
	const (
		deleteAllCars     = "DELETE FROM ocars"
		deleteAllPolicies = "DELETE FROM opolicies"
		deleteAllOrders   = "DELETE FROM orders"
	)

	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.Exec(deleteAllCars); err != nil {
		t.Fatal(err)
	}

	if _, err := tx.Exec(deleteAllPolicies); err != nil {
		t.Fatal(err)
	}

	if _, err := tx.Exec(deleteAllOrders); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
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

func TestOrderRepositorySqlx_FindOne(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	orders := []domain.Order{*newOrderFixture()}
	ClearDB(t, db)
	InitDB(t, db, orders)

	defer ClearDB(t, db)

	repo := NewOrderRepositorySqlx(context.Background(), db, events.NewEventDispatcher())

	testCases := []struct {
		name        string
		idArg       string
		wantIsOrder bool
		wantError   error
	}{
		{
			name:        "correct input",
			idArg:       orders[0].ID,
			wantIsOrder: true,
			wantError:   nil,
		},
		{
			name:        "incorrect id input",
			idArg:       "invalid-id",
			wantIsOrder: false,
			wantError:   application.ErrNotFoundOrder,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := repo.FindOne(tc.idArg)

			// if !reflect.DeepEqual(orders[0], s) {
			// 	t.Error("unequal category", orders[0], s)
			// }

			if reflect.ValueOf(s).IsNil() == tc.wantIsOrder {
				t.Error("unexpected result", s)
			}

			if !errors.Is(err, tc.wantError) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestOrderRepositorySqlx_Save(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	categories := []domain.Order{*newOrderFixture()}
	ClearDB(t, db)
	InitDB(t, db, categories)

	defer ClearDB(t, db)

	updatedOrder := *newOrderFixture()
	updatedOrder.Status = domain.Canceled

	dispatcher := &dispatcherMock{calls: make(map[string]uint)}
	repo := NewOrderRepositorySqlx(context.Background(), db, dispatcher)

	testCases := []struct {
		name          string
		categoryArg   domain.Order
		wantIsOrder   bool
		wantError     error
		wantDispErr   error
		wantDispCalls uint
	}{
		{
			name:          "correct input",
			categoryArg:   *newOrderFixture(),
			wantIsOrder:   true,
			wantError:     nil,
			wantDispErr:   nil,
			wantDispCalls: 1,
		},
		{
			name:          "correct update order input",
			categoryArg:   updatedOrder,
			wantError:     nil,
			wantDispErr:   nil,
			wantDispCalls: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dispatcher.expectedDispatchErr = tc.wantDispErr
			dispatcher.calls["Dispatch"] = 0

			err := repo.Save(tc.categoryArg)

			if !errors.Is(err, tc.wantError) {
				t.Error(err)
			}

			if dispatcher.calls["Dispatch"] != tc.wantDispCalls {
				t.Error("invalid dispatcher call", dispatcher.calls["Dispatch"])
			}
		})
	}
}
