package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
)

type carController struct {
	carUC application.CarUseCase
}

func NewCarController(carUC application.CarUseCase) *carController {
	return &carController{carUC}
}

func (c *carController) SearchCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var params application.SearchCarParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{error: %v}", application.ErrInvalidEntity)
		return
	}
	cars := c.carUC.SearchCars(params)

	w.WriteHeader(http.StatusOK)
	json, err := json.Marshal(cars)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{error: %v}", err)
		return
	}
	w.Write(json)
}

func (c *carController) GetCarById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	car, err := c.carUC.GetCarById(vars["id"])

	switch err {
	case application.ErrInvalidId:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrNotFoundCar:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		json, _ := json.Marshal(car)
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *carController) CreateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var params struct {
		Age       uint16 `json:"age"`
		KM        uint64 `json:"km"`
		Plate     string `json:"plate"`
		Document  string `json:"document"`
		StationId string `json:"stationId"`
		Model     string `json:"model"`
		Make      string `json:"make"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, application.ErrInvalidEntity)
		return
	}
	err := c.carUC.AddCar(params.Age, params.KM, params.Plate, params.Document, params.StationId, params.Model, params.Make)
	switch err {
	case application.ErrInvalidEntity, application.ErrStationMaxCapacity:
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

func (c *carController) DeleteCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	err := c.carUC.DeleteCar(vars["id"])

	switch err {
	case application.ErrInvalidId, application.ErrCarNotInMaintenance:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidCar:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *carController) UpdateCarToMaintenance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	var params struct {
		StationId string `json:"stationId"`
		KM        uint64 `json:"km"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}
	err := c.carUC.MoveCarToMaintenance(vars["id"], params.StationId, params.KM)
	switch err {
	case application.ErrInvalidId, application.ErrInvalidMaintenance:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidCar:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *carController) UpdateCarToPark(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	var params struct {
		StationId string `json:"stationId"`
		KM        uint64 `json:"km"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}
	err := c.carUC.ParkCar(vars["id"], params.StationId, params.KM)
	switch err {
	case application.ErrInvalidId, application.ErrStationMaxCapacity, application.ErrInvalidEntity, application.ErrInvalidPark:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidCar:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *carController) UpdateCarToTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	var params struct {
		StationId string `json:"stationId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}
	err := c.carUC.TransferCar(vars["id"], params.StationId)
	switch err {
	case application.ErrInvalidId, application.ErrStationMaxCapacity, application.ErrInvalidEntity, application.ErrInvalidTransfer:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidCar:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}
