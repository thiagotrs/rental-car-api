package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/adapters/repository"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/application"
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
		Status:    domain.Parked,
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
	o, _ := domain.NewOrder(
		time.Now(),
		time.Now().Add(time.Hour*24*5),
		*newCarFixture(),
		"83369771-f9a4-48b7-b87b-463f19f7b187",
		"2520aade-a397-4e3c-a589-39c6ae5c2eff",
		*newPolicyFixture(),
	)
	return o
}

type orderOrderServiceMock struct {
	expectedGetPolicy    *domain.Policy
	expectedGetPolicyErr error
	expectedGetCar       *domain.Car
	expectedGetCarErr    error
}

func (m *orderOrderServiceMock) GetPolicy(categoryId, modelId, policyId string) (*domain.Policy, error) {
	return m.expectedGetPolicy, m.expectedGetPolicyErr
}

func (m *orderOrderServiceMock) GetCar(stationId, modelId string) (*domain.Car, error) {
	return m.expectedGetCar, m.expectedGetCarErr
}

func TestOrderController_GetOrderById(t *testing.T) {
	orders := []domain.Order{*newOrderFixture(), *newOrderFixture()}
	orderRepo := repository.NewOrderRepositoryInMemory(orders)
	orderSvc := &orderOrderServiceMock{}
	orderUC := application.NewOrderUseCase(orderRepo, orderSvc)
	orderController := NewOrderController(orderUC)

	testCases := []struct {
		name           string
		idArg          string
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          orders[0].ID,
			wantStatusCode: http.StatusOK,
			wantBody:       orders[0],
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
			wantBody:       map[string]string{"error": application.ErrNotFoundOrder.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/orders/"+tc.idArg, nil)
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/orders/{id}", orderController.GetOrderById).Methods("GET")
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

func TestOrderController_CreateOrder(t *testing.T) {
	newOrder := newOrderFixture()
	orders := []domain.Order{*newOrderFixture(), *newOrderFixture()}
	orderRepo := repository.NewOrderRepositoryInMemory(orders)
	orderSvc := &orderOrderServiceMock{
		expectedGetPolicy: newPolicyFixture(),
		expectedGetCar:    newCarFixture(),
	}
	orderUC := application.NewOrderUseCase(orderRepo, orderSvc)
	orderController := NewOrderController(orderUC)

	type params struct {
		DateReservFrom time.Time
		DateReservTo   time.Time
		StationFromId  string
		StationToId    string
		CategoryId     string
		CarModel       string
		PolicyId       string
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
				DateReservFrom: newOrder.DateReservFrom,
				DateReservTo:   newOrder.DateReservTo,
				StationFromId:  newOrder.StationFromId,
				StationToId:    newOrder.StationToId,
				CategoryId:     newOrder.Policy.CategoryId,
				CarModel:       newOrder.Policy.CarModel,
				PolicyId:       newOrder.Policy.ID,
			},
			wantBody: nil,
		},
		{
			name:           "incorrect station id body req",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				DateReservFrom: newOrder.DateReservFrom,
				DateReservTo:   newOrder.DateReservFrom.Add(time.Hour * -5),
				StationFromId:  newOrder.StationFromId,
				StationToId:    newOrder.StationToId,
				CategoryId:     newOrder.Policy.CategoryId,
				CarModel:       newOrder.Policy.CarModel,
				PolicyId:       newOrder.Policy.ID,
			},
			wantBody: map[string]string{"error": application.ErrInvalidEntity.Error()},
		},
		{
			name:           "incorrect malformed body req",
			wantStatusCode: http.StatusBadRequest,
			bodyArg:        "",
			wantBody:       map[string]string{"error": application.ErrInvalidEntity.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)
			t.Log(string(reqJSON))

			req := httptest.NewRequest("POST", "/orders/", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/orders/", orderController.CreateOrder).Methods("POST")
			router.ServeHTTP(res, req)

			if res.Code != tc.wantStatusCode {
				t.Error("wrong response code", res.Code, tc.wantStatusCode)
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

func TestOrderController_UpdateToComfirmOrder(t *testing.T) {
	orders := []domain.Order{*newOrderFixture(), *newOrderFixture()}
	orderRepo := repository.NewOrderRepositoryInMemory(orders)
	orderSvc := &orderOrderServiceMock{
		expectedGetPolicy: newPolicyFixture(),
		expectedGetCar:    newCarFixture(),
	}
	orderUC := application.NewOrderUseCase(orderRepo, orderSvc)
	orderController := NewOrderController(orderUC)

	type params struct {
		DateFrom time.Time
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
			idArg:          orders[0].ID,
			wantStatusCode: http.StatusNoContent,
			bodyArg: params{
				DateFrom: time.Now().Add(time.Hour),
			},
			wantBody: nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				DateFrom: time.Now().Add(time.Hour),
			},
			wantBody: map[string]string{"error": application.ErrInvalidOrder.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				DateFrom: time.Now().Add(time.Hour),
			},
			wantBody: map[string]string{"error": application.ErrInvalidOrder.Error()},
		},
		{
			name:           "incorrect date from body",
			idArg:          orders[0].ID,
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				DateFrom: time.Now().Add(time.Hour * -10),
			},
			wantBody: map[string]string{"error": application.ErrInvalidOrder.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("PUT", "/orders/"+tc.idArg+"/confirm", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/orders/{id}/confirm", orderController.UpdateToComfirmOrder).Methods("PUT")
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

func TestOrderController_UpdateToCloseOrder(t *testing.T) {
	orders := []domain.Order{*newOrderFixture(), *newOrderFixture()}
	orders[0].Status = domain.Confirmed
	dateFrom := time.Now().Add(time.Hour)
	orders[0].DateFrom = &dateFrom
	orderRepo := repository.NewOrderRepositoryInMemory(orders)
	orderSvc := &orderOrderServiceMock{
		expectedGetPolicy: newPolicyFixture(),
		expectedGetCar:    newCarFixture(),
	}
	orderUC := application.NewOrderUseCase(orderRepo, orderSvc)
	orderController := NewOrderController(orderUC)

	type params struct {
		Discount float32
		Tax      float32
		DateTo   time.Time
		KM       uint64
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
			idArg:          orders[0].ID,
			wantStatusCode: http.StatusNoContent,
			bodyArg: params{
				Discount: 10,
				Tax:      1.5,
				DateTo:   orders[0].DateReservTo.Add(time.Hour),
				KM:       orders[0].Car.InitialKM + 50,
			},
			wantBody: nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				Discount: 10,
				Tax:      1.5,
				DateTo:   orders[0].DateReservTo.Add(time.Hour),
				KM:       orders[0].Car.InitialKM + 50,
			},
			wantBody: map[string]string{"error": application.ErrInvalidOrder.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				Discount: 10,
				Tax:      1.5,
				DateTo:   orders[0].DateReservTo.Add(time.Hour),
				KM:       orders[0].Car.InitialKM + 50,
			},
			wantBody: map[string]string{"error": application.ErrInvalidOrder.Error()},
		},
		{
			name:           "incorrect km body",
			idArg:          orders[0].ID,
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				Discount: 10,
				Tax:      1.5,
				DateTo:   orders[0].DateReservTo.Add(time.Hour),
				KM:       orders[0].Car.InitialKM - 50,
			},
			wantBody: map[string]string{"error": application.ErrInvalidOrder.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("PUT", "/orders/"+tc.idArg+"/close", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/orders/{id}/close", orderController.UpdateToCloseOrder).Methods("PUT")
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

func TestOrderController_UpdateToCancelOrder(t *testing.T) {
	orders := []domain.Order{*newOrderFixture(), *newOrderFixture()}
	orderRepo := repository.NewOrderRepositoryInMemory(orders)
	orderSvc := &orderOrderServiceMock{
		expectedGetPolicy: newPolicyFixture(),
		expectedGetCar:    newCarFixture(),
	}
	orderUC := application.NewOrderUseCase(orderRepo, orderSvc)
	orderController := NewOrderController(orderUC)

	testCases := []struct {
		name           string
		idArg          string
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          orders[0].ID,
			wantStatusCode: http.StatusNoContent,
			wantBody:       nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       map[string]string{"error": application.ErrInvalidOrder.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       map[string]string{"error": application.ErrInvalidOrder.Error()},
		},
		{
			name:           "incorrect km body",
			idArg:          orders[0].ID,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       map[string]string{"error": application.ErrInvalidOrder.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/orders/"+tc.idArg+"/cancel", nil)
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/orders/{id}/cancel", orderController.UpdateToCancelOrder).Methods("PUT")
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
