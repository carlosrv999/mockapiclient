// users_test.go

package mockclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetUsers tests the GetUsers method
func TestGetUsers(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user" {
			t.Errorf("expected request to /user, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": "1", "name": "John", "lastName": "Doe", "address": "123 Main St", "favoriteDogBreed": "Labrador", "createdAt": "2024-08-05T10:00:00Z"}, {"id": "2", "name": "Jane", "lastName": "Smith", "address": "456 Elm St", "favoriteDogBreed": "Beagle", "createdAt": "2024-08-05T11:00:00Z"}]`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	users, err := client.GetUsers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
	if users[0].Name != "John" || users[1].Name != "Jane" {
		t.Errorf("expected user names 'John' and 'Jane', got %s and %s", users[0].Name, users[1].Name)
	}
	if users[0].FavoriteDogBreed != "Labrador" || users[1].FavoriteDogBreed != "Beagle" {
		t.Errorf("expected favorite dog breeds 'Labrador' and 'Beagle', got %s and %s", users[0].FavoriteDogBreed, users[1].FavoriteDogBreed)
	}
}

// TestGetUsers_InvalidResponse tests the GetUsers method with an invalid response
func TestGetUsers_InvalidResponse(t *testing.T) {
	// Create a mock server with an invalid JSON response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	_, err := client.GetUsers()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

// TestGetUsers_ErrorResponse tests the GetUsers method with a server error response
func TestGetUsers_ErrorResponse(t *testing.T) {
	// Create a mock server that returns a 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	_, err := client.GetUsers()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

// TestGetUserByID tests the GetUserByID method for a successful response
func TestGetUserByID(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/user/1"
		if r.URL.Path != expectedPath {
			t.Errorf("expected request to %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "1", "name": "John", "lastName": "Doe", "address": "123 Main St", "favoriteDogBreed": "Labrador", "createdAt": "2024-08-05T10:00:00Z"}`)) // id as a string
	}))
	defer server.Close()

	client, err := NewClient(&server.URL)
	if err != nil {
		t.Fatalf("unexpected error while creating client: %v", err)
	}

	user, err := client.GetUserByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID to be 1, got %d", user.ID)
	}
	if user.Name != "John" {
		t.Errorf("expected user name to be 'John', got '%s'", user.Name)
	}
	if user.FavoriteDogBreed != "Labrador" {
		t.Errorf("expected favorite dog breed to be 'Labrador', got '%s'", user.FavoriteDogBreed)
	}

	// Additional debug logging to help identify errors
	t.Logf("User Retrieved: %+v", user)
}

// TestGetUserByID_NotFound tests the GetUserByID method for a user not found scenario
func TestGetUserByID_NotFound(t *testing.T) {
	// Create a mock server that returns a 404 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "user not found"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	_, err := client.GetUserByID(999) // Assume 999 is a non-existing user ID
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	expectedError := "error: status: 404, body: {\"error\": \"user not found\"}"
	if err.Error() != expectedError {
		t.Errorf("expected error message %q, got %q", expectedError, err.Error())
	}
}

// TestGetUserByID_InvalidResponse tests the GetUserByID method for an invalid JSON response
func TestGetUserByID_InvalidResponse(t *testing.T) {
	// Create a mock server with an invalid JSON response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	_, err := client.GetUserByID(1)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

// TestGetUserByID_ServerError tests the GetUserByID method for a server error response
func TestGetUserByID_ServerError(t *testing.T) {
	// Create a mock server that returns a 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	_, err := client.GetUserByID(1)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	expectedError := "error: status: 500, body: {\"error\": \"internal server error\"}"
	if err.Error() != expectedError {
		t.Errorf("expected error message %q, got %q", expectedError, err.Error())
	}
}

// TestUserUnmarshalJSON_Valid tests the UnmarshalJSON method for valid JSON input
func TestUserUnmarshalJSON_Valid(t *testing.T) {
	jsonData := `{"id": "1", "name": "John", "lastName": "Doe", "address": "123 Main St", "favoriteDogBreed": "Labrador", "createdAt": "2024-08-05T10:00:00Z"}`
	var user User

	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID to be 1, got %d", user.ID)
	}
	if user.Name != "John" {
		t.Errorf("expected user name to be 'John', got '%s'", user.Name)
	}
	if user.FavoriteDogBreed != "Labrador" {
		t.Errorf("expected favorite dog breed to be 'Labrador', got '%s'", user.FavoriteDogBreed)
	}
}

// TestUserUnmarshalJSON_InvalidID tests the UnmarshalJSON method for JSON with an invalid ID
func TestUserUnmarshalJSON_InvalidID(t *testing.T) {
	jsonData := `{"id": "abc", "name": "John", "lastName": "Doe", "address": "123 Main St", "favoriteDogBreed": "Labrador", "createdAt": "2024-08-05T10:00:00Z"}`
	var user User

	err := json.Unmarshal([]byte(jsonData), &user)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	expectedError := "invalid id: abc"
	if err.Error() != expectedError {
		t.Errorf("expected error message %q, got %q", expectedError, err.Error())
	}
}

// TestUserUnmarshalJSON_MissingID tests the UnmarshalJSON method for JSON missing the ID
func TestUserUnmarshalJSON_MissingID(t *testing.T) {
	jsonData := `{"name": "John", "lastName": "Doe", "address": "123 Main St", "favoriteDogBreed": "Labrador", "createdAt": "2024-08-05T10:00:00Z"}`
	var user User

	err := json.Unmarshal([]byte(jsonData), &user)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	expectedError := "invalid id: "
	if err.Error() != expectedError {
		t.Errorf("expected error message %q, got %q", expectedError, err.Error())
	}
}

// TestUserMarshalJSON tests the MarshalJSON method for valid User struct
func TestUserMarshalJSON(t *testing.T) {
	user := User{
		ID:               1,
		Name:             "John",
		LastName:         "Doe",
		Address:          "123 Main St",
		FavoriteDogBreed: "Labrador",
		CreatedAt:        "2024-08-05T10:00:00Z",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedJSON := `{"id":"1","name":"John","lastName":"Doe","address":"123 Main St","favoriteDogBreed":"Labrador","createdAt":"2024-08-05T10:00:00Z"}`
	if string(data) != expectedJSON {
		t.Errorf("expected JSON %s, got %s", expectedJSON, string(data))
	}
}

// TestUserMarshalUnmarshal tests round-trip marshal and unmarshal
func TestUserMarshalUnmarshal(t *testing.T) {
	originalUser := User{
		ID:               2,
		Name:             "Alice",
		LastName:         "Smith",
		Address:          "456 Oak St",
		FavoriteDogBreed: "Golden Retriever",
		CreatedAt:        "2024-08-06T10:00:00Z",
	}

	// Marshal the original user to JSON
	data, err := json.Marshal(originalUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Unmarshal the JSON back to a User struct
	var unmarshalledUser User
	err = json.Unmarshal(data, &unmarshalledUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Compare the original and unmarshalled User structs
	if originalUser != unmarshalledUser {
		t.Errorf("expected unmarshalled user to be %+v, got %+v", originalUser, unmarshalledUser)
	}
}

// TestCreateUser tests the CreateUser method for a successful user creation
func TestCreateUser(t *testing.T) {
	// Setup a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/user" {
			t.Errorf("expected request to /user, got %s", r.URL.Path)
		}

		// Check the request body
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			t.Fatalf("expected no error decoding request body, got %v", err)
		}

		// Verify that the user data is correct
		if user.Name != "Alice" || user.LastName != "Smith" {
			t.Errorf("unexpected user data: %+v", user)
		}

		// Simulate a successful response
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "1", "name": "Alice", "lastName": "Smith", "address": "456 Elm St", "favoriteDogBreed": "Beagle", "createdAt": "2024-08-06T10:00:00Z"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	newUser := &User{
		Name:             "Alice",
		LastName:         "Smith",
		Address:          "456 Elm St",
		FavoriteDogBreed: "Beagle",
	}

	createdUser, err := client.CreateUser(newUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the returned user
	if createdUser.ID != 1 {
		t.Errorf("expected user ID to be 1, got %d", createdUser.ID)
	}
	if createdUser.Name != "Alice" {
		t.Errorf("expected user name to be 'Alice', got '%s'", createdUser.Name)
	}
	if createdUser.FavoriteDogBreed != "Beagle" {
		t.Errorf("expected favorite dog breed to be 'Beagle', got '%s'", createdUser.FavoriteDogBreed)
	}
}

// TestUpdateUser tests the UpdateUser method for a successful user update
func TestUpdateUser(t *testing.T) {
	// Setup a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH method, got %s", r.Method)
		}
		if r.URL.Path != "/user/1" {
			t.Errorf("expected request to /user/1, got %s", r.URL.Path)
		}

		// Check the request body
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			t.Fatalf("expected no error decoding request body, got %v", err)
		}

		// Verify that the user data is correct
		if user.ID != 1 || user.Name != "Bob" {
			t.Errorf("unexpected user data: %+v", user)
		}

		// Simulate a successful response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "1", "name": "Bob", "lastName": "Smith", "address": "123 Main St", "favoriteDogBreed": "Labrador", "createdAt": "2024-08-06T10:00:00Z"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	updatedUser := &User{
		ID:               1,
		Name:             "Bob",
		LastName:         "Smith",
		Address:          "123 Main St",
		FavoriteDogBreed: "Labrador",
	}

	user, err := client.UpdateUser(updatedUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the returned user
	if user.ID != 1 {
		t.Errorf("expected user ID to be 1, got %d", user.ID)
	}
	if user.Name != "Bob" {
		t.Errorf("expected user name to be 'Bob', got '%s'", user.Name)
	}
	if user.FavoriteDogBreed != "Labrador" {
		t.Errorf("expected favorite dog breed to be 'Labrador', got '%s'", user.FavoriteDogBreed)
	}
}

// TestDeleteUser tests the DeleteUser method for a successful user deletion
func TestDeleteUser(t *testing.T) {
	// Setup a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		if r.URL.Path != "/user/1" {
			t.Errorf("expected request to /user/1, got %s", r.URL.Path)
		}

		// Simulate a successful response
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	err := client.DeleteUser(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
