// client_test.go

package mockclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewClient tests the NewClient function
func TestNewClient(t *testing.T) {
	hostURL := "http://example.com"
	client, err := NewClient(&hostURL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client.HostURL != hostURL {
		t.Errorf("expected host URL to be %v, got %v", hostURL, client.HostURL)
	}
	if client.HTTPClient.Timeout != 10*time.Second {
		t.Errorf("expected timeout to be 10s, got %v", client.HTTPClient.Timeout)
	}
}

// TestNewClient_NilURL tests NewClient with a nil hostURL
func TestNewClient_NilURL(t *testing.T) {
	_, err := NewClient(nil)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	expectedError := "hostURL is required"
	if err.Error() != expectedError {
		t.Errorf("expected error message to be %q, got %q", expectedError, err.Error())
	}
}

// TestDoRequest tests the doRequest function
func TestDoRequest(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	req, _ := http.NewRequest("GET", server.URL, nil)
	body, err := client.doRequest(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expectedBody := `{"message": "success"}`
	if string(body) != expectedBody {
		t.Errorf("expected body to be %q, got %q", expectedBody, string(body))
	}
}

// TestDoRequest_NonOKResponse tests doRequest with a non-OK status
func TestDoRequest_NonOKResponse(t *testing.T) {
	// Create a mock server with a 404 response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	req, _ := http.NewRequest("GET", server.URL, nil)
	_, err := client.doRequest(req)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	expectedError := "error: status: 404, body: {\"error\": \"not found\"}"
	if err.Error() != expectedError {
		t.Errorf("expected error message to be %q, got %q", expectedError, err.Error())
	}
}

// TestDoRequest_InvalidBody tests doRequest with an invalid response body
func TestDoRequest_InvalidBody(t *testing.T) {
	// Create a mock server with a malformed body
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{`)) // Incomplete JSON
	}))
	defer server.Close()

	client, _ := NewClient(&server.URL)
	req, _ := http.NewRequest("GET", server.URL, nil)
	body, err := client.doRequest(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expectedBody := `{`
	if string(body) != expectedBody {
		t.Errorf("expected body to be %q, got %q", expectedBody, string(body))
	}
}
