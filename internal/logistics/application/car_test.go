package application

import (
	"errors"
	"reflect"
	"testing"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
)

func newCarFixture() *domain.Car {
	return &domain.Car{
		ID:        "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
		Age:       2020,
		Plate:     "KST-9016",
		Document:  "abc.123.op-x",
		Model:     "Uno",
		Make:      "FIAT",
		StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
		KM:        12000,
		Status:    domain.Parked,
	}
}

type carRepositoryMock struct {
	expectedFindAllCars []domain.Car
	expectedFindOneCar  *domain.Car
	expectedFindOneErr  error
	expectedSaveErr     error
	expectedDeleteErr   error
	calls               map[string]uint
}

func (m *carRepositoryMock) Find(search SearchCarParams) []domain.Car {
	m.calls["FindAll"] = m.calls["FindAll"] + 1
	return m.expectedFindAllCars
}

func (m *carRepositoryMock) FindOne(id string) (*domain.Car, error) {
	m.calls["FindOne"] = m.calls["FindOne"] + 1
	return m.expectedFindOneCar, m.expectedFindOneErr
}

func (m *carRepositoryMock) Save(car domain.Car) error {
	m.calls["Save"] = m.calls["Save"] + 1
	return m.expectedSaveErr
}

func (m *carRepositoryMock) Delete(id string) error {
	m.calls["Delete"] = m.calls["Delete"] + 1
	return m.expectedDeleteErr
}

func TestCarUseCase_GetCarById(t *testing.T) {
	newCar := newCarFixture()

	type setup struct {
		repoCar *domain.Car
		repoErr error
	}

	type args struct {
		id string
	}

	type want struct {
		car   *domain.Car
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
				repoCar: newCar,
				repoErr: nil,
			},
			args: args{
				id: newCar.ID,
			},
			want: want{
				car:   newCar,
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
				car:   nil,
				err:   ErrInvalidId,
				calls: 0,
			},
		},
		{
			name: "not found car",
			setup: setup{
				repoCar: nil,
				repoErr: ErrNotFoundCar,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				car:   nil,
				err:   ErrNotFoundCar,
				calls: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			carRepo := &carRepositoryMock{
				expectedFindOneCar: tc.setup.repoCar,
				expectedFindOneErr: tc.setup.repoErr,
				calls:              make(map[string]uint),
			}
			stationRepo := &stationRepositoryMock{}
			carUC := NewCarUseCase(carRepo, stationRepo)
			cars, err := carUC.GetCarById(tc.args.id)

			if carRepo.calls["FindOne"] != tc.want.calls {
				t.Error("invalid repo call", carRepo.calls["FindAll"])
			}

			if !reflect.DeepEqual(cars, tc.want.car) {
				t.Error("unequal car", cars)
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error")
			}
		})
	}
}

func TestCarUseCase_AddCar(t *testing.T) {
	newStation := newStationFixture()

	maxCapStation := newStationFixture()
	maxCapStation.Idle = maxCapStation.Capacity

	type setup struct {
		repoStation    *domain.Station
		repoStationErr error
		repoCarErr     error
	}

	type args struct {
		age       uint16
		km        uint64
		plate     string
		document  string
		stationId string
		model     string
		make      string
	}

	type want struct {
		err          error
		carCalls     uint
		stationCalls uint
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
				repoStation:    newStation,
				repoStationErr: nil,
				repoCarErr:     nil,
			},
			args: args{
				age:       2020,
				km:        12000,
				plate:     "KST-9016",
				document:  "abc.123.op-x",
				stationId: newStation.ID,
				model:     "Uno",
				make:      "FIAT",
			},
			want: want{
				err:          nil,
				carCalls:     1,
				stationCalls: 1,
			},
		},
		{
			name: "incorrect station id input",
			setup: setup{
				repoStation:    nil,
				repoStationErr: ErrInvalidEntity,
				repoCarErr:     nil,
			},
			args: args{
				age:       2020,
				km:        12000,
				plate:     "KST-9016",
				document:  "abc.123.op-x",
				stationId: "invalid-id",
				model:     "Uno",
				make:      "FIAT",
			},
			want: want{
				err:          ErrInvalidEntity,
				carCalls:     0,
				stationCalls: 1,
			},
		},
		{
			name: "incorrect plate input",
			setup: setup{
				repoStation:    newStation,
				repoStationErr: nil,
				repoCarErr:     nil,
			},
			args: args{
				age:       2020,
				km:        12000,
				plate:     "",
				document:  "abc.123.op-x",
				stationId: newStation.ID,
				model:     "Uno",
				make:      "FIAT",
			},
			want: want{
				err:          ErrInvalidEntity,
				carCalls:     0,
				stationCalls: 1,
			},
		},
		{
			name: "incorrect plate input",
			setup: setup{
				repoStation:    maxCapStation,
				repoStationErr: nil,
				repoCarErr:     nil,
			},
			args: args{
				age:       2020,
				km:        12000,
				plate:     "KST-9016",
				document:  "abc.123.op-x",
				stationId: maxCapStation.ID,
				model:     "Uno",
				make:      "FIAT",
			},
			want: want{
				err:          ErrStationMaxCapacity,
				carCalls:     0,
				stationCalls: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			carRepo := &carRepositoryMock{
				expectedSaveErr: tc.setup.repoCarErr,
				calls:           make(map[string]uint),
			}
			stationRepo := &stationRepositoryMock{
				expectedFindOneStation: tc.setup.repoStation,
				expectedFindOneErr:     tc.setup.repoStationErr,
				calls:                  make(map[string]uint),
			}
			carUC := NewCarUseCase(carRepo, stationRepo)
			err := carUC.AddCar(tc.args.age, tc.args.km, tc.args.plate, tc.args.document, tc.args.stationId, tc.args.model, tc.args.make)

			if stationRepo.calls["FindOne"] != tc.want.stationCalls {
				t.Error("invalid repo call", stationRepo.calls["FindOne"])
			}
			if carRepo.calls["Save"] != tc.want.carCalls {
				t.Error("invalid repo call", carRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestCarUseCase_DeleteCar(t *testing.T) {
	newCar := newCarFixture()
	newCar.Status = domain.Maintenance

	activeCar := newCarFixture()

	type setup struct {
		repoFindCar *domain.Car
		repoFindErr error
		repoDelErr  error
	}

	type args struct {
		id string
	}

	type want struct {
		err         error
		findCalls   uint
		deleteCalls uint
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
				repoFindCar: newCar,
				repoFindErr: nil,
				repoDelErr:  nil,
			},
			args: args{
				id: newCar.ID,
			},
			want: want{
				err:         nil,
				findCalls:   1,
				deleteCalls: 1,
			},
		},
		{
			name: "incorrect id input",
			setup: setup{
				repoFindCar: nil,
				repoFindErr: nil,
				repoDelErr:  nil,
			},
			args: args{
				id: "invalid-id",
			},
			want: want{
				err:         ErrInvalidId,
				findCalls:   0,
				deleteCalls: 0,
			},
		},
		{
			name: "not found car",
			setup: setup{
				repoFindCar: nil,
				repoFindErr: ErrInvalidCar,
				repoDelErr:  nil,
			},
			args: args{
				id: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				err:         ErrInvalidCar,
				findCalls:   1,
				deleteCalls: 0,
			},
		},
		{
			name: "incorrect car status",
			setup: setup{
				repoFindCar: activeCar,
				repoFindErr: nil,
				repoDelErr:  nil,
			},
			args: args{
				id: activeCar.ID,
			},
			want: want{
				err:         ErrCarNotInMaintenance,
				findCalls:   1,
				deleteCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			carRepo := &carRepositoryMock{
				expectedFindOneCar: tc.setup.repoFindCar,
				expectedFindOneErr: tc.setup.repoFindErr,
				expectedDeleteErr:  tc.setup.repoDelErr,
				calls:              make(map[string]uint),
			}
			stationRepo := &stationRepositoryMock{}
			carUC := NewCarUseCase(carRepo, stationRepo)
			err := carUC.DeleteCar(tc.args.id)

			if carRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", carRepo.calls["FindOne"])
			}

			if carRepo.calls["Delete"] != tc.want.deleteCalls {
				t.Error("invalid repo call", carRepo.calls["Delete"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err, tc.want.err)
			}
		})
	}
}

func TestCarUseCase_MoveCarToMaintenance(t *testing.T) {
	newCar := newCarFixture()

	underFixCar := newCarFixture()
	underFixCar.Status = domain.Maintenance

	type setup struct {
		repoFindCar *domain.Car
		repoFindErr error
		repoSaveErr error
	}

	type args struct {
		id        string
		stationId string
		km        uint64
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
				repoFindCar: newCar,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				id:        newCar.ID,
				stationId: newCar.StationId,
				km:        newCar.KM,
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
				repoFindCar: nil,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				id:        "invalid-id",
				stationId: newCar.StationId,
				km:        newCar.KM,
			},
			want: want{
				err:       ErrInvalidId,
				findCalls: 0,
				saveCalls: 0,
			},
		},
		{
			name: "not found car",
			setup: setup{
				repoFindCar: nil,
				repoFindErr: ErrInvalidCar,
				repoSaveErr: nil,
			},
			args: args{
				id:        "35098f2d-6351-4509-87a2-896bab961a25",
				stationId: newCar.StationId,
				km:        newCar.KM,
			},
			want: want{
				err:       ErrInvalidCar,
				findCalls: 1,
				saveCalls: 0,
			},
		},
		{
			name: "incorrect status car",
			setup: setup{
				repoFindCar: underFixCar,
				repoFindErr: nil,
				repoSaveErr: nil,
			},
			args: args{
				id:        underFixCar.ID,
				stationId: underFixCar.StationId,
				km:        underFixCar.KM,
			},
			want: want{
				err:       ErrInvalidMaintenance,
				findCalls: 1,
				saveCalls: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			carRepo := &carRepositoryMock{
				expectedFindOneCar: tc.setup.repoFindCar,
				expectedFindOneErr: tc.setup.repoFindErr,
				expectedSaveErr:    tc.setup.repoSaveErr,
				calls:              make(map[string]uint),
			}
			stationRepo := &stationRepositoryMock{}
			carUC := NewCarUseCase(carRepo, stationRepo)
			err := carUC.MoveCarToMaintenance(tc.args.id, tc.args.stationId, tc.args.km)

			if carRepo.calls["FindOne"] != tc.want.findCalls {
				t.Error("invalid repo call", carRepo.calls["FindOne"])
			}

			if carRepo.calls["Save"] != tc.want.saveCalls {
				t.Error("invalid repo call", carRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err, tc.want.err)
			}
		})
	}
}

func TestCarUseCase_ParkCar(t *testing.T) {
	newCar := newCarFixture()
	newCar.Status = domain.Transfer
	newStation := newStationFixture()
	newStation.ID = newCar.StationId

	maxCapStation := newStationFixture()
	maxCapStation.Idle = maxCapStation.Capacity

	type setup struct {
		repoFindStation    *domain.Station
		repoFindErrStation error
		repoFindCar        *domain.Car
		repoFindCarErr     error
		repoSaveCarErr     error
	}

	type args struct {
		id        string
		stationId string
		km        uint64
	}

	type want struct {
		err              error
		findStationCalls uint
		findCarCalls     uint
		saveCarCalls     uint
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
				repoFindStation:    newStation,
				repoFindErrStation: nil,
				repoFindCar:        newCar,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        newCar.ID,
				stationId: newCar.StationId,
				km:        newCar.KM,
			},
			want: want{
				err:              nil,
				findStationCalls: 1,
				findCarCalls:     1,
				saveCarCalls:     1,
			},
		},
		{
			name: "incorrect id input",
			setup: setup{
				repoFindStation:    nil,
				repoFindErrStation: nil,
				repoFindCar:        nil,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        "invalid-id",
				stationId: newCar.StationId,
				km:        newCar.KM,
			},
			want: want{
				err:              ErrInvalidId,
				findStationCalls: 0,
				findCarCalls:     0,
				saveCarCalls:     0,
			},
		},
		{
			name: "not found station id",
			setup: setup{
				repoFindStation:    nil,
				repoFindErrStation: ErrInvalidEntity,
				repoFindCar:        nil,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        newCar.ID,
				stationId: "35098f2d-6351-4509-87a2-896bab961a25",
				km:        newCar.KM,
			},
			want: want{
				err:              ErrInvalidEntity,
				findStationCalls: 1,
				findCarCalls:     0,
				saveCarCalls:     0,
			},
		},
		{
			name: "incorrect max capacity car",
			setup: setup{
				repoFindStation:    maxCapStation,
				repoFindErrStation: nil,
				repoFindCar:        nil,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        newCar.ID,
				stationId: newCar.StationId,
				km:        newCar.KM,
			},
			want: want{
				err:              ErrStationMaxCapacity,
				findStationCalls: 1,
				findCarCalls:     0,
				saveCarCalls:     0,
			},
		},
		{
			name: "not found car id",
			setup: setup{
				repoFindStation:    newStation,
				repoFindErrStation: nil,
				repoFindCar:        nil,
				repoFindCarErr:     ErrInvalidCar,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        "35098f2d-6351-4509-87a2-896bab961a25",
				stationId: newCar.StationId,
				km:        newCar.KM,
			},
			want: want{
				err:              ErrInvalidCar,
				findStationCalls: 1,
				findCarCalls:     1,
				saveCarCalls:     0,
			},
		},
		{
			name: "incorrect kilometrage input",
			setup: setup{
				repoFindStation:    newStation,
				repoFindErrStation: nil,
				repoFindCar:        newCar,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        newCar.ID,
				stationId: newCar.StationId,
				km:        newCar.KM - 10,
			},
			want: want{
				err:              ErrInvalidPark,
				findStationCalls: 1,
				findCarCalls:     1,
				saveCarCalls:     0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			carRepo := &carRepositoryMock{
				expectedFindOneCar: tc.setup.repoFindCar,
				expectedFindOneErr: tc.setup.repoFindCarErr,
				expectedSaveErr:    tc.setup.repoSaveCarErr,
				calls:              make(map[string]uint),
			}
			stationRepo := &stationRepositoryMock{
				expectedFindOneStation: tc.setup.repoFindStation,
				expectedFindOneErr:     tc.setup.repoFindErrStation,
				calls:                  make(map[string]uint),
			}
			carUC := NewCarUseCase(carRepo, stationRepo)
			err := carUC.ParkCar(tc.args.id, tc.args.stationId, tc.args.km)

			if stationRepo.calls["FindOne"] != tc.want.findStationCalls {
				t.Error("invalid repo call", carRepo.calls["FindOne"])
			}

			if carRepo.calls["FindOne"] != tc.want.findCarCalls {
				t.Error("invalid repo call", carRepo.calls["FindOne"])
			}

			if carRepo.calls["Save"] != tc.want.saveCarCalls {
				t.Error("invalid repo call", carRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err, tc.want.err)
			}
		})
	}
}

func TestCarUseCase_TransferCar(t *testing.T) {
	newCar := newCarFixture()
	newStation := newStationFixture()
	newStation.ID = newCar.StationId

	maxCapStation := newStationFixture()
	maxCapStation.Idle = maxCapStation.Capacity

	type setup struct {
		repoFindStation    *domain.Station
		repoFindErrStation error
		repoFindCar        *domain.Car
		repoFindCarErr     error
		repoSaveCarErr     error
	}

	type args struct {
		id        string
		stationId string
	}

	type want struct {
		err              error
		findStationCalls uint
		findCarCalls     uint
		saveCarCalls     uint
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
				repoFindStation:    newStation,
				repoFindErrStation: nil,
				repoFindCar:        newCar,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        newCar.ID,
				stationId: newStation.ID,
			},
			want: want{
				err:              nil,
				findStationCalls: 1,
				findCarCalls:     1,
				saveCarCalls:     1,
			},
		},
		{
			name: "incorrect id input",
			setup: setup{
				repoFindStation:    nil,
				repoFindErrStation: nil,
				repoFindCar:        nil,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        "invalid-id",
				stationId: newStation.ID,
			},
			want: want{
				err:              ErrInvalidId,
				findStationCalls: 0,
				findCarCalls:     0,
				saveCarCalls:     0,
			},
		},
		{
			name: "not found station id",
			setup: setup{
				repoFindStation:    nil,
				repoFindErrStation: ErrInvalidEntity,
				repoFindCar:        nil,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        newCar.ID,
				stationId: "35098f2d-6351-4509-87a2-896bab961a25",
			},
			want: want{
				err:              ErrInvalidEntity,
				findStationCalls: 1,
				findCarCalls:     0,
				saveCarCalls:     0,
			},
		},
		{
			name: "incorrect max capacity car",
			setup: setup{
				repoFindStation:    maxCapStation,
				repoFindErrStation: nil,
				repoFindCar:        nil,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        newCar.ID,
				stationId: maxCapStation.ID,
			},
			want: want{
				err:              ErrStationMaxCapacity,
				findStationCalls: 1,
				findCarCalls:     0,
				saveCarCalls:     0,
			},
		},
		{
			name: "not found car id",
			setup: setup{
				repoFindStation:    newStation,
				repoFindErrStation: nil,
				repoFindCar:        nil,
				repoFindCarErr:     ErrInvalidCar,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        "35098f2d-6351-4509-87a2-896bab961a25",
				stationId: newStation.ID,
			},
			want: want{
				err:              ErrInvalidCar,
				findStationCalls: 1,
				findCarCalls:     1,
				saveCarCalls:     0,
			},
		},
		{
			name: "incorrect car status",
			setup: setup{
				repoFindStation:    newStation,
				repoFindErrStation: nil,
				repoFindCar:        newCar,
				repoFindCarErr:     nil,
				repoSaveCarErr:     nil,
			},
			args: args{
				id:        newCar.ID,
				stationId: newStation.ID,
			},
			want: want{
				err:              ErrInvalidTransfer,
				findStationCalls: 1,
				findCarCalls:     1,
				saveCarCalls:     0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			carRepo := &carRepositoryMock{
				expectedFindOneCar: tc.setup.repoFindCar,
				expectedFindOneErr: tc.setup.repoFindCarErr,
				expectedSaveErr:    tc.setup.repoSaveCarErr,
				calls:              make(map[string]uint),
			}
			stationRepo := &stationRepositoryMock{
				expectedFindOneStation: tc.setup.repoFindStation,
				expectedFindOneErr:     tc.setup.repoFindErrStation,
				calls:                  make(map[string]uint),
			}
			carUC := NewCarUseCase(carRepo, stationRepo)
			err := carUC.TransferCar(tc.args.id, tc.args.stationId)

			if stationRepo.calls["FindOne"] != tc.want.findStationCalls {
				t.Error("invalid repo call", carRepo.calls["FindOne"])
			}

			if carRepo.calls["FindOne"] != tc.want.findCarCalls {
				t.Error("invalid repo call", carRepo.calls["FindOne"])
			}

			if carRepo.calls["Save"] != tc.want.saveCarCalls {
				t.Error("invalid repo call", carRepo.calls["Save"])
			}

			if !errors.Is(err, tc.want.err) {
				t.Error("unexpected error", err, tc.want.err)
			}
		})
	}
}
