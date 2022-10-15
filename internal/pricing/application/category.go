package application

import (
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/domain"
)

type CategoryUseCase interface {
	GetCategories() []domain.Category
	GetCategoryById(id string) (*domain.Category, error)
	AddCategory(name, description string) error
	DeleteCategory(id string) error
	AddModelInCategory(categoryId, modelId string) error
	DeleteModelInCategory(categoryId, modelId string) error
	AddPolicyInCategory(categoryId, name string, price float32, unit, minUnit uint) error
	DeletePolicyInCategory(categoryId, policyId string) error
}

type categoryUseCase struct {
	categoryRepo CategoryRepository
}

func NewCategoryUseCase(categoryRepo CategoryRepository) *categoryUseCase {
	return &categoryUseCase{categoryRepo}
}

func (uc categoryUseCase) GetCategories() []domain.Category {
	return uc.categoryRepo.FindAll()
}

func (uc categoryUseCase) GetCategoryById(id string) (*domain.Category, error) {
	if err := validation.ValidId(id); err != nil {
		return nil, ErrInvalidId
	}

	category, err := uc.categoryRepo.FindOne(id)

	if err != nil {
		return nil, ErrNotFoundCategory
	}

	return category, nil
}

func (uc categoryUseCase) AddCategory(name, description string) error {
	newCategory, err := domain.NewCategory(name, description, []string{}, []domain.Policy{})

	if err != nil {
		return ErrInvalidEntity
	}

	if err := uc.categoryRepo.Save(*newCategory); err != nil {
		return ErrInvalidCategory
	}

	return nil
}

func (uc categoryUseCase) DeleteCategory(id string) error {
	if err := validation.ValidId(id); err != nil {
		return ErrInvalidId
	}

	if err := uc.categoryRepo.Delete(id); err != nil {
		return ErrInvalidCategory
	}

	return nil
}

func (uc categoryUseCase) AddModelInCategory(categoryId, modelName string) error {
	category, err := uc.categoryRepo.FindOne(categoryId)

	if err != nil {
		return ErrInvalidCategory
	}

	if err := category.AddModel(modelName); err != nil {
		return ErrInvalidModel
	}

	if err := uc.categoryRepo.Save(*category); err != nil {
		return ErrInvalidCategory
	}

	return nil
}

func (uc categoryUseCase) DeleteModelInCategory(categoryId, modelName string) error {
	category, err := uc.categoryRepo.FindOne(categoryId)

	if err != nil {
		return ErrInvalidCategory
	}

	if err := category.DelModel(modelName); err != nil {
		return ErrInvalidModel
	}

	if err := uc.categoryRepo.Save(*category); err != nil {
		return ErrInvalidCategory
	}

	return nil
}

func (uc categoryUseCase) AddPolicyInCategory(categoryId, name string, price float32, unit, minUnit uint) error {
	category, err := uc.categoryRepo.FindOne(categoryId)

	if err != nil {
		return ErrInvalidCategory
	}

	policy, err := domain.NewPolicy(name, price, domain.Unit(unit), minUnit)

	if err != nil {
		return ErrInvalidEntity
	}

	if err := category.AddPolicy(*policy); err != nil {
		return ErrInvalidPolicy
	}

	if err := uc.categoryRepo.Save(*category); err != nil {
		return ErrInvalidCategory
	}

	return nil
}

func (uc categoryUseCase) DeletePolicyInCategory(categoryId, policyId string) error {
	category, err := uc.categoryRepo.FindOne(categoryId)

	if err != nil {
		return ErrInvalidCategory
	}

	if err := category.DelPolicy(policyId); err != nil {
		return ErrInvalidPolicy
	}

	if err := uc.categoryRepo.Save(*category); err != nil {
		return ErrInvalidCategory
	}

	return nil
}
