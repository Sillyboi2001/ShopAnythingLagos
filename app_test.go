package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
  "encoding/json"

	"github.com/gin-gonic/gin"
  "github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	registerRoutes(router)
	return router
}

func registerRoutes(router *gin.Engine) {
	router.POST("/products", createProduct)
	router.GET("/products", getProducts)
  router.GET("/products/:id", getProductById)
  router.DELETE("/products/:id", deleteProduct)
  router.PATCH("/products/:id", updateProduct)
}

// TestCreateProduct tests the createProduct endpoint
func TestCreateProduct(t *testing.T) {
	router := SetupRouter()

  // Define a product JSON payload
  var jsonStr = []byte(`{"name":"Test Product","description":"This is a test","price":100}`)
  var jsonStr1 = []byte(`{"name":"Test Product","description":"This is a test","price":"100"}`)

	// Create a new HTTP request with the product payload
	req, _ := http.NewRequest("POST", "/products?merchantId=silas", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

  // Testing for errors
  // Testing for invalid inputs
  req1, _ := http.NewRequest("POST", "/products?merchantId=silas", bytes.NewBuffer(jsonStr1))
	req.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w, req1)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w1.Code)
	}

  // Testing for missing merchant ID
  wReq := httptest.NewRecorder()
  reqGet2, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonStr))
  router.ServeHTTP(wReq, reqGet2)
  assert.Equal(t, http.StatusBadRequest, wReq.Code)
  assert.Contains(t, wReq.Body.String(), "Merchant ID is required")
}

// TestGetProduct tests the getProducts endpoint
func TestGetProduct(t *testing.T) {
	router := SetupRouter()

  // Define a product JSON payload
  var jsonStr = []byte(`{"name":"Test Product","description":"This is a test","price":100}, {"name":"Test Product1","description":"This is a test1","price":100},`)

	// Create a new HTTP request with the product payload
	req, _ := http.NewRequest("POST", "/products?merchantId=silas", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

  // Get the product
  wGet := httptest.NewRecorder()
  reqGet, _ := http.NewRequest("GET", "/products?merchantId=silas", nil)
  router.ServeHTTP(wGet, reqGet)
  assert.Equal(t, http.StatusOK, wGet.Code)
  assert.Contains(t, wGet.Body.String(), "Test Product") 

  // Empty product
  wGet1 := httptest.NewRecorder()
  reqGet1, _ := http.NewRequest("GET", "/products?merchantId=sil", nil)
  router.ServeHTTP(wGet1, reqGet1)
  assert.Equal(t, http.StatusNotFound, wGet1.Code)
  assert.Contains(t, wGet1.Body.String(), "No products found for this merchant")

  // No merchant ID
  wGet2 := httptest.NewRecorder()
  reqGet2, _ := http.NewRequest("GET", "/products", nil)
  router.ServeHTTP(wGet2, reqGet2)
  assert.Equal(t, http.StatusBadRequest, wGet2.Code)
  assert.Contains(t, wGet2.Body.String(), "Merchant ID is required")
}

// TestGetProductByID tests the getProductById endpoint
func TestGetProductByID(t *testing.T) {
  router := SetupRouter()

  // Create the product
  var jsonStr = []byte(`{"name":"Test Product","description":"This is a test","price":100}`)
  reqCreate, _ := http.NewRequest("POST", "/products?merchantId=silas", bytes.NewBuffer(jsonStr))
  reqCreate.Header.Set("Content-Type", "application/json")
  wCreate := httptest.NewRecorder()
  router.ServeHTTP(wCreate, reqCreate)
  assert.Equal(t, http.StatusCreated, wCreate.Code)

  // Extract the product ID
  var createdProduct map[string]interface{}
  err := json.Unmarshal(wCreate.Body.Bytes(), &createdProduct)
  if err != nil {
      t.Fatalf("Failed to unmarshal response: %v", err)
  }
  productID, ok := createdProduct["id"].(string)
  if !ok {
      t.Fatalf("Product ID not found or invalid")
  }

  // Get the product by ID
  wGetByID := httptest.NewRecorder()
  reqGetByID, _ := http.NewRequest("GET", "/products/"+productID+"?merchantId=silas", nil)
  router.ServeHTTP(wGetByID, reqGetByID)
  assert.Equal(t, http.StatusOK, wGetByID.Code)
  assert.Contains(t, wGetByID.Body.String(), "Test Product")

  // Testing for unauthorized access
  reqGet, _ := http.NewRequest("GET", "/products/"+productID+"?merchantId=si", nil)
  wGet := httptest.NewRecorder()
  router.ServeHTTP(wGet, reqGet)
  assert.Equal(t, http.StatusUnauthorized, wGet.Code)
  assert.Contains(t, wGet.Body.String(), "Unauthorized")

  // Testing for wrong id
  reqGet1, _ := http.NewRequest("GET", "/products/"+productID+"abc?merchantId=silas", nil)
  wGet1 := httptest.NewRecorder()
  router.ServeHTTP(wGet1, reqGet1)
  assert.Equal(t, http.StatusNotFound, wGet1.Code)
  assert.Contains(t, wGet1.Body.String(), "Product not found")
}

// TestDeleteProductByID tests the deleteProduct endpoint
func TestDeleteProductByID(t *testing.T) {
  router := SetupRouter()

  // Create the product
  var jsonStr = []byte(`{"name":"Test Product","description":"This is a test","price":100}`)
  reqCreate, _ := http.NewRequest("POST", "/products?merchantId=silas", bytes.NewBuffer(jsonStr))
  reqCreate.Header.Set("Content-Type", "application/json")
  wCreate := httptest.NewRecorder()
  router.ServeHTTP(wCreate, reqCreate)
  assert.Equal(t, http.StatusCreated, wCreate.Code)

  // Extract the product ID
  var createdProduct map[string]interface{}
  err := json.Unmarshal(wCreate.Body.Bytes(), &createdProduct)
  if err != nil {
      t.Fatalf("Failed to unmarshal response: %v", err)
  }
  productID, ok := createdProduct["id"].(string)
  if !ok {
      t.Fatalf("Product ID not found or invalid")
  }

  // Testing for unauthorized access
  reqDel1, _ := http.NewRequest("DELETE", "/products/"+productID+"?merchantId=si", nil)
  wDel1 := httptest.NewRecorder()
  router.ServeHTTP(wDel1, reqDel1)
  assert.Equal(t, http.StatusUnauthorized, wDel1.Code)
  assert.Contains(t, wDel1.Body.String(), "Unauthorized")

  // Testing for wrong id
  reqDel, _ := http.NewRequest("DELETE", "/products/"+productID+"abc?merchantId=silas", nil)
  wDel := httptest.NewRecorder()
  router.ServeHTTP(wDel, reqDel)
  assert.Equal(t, http.StatusNotFound, wDel.Code)
  assert.Contains(t, wDel.Body.String(), "Product not found")

  // Delete the product by ID
  wGetByID := httptest.NewRecorder()
  reqGetByID, _ := http.NewRequest("DELETE", "/products/"+productID+"?merchantId=silas", nil)
  router.ServeHTTP(wGetByID, reqGetByID)
  assert.Equal(t, http.StatusOK, wGetByID.Code)
  assert.Contains(t, wGetByID.Body.String(), "Product deleted")
}

// TestUpdateProductByID tests the updateProduct endpoint
func TestUpdateProductByID(t *testing.T) {
  router := SetupRouter()

  // Step 1: Create the product
  var jsonStr = []byte(`{"name":"Test Product","description":"This is a test","price":100}`)
  reqCreate, _ := http.NewRequest("POST", "/products?merchantId=silas", bytes.NewBuffer(jsonStr))
  reqCreate.Header.Set("Content-Type", "application/json")
  wCreate := httptest.NewRecorder()
  router.ServeHTTP(wCreate, reqCreate)
  assert.Equal(t, http.StatusCreated, wCreate.Code)

  // Extract the product ID
  var createdProduct map[string]interface{}
  err := json.Unmarshal(wCreate.Body.Bytes(), &createdProduct)
  if err != nil {
      t.Fatalf("Failed to unmarshal response: %v", err)
  }
  productID, ok := createdProduct["id"].(string)
  if !ok {
      t.Fatalf("Product ID not found or invalid")
  }

  // Update the product
  var updateJsonStr = []byte(`{"name":"Updated Product","description":"Updated Description","price":150}`)
  reqUpdate, _ := http.NewRequest("PATCH", "/products/"+productID+"?merchantId=silas", bytes.NewBuffer(updateJsonStr))
  reqUpdate.Header.Set("Content-Type", "application/json")

  wUpdate := httptest.NewRecorder()
  router.ServeHTTP(wUpdate, reqUpdate)
  assert.Equal(t, http.StatusOK, wUpdate.Code)

  // Verify the update
  var updatedProduct Product
	err = json.Unmarshal(wUpdate.Body.Bytes(), &updatedProduct)
	if err != nil {
		t.Fatalf("Failed to unmarshal response from update: %v", err)
	}
	assert.Equal(t, "Updated Product", updatedProduct.NAME, "Product name did not update correctly")
	assert.Equal(t, "Updated Description", updatedProduct.DESCRIPTION, "Product description did not update correctly")
	assert.Equal(t, 150, updatedProduct.PRICE, "Product price did not update correctly")

  // Checking for errors
  // Testing for bad request
  var updateJsonStr1 = []byte(`{"name":"Updated Product","description": "Updated Description","price":"150"}`)
  reqUpdate1, _ := http.NewRequest("PATCH", "/products/"+productID+"?merchantId=silas", bytes.NewBuffer(updateJsonStr1))
  reqUpdate.Header.Set("Content-Type", "application/json")

  wUpdate1 := httptest.NewRecorder()
  router.ServeHTTP(wUpdate1, reqUpdate1)
  assert.Equal(t, http.StatusBadRequest, wUpdate1.Code)
  assert.Contains(t, wUpdate1.Body.String(), "Bad request")

  // Testing for unauthorized access
  var updateJsonStr2 = []byte(`{"name":"Updated Product","description": "Updated Description","price":150}`)
  reqUpdate2, _ := http.NewRequest("PATCH", "/products/"+productID, bytes.NewBuffer(updateJsonStr2))
  reqUpdate.Header.Set("Content-Type", "application/json")

  wUpdate2 := httptest.NewRecorder()
  router.ServeHTTP(wUpdate2, reqUpdate2)
  assert.Equal(t, http.StatusUnauthorized, wUpdate2.Code)
  assert.Contains(t, wUpdate2.Body.String(), "Unauthorized")

  // Testing for wrong id
  var updateJsonStr3 = []byte(`{"name":"Updated Product","description": "Updated Description","price":150}`)
  reqUpdate3, _ := http.NewRequest("PATCH", "/products/"+productID+"abc?merchantId=silas", bytes.NewBuffer(updateJsonStr3))
  reqUpdate.Header.Set("Content-Type", "application/json")

  wUpdate3 := httptest.NewRecorder()
  router.ServeHTTP(wUpdate3, reqUpdate3)
  assert.Equal(t, http.StatusNotFound, wUpdate3.Code)
  assert.Contains(t, wUpdate3.Body.String(), "Product not found")
}

