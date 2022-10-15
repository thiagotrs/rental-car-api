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

func newStationFixture() *domain.Station {
	s, _ := domain.NewStation("Station 1", "Farway Av.", "45, Ap. 50", "Polar", "Nort City", "20778990", 100, 0)
	return s
}

func TestStationController_GetStations(t *testing.T) {
	stations := []domain.Station{*newStationFixture(), *newStationFixture()}

	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationController := NewStationController(stationUC)

	req := httptest.NewRequest("GET", "/stations", nil)
	res := httptest.NewRecorder()

	stationController.GetStations(res, req)

	if res.Code != http.StatusOK {
		t.Error("wrong response code", res.Code)
	}

	json, _ := json.Marshal(stations)
	expectedBody := strings.Trim(string(json), "\n")

	if strings.Trim(res.Body.String(), "\n") != expectedBody {
		t.Error("wrong response body", strings.Trim(res.Body.String(), "\n"))
	}
}

func TestStationController_GetStationById(t *testing.T) {
	stations := []domain.Station{*newStationFixture(), *newStationFixture()}
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationController := NewStationController(stationUC)

	testCases := []struct {
		name           string
		idArg          string
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          stations[0].ID,
			wantStatusCode: http.StatusOK,
			wantBody:       stations[0],
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
			wantBody:       map[string]string{"error": application.ErrNotFoundStation.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/stations/"+tc.idArg, nil)
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/stations/{id}", stationController.GetStationById).Methods("GET")
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

func TestStationController_CreateStation(t *testing.T) {
	stations := []domain.Station{*newStationFixture()}
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationController := NewStationController(stationUC)

	type params struct {
		Name, Address, Complement, State, City, Cep string
		Capacity, Idle                              uint
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
				Name:       "Station 1",
				Address:    "Farway Av.",
				Complement: "45, Ap. 50",
				State:      "Polar",
				City:       "Nort City",
				Cep:        "20778990",
				Capacity:   100,
				Idle:       0,
			},
			wantBody: nil,
		},
		{
			name:           "incorrect capacity req",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				Name:       "Station 1",
				Address:    "Farway Av.",
				Complement: "45, Ap. 50",
				State:      "Polar",
				City:       "Nort City",
				Cep:        "20778990",
				Capacity:   0,
				Idle:       0,
			},
			wantBody: map[string]string{"error": application.ErrInvalidEntity.Error()},
		},
		{
			name:           "incorrect malformed req",
			wantStatusCode: http.StatusBadRequest,
			bodyArg:        "",
			wantBody:       map[string]string{"error": application.ErrInvalidEntity.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("POST", "/stations/", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/stations/", stationController.CreateStation).Methods("POST")
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

func TestStationController_DeleteStation(t *testing.T) {
	noZeroIdleStation := *newStationFixture()
	noZeroIdleStation.Idle = 5
	stations := []domain.Station{*newStationFixture(), noZeroIdleStation}
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationController := NewStationController(stationUC)

	testCases := []struct {
		name           string
		idArg          string
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          stations[0].ID,
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
			wantBody:       map[string]string{"error": application.ErrInvalidStation.Error()},
		},
		{
			name:           "incorrect station req",
			idArg:          stations[1].ID,
			wantStatusCode: http.StatusNotFound,
			wantBody:       map[string]string{"error": application.ErrStationHasCars.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/stations/"+tc.idArg, nil)
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/stations/{id}", stationController.DeleteStation).Methods("DELETE")
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

func TestStationController_UpdateStationCapacity(t *testing.T) {
	noZeroIdleStation := *newStationFixture()
	noZeroIdleStation.Idle = 5
	stations := []domain.Station{*newStationFixture(), noZeroIdleStation}
	stationRepo := repository.NewStationRepositoryInMemory(stations)
	stationUC := application.NewStationUseCase(stationRepo)
	stationController := NewStationController(stationUC)

	testCases := []struct {
		name           string
		idArg          string
		bodyArg        interface{}
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          stations[0].ID,
			wantStatusCode: http.StatusNoContent,
			bodyArg:        map[string]uint{"capacity": 150},
			wantBody:       nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusBadRequest,
			bodyArg:        map[string]uint{"capacity": 150},
			wantBody:       map[string]string{"error": application.ErrInvalidId.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusNotFound,
			bodyArg:        map[string]uint{"capacity": 150},
			wantBody:       map[string]string{"error": application.ErrInvalidStation.Error()},
		},
		{
			name:           "incorrect capacity req",
			idArg:          stations[0].ID,
			wantStatusCode: http.StatusBadRequest,
			bodyArg:        map[string]uint{"capacity": 0},
			wantBody:       map[string]string{"error": application.ErrInvalidCapacity.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("PUT", "/stations/"+tc.idArg+"/capacity", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/stations/{id}/capacity", stationController.UpdateStationCapacity).Methods("PUT")
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
