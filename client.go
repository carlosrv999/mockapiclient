package mockclient

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client
type Client struct {
	HostURL    string
	HTTPClient *http.Client
}

// NewClient
func NewClient(hostURL *string) (*Client, error) {
	if hostURL == nil {
		return nil, fmt.Errorf("hostURL is required")
	}

	return &Client{
		HostURL:    *hostURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// doRequest
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("error: status: %d, body: %s", res.StatusCode, body)
	}

	return body, nil
}
