package repository

import (
	"sync"

	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/domain"
)

type categoryRepositoryInMemory struct {
	categories map[string]domain.Category
	*sync.RWMutex
}

func NewCategoryRepositoryInMemory(categories []domain.Category) *categoryRepositoryInMemory {
	categoriesMap := make(map[string]domain.Category)
	for _, v := range categories {
		categoriesMap[v.ID] = v
	}
	return &categoryRepositoryInMemory{categoriesMap, &sync.RWMutex{}}
}

func (repo categoryRepositoryInMemory) FindAll() []domain.Category {
	repo.Lock()
	defer repo.Unlock()

	categories := []domain.Category{}
	for k := range repo.categories {
		categories = append(categories, repo.categories[k])
	}

	return categories
}

func (repo categoryRepositoryInMemory) FindOne(id string) (*domain.Category, error) {
	repo.Lock()
	defer repo.Unlock()

	s, exists := repo.categories[id]
	if !exists {
		return &s, application.ErrNotFoundCategory
	}

	return &s, nil
}

func (repo *categoryRepositoryInMemory) Save(category domain.Category) error {
	repo.Lock()
	defer repo.Unlock()

	repo.categories[category.ID] = category

	return nil
}

func (repo *categoryRepositoryInMemory) Delete(id string) error {
	repo.Lock()
	defer repo.Unlock()

	if _, exists := repo.categories[id]; !exists {
		return application.ErrInvalidCategory
	}

	delete(repo.categories, id)

	return nil
}
