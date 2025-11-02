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

// =============
// ResponseData represents the structure you expect from the API (list of areas)
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

// =============
// Explore pokemons in area
type LocationDetail struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

// Each encounter in the list
type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

// The pokemon name
type Pokemon struct {
	Name string `json:"name"`
}

// =============
// Get pokemon details info
type PokemonDetails struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

// FetchItems makes a GET request to the specified URL (GETS LIST OF AREAS)
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

// GET request to the specified URL (GETS LIST OF POKEMONS IN THE AREA)
func (c *Client) FetchLocationDetail(url string) (*LocationDetail, error) {
	 // Try cache first
	if cachedData, found := c.cache.Get(url); found {
		var data LocationDetail
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
	var data LocationDetail
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GET request to the specified URL (GETS POKEMON DATA)
func (c *Client) FetchPokemon(url string) (*PokemonDetails, error) {
	 // Try cache first
	if cachedData, found := c.cache.Get(url); found {
		var data PokemonDetails
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
	var data PokemonDetails
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GENERIC FETCH
// func (c *Client) Fetch[T any](url string) (*T, error) {
// 	// Try cache first
// 	if cachedData, found := c.cache.Get(url); found {
// 		var data T
// 		err := json.Unmarshal(cachedData, &data) // cachedData - bytes
// 		if err != nil {
// 			return nil, err
// 		}
// 		// fmt.Println("===== Using cached data =====")
// 		return &data, nil
// 	}

// 	// fmt.Println("===== New Fetch =====")

// 	// Step 1: Make HTTP request
// 	resp, err := c.httpClient.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close() // Always close the body!

// 	// Step 2: Read response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	c.cache.Add(url, body)

// 	// Step 3: Unmarshal JSON into struct
// 	var data T
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &data, nil
// }