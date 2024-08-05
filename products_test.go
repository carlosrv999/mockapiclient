// products_test.go

package mockclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestProductUnmarshalJSON tests the custom UnmarshalJSON function for the Product struct
func TestProductUnmarshalJSON(t *testing.T) {
	jsonData := `{"id": "1", "name": "Test Product", "price": "19.99", "stock": "100", "createdAt": "2024-08-05T10:00:00Z", "type": "Gadget", "department": "Electronics"}`
	var product Product
	err := json.Unmarshal([]byte(jsonData), &product)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if product.ID != 1 || product.Price != 19.99 || product.Stock != 100 {
		t.Errorf("expected product with ID 1, Price 19.99, Stock 100; got ID %d, Price %f, Stock %d", product.ID, product.Price, product.Stock)
	}
}

// TestProductMarshalJSON tests the custom MarshalJSON function for the Product struct
func TestProductMarshalJSON(t *testing.T) {
	product := Product{
		ID:         1,
		Name:       "Test Product",
		Price:      19.99,
		Stock:      100,
		CreatedAt:  "2024-08-05T10:00:00Z",
		Type:       "Gadget",
		Department: "Electronics",
	}

	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("expected no error unmarshalling JSON, got %v", err)
	}

	expected := map[string]interface{}{
		"id":         "1",
		"name":       "Test Product",
		"price":      "19.990000",
		"stock":      "100",
		"createdAt":  "2024-08-05T10:00:00Z",
		"type":       "Gadget",
		"department": "Electronics",
	}

	if !equalMaps(result, expected) {
		t.Errorf("expected JSON %v, got %v", expected, result)
	}
}

// Helper function to compare two maps
func equalMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if vb, ok := b[k]; !ok || vb != v {
			return false
		}
	}
	return true
}

// TestGetProducts tests the GetProducts method
func TestGetProducts(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products" {
			t.Errorf("expected request to /products, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": "1", "name": "Product 1", "price": "10.00", "stock": "50"}, {"id": "2", "name": "Product 2", "price": "20.00", "stock": "30"}]`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	products, err := client.GetProducts()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
	if products[0].Name != "Product 1" || products[1].Name != "Product 2" {
		t.Errorf("expected product names 'Product 1' and 'Product 2', got %s and %s", products[0].Name, products[1].Name)
	}
}

// TestGetProductByID tests the GetProductByID method
func TestGetProductByID(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products/1" {
			t.Errorf("expected request to /products/1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "1", "name": "Product 1", "price": "10.00", "stock": "50"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	product, err := client.GetProductByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.Name != "Product 1" || product.Price != 10.00 {
		t.Errorf("expected product with name 'Product 1' and price 10.00, got %s and %f", product.Name, product.Price)
	}
}

// TestCreateProduct tests the CreateProduct method
func TestCreateProduct(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/products" {
			t.Errorf("expected request to /products, got %s", r.URL.Path)
		}

		var prod Product
		err := json.NewDecoder(r.Body).Decode(&prod)
		if err != nil {
			t.Fatalf("expected no error decoding request body, got %v", err)
		}

		if prod.Name != "New Product" {
			t.Errorf("expected product name 'New Product', got %s", prod.Name)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "3", "name": "New Product", "price": "15.00", "stock": "100"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	newProduct := &Product{Name: "New Product", Price: 15.00, Stock: 100}
	createdProduct, err := client.CreateProduct(newProduct)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdProduct.ID != 3 || createdProduct.Name != "New Product" {
		t.Errorf("expected created product with ID 3 and name 'New Product', got ID %d and name %s", createdProduct.ID, createdProduct.Name)
	}
}

// TestUpdateProduct tests the UpdateProduct method
func TestUpdateProduct(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH method, got %s", r.Method)
		}
		if r.URL.Path != "/products/1" {
			t.Errorf("expected request to /products/1, got %s", r.URL.Path)
		}

		var prod Product
		err := json.NewDecoder(r.Body).Decode(&prod)
		if err != nil {
			t.Fatalf("expected no error decoding request body, got %v", err)
		}

		if prod.ID != 1 || prod.Name != "Updated Product" {
			t.Errorf("expected product with ID 1 and name 'Updated Product', got ID %d and name %s", prod.ID, prod.Name)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "1", "name": "Updated Product", "price": "25.00", "stock": "75"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	updatedProduct := &Product{ID: 1, Name: "Updated Product", Price: 25.00, Stock: 75}
	product, err := client.UpdateProduct(updatedProduct)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.Price != 25.00 || product.Stock != 75 {
		t.Errorf("expected updated product with price 25.00 and stock 75, got price %f and stock %d", product.Price, product.Stock)
	}
}

// TestDeleteProduct tests the DeleteProduct method
func TestDeleteProduct(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		if r.URL.Path != "/products/1" {
			t.Errorf("expected request to /products/1, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	err := client.DeleteProduct(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
