package domain

import (
	"errors"
	"reflect"
	"testing"
)

func newCarFixture() *Car {
	return &Car{
		ID:        "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		Age:       2020,
		Plate:     "KST-9016",
		Document:  "abc.123.op-x",
		Model:     "Uno",
		Make:      "FIAT",
		StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
		KM:        12000,
		Status:    Parked,
	}
}

func TestNewCar(t *testing.T) {
	type args struct {
		age       uint16
		plate     string
		document  string
		model     string
		make      string
		stationId string
		km        uint64
	}

	type want struct {
		isCar bool
		err   error
	}

	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct input",
			args: args{
				age:       2020,
				plate:     "KST-9016",
				document:  "abc.123.op-x",
				model:     "Uno",
				make:      "fiat",
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
				km:        12000,
			},
			want: want{
				isCar: true,
				err:   nil,
			},
		},
		{
			name: "incorrect owner station id input",
			args: args{
				age:       2020,
				plate:     "KST-9016",
				document:  "abc.123.op-x",
				model:     "Uno",
				make:      "fiat",
				stationId: "incorrect-id",
				km:        12000,
			},
			want: want{
				isCar: false,
				err:   ErrInvalidEntity,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCar(
				tc.args.age,
				tc.args.km,
				tc.args.plate,
				tc.args.document,
				tc.args.stationId,
				tc.args.model,
				tc.args.make,
			)

			if reflect.ValueOf(c).IsNil() == tc.want.isCar {
				t.Error("unexpected result")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestCar_ToMaintenance(t *testing.T) {
	type init struct {
		stationId string
		status    CarStatus
		km        uint64
	}

	type args struct {
		stationId string
		km        uint64
	}

	type want struct {
		stationId string
		status    CarStatus
		km        uint64
		err       error
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
				status:    Parked,
				km:        12005,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
			},
			args: args{
				km:        12005,
				stationId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			want: want{
				status:    Maintenance,
				km:        12005,
				stationId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
				err:       nil,
			},
		},
		{
			name: "incorrect km input",
			init: init{
				status:    Parked,
				km:        12005,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
			},
			args: args{
				km:        12000,
				stationId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			want: want{
				status:    Parked,
				km:        12005,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
				err:       ErrInvalidMaintenance,
			},
		},
		{
			name: "incorrect car status",
			init: init{
				status:    Maintenance,
				km:        12005,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
			},
			args: args{
				km:        12005,
				stationId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			want: want{
				status:    Maintenance,
				km:        12005,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
				err:       ErrInvalidMaintenance,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newCar := newCarFixture()
			newCar.Status = Parked
			newCar.KM = tc.init.km
			newCar.StationId = tc.init.stationId
			newCar.Status = tc.init.status

			err := newCar.ToMaintenance(tc.args.stationId, tc.args.km)

			if newCar.Status != tc.want.status {
				t.Error("unexpected status value")
			}

			if newCar.KM != tc.want.km {
				t.Error("unexpected kilometrage value")
			}

			if newCar.StationId != tc.want.stationId {
				t.Error("unexpected StationId value")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCar_Transfer(t *testing.T) {
	type init struct {
		status    CarStatus
		stationId string
	}

	type args struct {
		stationId string
	}

	type want struct {
		status    CarStatus
		stationId string
		err       error
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
				status:    Parked,
				stationId: "11ab50ac-d649-4fdd-b5bb-d9c1ac2fdfdd",
			},
			args: args{
				stationId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			want: want{
				status:    Transfer,
				stationId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
				err:       nil,
			},
		},
		{
			name: "incorrect car status",
			init: init{
				status:    Transit,
				stationId: "11ab50ac-d649-4fdd-b5bb-d9c1ac2fdfdd",
			},
			args: args{
				stationId: "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			},
			want: want{
				status:    Transit,
				stationId: "11ab50ac-d649-4fdd-b5bb-d9c1ac2fdfdd",
				err:       ErrInvalidTransfer,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newCar := newCarFixture()
			newCar.StationId = tc.init.stationId
			newCar.Status = tc.init.status

			err := newCar.Transfer(tc.args.stationId)

			if newCar.Status != tc.want.status {
				t.Error("unexpected status value")
			}

			if newCar.StationId != tc.want.stationId {
				t.Error("unexpected currStationId value")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCar_Park(t *testing.T) {
	type init struct {
		stationId string
		status    CarStatus
		km        uint64
	}

	type args struct {
		stationId string
		km        uint64
	}

	type want struct {
		stationId string
		status    CarStatus
		km        uint64
		err       error
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
				status:    Transfer,
				km:        12050,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
			},
			args: args{
				km:        12050,
				stationId: "11ab50ac-d649-4fdd-b5bb-d9c1ac2fdfdd",
			},
			want: want{
				status:    Parked,
				km:        12050,
				stationId: "11ab50ac-d649-4fdd-b5bb-d9c1ac2fdfdd",
				err:       nil,
			},
		},
		{
			name: "incorrect km input",
			init: init{
				status:    Transit,
				km:        12050,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
			},
			args: args{
				km:        12000,
				stationId: "11ab50ac-d649-4fdd-b5bb-d9c1ac2fdfdd",
			},
			want: want{
				status:    Transit,
				km:        12050,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
				err:       ErrInvalidPark,
			},
		},
		{
			name: "incorrect park",
			init: init{
				status:    Parked,
				km:        12000,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
			},
			args: args{
				km:        12050,
				stationId: "11ab50ac-d649-4fdd-b5bb-d9c1ac2fdfdd",
			},
			want: want{
				status:    Parked,
				km:        12000,
				stationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
				err:       ErrInvalidPark,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newCar := newCarFixture()
			newCar.Status = tc.init.status
			newCar.StationId = tc.init.stationId
			newCar.KM = tc.init.km

			err := newCar.Park(tc.args.stationId, tc.args.km)

			if newCar.Status != tc.want.status {
				t.Error("unexpected status value")
			}

			if newCar.KM != tc.want.km {
				t.Error("unexpected kilometrage value")
			}

			if newCar.StationId != tc.want.stationId {
				t.Error("unexpected stationId value")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}
