package application

import (
	"errors"
	"reflect"
	"testing"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
)

func newStationFixture() *domain.Station {
	return &domain.Station{
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

type stationRepositoryMock struct {
	expectedFindAllStations []domain.Station
	expectedFindOneStation  *domain.Station
	expectedFindOneErr      error
	expectedSaveErr         error
	expectedDeleteErr       error
	calls                   map[string]uint
}

func (m *stationRepositoryMock) FindAll() []domain.Station {
	m.calls["FindAll"] = m.calls["FindAll"] + 1
	return m.expectedFindAllStations
}

func (m *stationRepositoryMock) FindOne(id string) (*domain.Station, error) {
	m.calls["FindOne"] = m.calls["FindOne"] + 1
	return m.expectedFindOneStation, m.expectedFindOneErr
}

func (m *stationRepositoryMock) Save(station domain.Station) error {
	m.calls["Save"] = m.calls["Save"] + 1
	return m.expectedSaveErr
}

func (m *stationRepositoryMock) Delete(id string) error {
	m.calls["Delete"] = m.calls["Delete"] + 1
	return m.expectedDeleteErr
}

func TestStationUseCase_GetStations(t *testing.T) {
	newStation := newStationFixture()

	type setup struct {
		repoStations []domain.Station
	}

	type want struct {
		stations []domain.Station
	}

	testCases := []struct {
		name  string
		setup setup
		want  want
	}{
		{
			name: "correct 1",
			setup: setup{
				repoStations: []domain.Station{},
			},
			want: want{
				stations: []domain.Station{},
			},
		},
		{
			name: "correct 2",
			setup: setup{
				repoStations: []domain.Station{*newStation},
			},
			want: want{
				stations: []domain.Station{*newStation},
			},
		},
		{
			name: "unexpected error",
			setup: setup{
				repoStations: nil,
			},
			want: want{
				stations: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stationRepo := &stationRepositoryMock{
				expectedFindAllStations: tc.setup.repoStations,
				calls:                   make(map[string]uint),
			}
			stationUC := NewStationUseCase(stationRepo)
			stations := stationUC.GetStations()

			if stationRepo.calls["FindAll"] != 1 {
				t.Error("invalid repo call", stationRepo.calls["FindAll"])
			}

			if !reflect.DeepEqual(stations, tc.want.stations) {
				t.Error("unequal station", stations)
			}
		})
	}
}

func TestStationUseCase_GetStationById(t *testing.T) {
	newStation := newStationFixture()

	type setup struct {
		repoStation *domain.Station
		repoErr     error
	}

	type args struct {
		id string
	}

	type want struct {
		station *domain.Station
		err     error
		calls   uint
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
				repoStation: newStation,
				repoErr:     nil,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				station: newStation,
				err:     nil,
				calls:   1,
			},
		},
		{
			name:  "incorrect id input",
			setup: setup{},
			args: args{
				id: "invalid-id",
			},
			want: want{
				station: nil,
				err:     ErrInvalidId,
				calls:   0,
			},
		},
		{
			name: "not found station",
			setup: setup{
				repoStation: nil,
				repoErr:     ErrNotFoundStation,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				station: nil,
				err:     ErrNotFoundStation,
				calls:   1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stationRepo := &stationRepositoryMock{
				expectedFindOneStation: tc.setup.repoStation,
				expectedFindOneErr:     tc.setup.repoErr,
				calls:                  make(map[string]uint),
			}
			stationUC := NewStationUseCase(stationRepo)
			stations, err := stationUC.GetStationById(tc.args.id)

			if stationRepo.calls["FindOne"] != tc.want.calls {
				t.Error("invalid repo call", stationRepo.calls["FindOne"])
			}

			if !reflect.DeepEqual(stations, tc.want.station) {
				t.Error("unequal station", stations)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestStationUseCase_AddStation(t *testing.T) {
	type setup struct {
		repoSaveErr error
	}

	type args struct {
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
		err       error
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
				repoSaveErr: nil,
			},
			args: args{
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
				err:       nil,
				saveCalls: 1,
			},
		},
		{
			name:  "incorrect capacity input",
			setup: setup{},
			args: args{
				name:       "Station 1",
				address:    "Farway Av.",
				complement: "45, Ap. 50",
				state:      "Polar",
				city:       "Nort City",
				cep:        "20778990",
				capacity:   0,
				idle:       0,
			},
			want: want{
				err:       ErrInvalidEntity,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stationRepo := &stationRepositoryMock{
				expectedSaveErr: tc.setup.repoSaveErr,
				calls:           make(map[string]uint),
			}
			stationUC := NewStationUseCase(stationRepo)
			err := stationUC.AddStation(
				tc.args.name,
				tc.args.address,
				tc.args.complement,
				tc.args.state,
				tc.args.city,
				tc.args.cep,
				tc.args.capacity,
				tc.args.idle,
			)

			if stationRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", stationRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", tc.want.err)
			}
		})
	}
}

func TestStationUseCase_DeleteStation(t *testing.T) {
	newStation := newStationFixture()

	type setup struct {
		repoFindStation *domain.Station
		repoFindErr     error
		repoDelErr      error
		idle            uint
	}

	type args struct {
		id string
	}

	type want struct {
		err       error
		findCalls uint
		delCalls  uint
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
				repoFindStation: newStation,
				repoFindErr:     nil,
				repoDelErr:      nil,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				err:       nil,
				findCalls: 1,
				delCalls:  1,
			},
		},
		{
			name:  "incorrect id input",
			setup: setup{},
			args: args{
				id: "invalid-id",
			},
			want: want{
				err:       ErrInvalidId,
				findCalls: 0,
				delCalls:  0,
			},
		},
		{
			name: "not found station",
			setup: setup{
				repoFindStation: nil,
				repoFindErr:     ErrInvalidStation,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				err:       ErrInvalidStation,
				findCalls: 1,
				delCalls:  0,
			},
		},
		{
			name: "station has cars",
			setup: setup{
				repoFindStation: newStation,
				repoFindErr:     nil,
				idle:            5,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				err:       ErrStationHasCars,
				findCalls: 1,
				delCalls:  0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init
			if tc.setup.repoFindStation != nil {
				tc.setup.repoFindStation.Idle = tc.setup.idle
			}
			stationRepo := &stationRepositoryMock{
				expectedFindOneStation: tc.setup.repoFindStation,
				expectedFindOneErr:     tc.setup.repoFindErr,
				expectedDeleteErr:      tc.setup.repoDelErr,
				calls:                  make(map[string]uint),
			}
			stationUC := NewStationUseCase(stationRepo)
			err := stationUC.DeleteStation(tc.args.id)

			if stationRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", stationRepo.calls["FindAll"])
			}

			if stationRepo.calls["Delete"] != tc.want.delCalls {
				t.Error("invalid repo call", stationRepo.calls["Delete"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestStationUseCase_ChangeStationCapacity(t *testing.T) {
	newStation := newStationFixture()

	type setup struct {
		repoFindStation *domain.Station
		repoFindErr     error
		repoSaveErr     error
		idle            uint
	}

	type args struct {
		id       string
		capacity uint
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
				repoFindStation: newStation,
				repoFindErr:     nil,
				repoSaveErr:     nil,
			},
			args: args{
				id:       "35098f2d-6351-4509-87a2-896bab961a25",
				capacity: 150,
			},
			want: want{
				err:       nil,
				findCalls: 1,
				saveCalls: 1,
			},
		},
		{
			name:  "incorrect id input",
			setup: setup{},
			args: args{
				id:       "invalid-id",
				capacity: 150,
			},
			want: want{
				err:       ErrInvalidId,
				findCalls: 0,
				saveCalls: 0,
			},
		},
		{
			name: "not found station",
			setup: setup{
				repoFindStation: nil,
				repoFindErr:     ErrInvalidStation,
			},
			args: args{
				id:       "35098f2d-6351-4509-87a2-896bab961a25",
				capacity: 150,
			},
			want: want{
				err:       ErrInvalidStation,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "invalid capacity input",
			setup: setup{
				repoFindStation: newStation,
				repoFindErr:     nil,
				idle:            95,
			},
			args: args{
				id:       "35098f2d-6351-4509-87a2-896bab961a25",
				capacity: 90,
			},
			want: want{
				err:       ErrInvalidCapacity,
				findCalls: 1,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init
			if tc.setup.repoFindStation != nil {
				tc.setup.repoFindStation.Idle = tc.setup.idle
			}
			stationRepo := &stationRepositoryMock{
				expectedFindOneStation: tc.setup.repoFindStation,
				expectedFindOneErr:     tc.setup.repoFindErr,
				expectedSaveErr:        tc.setup.repoSaveErr,
				calls:                  make(map[string]uint),
			}
			stationUC := NewStationUseCase(stationRepo)
			err := stationUC.ChangeStationCapacity(tc.args.id, tc.args.capacity)

			if stationRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", stationRepo.calls["FindAll"])
			}

			if stationRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", stationRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", tc.want.err)
			}
		})
	}
}
