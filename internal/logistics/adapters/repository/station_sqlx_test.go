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
)

func newStationFixture() *domain.Station {
	s, _ := domain.NewStation("Station 1", "Farway Av.", "45, Ap. 50", "Polar", "Nort City", "20778990", 100, 0)
	return s
}

func newInvalidStationFixture() *domain.Station {
	return &domain.Station{
		ID:         "",
		Name:       "Station 1000",
		Address:    "Farway Av.",
		Complement: "45, Ap. 50",
		State:      "Polar",
		City:       "Nort City",
		Cep:        "20778990",
		Capacity:   100,
		Idle:       0,
	}
}

func GetDBConn(t *testing.T) *sqlx.DB {
	t.Helper()
	// connStr := "postgres://postgres:admin1234@localhost/?sslmode=disable"
	// db, err := sqlx.Open("postgres", connStr)
	db, err := sqlx.Open("sqlite3", "data/logistics_test.db")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func InitDB(t *testing.T, db *sqlx.DB, stations []domain.Station) {
	t.Helper()
	const saveStation = "INSERT INTO stations VALUES (:id, :name, :address, :complement, :state, :city, :cep, :capacity, :idle)"

	for _, s := range stations {
		if _, err := db.NamedExec(saveStation, s); err != nil {
			t.Fatal(err)
		}
	}
}

func ClearDB(t *testing.T, db *sqlx.DB) {
	t.Helper()
	const deleteAllStations = "DELETE FROM stations"

	if _, err := db.Exec(deleteAllStations); err != nil {
		t.Fatal(err)
	}
}

func TestStationRepositorySqlx_FindAll(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	stations := []domain.Station{*newStationFixture()}
	InitDB(t, db, stations)

	defer ClearDB(t, db)

	repo := NewStationRepositorySqlx(context.Background(), db)

	t.Run("newRepo", func(t *testing.T) {
		s := repo.FindAll()

		if len(stations) != len(s) {
			t.Error("unequal stations", stations, s)
		}

		if !reflect.DeepEqual(stations, s) {
			t.Error("unequal stations", stations, s)
		}
	})
}

func TestStationRepositorySqlx_FindOne(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	stations := []domain.Station{*newStationFixture()}
	InitDB(t, db, stations)

	defer ClearDB(t, db)

	repo := NewStationRepositorySqlx(context.Background(), db)

	testCases := []struct {
		name          string
		idArg         string
		wantIsStation bool
		wantError     error
	}{
		{
			name:          "correct input",
			idArg:         stations[0].ID,
			wantIsStation: true,
			wantError:     nil,
		},
		{
			name:          "incorrect id input",
			idArg:         "invalid-id",
			wantIsStation: false,
			wantError:     application.ErrNotFoundStation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := repo.FindOne(tc.idArg)

			// if !reflect.DeepEqual(stations[0], s) {
			// 	t.Error("unequal station", stations[0], s)
			// }

			if reflect.ValueOf(s).IsNil() == tc.wantIsStation {
				t.Error("unexpected result", s)
			}

			if !errors.Is(err, tc.wantError) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestStationRepositorySqlx_Save(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	stations := []domain.Station{*newStationFixture()}
	InitDB(t, db, stations)

	defer ClearDB(t, db)

	repo := NewStationRepositorySqlx(context.Background(), db)

	testCases := []struct {
		name          string
		stationArg    domain.Station
		wantIsStation bool
		wantError     error
	}{
		{
			name:          "correct input",
			stationArg:    *newStationFixture(),
			wantIsStation: true,
			wantError:     nil,
		},
		{
			name:       "correct update input",
			stationArg: stations[0],
			wantError:  nil,
		},
		{
			name:       "incorrect station input",
			stationArg: *newInvalidStationFixture(),
			wantError:  application.ErrInvalidStation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Save(tc.stationArg)

			if !errors.Is(err, tc.wantError) {
				t.Error(err)
			}
		})
	}
}

func TestStationRepositorySqlx_Delete(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	stations := []domain.Station{*newStationFixture()}
	InitDB(t, db, stations)

	defer ClearDB(t, db)

	repo := NewStationRepositorySqlx(context.Background(), db)

	testCases := []struct {
		name      string
		idArg     string
		wantError error
	}{
		{
			name:      "correct input",
			idArg:     stations[0].ID,
			wantError: nil,
		},
		{
			name:      "incorrect id input",
			idArg:     "invalid-id",
			wantError: application.ErrInvalidStation,
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
