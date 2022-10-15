package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

const (
	findOrder = `SELECT id, "dateFrom", "dateTo", "dateReservFrom", "dateReservTo", status, "stationFromId", "stationToId", discount,	tax FROM orders WHERE id = $1 LIMIT 1`

	findCarByOrder = `
	SELECT id, age, plate, document, "carModel", "initialKM", "finalKM", status, "stationId" FROM ocars 
	WHERE "orderId" = $1 LIMIT 1`

	findPolicyByOrder = `
	SELECT id, name, price, unit, "minUnit", "carModel", "categoryId" FROM opolicies 
	WHERE "orderId" = $1 LIMIT 1`

	upsertOrder = `
	INSERT INTO orders 
	VALUES (:id, :dateFrom, :dateTo, :dateReservFrom, :dateReservTo, :status, :stationFromId, :stationToId, :discount, :tax) 
	ON CONFLICT(id) DO 
	UPDATE SET "dateFrom" = :dateFrom, "dateTo" = :dateTo, "dateReservFrom" = :dateReservFrom, "dateReservTo" = :dateReservTo, status = :status, "stationFromId" = :stationFromId, "stationToId" = :stationToId, discount = :discount, tax = :tax 
	WHERE orders.id = :id`

	upsertCarOrder = `
	INSERT INTO ocars (id, "orderId", age, plate, document, "carModel", "initialKM", "finalKM", status, "stationId") 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT(id, "orderId") DO 
	UPDATE SET age = $3, plate = $4, document = $5, "carModel" = $6, "initialKM" = $7, "finalKM" = $8, status = $9, "stationId" = $10 
	WHERE ocars.id = $1 AND ocars."orderId" = $2`

	upsertPolicyOrder = `
	INSERT INTO opolicies (id, "orderId", name, price, unit, "minUnit", "carModel", "categoryId") 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	ON CONFLICT(id, "orderId") DO 
	UPDATE SET name = $3, price = $4, unit = $5, "minUnit" = $6, "carModel" = $7, "categoryId" = $8 
	WHERE opolicies.id = $1 AND opolicies."orderId" = $2`
)

type orderRepositorySqlx struct {
	ctx  context.Context
	DB   *sqlx.DB
	disp events.Dispatcher
}

func NewOrderRepositorySqlx(ctx context.Context, DB *sqlx.DB, disp events.Dispatcher) *orderRepositorySqlx {
	return &orderRepositorySqlx{ctx, DB, disp}
}

func (repo *orderRepositorySqlx) FindOne(id string) (*domain.Order, error) {
	var order domain.Order

	if err := repo.DB.GetContext(repo.ctx, &order, findOrder, id); err != nil {
		return nil, application.ErrNotFoundOrder
	}

	if err := repo.DB.GetContext(repo.ctx, &order.Car, findCarByOrder, order.ID); err != nil {
		return nil, application.ErrNotFoundOrder
	}
	if err := repo.DB.GetContext(repo.ctx, &order.Policy, findPolicyByOrder, order.ID); err != nil {
		return nil, application.ErrNotFoundOrder
	}

	return &order, nil
}

func (repo *orderRepositorySqlx) Save(order domain.Order) error {
	tx, err := repo.DB.BeginTxx(repo.ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.NamedExecContext(repo.ctx, upsertOrder, order)
	if err != nil {
		tx.Rollback()
		return err
	}
	n, err := result.RowsAffected()
	if err != nil || n == 0 {
		tx.Rollback()
		return err
	}

	result, err = tx.ExecContext(
		repo.ctx,
		upsertCarOrder,
		order.Car.ID,
		order.ID,
		order.Car.Age,
		order.Car.Plate,
		order.Car.Document,
		order.Car.CarModel,
		order.Car.InitialKM,
		order.Car.FinalKM,
		order.Car.Status,
		order.Car.StationId)
	if err != nil {
		tx.Rollback()
		return err
	}
	n, err = result.RowsAffected()
	if err != nil || n == 0 {
		tx.Rollback()
		return err
	}

	result, err = tx.ExecContext(
		repo.ctx,
		upsertPolicyOrder,
		order.Policy.ID,
		order.ID,
		order.Policy.Name,
		order.Policy.Price,
		order.Policy.Unit,
		order.Policy.MinUnit,
		order.Policy.CarModel,
		order.Policy.CategoryId)
	if err != nil {
		tx.Rollback()
		return err
	}
	n, err = result.RowsAffected()
	if err != nil || n == 0 {
		tx.Rollback()
		return err
	}

	if len(order.Events) > 0 {
		if err := repo.disp.Dispatch(order.Events); err != nil {
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
