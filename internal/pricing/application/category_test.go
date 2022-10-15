package application

import (
	"errors"
	"reflect"
	"testing"

	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/domain"
)

func newCategoryFixture() *domain.Category {
	return &domain.Category{
		ID:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		Name:        "Basic",
		Description: "basic cars",
		CarModels:   []string{"UNO", "MERIVA"},
		Policies: []domain.Policy{{
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
	}
}

type categoryRepositoryMock struct {
	expectedFindAllCategories []domain.Category
	expectedFindOneCategory   *domain.Category
	expectedFindOneErr        error
	expectedSaveErr           error
	expectedDeleteErr         error
	calls                     map[string]uint
}

func (m *categoryRepositoryMock) FindAll() []domain.Category {
	m.calls["FindAll"] = m.calls["FindAll"] + 1
	return m.expectedFindAllCategories
}

func (m *categoryRepositoryMock) FindOne(id string) (*domain.Category, error) {
	m.calls["FindOne"] = m.calls["FindOne"] + 1
	return m.expectedFindOneCategory, m.expectedFindOneErr
}

func (m *categoryRepositoryMock) Save(category domain.Category) error {
	m.calls["Save"] = m.calls["Save"] + 1
	return m.expectedSaveErr
}

func (m *categoryRepositoryMock) Delete(id string) error {
	m.calls["Delete"] = m.calls["Delete"] + 1
	return m.expectedDeleteErr
}

func TestCategoryUseCase_GetCategories(t *testing.T) {
	newCategory := newCategoryFixture()

	type setup struct {
		repoCategories []domain.Category
	}

	type want struct {
		categories []domain.Category
	}

	testCases := []struct {
		name  string
		setup setup
		want  want
	}{
		{
			name: "correct 1",
			setup: setup{
				repoCategories: []domain.Category{},
			},
			want: want{
				categories: []domain.Category{},
			},
		},
		{
			name: "correct 2",
			setup: setup{
				repoCategories: []domain.Category{*newCategory},
			},
			want: want{
				categories: []domain.Category{*newCategory},
			},
		},
		{
			name: "unexpected error",
			setup: setup{
				repoCategories: nil,
			},
			want: want{
				categories: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := &categoryRepositoryMock{
				expectedFindAllCategories: tc.setup.repoCategories,
				calls:                     make(map[string]uint),
			}
			categoryUC := NewCategoryUseCase(categoryRepo)
			categories := categoryUC.GetCategories()

			if categoryRepo.calls["FindAll"] != 1 {
				t.Error("invalid repo call", categoryRepo.calls["FindAll"])
			}

			if !reflect.DeepEqual(categories, tc.want.categories) {
				t.Error("unequal category", categories)
			}
		})
	}
}

func TestCategoryUseCase_GetCategoryById(t *testing.T) {
	newCategory := newCategoryFixture()

	type setup struct {
		repoCategory *domain.Category
		repoErr      error
	}

	type args struct {
		id string
	}

	type want struct {
		category *domain.Category
		err      error
		calls    uint
	}

	testCases := []struct {
		name  string
		setup setup
		args  args
		want  want
	}{
		{
			name: "correct input",
			setup: setup{
				repoCategory: newCategory,
				repoErr:      nil,
			},
			args: args{
				id: newCategory.ID,
			},
			want: want{
				category: newCategory,
				err:      nil,
				calls:    1,
			},
		},
		{
			name:  "incorrect id input",
			setup: setup{},
			args: args{
				id: "invalid-id",
			},
			want: want{
				category: nil,
				err:      ErrInvalidId,
				calls:    0,
			},
		},
		{
			name: "not found category",
			setup: setup{
				repoCategory: nil,
				repoErr:      ErrNotFoundCategory,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				category: nil,
				err:      ErrNotFoundCategory,
				calls:    1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := &categoryRepositoryMock{
				expectedFindOneCategory: tc.setup.repoCategory,
				expectedFindOneErr:      tc.setup.repoErr,
				calls:                   make(map[string]uint),
			}
			categoryUC := NewCategoryUseCase(categoryRepo)
			category, err := categoryUC.GetCategoryById(tc.args.id)

			if categoryRepo.calls["FindOne"] != tc.want.calls {
				t.Error("invalid repo call", categoryRepo.calls["FindOne"])
			}

			if !reflect.DeepEqual(category, tc.want.category) {
				t.Error("unequal category", category)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCategoryUseCase_DeleteCategory(t *testing.T) {
	newCategory := newCategoryFixture()

	type setup struct {
		repoDelErr error
	}

	type args struct {
		id string
	}

	type want struct {
		err      error
		delCalls uint
	}

	testCases := []struct {
		name  string
		setup setup
		args  args
		want  want
	}{
		{
			name: "correct input",
			setup: setup{
				repoDelErr: nil,
			},
			args: args{
				id: newCategory.ID,
			},
			want: want{
				err:      nil,
				delCalls: 1,
			},
		},
		{
			name:  "incorrect id input",
			setup: setup{},
			args: args{
				id: "invalid-id",
			},
			want: want{
				err:      ErrInvalidId,
				delCalls: 0,
			},
		},
		{
			name: "not found category",
			setup: setup{
				repoDelErr: ErrInvalidCategory,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				err:      ErrInvalidCategory,
				delCalls: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := &categoryRepositoryMock{
				expectedDeleteErr: tc.setup.repoDelErr,
				calls:             make(map[string]uint),
			}
			categoryUC := NewCategoryUseCase(categoryRepo)
			err := categoryUC.DeleteCategory(tc.args.id)

			if categoryRepo.calls["Delete"] != tc.want.delCalls {
				t.Error("invalid repo call", categoryRepo.calls["Delete"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCategoryUseCase_AddCategory(t *testing.T) {
	type setup struct {
		repoSaveErr error
	}

	type args struct {
		name        string
		description string
	}

	type want struct {
		err       error
		saveCalls uint
	}

	testCases := []struct {
		name  string
		setup setup
		args  args
		want  want
	}{
		{
			name: "correct input",
			setup: setup{
				repoSaveErr: nil,
			},
			args: args{
				name:        "Premium",
				description: "premium cars",
			},
			want: want{
				err:       nil,
				saveCalls: 1,
			},
		},
		{
			name: "incorrect input",
			setup: setup{
				repoSaveErr: nil,
			},
			args: args{
				name:        "",
				description: "",
			},
			want: want{
				err:       ErrInvalidEntity,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := &categoryRepositoryMock{
				expectedSaveErr: tc.setup.repoSaveErr,
				calls:           make(map[string]uint),
			}
			categoryUC := NewCategoryUseCase(categoryRepo)
			err := categoryUC.AddCategory(
				tc.args.name,
				tc.args.description,
			)

			if categoryRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", categoryRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", tc.want.err, err)
			}
		})
	}
}

func TestCategoryUseCase_AddModelInCategory(t *testing.T) {
	newCategory := newCategoryFixture()

	type setup struct {
		repoFindOne *domain.Category
		repoFindErr error
		repoSaveErr error
	}

	type args struct {
		categoryId string
		modelName  string
	}

	type want struct {
		err       error
		findCalls uint
		saveCalls uint
	}

	testCases := []struct {
		name  string
		setup setup
		args  args
		want  want
	}{
		{
			name: "correct input",
			setup: setup{
				repoFindOne: newCategory,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: newCategory.ID,
				modelName:  "COROLA",
			},
			want: want{
				err:       nil,
				findCalls: 1,
				saveCalls: 1,
			},
		},
		{
			name: "incorrect category input",
			setup: setup{
				repoFindOne: nil,
				repoFindErr: ErrInvalidCategory,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: "invalid-id",
				modelName:  "COROLA",
			},
			want: want{
				err:       ErrInvalidCategory,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect model input",
			setup: setup{
				repoFindOne: newCategory,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: newCategory.ID,
				modelName:  "UNO",
			},
			want: want{
				err:       ErrInvalidModel,
				findCalls: 1,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := &categoryRepositoryMock{
				expectedFindOneCategory: tc.setup.repoFindOne,
				expectedFindOneErr:      tc.setup.repoFindErr,
				expectedSaveErr:         tc.setup.repoSaveErr,
				calls:                   make(map[string]uint),
			}
			categoryUC := NewCategoryUseCase(categoryRepo)
			err := categoryUC.AddModelInCategory(
				tc.args.categoryId,
				tc.args.modelName,
			)

			if categoryRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", categoryRepo.calls["FindOne"])
			}

			if categoryRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", categoryRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", tc.want.err, err)
			}
		})
	}
}

func TestCategoryUseCase_DeleteModelInCategory(t *testing.T) {
	newCategory := newCategoryFixture()

	type setup struct {
		repoFindOne *domain.Category
		repoFindErr error
		repoSaveErr error
	}

	type args struct {
		categoryId string
		modelName  string
	}

	type want struct {
		err       error
		findCalls uint
		saveCalls uint
	}

	testCases := []struct {
		name  string
		setup setup
		args  args
		want  want
	}{
		{
			name: "correct input",
			setup: setup{
				repoFindOne: newCategory,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: newCategory.ID,
				modelName:  "UNO",
			},
			want: want{
				err:       nil,
				findCalls: 1,
				saveCalls: 1,
			},
		},
		{
			name: "incorrect category input",
			setup: setup{
				repoFindOne: nil,
				repoFindErr: ErrInvalidCategory,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: "invalid-id",
				modelName:  "UNO",
			},
			want: want{
				err:       ErrInvalidCategory,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect model input",
			setup: setup{
				repoFindOne: newCategory,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: newCategory.ID,
				modelName:  "COROLA",
			},
			want: want{
				err:       ErrInvalidModel,
				findCalls: 1,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := &categoryRepositoryMock{
				expectedFindOneCategory: tc.setup.repoFindOne,
				expectedFindOneErr:      tc.setup.repoFindErr,
				expectedSaveErr:         tc.setup.repoSaveErr,
				calls:                   make(map[string]uint),
			}
			categoryUC := NewCategoryUseCase(categoryRepo)
			err := categoryUC.DeleteModelInCategory(
				tc.args.categoryId,
				tc.args.modelName,
			)

			if categoryRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", categoryRepo.calls["FindOne"])
			}

			if categoryRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", categoryRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", tc.want.err, err)
			}
		})
	}
}

func TestCategoryUseCase_AddPolicyInCategory(t *testing.T) {
	newCategory := newCategoryFixture()

	type setup struct {
		repoFindOne *domain.Category
		repoFindErr error
		repoSaveErr error
	}

	type args struct {
		categoryId, name string
		price            float32
		unit, minUnit    uint
	}

	type want struct {
		err       error
		findCalls uint
		saveCalls uint
	}

	testCases := []struct {
		name  string
		setup setup
		args  args
		want  want
	}{
		{
			name: "correct input",
			setup: setup{
				repoFindOne: newCategory,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: newCategory.ID,
				name:       "policy 1",
				price:      1.5,
				unit:       uint(domain.PerKM),
				minUnit:    20,
			},
			want: want{
				err:       nil,
				findCalls: 1,
				saveCalls: 1,
			},
		},
		{
			name: "incorrect category input",
			setup: setup{
				repoFindOne: nil,
				repoFindErr: ErrInvalidCategory,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: "invalid-id",
				name:       "policy 1",
				price:      1.5,
				unit:       uint(domain.PerKM),
				minUnit:    20,
			},
			want: want{
				err:       ErrInvalidCategory,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect invalid policy input",
			setup: setup{
				repoFindOne: newCategory,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: newCategory.ID,
				name:       "policy 1",
				price:      0,
				unit:       uint(domain.PerKM),
				minUnit:    20,
			},
			want: want{
				err:       ErrInvalidEntity,
				findCalls: 1,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := &categoryRepositoryMock{
				expectedFindOneCategory: tc.setup.repoFindOne,
				expectedFindOneErr:      tc.setup.repoFindErr,
				expectedSaveErr:         tc.setup.repoSaveErr,
				calls:                   make(map[string]uint),
			}
			categoryUC := NewCategoryUseCase(categoryRepo)
			err := categoryUC.AddPolicyInCategory(
				tc.args.categoryId,
				tc.args.name,
				tc.args.price,
				tc.args.unit,
				tc.args.minUnit,
			)

			if categoryRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", categoryRepo.calls["FindOne"])
			}

			if categoryRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", categoryRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", tc.want.err, err)
			}
		})
	}
}

func TestCategoryUseCase_DeletePolicyInCategory(t *testing.T) {
	newCategory := newCategoryFixture()

	type setup struct {
		repoFindOne *domain.Category
		repoFindErr error
		repoSaveErr error
	}

	type args struct {
		categoryId string
		policyId   string
	}

	type want struct {
		err       error
		findCalls uint
		saveCalls uint
	}

	testCases := []struct {
		name  string
		setup setup
		args  args
		want  want
	}{
		{
			name: "correct input",
			setup: setup{
				repoFindOne: newCategory,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: newCategory.ID,
				policyId:   newCategory.Policies[0].ID,
			},
			want: want{
				err:       nil,
				findCalls: 1,
				saveCalls: 1,
			},
		},
		{
			name: "incorrect category input",
			setup: setup{
				repoFindOne: nil,
				repoFindErr: ErrInvalidCategory,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: "invalid-id",
				policyId:   newCategory.Policies[0].ID,
			},
			want: want{
				err:       ErrInvalidCategory,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect invalid policy input",
			setup: setup{
				repoFindOne: newCategory,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				categoryId: newCategory.ID,
				policyId:   "invalid-id",
			},
			want: want{
				err:       ErrInvalidPolicy,
				findCalls: 1,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categoryRepo := &categoryRepositoryMock{
				expectedFindOneCategory: tc.setup.repoFindOne,
				expectedFindOneErr:      tc.setup.repoFindErr,
				expectedSaveErr:         tc.setup.repoSaveErr,
				calls:                   make(map[string]uint),
			}
			categoryUC := NewCategoryUseCase(categoryRepo)
			err := categoryUC.DeletePolicyInCategory(
				tc.args.categoryId,
				tc.args.policyId,
			)

			if categoryRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", categoryRepo.calls["FindOne"])
			}

			if categoryRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", categoryRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", tc.want.err, err)
			}
		})
	}
}
