package domain

import (
	"errors"
	"reflect"
	"testing"
)

func newStationFixture() *Station {
	return &Station{
		ID:         "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		Name:       "Station 1",
		Address:    "Farway Av.",
		Complement: "45, Ap. 50",
		State:      "Polar",
		City:       "Nort City",
		Cep:        "20778990",
		Capacity:   100,
		Idle:       0,
	}
}

func TestNewStation(t *testing.T) {
	type args struct {
		id         string
		name       string
		address    string
		complement string
		state      string
		city       string
		cep        string
		capacity   uint
		idle       uint
	}

	type want struct {
		isStation bool
		err       error
	}

	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct input",
			args: args{
				id:         "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
				name:       "Station 1",
				address:    "Farway Av.",
				complement: "45, Ap. 50",
				state:      "Polar",
				city:       "Nort City",
				cep:        "20778990",
				capacity:   100,
				idle:       0,
			},
			want: want{
				isStation: true,
				err:       nil,
			},
		},
		{
			name: "incorrect cep input",
			args: args{
				id:         "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
				name:       "Station 1",
				address:    "Farway Av.",
				complement: "45, Ap. 50",
				state:      "Polar",
				city:       "Nort City",
				cep:        "20778-990",
				capacity:   100,
				idle:       0,
			},
			want: want{
				isStation: false,
				err:       ErrInvalidEntity,
			},
		},
		{
			name: "incorrect capacity input",
			args: args{
				id:         "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
				name:       "Station 1",
				address:    "Farway Av.",
				complement: "45, Ap. 50",
				state:      "Polar",
				city:       "Nort City",
				cep:        "20778990",
				capacity:   30,
				idle:       50,
			},
			want: want{
				isStation: false,
				err:       ErrInvalidCapacity,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewStation(
				tc.name,
				tc.args.address,
				tc.args.complement,
				tc.args.state,
				tc.args.city,
				tc.args.cep,
				tc.args.capacity,
				tc.args.idle,
			)

			if reflect.ValueOf(c).IsNil() == tc.want.isStation {
				t.Error("unexpected result", c)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestStation_SetCapacity(t *testing.T) {
	newStation := newStationFixture()

	type init struct {
		capacity uint
		idle     uint
	}

	type args struct {
		capacity uint
	}

	type want struct {
		capacity uint
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
				capacity: 100,
				idle:     0,
			},
			args: args{
				capacity: 50,
			},
			want: want{
				capacity: 50,
				err:      nil,
			},
		},
		{
			name: "incorrect capacity input",
			init: init{
				capacity: 100,
				idle:     90,
			},
			args: args{
				capacity: 50,
			},
			want: want{
				capacity: 100,
				err:      ErrInvalidCapacity,
			},
		},
		{
			name: "incorrect capacity input",
			init: init{
				capacity: 1,
				idle:     0,
			},
			args: args{
				capacity: 0,
			},
			want: want{
				capacity: 1,
				err:      ErrInvalidCapacity,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set scenario
			newStation.Capacity = tc.init.capacity
			newStation.Idle = tc.init.idle

			// test
			err := newStation.SetCapacity(tc.args.capacity)

			if newStation.Capacity != tc.want.capacity {
				t.Error("unexpected capacity value")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestStation_SetIdle(t *testing.T) {
	newStation := newStationFixture()

	type init struct {
		capacity uint
		idle     uint
	}

	type args struct {
		idle uint
	}

	type want struct {
		idle uint
		err  error
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
				capacity: 100,
				idle:     0,
			},
			args: args{
				idle: 1,
			},
			want: want{
				idle: 1,
				err:  nil,
			},
		},
		{
			name: "incorrect current cars input",
			init: init{
				capacity: 50,
				idle:     40,
			},
			args: args{
				idle: 60,
			},
			want: want{
				idle: 40,
				err:  ErrInvalidIdle,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set scenario
			newStation.Capacity = tc.init.capacity
			newStation.Idle = tc.init.idle

			// test
			err := newStation.SetIdle(tc.args.idle)

			if newStation.Idle != tc.want.idle {
				t.Error("unexpected current cars value")
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}
