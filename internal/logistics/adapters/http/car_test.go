package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/adapters/repository"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"
)

func newCarFixture() *domain.Car {
	c, _ := domain.NewCar(2020, 12000, "KST-9016", "abc.123.op-x", "83369771-f9a4-48b7-b87b-463f19f7b187", "Uno", "FIAT")
	return c
}

func TestCarController_GetCarById(t *testing.T) {
	cars := []domain.Car{*newCarFixture(), *newCarFixture()}
	stations := []domain.Station{}

	carRepo := repository.NewCarRepositoryInMemory(cars)
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	carUC := application.NewCarUseCase(carRepo, stationRepo)
	carController := NewCarController(carUC)

	testCases := []struct {
		name           string
		idArg          string
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          cars[0].ID,
			wantStatusCode: http.StatusOK,
			wantBody:       cars[0],
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       map[string]string{"error": application.ErrInvalidId.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusNotFound,
			wantBody:       map[string]string{"error": application.ErrNotFoundCar.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/cars/"+tc.idArg, nil)
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/cars/{id}", carController.GetCarById).Methods("GET")
			router.ServeHTTP(res, req)

			if res.Code != tc.wantStatusCode {
				t.Error("wrong response code", res.Code)
			}

			json, _ := json.Marshal(tc.wantBody)
			expectedBody := strings.Trim(string(json), "\n")

			if strings.Trim(res.Body.String(), "\n") != expectedBody {
				t.Error("wrong response body", strings.Trim(res.Body.String(), "\n"), expectedBody)
			}
		})
	}
}

func TestCarController_DeleteCar(t *testing.T) {
	cars := []domain.Car{*newCarFixture(), *newCarFixture()}
	cars[0].Status = domain.Maintenance
	stations := []domain.Station{}

	carRepo := repository.NewCarRepositoryInMemory(cars)
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	carUC := application.NewCarUseCase(carRepo, stationRepo)
	carController := NewCarController(carUC)

	testCases := []struct {
		name           string
		idArg          string
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          cars[0].ID,
			wantStatusCode: http.StatusNoContent,
			wantBody:       nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       map[string]string{"error": application.ErrInvalidId.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusNotFound,
			wantBody:       map[string]string{"error": application.ErrInvalidCar.Error()},
		},
		{
			name:           "incorrect station req",
			idArg:          cars[1].ID,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       map[string]string{"error": application.ErrCarNotInMaintenance.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/cars/"+tc.idArg, nil)
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/cars/{id}", carController.DeleteCar).Methods("DELETE")
			router.ServeHTTP(res, req)

			if res.Code != tc.wantStatusCode {
				t.Error("wrong response code", res.Code)
			}

			var expectedBody string
			if tc.wantBody != nil {
				json, _ := json.Marshal(tc.wantBody)
				expectedBody = strings.Trim(string(json), "\n")
			}

			if strings.Trim(res.Body.String(), "\n") != expectedBody {
				t.Error("wrong response body", strings.Trim(res.Body.String(), "\n"), expectedBody)
			}
		})
	}
}

func TestCarController_CreateCar(t *testing.T) {
	stations := []domain.Station{*newStationFixture()}
	cars := []domain.Car{*newCarFixture()}
	cars[0].StationId = stations[0].ID
	carRepo := repository.NewCarRepositoryInMemory(cars)
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	carUC := application.NewCarUseCase(carRepo, stationRepo)
	carController := NewCarController(carUC)

	type params struct {
		Age                                     uint16
		KM                                      uint64
		Plate, Document, StationId, Model, Make string
	}

	testCases := []struct {
		name           string
		bodyArg        interface{}
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			wantStatusCode: http.StatusCreated,
			bodyArg: params{
				Age:       2020,
				KM:        12000,
				Plate:     "KST-9016",
				Document:  "abc.123.op-x",
				Model:     "Uno",
				Make:      "FIAT",
				StationId: stations[0].ID,
			},
			wantBody: nil,
		},
		{
			name:           "incorrect station id body req",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				Age:       2020,
				KM:        12000,
				Plate:     "KST-9016",
				Document:  "abc.123.op-x",
				Model:     "Uno",
				Make:      "FIAT",
				StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
			},
			wantBody: map[string]string{"error": application.ErrInvalidEntity.Error()},
		},
		{
			name:           "incorrect malformed body req",
			wantStatusCode: http.StatusBadRequest,
			bodyArg:        "",
			wantBody:       map[string]string{"error": application.ErrInvalidEntity.Error()},
		},
		{
			name:           "incorrect plate body req",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				Age:       2020,
				KM:        12000,
				Plate:     "",
				Document:  "abc.123.op-x",
				Model:     "Uno",
				Make:      "FIAT",
				StationId: "83369771-f9a4-48b7-b87b-463f19f7b187",
			},
			wantBody: map[string]string{"error": application.ErrInvalidEntity.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("POST", "/cars/", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/cars/", carController.CreateCar).Methods("POST")
			router.ServeHTTP(res, req)

			if res.Code != tc.wantStatusCode {
				t.Error("wrong response code", res.Code)
			}

			var expectedBody string
			if tc.wantBody != nil {
				json, _ := json.Marshal(tc.wantBody)
				expectedBody = strings.Trim(string(json), "\n")
			}

			if strings.Trim(res.Body.String(), "\n") != expectedBody {
				t.Error("wrong response body", strings.Trim(res.Body.String(), "\n"), expectedBody)
			}
		})
	}
}

func TestCarController_UpdateCarToMaintenance(t *testing.T) {
	stations := []domain.Station{*newStationFixture()}
	maintenanceCar := *newCarFixture()
	maintenanceCar.StationId = stations[0].ID
	maintenanceCar.Status = domain.Maintenance
	cars := []domain.Car{*newCarFixture(), maintenanceCar}
	cars[0].StationId = stations[0].ID
	carRepo := repository.NewCarRepositoryInMemory(cars)
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	carUC := application.NewCarUseCase(carRepo, stationRepo)
	carController := NewCarController(carUC)

	type params struct {
		StationId string
		KM        uint64
	}

	testCases := []struct {
		name           string
		idArg          string
		bodyArg        interface{}
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          cars[0].ID,
			wantStatusCode: http.StatusNoContent,
			bodyArg: params{
				StationId: cars[0].StationId,
				KM:        cars[0].KM + 5,
			},
			wantBody: nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				StationId: cars[0].StationId,
				KM:        cars[0].KM + 5,
			},
			wantBody: map[string]string{"error": application.ErrInvalidId.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusNotFound,
			bodyArg: params{
				StationId: cars[0].StationId,
				KM:        cars[0].KM + 5,
			},
			wantBody: map[string]string{"error": application.ErrInvalidCar.Error()},
		},
		{
			name:           "incorrect car status req",
			idArg:          maintenanceCar.ID,
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				StationId: maintenanceCar.StationId,
				KM:        maintenanceCar.KM + 5,
			},
			wantBody: map[string]string{"error": application.ErrInvalidMaintenance.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("PUT", "/cars/"+tc.idArg+"/maintenance", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/cars/{id}/maintenance", carController.UpdateCarToMaintenance).Methods("PUT")
			router.ServeHTTP(res, req)

			if res.Code != tc.wantStatusCode {
				t.Error("wrong response code", res.Code)
			}

			var expectedBody string
			if tc.wantBody != nil {
				json, _ := json.Marshal(tc.wantBody)
				expectedBody = strings.Trim(string(json), "\n")
			}

			if strings.Trim(res.Body.String(), "\n") != expectedBody {
				t.Error("wrong response body", strings.Trim(res.Body.String(), "\n"), expectedBody)
			}
		})
	}
}

func TestCarController_UpdateCarToPark(t *testing.T) {
	stations := []domain.Station{*newStationFixture()}
	parkedCar := *newCarFixture()
	parkedCar.StationId = stations[0].ID
	transferCar := *newCarFixture()
	transferCar.StationId = stations[0].ID
	transferCar.Status = domain.Transfer
	cars := []domain.Car{transferCar, parkedCar}
	carRepo := repository.NewCarRepositoryInMemory(cars)
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	carUC := application.NewCarUseCase(carRepo, stationRepo)
	carController := NewCarController(carUC)

	type params struct {
		StationId string
		KM        uint64
	}

	testCases := []struct {
		name           string
		idArg          string
		bodyArg        interface{}
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          transferCar.ID,
			wantStatusCode: http.StatusNoContent,
			bodyArg: params{
				StationId: transferCar.StationId,
				KM:        transferCar.KM + 5,
			},
			wantBody: nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				StationId: transferCar.StationId,
				KM:        transferCar.KM + 5,
			},
			wantBody: map[string]string{"error": application.ErrInvalidId.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusNotFound,
			bodyArg: params{
				StationId: transferCar.StationId,
				KM:        transferCar.KM + 5,
			},
			wantBody: map[string]string{"error": application.ErrInvalidCar.Error()},
		},
		{
			name:           "incorrect car status req",
			idArg:          parkedCar.ID,
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				StationId: parkedCar.StationId,
				KM:        parkedCar.KM + 5,
			},
			wantBody: map[string]string{"error": application.ErrInvalidPark.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("PUT", "/cars/"+tc.idArg+"/park", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/cars/{id}/park", carController.UpdateCarToPark).Methods("PUT")
			router.ServeHTTP(res, req)

			if res.Code != tc.wantStatusCode {
				t.Error("wrong response code", res.Code)
			}

			var expectedBody string
			if tc.wantBody != nil {
				json, _ := json.Marshal(tc.wantBody)
				expectedBody = strings.Trim(string(json), "\n")
			}

			if strings.Trim(res.Body.String(), "\n") != expectedBody {
				t.Error("wrong response body", strings.Trim(res.Body.String(), "\n"), expectedBody)
			}
		})
	}
}

func TestCarController_UpdateCarToTransfer(t *testing.T) {
	stations := []domain.Station{*newStationFixture()}
	parkedCar := *newCarFixture()
	parkedCar.StationId = stations[0].ID
	transferCar := *newCarFixture()
	transferCar.StationId = stations[0].ID
	transferCar.Status = domain.Transfer
	cars := []domain.Car{transferCar, parkedCar}
	carRepo := repository.NewCarRepositoryInMemory(cars)
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	carUC := application.NewCarUseCase(carRepo, stationRepo)
	carController := NewCarController(carUC)

	type params struct {
		StationId string
	}

	testCases := []struct {
		name           string
		idArg          string
		bodyArg        interface{}
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          parkedCar.ID,
			wantStatusCode: http.StatusNoContent,
			bodyArg: params{
				StationId: parkedCar.StationId,
			},
			wantBody: nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				StationId: parkedCar.StationId,
			},
			wantBody: map[string]string{"error": application.ErrInvalidId.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusNotFound,
			bodyArg: params{
				StationId: parkedCar.StationId,
			},
			wantBody: map[string]string{"error": application.ErrInvalidCar.Error()},
		},
		{
			name:           "incorrect car status req",
			idArg:          transferCar.ID,
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				StationId: transferCar.StationId,
			},
			wantBody: map[string]string{"error": application.ErrInvalidTransfer.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("PUT", "/cars/"+tc.idArg+"/transfer", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/cars/{id}/transfer", carController.UpdateCarToTransfer).Methods("PUT")
			router.ServeHTTP(res, req)

			if res.Code != tc.wantStatusCode {
				t.Error("wrong response code", res.Code)
			}

			var expectedBody string
			if tc.wantBody != nil {
				json, _ := json.Marshal(tc.wantBody)
				expectedBody = strings.Trim(string(json), "\n")
			}

			if strings.Trim(res.Body.String(), "\n") != expectedBody {
				t.Error("wrong response body", strings.Trim(res.Body.String(), "\n"), expectedBody)
			}
		})
	}
}
