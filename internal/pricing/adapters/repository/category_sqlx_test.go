package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	// _ "github.com/lib/pq"

	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/domain"
)

func newCategoryFixture() *domain.Category {
	c, _ := domain.NewCategory(
		"Basic",
		"basic cars",
		[]string{"UNO", "MERIVA"},
		[]domain.Policy{{
			ID:      "83369771-f9a4-48b7-b87b-463f19f7b187",
			Name:    "Promo 1",
			Price:   0.2,
			Unit:    domain.PerKM,
			MinUnit: 50,
		}, {
			ID:      "4202b708-a387-4bae-85ce-11cb7a95759d",
			Name:    "Promo 2",
			Price:   30.5,
			Unit:    domain.PerDay,
			MinUnit: 5,
		}},
	)
	return c
}

func GetDBConn(t *testing.T) *sqlx.DB {
	t.Helper()
	// connStr := "postgres://postgres:admin1234@localhost/?sslmode=disable"
	// db, err := sqlx.Open("postgres", connStr)
	db, err := sqlx.Open("sqlite3", "data/pricing_test.db")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func InitDB(t *testing.T, db *sqlx.DB, categories []domain.Category) {
	t.Helper()
	const (
		saveCategory         = `INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)`
		insertModelCategory  = `INSERT INTO cmodels (name, "categoryId") VALUES ($1, $2)`
		insertPolicyCategory = `INSERT INTO cpolicies (id, name, price, unit, "minUnit", "categoryId") VALUES ($1, $2, $3, $4, $5, $6)`
	)

	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}

	for _, category := range categories {
		if _, err := db.Exec(saveCategory, category.ID, category.Name, category.Description); err != nil {
			t.Fatal("CATEGORY", err)
		}

		for _, model := range category.CarModels {
			if _, err := tx.Exec(insertModelCategory, model, category.ID); err != nil {
				t.Fatal("MODEL", err)
			}
		}

		for _, policy := range category.Policies {
			if _, err := tx.Exec(insertPolicyCategory, policy.ID, policy.Name, policy.Price, policy.Unit, policy.MinUnit, category.ID); err != nil {
				t.Fatal("POLICY", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func ClearDB(t *testing.T, db *sqlx.DB) {
	t.Helper()
	const (
		deleteAllCategories = "DELETE FROM categories"
		deleteAllModels     = "DELETE FROM cmodels"
		deleteAllPolicies   = "DELETE FROM cpolicies"
	)

	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.Exec(deleteAllModels); err != nil {
		t.Fatal(err)
	}

	if _, err := tx.Exec(deleteAllPolicies); err != nil {
		t.Fatal(err)
	}

	if _, err := tx.Exec(deleteAllCategories); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestCategoryRepositorySqlx_FindAll(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	categories := []domain.Category{*newCategoryFixture()}
	InitDB(t, db, categories)

	defer ClearDB(t, db)

	repo := NewCategoryRepositorySqlx(context.Background(), db)

	t.Run("findAll", func(t *testing.T) {
		s := repo.FindAll()

		if len(categories) != len(s) {
			t.Error("unequal categories", categories, s)
		}

		if !reflect.DeepEqual(categories, s) {
			t.Error("unequal categories", categories, s)
		}
	})
}

func TestCategoryRepositorySqlx_FindOne(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	categories := []domain.Category{*newCategoryFixture()}
	InitDB(t, db, categories)

	defer ClearDB(t, db)

	repo := NewCategoryRepositorySqlx(context.Background(), db)

	testCases := []struct {
		name           string
		idArg          string
		wantIsCategory bool
		wantError      error
	}{
		{
			name:           "correct input",
			idArg:          categories[0].ID,
			wantIsCategory: true,
			wantError:      nil,
		},
		{
			name:           "incorrect id input",
			idArg:          "invalid-id",
			wantIsCategory: false,
			wantError:      application.ErrNotFoundCategory,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := repo.FindOne(tc.idArg)

			// if !reflect.DeepEqual(categories[0], s) {
			// 	t.Error("unequal category", categories[0], s)
			// }

			if reflect.ValueOf(s).IsNil() == tc.wantIsCategory {
				t.Error("unexpected result", s)
			}

			if !errors.Is(err, tc.wantError) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCategoryRepositorySqlx_Delete(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	categories := []domain.Category{*newCategoryFixture()}
	InitDB(t, db, categories)

	defer ClearDB(t, db)

	repo := NewCategoryRepositorySqlx(context.Background(), db)

	testCases := []struct {
		name      string
		idArg     string
		wantError error
	}{
		{
			name:      "correct input",
			idArg:     categories[0].ID,
			wantError: nil,
		},
		{
			name:      "incorrect id input",
			idArg:     "invalid-id",
			wantError: application.ErrInvalidCategory,
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

func TestCategoryRepositorySqlx_Save(t *testing.T) {
	db := GetDBConn(t)
	defer db.Close()

	categories := []domain.Category{*newCategoryFixture()}
	InitDB(t, db, categories)

	defer ClearDB(t, db)

	repo := NewCategoryRepositorySqlx(context.Background(), db)

	testCases := []struct {
		name           string
		categoryArg    domain.Category
		wantIsCategory bool
		wantError      error
	}{
		{
			name:           "correct input",
			categoryArg:    *newCategoryFixture(),
			wantIsCategory: true,
			wantError:      nil,
		},
		{
			name:        "correct update input",
			categoryArg: categories[0],
			wantError:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Save(tc.categoryArg)

			if !errors.Is(err, tc.wantError) {
				t.Error(err)
			}
		})
	}
}
