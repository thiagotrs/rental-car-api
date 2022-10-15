package application

import "github.com/thiagotrs/rentalcar-ddd/internal/pricing/domain"

type CategoryReadRepository interface {
	FindAll() []domain.Category
	FindOne(id string) (*domain.Category, error)
}

type CategoryWriteRepository interface {
	Save(category domain.Category) error
	Delete(id string) error
}

type CategoryRepository interface {
	CategoryReadRepository
	CategoryWriteRepository
}
