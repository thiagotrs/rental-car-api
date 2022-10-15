package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/application"
)

type categoryController struct {
	categoryUC application.CategoryUseCase
}

func NewCategoryController(categoryUC application.CategoryUseCase) *categoryController {
	return &categoryController{categoryUC}
}

func (c *categoryController) GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	categories := c.categoryUC.GetCategories()
	w.WriteHeader(http.StatusOK)
	json, _ := json.Marshal(categories)
	w.Write(json)
}

func (c *categoryController) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	category, err := c.categoryUC.GetCategoryById(vars["id"])

	switch err {
	case application.ErrInvalidId:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrNotFoundCategory:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		json, _ := json.Marshal(category)
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *categoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var params struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, application.ErrInvalidEntity)
		return
	}
	err := c.categoryUC.AddCategory(params.Name, params.Description)
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

func (c *categoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	err := c.categoryUC.DeleteCategory(vars["id"])
	switch err {
	case application.ErrInvalidId:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidCategory:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *categoryController) UpdateAddModelInCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	var params struct {
		ModelName string `json:"modelName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, application.ErrInvalidEntity)
		return
	}
	err := c.categoryUC.AddModelInCategory(vars["id"], params.ModelName)
	switch err {
	case application.ErrInvalidId, application.ErrInvalidModel:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidCategory:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *categoryController) UpdateDelModelInCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	var params struct {
		ModelName string `json:"modelName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, application.ErrInvalidEntity)
		return
	}
	err := c.categoryUC.DeleteModelInCategory(vars["id"], params.ModelName)
	switch err {
	case application.ErrInvalidId, application.ErrInvalidModel:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidCategory:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *categoryController) UpdateAddPolicyInCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	var params struct {
		Name    string  `json:"name"`
		Price   float32 `json:"price"`
		Unit    uint    `json:"unit"`
		MinUnit uint    `json:"minUnit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, application.ErrInvalidEntity)
		return
	}
	err := c.categoryUC.AddPolicyInCategory(vars["id"], params.Name, params.Price, params.Unit, params.MinUnit)
	switch err {
	case application.ErrInvalidId, application.ErrInvalidEntity, application.ErrInvalidPolicy:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidCategory:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}

func (c *categoryController) UpdateDelPolicyInCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	err := c.categoryUC.DeletePolicyInCategory(vars["id"], vars["policyId"])
	switch err {
	case application.ErrInvalidId, application.ErrInvalidPolicy:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case application.ErrInvalidCategory:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"%v"}`, err)
	}
}
