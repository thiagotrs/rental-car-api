package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
)

type stationController struct {
	stationUC application.StationUseCase
}

func NewStationController(stationUC application.StationUseCase) *stationController {
	return &stationController{stationUC}
}

func (c *stationController) GetStations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	stations := c.stationUC.GetStations()
	w.WriteHeader(http.StatusOK)
	json, _ := json.Marshal(stations)
	w.Write(json)
}

func (c *stationController) GetStationById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	station, err := c.stationUC.GetStationById(vars["id"])

	switch err {
	case application.ErrInvalidId:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrNotFoundStation:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		json, _ := json.Marshal(station)
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *stationController) CreateStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var params struct {
		Name       string `json:"name"`
		Address    string `json:"address"`
		Complement string `json:"complement"`
		State      string `json:"state"`
		City       string `json:"city"`
		Cep        string `json:"cep"`
		Capacity   uint   `json:"capacity"`
		Idle       uint   `json:"idle"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, application.ErrInvalidEntity)
		return
	}
	err := c.stationUC.AddStation(params.Name, params.Address, params.Complement, params.State, params.City, params.Cep, params.Capacity, params.Idle)
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

func (c *stationController) DeleteStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	err := c.stationUC.DeleteStation(vars["id"])
	switch err {
	case application.ErrInvalidId:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidStation, application.ErrStationHasCars:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *stationController) UpdateStationCapacity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	var params struct {
		Capacity uint `json:"capacity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
		return
	}
	err := c.stationUC.ChangeStationCapacity(vars["id"], params.Capacity)
	switch err {
	case application.ErrInvalidId, application.ErrInvalidCapacity:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidStation:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}
