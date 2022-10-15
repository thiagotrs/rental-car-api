package domain

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func newCarFixture() *Car {
	return &Car{
		ID:        "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		Age:       2020,
		Plate:     "KST-9016",
		Document:  "abc.123.op-x",
		CarModel:  "UNO",
		InitialKM: 12000,
		Status:    Parked,
		StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
	}
}

func newPolicyFixture() *Policy {
	return &Policy{
		ID:         "5ecf09ce-8c41-4faa-a4e5-824af9c80892",
		Name:       "Promo default",
		Price:      30.5,
		Unit:       PerDay,
		MinUnit:    5,
		CarModel:   "UNO",
		CategoryId: "479ab9e7-ad16-4864-8e49-29b15e4b390e",
	}
}

func newOrderFixture() *Order {
	return &Order{
		ID:             "c6f31fdd-a77a-464b-9475-2d12441963a6",
		DateReservFrom: time.Now(),
		DateReservTo:   time.Now().Add(time.Hour * 24 * 5),
		Status:         Opened,
		Car:            *newCarFixture(),
		StationFromId:  "83369771-f9a4-48b7-b87b-463f19f7b187",
		StationToId:    "2520aade-a397-4e3c-a589-39c6ae5c2eff",
		Policy:         *newPolicyFixture(),
	}
}

func TestNewOrder(t *testing.T) {
	otherCar := *newCarFixture()
	otherCar.CarModel = "PORCHE"

	otherCar2 := *newCarFixture()
	otherCar2.Status = Maintenance

	type args struct {
		dateReservFrom time.Time
		dateReservTo   time.Time
		car            Car
		stationFromId  string
		stationToId    string
		policy         Policy
	}

	type want struct {
		isOrder bool
		err     error
	}

	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct input",
			args: args{
				dateReservFrom: time.Now(),
				dateReservTo:   time.Now().Add(time.Hour * 24 * 5),
				car:            *newCarFixture(),
				stationFromId:  "83369771-f9a4-48b7-b87b-463f19f7b187",
				stationToId:    "2520aade-a397-4e3c-a589-39c6ae5c2eff",
				policy:         *newPolicyFixture(),
			},
			want: want{
				isOrder: true,
				err:     nil,
			},
		},
		{
			name: "incorrect date reserve to input",
			args: args{
				dateReservFrom: time.Now(),
				dateReservTo:   time.Now().Add(time.Hour * -5),
				car:            *newCarFixture(),
				stationFromId:  "83369771-f9a4-48b7-b87b-463f19f7b187",
				stationToId:    "2520aade-a397-4e3c-a589-39c6ae5c2eff",
				policy:         *newPolicyFixture(),
			},
			want: want{
				isOrder: false,
				err:     ErrInvalidReservedDate,
			},
		},
		{
			name: "incorrect station input",
			args: args{
				dateReservFrom: time.Now(),
				dateReservTo:   time.Now().Add(time.Hour * 24 * 5),
				car:            *newCarFixture(),
				stationFromId:  "7621d238-cf12-4570-8ed1-6c0a38b76b4d",
				stationToId:    "2520aade-a397-4e3c-a589-39c6ae5c2eff",
				policy:         *newPolicyFixture(),
			},
			want: want{
				isOrder: false,
				err:     ErrInvalidCarStation,
			},
		},
		{
			name: "incorrect model input",
			args: args{
				dateReservFrom: time.Now(),
				dateReservTo:   time.Now().Add(time.Hour * 24 * 5),
				car:            otherCar,
				stationFromId:  "83369771-f9a4-48b7-b87b-463f19f7b187",
				stationToId:    "2520aade-a397-4e3c-a589-39c6ae5c2eff",
				policy:         *newPolicyFixture(),
			},
			want: want{
				isOrder: false,
				err:     ErrInvalidCarStation,
			},
		},
		{
			name: "incorrect car status input",
			args: args{
				dateReservFrom: time.Now(),
				dateReservTo:   time.Now().Add(time.Hour * 24 * 5),
				car:            otherCar2,
				stationFromId:  "83369771-f9a4-48b7-b87b-463f19f7b187",
				stationToId:    "2520aade-a397-4e3c-a589-39c6ae5c2eff",
				policy:         *newPolicyFixture(),
			},
			want: want{
				isOrder: false,
				err:     ErrInvalidReserve,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewOrder(
				tc.args.dateReservFrom,
				tc.args.dateReservTo,
				tc.args.car,
				tc.args.stationFromId,
				tc.args.stationToId,
				tc.args.policy,
			)

			if reflect.ValueOf(c).IsNil() == tc.want.isOrder {
				t.Error("unexpected result", c)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestOrder_Confirm(t *testing.T) {
	dateReservFrom := time.Now()
	dateFrom := time.Now().Add(time.Hour)

	type init struct {
		orderStatus    OrderStatus
		dateReservFrom time.Time
	}

	type args struct {
		dateFrom time.Time
	}

	type want struct {
		err      error
		status   OrderStatus
		dateFrom *time.Time
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
				orderStatus:    Opened,
				dateReservFrom: dateReservFrom,
			},
			args: args{
				dateFrom: dateFrom,
			},
			want: want{
				err:      nil,
				status:   Confirmed,
				dateFrom: &dateFrom,
			},
		},
		{
			name: "incorrect order status",
			init: init{
				orderStatus:    Confirmed,
				dateReservFrom: dateReservFrom,
			},
			args: args{
				dateFrom: dateFrom,
			},
			want: want{
				err:      ErrClose,
				status:   Confirmed,
				dateFrom: nil,
			},
		},
		{
			name: "incorrect date from input",
			init: init{
				orderStatus:    Opened,
				dateReservFrom: dateReservFrom,
			},
			args: args{
				dateFrom: time.Now().Add(time.Hour * 24 * 10),
			},
			want: want{
				err:      ErrIvalidConfirmDate,
				status:   Opened,
				dateFrom: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newOrder := newOrderFixture()
			newOrder.Car.Status = Reserved
			newOrder.Status = tc.init.orderStatus
			newOrder.DateReservFrom = tc.init.dateReservFrom

			err := newOrder.Confirm(tc.args.dateFrom)

			if newOrder.Status != tc.want.status {
				t.Error("unexpected result", newOrder.Status)
			}

			if !reflect.DeepEqual(newOrder.DateFrom, tc.want.dateFrom) {
				t.Error("unexpected result", newOrder.DateFrom)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestOrder_Cancel(t *testing.T) {
	type init struct {
		orderStatus OrderStatus
	}

	type want struct {
		err    error
		status OrderStatus
	}

	testCases := []struct {
		name string
		init init
		want want
	}{
		{
			name: "correct input",
			init: init{
				orderStatus: Opened,
			},
			want: want{
				err:    nil,
				status: Canceled,
			},
		},
		{
			name: "incorrect order status",
			init: init{
				orderStatus: Canceled,
			},
			want: want{
				err:    ErrClose,
				status: Canceled,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newOrder := newOrderFixture()
			newOrder.Status = tc.init.orderStatus
			newOrder.Car.Status = Reserved

			err := newOrder.Cancel()

			if newOrder.Status != tc.want.status {
				t.Error("unexpected result", newOrder.Status)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestOrder_Close(t *testing.T) {
	dateTo := time.Now().Add(time.Hour * 24 * 6)

	type init struct {
		orderStatus OrderStatus
	}

	type args struct {
		discount float32
		tax      float32
		dateTo   time.Time
		finalKM  uint64
	}

	type want struct {
		err      error
		status   OrderStatus
		dateTo   *time.Time
		discount float32
		tax      float32
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
				orderStatus: Confirmed,
			},
			args: args{
				discount: 10.5,
				tax:      1.2,
				dateTo:   dateTo,
				finalKM:  12050,
			},
			want: want{
				err:      nil,
				status:   Closed,
				dateTo:   &dateTo,
				discount: 10.5,
				tax:      1.2,
			},
		},
		{
			name: "incorrect order status",
			init: init{
				orderStatus: Opened,
			},
			args: args{
				discount: 10.5,
				tax:      1.2,
				dateTo:   dateTo,
				finalKM:  12050,
			},
			want: want{
				err:      ErrClose,
				status:   Opened,
				dateTo:   nil,
				discount: 0,
				tax:      0,
			},
		},
		{
			name: "incorrect date to input",
			init: init{
				orderStatus: Confirmed,
			},
			args: args{
				discount: 10.5,
				tax:      1.2,
				dateTo:   time.Now().Add(time.Hour * -2),
				finalKM:  12050,
			},
			want: want{
				err:      ErrIvalidCloseDate,
				status:   Confirmed,
				dateTo:   nil,
				discount: 0,
				tax:      0,
			},
		},
		{
			name: "incorrect tax input",
			init: init{
				orderStatus: Confirmed,
			},
			args: args{
				discount: 10.5,
				tax:      -1.2,
				dateTo:   dateTo,
				finalKM:  12050,
			},
			want: want{
				err:      ErrIvalidCloseTax,
				status:   Confirmed,
				dateTo:   nil,
				discount: 0,
				tax:      0,
			},
		},
		{
			name: "incorrect discount input",
			init: init{
				orderStatus: Confirmed,
			},
			args: args{
				discount: -10.5,
				tax:      1.2,
				dateTo:   dateTo,
				finalKM:  12050,
			},
			want: want{
				err:      ErrIvalidCloseDiscount,
				status:   Confirmed,
				dateTo:   nil,
				discount: 0,
				tax:      0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newOrder := newOrderFixture()
			newOrder.Car.Status = Reserved
			dateFrom := newOrder.DateReservFrom.Add(time.Hour)
			newOrder.DateFrom = &dateFrom
			newOrder.Status = tc.init.orderStatus

			err := newOrder.Close(tc.args.discount, tc.args.tax, tc.args.dateTo, tc.args.finalKM)

			if newOrder.Status != tc.want.status {
				t.Error("unexpected result", newOrder.Status)
			}

			if !reflect.DeepEqual(newOrder.DateTo, tc.want.dateTo) {
				t.Error("unexpected result", newOrder.DateTo)
			}

			if newOrder.Discount != tc.want.discount {
				t.Error("unexpected result", newOrder.Discount)
			}

			if newOrder.Tax != tc.want.tax {
				t.Error("unexpected result", newOrder.Tax)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}
