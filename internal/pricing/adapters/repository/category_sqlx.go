package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/domain"
)

const (
	findCategories = `SELECT * FROM categories`
	findCategory   = `SELECT * FROM categories WHERE id = $1 LIMIT 1`

	upsertCategory = `
	INSERT INTO categories VALUES (:id, :name, :description) 
	ON CONFLICT(id) DO UPDATE SET name = :name, description = :description WHERE categories.id = :id`
	deleteCategory = `DELETE FROM categories WHERE id = $1`

	insertModelCategory = `INSERT INTO cmodels (name, "categoryId") VALUES ($1, $2)`
	deleteModels        = `DELETE FROM cmodels WHERE "categoryId" = $1`

	insertPolicyCategory = `
	INSERT INTO cpolicies (id, name, price, unit, "minUnit", "categoryId") VALUES ($1, $2, $3, $4, $5, $6) 
	ON CONFLICT(id) DO UPDATE SET name = $2, price = $3, unit = $4, "minUnit" = $5 WHERE cpolicies.id = $1`
	deletePolicies = `DELETE FROM cpolicies WHERE "categoryId" = $1`

	findModelsByCategory   = `SELECT name FROM cmodels WHERE "categoryId" = $1`
	findPoliciesByCategory = `SELECT id, name, price, unit, "minUnit" FROM cpolicies WHERE "categoryId" = $1`
)

type categoryRepositorySqlx struct {
	ctx context.Context
	DB  *sqlx.DB
}

func NewCategoryRepositorySqlx(ctx context.Context, DB *sqlx.DB) *categoryRepositorySqlx {
	return &categoryRepositorySqlx{ctx, DB}
}

func (repo *categoryRepositorySqlx) FindAll() []domain.Category {
	categories := []domain.Category{}

	if err := repo.DB.SelectContext(repo.ctx, &categories, findCategories); err != nil {
		return categories
	}

	for i, c := range categories {
		repo.DB.SelectContext(repo.ctx, &c.CarModels, findModelsByCategory, c.ID)
		categories[i].CarModels = c.CarModels
	}

	for i, c := range categories {
		repo.DB.SelectContext(repo.ctx, &c.Policies, findPoliciesByCategory, c.ID)
		categories[i].Policies = c.Policies
	}

	return categories
}

func (repo *categoryRepositorySqlx) FindOne(id string) (*domain.Category, error) {
	var category domain.Category

	if err := repo.DB.GetContext(repo.ctx, &category, findCategory, id); err != nil {
		return nil, application.ErrNotFoundCategory
	}

	repo.DB.SelectContext(repo.ctx, &category.CarModels, findModelsByCategory, category.ID)
	repo.DB.SelectContext(repo.ctx, &category.Policies, findPoliciesByCategory, category.ID)

	return &category, nil
}

func (repo *categoryRepositorySqlx) Save(category domain.Category) error {
	tx, err := repo.DB.BeginTxx(repo.ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.NamedExecContext(repo.ctx, upsertCategory, category)
	if err != nil {
		tx.Rollback()
		return err
	}
	n, err := result.RowsAffected()
	if err != nil || n == 0 {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(repo.ctx, deleteModels, category.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(repo.ctx, deletePolicies, category.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, m := range category.CarModels {
		result, err = tx.ExecContext(repo.ctx, insertModelCategory, m, category.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		n, err = result.RowsAffected()
		if err != nil || n == 0 {
			tx.Rollback()
			return err
		}
	}

	for _, p := range category.Policies {
		result, err = tx.ExecContext(
			repo.ctx,
			insertPolicyCategory,
			p.ID,
			p.Name,
			p.Price,
			p.Unit,
			p.MinUnit,
			category.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		n, err = result.RowsAffected()
		if err != nil || n == 0 {
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

func (repo *categoryRepositorySqlx) Delete(id string) error {
	tx, err := repo.DB.BeginTxx(repo.ctx, nil)
	if err != nil {
		return application.ErrInvalidCategory
	}

	result, err := tx.ExecContext(repo.ctx, deleteCategory, id)
	if err != nil {
		tx.Rollback()
		return application.ErrInvalidCategory
	}
	n, err := result.RowsAffected()
	if err != nil || n == 0 {
		tx.Rollback()
		return application.ErrInvalidCategory
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return application.ErrInvalidCategory
	}

	return nil
}
