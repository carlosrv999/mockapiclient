package mockclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	Address          string `json:"address,omitempty"`
	FavoriteDogBreed string `json:"favoriteDogBreed,omitempty"`
	CreatedAt        string `json:"createdAt,omitempty"`
}

// GetUsers - Returns a list of users (no auth required)
func (c *Client) GetUsers() ([]User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	users := []User{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID - Returns a user by ID (no auth required)
func (c *Client) GetUserByID(id int) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Implementing custom unmarshaler for User
func (u *User) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to hold JSON data with ID as a string
	type Alias User
	temp := &struct {
		ID string `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	// Unmarshal JSON into the temporary struct
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Convert ID from string to int
	id, err := strconv.Atoi(temp.ID)
	if err != nil {
		return fmt.Errorf("invalid id: %s", temp.ID)
	}
	u.ID = id
	return nil
}

// Implementing custom marshaler for User
func (u User) MarshalJSON() ([]byte, error) {
	// Define a temporary struct to hold JSON data with ID as a string
	type Alias User
	temp := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    strconv.Itoa(u.ID), // Convert ID from int to string
		Alias: (*Alias)(&u),
	}

	// Marshal the temporary struct to JSON
	return json.Marshal(temp)
}

// CreateUser - Creates a new user (no auth required)
func (c *Client) CreateUser(user *User) (*User, error) {
	rb, err := json.Marshal(*user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newUser := User{}

	err = json.Unmarshal(body, &newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// UpdateUser - Updates an existing user (no auth required)
func (c *Client) UpdateUser(user *User) (*User, error) {
	rb, err := json.Marshal(*user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/user/%d", c.HostURL, user.ID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedUser := User{}

	err = json.Unmarshal(body, &updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// DeleteUser - Deletes a user by ID (no auth required)
func (c *Client) DeleteUser(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user/%d", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
