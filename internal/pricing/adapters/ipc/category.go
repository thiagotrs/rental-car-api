package ipc

import (
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/ipc"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/domain"
)

type categoryIPC struct {
	categoryUC application.CategoryUseCase
}

func NewCategoryIPC(carUC application.CategoryUseCase) *categoryIPC {
	return &categoryIPC{carUC}
}

func (uc categoryIPC) GetPolicy(categoryId, carModel, policyId string) (*ipc.PolicyData, error) {
	category, err := uc.categoryUC.GetCategoryById(categoryId)
	if err != nil {
		return nil, application.ErrNotFoundPolicy
	}

	if !category.IsModelAvailable(carModel) {
		return nil, application.ErrNotFoundPolicy
	}

	if !category.IsPolicyAvailable(policyId) {
		return nil, application.ErrNotFoundPolicy
	}

	var policy *domain.Policy
	for _, p := range category.Policies {
		if p.ID == policyId {
			policy = &p
		}
	}

	if policy == nil {
		return nil, application.ErrNotFoundPolicy
	}

	policyData := &ipc.PolicyData{
		ID:      policy.ID,
		Name:    policy.Name,
		Price:   policy.Price,
		Unit:    uint(policy.Unit),
		MinUnit: policy.MinUnit,
	}

	return policyData, nil
}
