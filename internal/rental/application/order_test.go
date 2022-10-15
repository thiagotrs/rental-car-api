package application

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

func newCarFixture() *domain.Car {
	return &domain.Car{
		ID:        "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		Age:       2020,
		Plate:     "KST-9016",
		Document:  "abc.123.op-x",
		CarModel:  "UNO",
		InitialKM: 12000,
		Status:    domain.Reserved,
		StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
	}
}

func newPolicyFixture() *domain.Policy {
	return &domain.Policy{
		ID:         "5ecf09ce-8c41-4faa-a4e5-824af9c80892",
		Name:       "Promo default",
		Price:      30.5,
		Unit:       domain.PerDay,
		MinUnit:    5,
		CarModel:   "UNO",
		CategoryId: "479ab9e7-ad16-4864-8e49-29b15e4b390e",
	}
}

func newOrderFixture() *domain.Order {
	return &domain.Order{
		ID:             "c6f31fdd-a77a-464b-9475-2d12441963a6",
		DateReservFrom: time.Now(),
		DateReservTo:   time.Now().Add(time.Hour * 24 * 5),
		Status:         domain.Opened,
		Car:            *newCarFixture(),
		StationFromId:  "83369771-f9a4-48b7-b87b-463f19f7b187",
		StationToId:    "2520aade-a397-4e3c-a589-39c6ae5c2eff",
		Policy:         *newPolicyFixture(),
	}
}

type orderRepositoryMock struct {
	expectedFindAllOrders []domain.Order
	expectedFindOneOrder  *domain.Order
	expectedFindOneErr    error
	expectedSaveErr       error
	expectedDeleteErr     error
	calls                 map[string]uint
}

func (m *orderRepositoryMock) FindAll() []domain.Order {
	m.calls["FindAll"] = m.calls["FindAll"] + 1
	return m.expectedFindAllOrders
}

func (m *orderRepositoryMock) FindOne(id string) (*domain.Order, error) {
	m.calls["FindOne"] = m.calls["FindOne"] + 1
	return m.expectedFindOneOrder, m.expectedFindOneErr
}

func (m *orderRepositoryMock) Save(order domain.Order) error {
	m.calls["Save"] = m.calls["Save"] + 1
	return m.expectedSaveErr
}

func (m *orderRepositoryMock) Delete(id string) error {
	m.calls["Delete"] = m.calls["Delete"] + 1
	return m.expectedDeleteErr
}

type orderOrderServiceMock struct {
	expectedGetPolicy    *domain.Policy
	expectedGetPolicyErr error
	expectedGetCar       *domain.Car
	expectedGetCarErr    error
	calls                map[string]uint
}

func (m *orderOrderServiceMock) GetPolicy(categoryId, modelId, policyId string) (*domain.Policy, error) {
	m.calls["GetPolicy"] = m.calls["GetPolicy"] + 1
	return m.expectedGetPolicy, m.expectedGetPolicyErr
}

func (m *orderOrderServiceMock) GetCar(stationId, modelId string) (*domain.Car, error) {
	m.calls["GetCar"] = m.calls["GetCar"] + 1
	return m.expectedGetCar, m.expectedGetCarErr
}

func TestOrderUseCase_GetById(t *testing.T) {
	newOrder := newOrderFixture()

	type setup struct {
		repoOrder *domain.Order
		repoErr   error
	}

	type args struct {
		id string
	}

	type want struct {
		order *domain.Order
		err   error
		calls uint
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
				repoOrder: newOrder,
				repoErr:   nil,
			},
			args: args{
				id: newOrder.ID,
			},
			want: want{
				order: newOrder,
				err:   nil,
				calls: 1,
			},
		},
		{
			name:  "incorrect id input",
			setup: setup{},
			args: args{
				id: "invalid-id",
			},
			want: want{
				order: nil,
				err:   ErrInvalidId,
				calls: 0,
			},
		},
		{
			name: "not found order",
			setup: setup{
				repoOrder: nil,
				repoErr:   ErrNotFoundOrder,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				order: nil,
				err:   ErrNotFoundOrder,
				calls: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderRepo := &orderRepositoryMock{
				expectedFindOneOrder: tc.setup.repoOrder,
				expectedFindOneErr:   tc.setup.repoErr,
				calls:                make(map[string]uint),
			}
			orderSvc := &orderOrderServiceMock{}
			orderUC := NewOrderUseCase(orderRepo, orderSvc)
			order, err := orderUC.GetById(tc.args.id)

			if orderRepo.calls["FindOne"] != tc.want.calls {
				t.Error("invalid repo call", orderRepo.calls["FindOne"])
			}

			if !reflect.DeepEqual(order, tc.want.order) {
				t.Error("unequal order", order)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestOrderUseCase_Open(t *testing.T) {
	newOrder := newOrderFixture()
	carReserved := newCarFixture()
	carReserved.Status = domain.Parked

	type setup struct {
		repoGetPolicy    *domain.Policy
		repoGetPolicyErr error
		repoGetCar       *domain.Car
		repoGetCarErr    error
		repoSaveErr      error
	}

	type args struct {
		dateReservFrom, dateReservTo                               time.Time
		stationFromId, stationToId, categoryId, carModel, policyId string
	}

	type want struct {
		err            error
		getPolicyCalls uint
		getCarCalls    uint
		saveCalls      uint
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
				repoGetPolicy:    newPolicyFixture(),
				repoGetPolicyErr: nil,
				repoGetCar:       carReserved,
				repoGetCarErr:    nil,
				repoSaveErr:      nil,
			},
			args: args{
				dateReservFrom: newOrder.DateReservFrom,
				dateReservTo:   newOrder.DateReservTo,
				stationFromId:  newOrder.StationFromId,
				stationToId:    newOrder.StationToId,
				categoryId:     newOrder.Policy.CategoryId,
				carModel:       newOrder.Policy.CarModel,
				policyId:       newOrder.Policy.ID,
			},
			want: want{
				err:            nil,
				getPolicyCalls: 1,
				getCarCalls:    1,
				saveCalls:      1,
			},
		},
		{
			name: "incorrect policy input",
			setup: setup{
				repoGetPolicy:    nil,
				repoGetPolicyErr: ErrInvalidPolicy,
				repoGetCar:       nil,
				repoGetCarErr:    nil,
				repoSaveErr:      nil,
			},
			args: args{
				dateReservFrom: newOrder.DateReservFrom,
				dateReservTo:   newOrder.DateReservTo,
				stationFromId:  newOrder.StationFromId,
				stationToId:    newOrder.StationToId,
				categoryId:     "df43a454-4d84-4094-b0b0-7023817aed2a",
				carModel:       newOrder.Policy.CarModel,
				policyId:       newOrder.Policy.ID,
			},
			want: want{
				err:            ErrInvalidPolicy,
				getPolicyCalls: 1,
				getCarCalls:    0,
				saveCalls:      0,
			},
		},
		{
			name: "incorrect car input",
			setup: setup{
				repoGetPolicy:    newPolicyFixture(),
				repoGetPolicyErr: nil,
				repoGetCar:       nil,
				repoGetCarErr:    ErrInvalidCar,
				repoSaveErr:      nil,
			},
			args: args{
				dateReservFrom: newOrder.DateReservFrom,
				dateReservTo:   newOrder.DateReservTo,
				stationFromId:  "df43a454-4d84-4094-b0b0-7023817aed2a",
				stationToId:    newOrder.StationToId,
				categoryId:     newOrder.Policy.CategoryId,
				carModel:       newOrder.Policy.CarModel,
				policyId:       newOrder.Policy.ID,
			},
			want: want{
				err:            ErrInvalidCar,
				getPolicyCalls: 1,
				getCarCalls:    1,
				saveCalls:      0,
			},
		},
		{
			name: "incorrect date reserve from input",
			setup: setup{
				repoGetPolicy:    newPolicyFixture(),
				repoGetPolicyErr: nil,
				repoGetCar:       newCarFixture(),
				repoGetCarErr:    nil,
				repoSaveErr:      nil,
			},
			args: args{
				dateReservFrom: newOrder.DateReservFrom,
				dateReservTo:   newOrder.DateReservFrom.Add(time.Hour * -1),
				stationFromId:  newOrder.StationFromId,
				stationToId:    newOrder.StationToId,
				categoryId:     newOrder.Policy.CategoryId,
				carModel:       newOrder.Policy.CarModel,
				policyId:       newOrder.Policy.ID,
			},
			want: want{
				err:            domain.ErrInvalidReservedDate,
				getPolicyCalls: 1,
				getCarCalls:    1,
				saveCalls:      0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderRepo := &orderRepositoryMock{
				expectedSaveErr: tc.setup.repoSaveErr,
				calls:           make(map[string]uint),
			}
			orderSvc := &orderOrderServiceMock{
				expectedGetPolicy:    tc.setup.repoGetPolicy,
				expectedGetPolicyErr: tc.setup.repoGetPolicyErr,
				expectedGetCar:       tc.setup.repoGetCar,
				expectedGetCarErr:    tc.setup.repoGetCarErr,
				calls:                make(map[string]uint),
			}
			orderUC := NewOrderUseCase(orderRepo, orderSvc)
			err := orderUC.Open(
				tc.args.dateReservFrom,
				tc.args.dateReservTo,
				tc.args.stationFromId,
				tc.args.stationToId,
				tc.args.categoryId,
				tc.args.carModel,
				tc.args.policyId)

			if orderSvc.calls["GetPolicy"] != tc.want.getPolicyCalls {
				t.Error("invalid repo call", orderSvc.calls["GetPolicy"])
			}

			if orderSvc.calls["GetCar"] != tc.want.getCarCalls {
				t.Error("invalid repo call", orderSvc.calls["GetCar"])
			}

			if orderRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", orderRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestOrderUseCase_Confirm(t *testing.T) {
	newOrder := newOrderFixture()
	newOrder.Status = domain.Confirmed

	type setup struct {
		repoFindOrder    *domain.Order
		repoFindOrderErr error
		repoSaveErr      error
	}

	type args struct {
		id       string
		dateFrom time.Time
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
				repoFindOrder:    newOrderFixture(),
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id:       newOrderFixture().ID,
				dateFrom: time.Now().Add(time.Hour),
			},
			want: want{
				err:       nil,
				findCalls: 1,
				saveCalls: 1,
			},
		},
		{
			name: "incorrect id input",
			setup: setup{
				repoFindOrder:    nil,
				repoFindOrderErr: ErrInvalidOrder,
				repoSaveErr:      nil,
			},
			args: args{
				id:       "invalid-id",
				dateFrom: time.Now().Add(time.Hour),
			},
			want: want{
				err:       ErrInvalidOrder,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect order status",
			setup: setup{
				repoFindOrder:    newOrder,
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id:       newOrder.ID,
				dateFrom: time.Now().Add(time.Hour),
			},
			want: want{
				err:       domain.ErrClose,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect date from input",
			setup: setup{
				repoFindOrder:    newOrderFixture(),
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id:       newOrderFixture().ID,
				dateFrom: time.Now().Add(time.Hour * -5),
			},
			want: want{
				err:       domain.ErrIvalidConfirmDate,
				findCalls: 1,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderRepo := &orderRepositoryMock{
				expectedFindOneOrder: tc.setup.repoFindOrder,
				expectedFindOneErr:   tc.setup.repoFindOrderErr,
				expectedSaveErr:      tc.setup.repoSaveErr,
				calls:                make(map[string]uint),
			}
			orderSvc := &orderOrderServiceMock{}
			orderUC := NewOrderUseCase(orderRepo, orderSvc)
			err := orderUC.Confirm(tc.args.id, tc.args.dateFrom)

			if orderRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", orderSvc.calls["FindOne"])
			}

			if orderRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", orderRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestOrderUseCase_Cancel(t *testing.T) {
	newOrder := newOrderFixture()
	newOrder.Status = domain.Canceled

	type setup struct {
		repoFindOrder    *domain.Order
		repoFindOrderErr error
		repoSaveErr      error
	}

	type args struct {
		id string
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
				repoFindOrder:    newOrderFixture(),
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id: newOrderFixture().ID,
			},
			want: want{
				err:       nil,
				findCalls: 1,
				saveCalls: 1,
			},
		},
		{
			name: "incorrect id input",
			setup: setup{
				repoFindOrder:    nil,
				repoFindOrderErr: ErrInvalidOrder,
				repoSaveErr:      nil,
			},
			args: args{
				id: "invalid-id",
			},
			want: want{
				err:       ErrInvalidOrder,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect order status",
			setup: setup{
				repoFindOrder:    newOrder,
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id: newOrder.ID,
			},
			want: want{
				err:       domain.ErrClose,
				findCalls: 1,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderRepo := &orderRepositoryMock{
				expectedFindOneOrder: tc.setup.repoFindOrder,
				expectedFindOneErr:   tc.setup.repoFindOrderErr,
				expectedSaveErr:      tc.setup.repoSaveErr,
				calls:                make(map[string]uint),
			}
			orderSvc := &orderOrderServiceMock{}
			orderUC := NewOrderUseCase(orderRepo, orderSvc)
			err := orderUC.Cancel(tc.args.id)

			if orderRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", orderSvc.calls["FindOne"])
			}

			if orderRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", orderRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestOrderUseCase_Close(t *testing.T) {
	newOrder := newOrderFixture()
	newOrder.Status = domain.Confirmed
	dateFrom := time.Now().Add(time.Hour)
	newOrder.DateFrom = &dateFrom

	confirmedOrder := newOrderFixture()
	confirmedOrder.Status = domain.Confirmed
	confirmedOrder.DateFrom = &dateFrom

	type setup struct {
		repoFindOrder    *domain.Order
		repoFindOrderErr error
		repoSaveErr      error
	}

	type args struct {
		id            string
		discount, tax float32
		dateTo        time.Time
		km            uint64
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
				repoFindOrder:    newOrder,
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id:       newOrder.ID,
				discount: 10,
				tax:      1.5,
				dateTo:   time.Now().Add(time.Hour),
				km:       12500,
			},
			want: want{
				err:       nil,
				findCalls: 1,
				saveCalls: 1,
			},
		},
		{
			name: "incorrect id input",
			setup: setup{
				repoFindOrder:    nil,
				repoFindOrderErr: ErrInvalidOrder,
				repoSaveErr:      nil,
			},
			args: args{
				id:       "invalid-id",
				discount: 10,
				tax:      1.5,
				dateTo:   time.Now().Add(time.Hour),
				km:       12500,
			},
			want: want{
				err:       ErrInvalidOrder,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect order status",
			setup: setup{
				repoFindOrder:    newOrderFixture(),
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id:       newOrderFixture().ID,
				discount: 10,
				tax:      1.5,
				dateTo:   time.Now().Add(time.Hour),
				km:       12500,
			},
			want: want{
				err:       domain.ErrClose,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect tax input",
			setup: setup{
				repoFindOrder:    confirmedOrder,
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id:       confirmedOrder.ID,
				discount: 10,
				tax:      1.5,
				dateTo:   time.Now().Add(time.Hour * -10),
				km:       12500,
			},
			want: want{
				err:       domain.ErrIvalidCloseDate,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect tax input",
			setup: setup{
				repoFindOrder:    confirmedOrder,
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id:       confirmedOrder.ID,
				discount: 10,
				tax:      -1.5,
				dateTo:   time.Now().Add(time.Hour),
				km:       12500,
			},
			want: want{
				err:       domain.ErrIvalidCloseTax,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect discount input",
			setup: setup{
				repoFindOrder:    confirmedOrder,
				repoFindOrderErr: nil,
				repoSaveErr:      nil,
			},
			args: args{
				id:       confirmedOrder.ID,
				discount: -10,
				tax:      1.5,
				dateTo:   time.Now().Add(time.Hour),
				km:       12500,
			},
			want: want{
				err:       domain.ErrIvalidCloseDiscount,
				findCalls: 1,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderRepo := &orderRepositoryMock{
				expectedFindOneOrder: tc.setup.repoFindOrder,
				expectedFindOneErr:   tc.setup.repoFindOrderErr,
				expectedSaveErr:      tc.setup.repoSaveErr,
				calls:                make(map[string]uint),
			}
			orderSvc := &orderOrderServiceMock{}
			orderUC := NewOrderUseCase(orderRepo, orderSvc)
			err := orderUC.Close(tc.args.id, tc.args.discount, tc.args.tax, tc.args.dateTo, tc.args.km)

			if orderRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", orderSvc.calls["FindOne"])
			}

			if orderRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", orderRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}
