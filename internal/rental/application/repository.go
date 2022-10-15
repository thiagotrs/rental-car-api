package application

import "github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"

type OrderReaderRepository interface {
	FindOne(id string) (*domain.Order, error)
}

type OrderWriterRepository interface {
	Save(order domain.Order) error
}

type OrderRepository interface {
	OrderReaderRepository
	OrderWriterRepository
}
