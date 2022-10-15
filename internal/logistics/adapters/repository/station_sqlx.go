package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

const (
	findStations  = "SELECT * FROM stations"
	findStation   = "SELECT * FROM stations WHERE id = $1 LIMIT 1"
	upsertStation = "INSERT INTO stations VALUES (:id, :name, :address, :complement, :state, :city, :cep, :capacity, :idle) ON CONFLICT(id) DO UPDATE SET name = :name, address = :address, complement = :complement, state = :state, city = :city, cep = :cep, capacity = :capacity, idle = :idle WHERE stations.id = :id"
	deleteStation = "DELETE FROM stations WHERE id = $1"
)

type stationRepositorySqlx struct {
	ctx context.Context
	DB  *sqlx.DB
}

func NewStationRepositorySqlx(ctx context.Context, DB *sqlx.DB) *stationRepositorySqlx {
	return &stationRepositorySqlx{ctx, DB}
}

func (repo *stationRepositorySqlx) FindAll() []domain.Station {
	stations := []domain.Station{}

	if err := repo.DB.SelectContext(repo.ctx, &stations, findStations); err != nil {
		return stations
	}

	return stations
}

func (repo *stationRepositorySqlx) FindOne(id string) (*domain.Station, error) {
	var station domain.Station

	if err := repo.DB.GetContext(repo.ctx, &station, findStation, id); err != nil {
		return nil, application.ErrNotFoundStation
	}

	return &station, nil
}

func (repo *stationRepositorySqlx) Save(station domain.Station) error {
	if err := validation.ValidateEntity(station); err != nil {
		return application.ErrInvalidStation
	}
	if _, err := repo.DB.NamedExecContext(repo.ctx, upsertStation, station); err != nil {
		return application.ErrInvalidStation
	}

	return nil
}

func (repo *stationRepositorySqlx) Delete(id string) error {
	r, err := repo.DB.ExecContext(repo.ctx, deleteStation, id)
	if err != nil {
		return application.ErrInvalidStation
	}
	n, err := r.RowsAffected()
	if err != nil || n == 0 {
		return application.ErrInvalidStation
	}
	return nil
}
