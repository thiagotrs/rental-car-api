package domain

import (
	"fmt"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/validation"
)

type Category struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" validate:"required" db:"name"`
	Description string   `json:"description" db:"description"`
	CarModels   []string `json:"carModels" validate:"required,dive,required"`
	Policies    []Policy `json:"policies" validate:"required,dive,required"`
}

func NewCategory(name, description string, carModels []string, policies []Policy) (*Category, error) {
	category := &Category{
		ID:          validation.NewId(),
		Name:        name,
		Description: description,
		CarModels:   carModels,
		Policies:    policies,
	}

	if err := validation.ValidateEntity(category); err != nil {
		return nil, fmt.Errorf("%w\n%v", ErrInvalidEntity, err)
	}

	return category, nil
}

func (c *Category) AddModel(modelName string) error {
	for _, m := range c.CarModels {
		if m == modelName {
			return ErrInvalidModel
		}
	}

	c.CarModels = append(c.CarModels, modelName)

	return nil
}

func (c *Category) DelModel(modelName string) error {
	for i, m := range c.CarModels {
		if m == modelName {
			c.CarModels = append(c.CarModels[:i], c.CarModels[i+1:]...)
			return nil
		}
	}

	return ErrInvalidModel
}

func (c *Category) IsModelAvailable(modelName string) bool {
	flag := false
	for _, m := range c.CarModels {
		if m == modelName {
			flag = true
		}
	}
	return flag
}

func (c *Category) AddPolicy(policy Policy) error {
	for _, p := range c.Policies {
		if policy.ID == p.ID {
			return ErrInvalidPolicy
		}
	}

	c.Policies = append(c.Policies, policy)

	return nil
}

func (c *Category) DelPolicy(policyId string) error {
	for i, p := range c.Policies {
		if p.ID == policyId {
			c.Policies = append(c.Policies[:i], c.Policies[i+1:]...)
			return nil
		}
	}

	return ErrInvalidPolicy
}

func (c *Category) IsPolicyAvailable(policyId string) bool {
	flag := false
	for _, m := range c.Policies {
		if policyId == m.ID {
			flag = true
		}
	}
	return flag
}
