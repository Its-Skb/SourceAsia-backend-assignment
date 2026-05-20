package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-assignment/internal/models"
	"backend-assignment/internal/routes"
	"backend-assignment/internal/storage"

	"github.com/gin-gonic/gin"
)

func TestCreateProduct_DuplicateSKU(t *testing.T) {

	// Reset test storage
	storage.Products = make(map[string]*models.Product)
	storage.ProductMediaStore = make(map[string]*models.ProductMedia)
	storage.SKUIndex = make(map[string]string)
	storage.ProductOrder = []string{}

	router := gin.Default()

	routes.RegisterRoutes(router)

	requestBody := []byte(`{
		"name":"Test Product",
		"sku":"TEST-SKU-001"
	}`)

	// First request
	req1, _ := http.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBuffer(requestBody),
	)

	req1.Header.Set("Content-Type", "application/json")

	resp1 := httptest.NewRecorder()

	router.ServeHTTP(resp1, req1)

	if resp1.Code != http.StatusCreated {
		t.Errorf("expected first request to succeed")
	}

	// Second request with same SKU
	req2, _ := http.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBuffer(requestBody),
	)

	req2.Header.Set("Content-Type", "application/json")

	resp2 := httptest.NewRecorder()

	router.ServeHTTP(resp2, req2)

	if resp2.Code != http.StatusConflict {
		t.Errorf("expected duplicate SKU conflict")
	}
}
