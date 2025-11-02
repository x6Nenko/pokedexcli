package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"github.com/x6Nenko/pokedexcli/internal/pokecache"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	cache      *pokecache.Cache
}

// NewClient creates a new API client
func NewClient(interval time.Duration, baseURL string) *Client {
	newCache := pokecache.NewCache(interval)
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{}, // Default HTTP client
		cache: newCache,
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
	 // Try cache first
	if cachedData, found := c.cache.Get(url); found {
		var data ResponseData
		err := json.Unmarshal(cachedData, &data) // cachedData - bytes
		if err != nil {
			return nil, err
		}
		// fmt.Println("===== Using cached data =====")
		return &data, nil
	}

	// fmt.Println("===== New Fetch =====")

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

	c.cache.Add(url, body)

	// Step 3: Unmarshal JSON into struct
	var data ResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}