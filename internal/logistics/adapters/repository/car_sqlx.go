package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

const (
	findCars  = `SELECT * FROM cars`
	findCar   = `SELECT * FROM cars WHERE id = $1 LIMIT 1`
	upsertCar = `
	INSERT INTO cars VALUES (:id, :age, :plate, :document, :model, :make, :stationId, :km, :status) 
	ON CONFLICT(id) DO UPDATE SET age = :age, plate = :plate, document = :document, model = :model, make = :make, "stationId" = :stationId, km = :km, status = :status 
	WHERE cars.id = :id`
	deleteCar = `DELETE FROM cars WHERE id = $1`
)

type carRepositorySqlx struct {
	ctx  context.Context
	DB   *sqlx.DB
	disp events.Dispatcher
}

func NewCarRepositorySqlx(ctx context.Context, DB *sqlx.DB, disp events.Dispatcher) *carRepositorySqlx {
	return &carRepositorySqlx{ctx, DB, disp}
}

func (repo *carRepositorySqlx) Find(search application.SearchCarParams) []domain.Car {
	cars := []domain.Car{}

	var args []string
	var values []interface{}
	count := 0

	if search.Age > 0 {
		count++
		args = append(args, fmt.Sprintf(`age = $%v`, count))
		values = append(values, search.Age)
	}
	if len(search.Plate) > 0 {
		count++
		args = append(args, fmt.Sprintf(`plate = $%v`, count))
		values = append(values, search.Plate)
	}
	if len(search.Document) > 0 {
		count++
		args = append(args, fmt.Sprintf(`document = $%v`, count))
		values = append(values, search.Document)
	}
	if len(search.Model) > 0 {
		count++
		args = append(args, fmt.Sprintf(`model = $%v`, count))
		values = append(values, search.Model)
	}
	if len(search.Make) > 0 {
		count++
		args = append(args, fmt.Sprintf(`make = $%v`, count))
		values = append(values, search.Make)
	}
	if len(search.StationId) > 0 {
		count++
		args = append(args, fmt.Sprintf(`"stationId" = $%v`, count))
		values = append(values, search.StationId)
	}
	if search.KM > 0 {
		count++
		args = append(args, fmt.Sprintf(`km = $%v`, count))
		values = append(values, search.KM)
	}
	if search.Status > 0 {
		count++
		args = append(args, fmt.Sprintf(`status = $%v`, count))
		values = append(values, search.Status)
	}

	findCarsWithFilter := findCars
	if len(args) > 0 {
		findCarsWithFilter = findCarsWithFilter + ` WHERE ` + strings.Join(args, ` AND `)
	}

	if err := repo.DB.SelectContext(repo.ctx, &cars, findCarsWithFilter, values...); err != nil {
		return cars
	}

	return cars
}

func (repo *carRepositorySqlx) FindOne(id string) (*domain.Car, error) {
	var car domain.Car

	if err := repo.DB.GetContext(repo.ctx, &car, findCar, id); err != nil {
		return nil, application.ErrNotFoundCar
	}

	return &car, nil
}

func (repo *carRepositorySqlx) Save(car domain.Car) error {
	if err := validation.ValidateEntity(car); err != nil {
		return application.ErrInvalidCar
	}

	tx, err := repo.DB.BeginTxx(repo.ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.NamedExecContext(repo.ctx, upsertCar, car); err != nil {
		return application.ErrInvalidCar
	}

	if len(car.Events) > 0 {
		if err := repo.disp.Dispatch(car.Events); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (repo *carRepositorySqlx) Delete(id string) error {
	r, err := repo.DB.ExecContext(repo.ctx, deleteCar, id)
	if err != nil {
		return application.ErrInvalidCar
	}
	n, err := r.RowsAffected()
	if err != nil || n == 0 {
		return application.ErrInvalidCar
	}
	return nil
}
