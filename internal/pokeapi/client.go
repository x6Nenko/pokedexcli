package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{}, // Default HTTP client
	}
}

// ResponseData represents the structure you expect from the API
type ResponseData struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Item  `json:"results"`
}

type Item struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// FetchItems makes a GET request to the specified URL
func (c *Client) FetchItems(url string) (*ResponseData, error) {
	// Step 1: Make HTTP request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Always close the body!

	// Step 2: Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Step 3: Unmarshal JSON into struct
	var data ResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}