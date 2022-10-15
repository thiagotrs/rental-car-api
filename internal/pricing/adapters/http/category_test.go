package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/adapters/repository"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/application"
	"github.com/thiagotrs/rentalcar-ddd/internal/pricing/domain"
)

func newCategoryFixture() *domain.Category {
	p1, _ := domain.NewPolicy("Promo 1", 0.2, domain.PerKM, 500)
	p2, _ := domain.NewPolicy("Promo 2", 30.5, domain.PerDay, 5)
	c, _ := domain.NewCategory(
		"Basic",
		"basic cars",
		[]string{"UNO", "MERIVA"},
		[]domain.Policy{*p1, *p2})
	return c
}

func TestCategoryController_GetCategories(t *testing.T) {
	categories := []domain.Category{*newCategoryFixture(), *newCategoryFixture()}
	stationRepo := repository.NewCategoryRepositoryInMemory(categories)
	stationUC := application.NewCategoryUseCase(stationRepo)
	stationController := NewCategoryController(stationUC)

	req := httptest.NewRequest("GET", "/categories", nil)
	res := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/categories", stationController.GetCategories).Methods("GET")
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("wrong response code", res.Code)
	}

	json, _ := json.Marshal(categories)
	expectedBody := strings.Trim(string(json), "\n")

	if strings.Trim(res.Body.String(), "\n") != expectedBody {
		t.Error("wrong response body", strings.Trim(res.Body.String(), "\n"), expectedBody)
	}
}

func TestCategoryController_GetCategoryById(t *testing.T) {
	categories := []domain.Category{*newCategoryFixture(), *newCategoryFixture()}
	categoryRepo := repository.NewCategoryRepositoryInMemory(categories)
	categoryUC := application.NewCategoryUseCase(categoryRepo)
	categoryController := NewCategoryController(categoryUC)

	testCases := []struct {
		name           string
		idArg          string
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          categories[0].ID,
			wantStatusCode: http.StatusOK,
			wantBody:       categories[0],
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
			wantBody:       map[string]string{"error": application.ErrNotFoundCategory.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/categories/"+tc.idArg, nil)
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/categories/{id}", categoryController.GetCategoryById).Methods("GET")
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

func TestCategoryController_CreateCategory(t *testing.T) {
	categories := []domain.Category{*newCategoryFixture(), *newCategoryFixture()}
	categoryRepo := repository.NewCategoryRepositoryInMemory(categories)
	categoryUC := application.NewCategoryUseCase(categoryRepo)
	categoryController := NewCategoryController(categoryUC)

	type params struct {
		Name        string
		Description string
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
				Name:        "Basic",
				Description: "basic cars",
			},
			wantBody: nil,
		},
		{
			name:           "incorrect station id body req",
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				Name:        "",
				Description: "",
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

			req := httptest.NewRequest("POST", "/categories/", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/categories/", categoryController.CreateCategory).Methods("POST")
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

func TestCategoryController_DeleteCategory(t *testing.T) {
	categories := []domain.Category{*newCategoryFixture(), *newCategoryFixture()}
	categoryRepo := repository.NewCategoryRepositoryInMemory(categories)
	categoryUC := application.NewCategoryUseCase(categoryRepo)
	categoryController := NewCategoryController(categoryUC)

	testCases := []struct {
		name           string
		idArg          string
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          categories[0].ID,
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
			wantBody:       map[string]string{"error": application.ErrInvalidCategory.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/category/"+tc.idArg, nil)
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/category/{id}", categoryController.DeleteCategory).Methods("DELETE")
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

func TestCategoryController_UpdateAddModelInCategory(t *testing.T) {
	categories := []domain.Category{*newCategoryFixture(), *newCategoryFixture()}
	categoryRepo := repository.NewCategoryRepositoryInMemory(categories)
	categoryUC := application.NewCategoryUseCase(categoryRepo)
	categoryController := NewCategoryController(categoryUC)

	type params struct {
		ModelName string
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
			idArg:          categories[0].ID,
			wantStatusCode: http.StatusNoContent,
			bodyArg: params{
				ModelName: "PALIO",
			},
			wantBody: nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusNotFound,
			bodyArg: params{
				ModelName: "PALIO",
			},
			wantBody: map[string]string{"error": application.ErrInvalidCategory.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusNotFound,
			bodyArg: params{
				ModelName: "PALIO",
			},
			wantBody: map[string]string{"error": application.ErrInvalidCategory.Error()},
		},
		{
			name:           "incorrect car status req",
			idArg:          categories[0].ID,
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				ModelName: "UNO",
			},
			wantBody: map[string]string{"error": application.ErrInvalidModel.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("PUT", "/categories/"+tc.idArg+"/model", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/categories/{id}/model", categoryController.UpdateAddModelInCategory).Methods("PUT")
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

func TestCategoryController_UpdateDelModelInCategory(t *testing.T) {
	categories := []domain.Category{*newCategoryFixture(), *newCategoryFixture()}
	categoryRepo := repository.NewCategoryRepositoryInMemory(categories)
	categoryUC := application.NewCategoryUseCase(categoryRepo)
	categoryController := NewCategoryController(categoryUC)

	type params struct {
		ModelName string
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
			idArg:          categories[0].ID,
			wantStatusCode: http.StatusNoContent,
			bodyArg: params{
				ModelName: "UNO",
			},
			wantBody: nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusNotFound,
			bodyArg: params{
				ModelName: "UNO",
			},
			wantBody: map[string]string{"error": application.ErrInvalidCategory.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusNotFound,
			bodyArg: params{
				ModelName: "UNO",
			},
			wantBody: map[string]string{"error": application.ErrInvalidCategory.Error()},
		},
		{
			name:           "incorrect car status req",
			idArg:          categories[0].ID,
			wantStatusCode: http.StatusBadRequest,
			bodyArg: params{
				ModelName: "PALIO",
			},
			wantBody: map[string]string{"error": application.ErrInvalidModel.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("DELETE", "/categories/"+tc.idArg+"/model", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/categories/{id}/model", categoryController.UpdateDelModelInCategory).Methods("DELETE")
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

func TestCategoryController_UpdateAddPolicyInCategory(t *testing.T) {
	categories := []domain.Category{*newCategoryFixture(), *newCategoryFixture()}
	categoryRepo := repository.NewCategoryRepositoryInMemory(categories)
	categoryUC := application.NewCategoryUseCase(categoryRepo)
	categoryController := NewCategoryController(categoryUC)

	type params struct {
		Name    string
		Price   float32
		Unit    uint
		MinUnit uint
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
			idArg:          categories[0].ID,
			wantStatusCode: http.StatusNoContent,
			bodyArg:        params{"Promo 1", 0.2, uint(domain.PerKM), 500},
			wantBody:       nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			wantStatusCode: http.StatusNotFound,
			bodyArg:        params{"Promo 1", 0.2, uint(domain.PerKM), 500},
			wantBody:       map[string]string{"error": application.ErrInvalidCategory.Error()},
		},
		{
			name:           "incorrect id req",
			idArg:          "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusNotFound,
			bodyArg:        params{"Promo 1", 0.2, uint(domain.PerKM), 500},
			wantBody:       map[string]string{"error": application.ErrInvalidCategory.Error()},
		},
		{
			name:           "incorrect policy req",
			idArg:          categories[0].ID,
			wantStatusCode: http.StatusBadRequest,
			bodyArg:        params{},
			wantBody:       map[string]string{"error": application.ErrInvalidEntity.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqJSON, _ := json.Marshal(tc.bodyArg)

			req := httptest.NewRequest("PUT", "/categories/"+tc.idArg+"/policy", bytes.NewBuffer(reqJSON))
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/categories/{id}/policy", categoryController.UpdateAddPolicyInCategory).Methods("PUT")
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

func TestCategoryController_UpdateDelPolicyInCategory(t *testing.T) {
	categories := []domain.Category{*newCategoryFixture(), *newCategoryFixture()}
	categoryRepo := repository.NewCategoryRepositoryInMemory(categories)
	categoryUC := application.NewCategoryUseCase(categoryRepo)
	categoryController := NewCategoryController(categoryUC)

	testCases := []struct {
		name           string
		idArg          string
		policyIdArg    string
		wantStatusCode int
		wantBody       interface{}
	}{
		{
			name:           "correct req",
			idArg:          categories[0].ID,
			policyIdArg:    categories[0].Policies[0].ID,
			wantStatusCode: http.StatusNoContent,
			wantBody:       nil,
		},
		{
			name:           "incorrect invalid id req",
			idArg:          "invalid-id",
			policyIdArg:    categories[0].Policies[0].ID,
			wantStatusCode: http.StatusNotFound,
			wantBody:       map[string]string{"error": application.ErrInvalidCategory.Error()},
		},
		{
			name:           "incorrect invalid policy id req",
			idArg:          categories[0].ID,
			policyIdArg:    "e4ce866a-f5b7-4774-8f9d-5eb74c3900cc",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       map[string]string{"error": application.ErrInvalidPolicy.Error()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/categories/"+tc.idArg+"/policy/"+tc.policyIdArg, nil)
			res := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/categories/{id}/policy/{policyId}", categoryController.UpdateDelPolicyInCategory).Methods("DELETE")
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
