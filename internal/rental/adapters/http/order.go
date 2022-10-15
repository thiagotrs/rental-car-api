package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/thiagotrs/rentalcar-ddd/internal/rental/application"
)

type orderController struct {
	orderUC application.OrderUseCase
}

func NewOrderController(orderUC application.OrderUseCase) *orderController {
	return &orderController{orderUC}
}

func (c *orderController) GetOrderById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	order, err := c.orderUC.GetById(vars["id"])

	switch err {
	case application.ErrInvalidId:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrNotFoundOrder:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		json, _ := json.Marshal(order)
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *orderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var params struct {
		DateReservFrom time.Time `json:"dateReservFrom"`
		DateReservTo   time.Time `json:"dateReservTo"`
		StationFromId  string    `json:"stationFromId"`
		StationToId    string    `json:"stationToId"`
		CategoryId     string    `json:"categoryId"`
		CarModel       string    `json:"carModel"`
		PolicyId       string    `json:"policyId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, application.ErrInvalidEntity)
		return
	}
	err := c.orderUC.Open(
		params.DateReservFrom, params.DateReservTo, params.StationFromId,
		params.StationToId, params.CategoryId, params.CarModel, params.PolicyId)

	switch err {
	case application.ErrInvalidEntity:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	case nil:
		w.WriteHeader(http.StatusCreated)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *orderController) UpdateToComfirmOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	var params struct {
		DateFrom time.Time `json:"dateFrom"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{error: %v}", err)
		return
	}
	err := c.orderUC.Confirm(vars["id"], params.DateFrom)

	switch err {
	case application.ErrInvalidOrder:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *orderController) UpdateToCloseOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	var params struct {
		Discount float32   `json:"discount"`
		Tax      float32   `json:"tax"`
		DateTo   time.Time `json:"dateTo"`
		KM       uint64    `json:"km"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{error: %v}", err)
		return
	}
	err := c.orderUC.Close(vars["id"], params.Discount, params.Tax, params.DateTo, params.KM)
	switch err {
	case application.ErrInvalidOrder:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *orderController) UpdateToCancelOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	err := c.orderUC.Cancel(vars["id"])
	switch err {
	case application.ErrInvalidOrder:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}
