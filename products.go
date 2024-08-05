package mockclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID         int     `json:"id,omitempty"`
	Name       string  `json:"name,omitempty"`
	Price      float64 `json:"price,omitempty"`
	Stock      int     `json:"stock,omitempty"`
	CreatedAt  string  `json:"createdAt,omitempty"`
	Type       string  `json:"type,omitempty"`
	Department string  `json:"department,omitempty"`
}

// UnmarshalJSON - Custom unmarshal for Product struct
func (p *Product) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to hold JSON data with ID, Price and Stock as a string
	type Alias Product
	temp := &struct {
		ID    string `json:"id"`
		Price string `json:"price"`
		Stock string `json:"stock"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	// Unmarshal JSON into the temporary struct
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Convert ID and Stock to int and Price to float64
	id, err := strconv.Atoi(temp.ID)
	if err != nil {
		return fmt.Errorf("invalid id: %s", temp.ID)
	}
	p.ID = id

	stock, err := strconv.Atoi(temp.Stock)
	if err != nil {
		return fmt.Errorf("invalid stock: %s", temp.Stock)
	}
	p.Stock = stock

	price, err := strconv.ParseFloat(temp.Price, 64)
	if err != nil {
		return fmt.Errorf("invalid price: %s", temp.Price)
	}
	p.Price = price
	return nil
}

// Implementinc custom marshaler for Product
func (p Product) MarshalJSON() ([]byte, error) {
	// Define a temporary struct to hold JSON data with ID, Price and Stock as a string
	type Alias Product
	temp := &struct {
		ID    string `json:"id"`
		Price string `json:"price"`
		Stock string `json:"stock"`
		*Alias
	}{
		ID:    strconv.Itoa(p.ID), // Convert ID to string
		Price: fmt.Sprintf("%f", p.Price),
		Stock: strconv.Itoa(p.Stock),
		Alias: (*Alias)(&p),
	}

	// Marshal the temporary struct to JSON
	return json.Marshal(temp)
}

func (c *Client) GetProducts() ([]Product, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/products", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	err = json.Unmarshal(body, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// GetProductByID - Get product by ID
func (c *Client) GetProductByID(id int) (*Product, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/products/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	product := &Product{}
	err = json.Unmarshal(body, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// CreateProduct - Create a new product
func (c *Client) CreateProduct(product *Product) (*Product, error) {
	rb, err := json.Marshal(*product)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/products", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newProduct := Product{}

	err = json.Unmarshal(body, &newProduct)
	if err != nil {
		return nil, err
	}

	return &newProduct, nil
}

// UpdateProduct - Update an existing product (no auth required)
func (c *Client) UpdateProduct(product *Product) (*Product, error) {
	rb, err := json.Marshal(*product)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/products/%d", c.HostURL, product.ID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedProduct := Product{}

	err = json.Unmarshal(body, &updatedProduct)
	if err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

// DeleteProduct - Delete a product by ID
func (c *Client) DeleteProduct(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/products/%d", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
