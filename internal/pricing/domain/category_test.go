package domain

import (
	"errors"
	"reflect"
	"testing"
)

func newPolicyFixture() *Policy {
	return &Policy{
		ID:      "5ecf09ce-8c41-4faa-a4e5-824af9c80892",
		Name:    "Promo default",
		Price:   30.5,
		Unit:    PerDay,
		MinUnit: 5,
	}
}

func newCategoryFixture() *Category {
	return &Category{
		ID:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		Name:        "Basic",
		Description: "basic cars",
		CarModels:   []string{"UNO", "MERIVA"},
		Policies: []Policy{{
			ID:      "83369771-f9a4-48b7-b87b-463f19f7b187",
			Name:    "Promo 1",
			Price:   0.2,
			Unit:    PerKM,
			MinUnit: 50,
		}, {
			ID:      "4202b708-a387-4bae-85ce-11cb7a95759d",
			Name:    "Promo 2",
			Price:   30.5,
			Unit:    PerDay,
			MinUnit: 5,
		}},
	}
}

func TestNewCategory(t *testing.T) {
	type args struct {
		name, description string
		carModels         []string
		policies          []Policy
	}

	type want struct {
		isCategory bool
		err        error
	}

	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct input",
			args: args{
				name:        "Category 1",
				description: "default category",
				carModels:   []string{"UNO", "MERIVA"},
				policies:    []Policy{*newPolicyFixture()},
			},
			want: want{
				isCategory: true,
				err:        nil,
			},
		},
		{
			name: "incorrect input",
			args: args{
				name:        "",
				description: "default category",
				carModels:   []string{"UNO", "MERIVA"},
				policies:    []Policy{*newPolicyFixture()},
			},
			want: want{
				isCategory: false,
				err:        ErrInvalidEntity,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCategory(tc.args.name, tc.args.description, tc.args.carModels, tc.args.policies)

			if reflect.ValueOf(c).IsNil() == tc.want.isCategory {
				t.Error("unexpected result", c)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestCategory_AddModel(t *testing.T) {
	newCategory := newCategoryFixture()

	type init struct {
		models []string
	}

	type args struct {
		model string
	}

	type want struct {
		models []string
		err    error
	}

	testCases := []struct {
		name string
		init init
		args args
		want want
	}{
		{
			name: "correct input",
			init: init{
				models: newCategory.CarModels,
			},
			args: args{
				model: "COROLA",
			},
			want: want{
				models: append(newCategory.CarModels, "COROLA"),
				err:    nil,
			},
		},
		{
			name: "incorrect model input",
			init: init{
				models: newCategory.CarModels,
			},
			args: args{
				model: newCategory.CarModels[0],
			},
			want: want{
				models: newCategory.CarModels,
				err:    ErrInvalidModel,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set scenario
			newCategory.CarModels = tc.init.models

			// test
			err := newCategory.AddModel(tc.args.model)

			if len(newCategory.CarModels) != len(tc.want.models) {
				t.Error("unexpected capacity value")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCategory_DelModel(t *testing.T) {
	newCategory := newCategoryFixture()

	type init struct {
		models []string
	}

	type args struct {
		model string
	}

	type want struct {
		models []string
		err    error
	}

	testCases := []struct {
		name string
		init init
		args args
		want want
	}{
		{
			name: "correct input",
			init: init{
				models: newCategory.CarModels,
			},
			args: args{
				model: newCategory.CarModels[0],
			},
			want: want{
				models: newCategory.CarModels[1:],
				err:    nil,
			},
		},
		{
			name: "incorrect model input",
			init: init{
				models: newCategory.CarModels,
			},
			args: args{
				model: "COROLA",
			},
			want: want{
				models: newCategory.CarModels,
				err:    ErrInvalidModel,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set scenario
			newCategory.CarModels = tc.init.models

			// test
			err := newCategory.DelModel(tc.args.model)

			if len(newCategory.CarModels) != len(tc.want.models) {
				t.Error("unexpected capacity value")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCategory_AddPolicy(t *testing.T) {
	newCategory := newCategoryFixture()

	type init struct {
		policies []Policy
	}

	type args struct {
		policy Policy
	}

	type want struct {
		policies []Policy
		err      error
	}

	testCases := []struct {
		name string
		init init
		args args
		want want
	}{
		{
			name: "correct input",
			init: init{
				policies: newCategory.Policies,
			},
			args: args{
				policy: *newPolicyFixture(),
			},
			want: want{
				policies: append(newCategory.Policies, *newPolicyFixture()),
				err:      nil,
			},
		},
		{
			name: "incorrect policy input",
			init: init{
				policies: newCategory.Policies,
			},
			args: args{
				policy: newCategory.Policies[0],
			},
			want: want{
				policies: newCategory.Policies,
				err:      ErrInvalidPolicy,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set scenario
			newCategory.Policies = tc.init.policies

			// test
			err := newCategory.AddPolicy(tc.args.policy)

			if len(newCategory.Policies) != len(tc.want.policies) {
				t.Error("unexpected capacity value")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCategory_DelPolicy(t *testing.T) {
	newCategory := newCategoryFixture()

	type init struct {
		policies []Policy
	}

	type args struct {
		policyId string
	}

	type want struct {
		policies []Policy
		err      error
	}

	testCases := []struct {
		name string
		init init
		args args
		want want
	}{
		{
			name: "correct input",
			init: init{
				policies: newCategory.Policies,
			},
			args: args{
				policyId: newCategory.Policies[0].ID,
			},
			want: want{
				policies: newCategory.Policies[1:],
				err:      nil,
			},
		},
		{
			name: "incorrect policyId input",
			init: init{
				policies: newCategory.Policies,
			},
			args: args{
				policyId: newPolicyFixture().ID,
			},
			want: want{
				policies: newCategory.Policies,
				err:      ErrInvalidPolicy,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set scenario
			newCategory.Policies = tc.init.policies

			// test
			err := newCategory.DelPolicy(tc.args.policyId)

			if len(newCategory.Policies) != len(tc.want.policies) {
				t.Error("unexpected capacity value")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}
