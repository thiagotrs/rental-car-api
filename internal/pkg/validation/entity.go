package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewId() string {
	return fmt.Sprintf("%v", uuid.New())
}

func ValidateEntity(e interface{}) error {
	validate := validator.New()

	if err := validate.Struct(e); err != nil {
		return err
	}

	return nil
}

func ValidId(id string) error {
	validate := validator.New()

	if err := validate.Var(id, "uuid4"); err != nil {
		return err
	}

	return nil
}
