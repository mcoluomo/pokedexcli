package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mcoluomo/pokedexcli/cache"
	"github.com/mcoluomo/pokedexcli/location"
	"github.com/mcoluomo/pokedexcli/pokemon"
)

// Client handles API communication
type Client struct {
	baseURL    string
	httpClient *http.Client
	cache      *cache.Cache
	next       string
	prev       string
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		baseURL: "https://pokeapi.co/api/v2",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: cache.NewCache(15 * time.Minute),
		next:  "https://pokeapi.co/api/v2/location-area/",
		prev:  "",
	}
}

// GetPokemon fetches a Pokemon and converts to domain model
func (c *Client) GetPokemon(name string) (pokemon.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, name)

	// Check cache first
	if data, exists := c.cache.Get(url); exists {
		return c.parsePokemonResponse(data)
	}

	// Make HTTP request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return pokemon.Pokemon{}, fmt.Errorf("failed to fetch Pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pokemon.Pokemon{}, fmt.Errorf("API returned status %d for Pokemon %s", resp.StatusCode, name)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return pokemon.Pokemon{}, fmt.Errorf("failed to read response: %w", err)
	}

	// Cache the response
	c.cache.Set(url, body)

	return c.parsePokemonResponse(body)
}

// GetLocationAreas fetches location areas
func (c *Client) GetLocationAreas() ([]location.LocationArea, error) {
	url := c.next
	if url == "" {
		return nil, fmt.Errorf("no more locations available")
	}

	// Check cache first
	if data, exists := c.cache.Get(url); exists {
		return c.parseLocationResponse(data)
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Cache the response
	c.cache.Set(url, body)

	areas, err := c.parseLocationResponse(body)
	if err != nil {
		return nil, err
	}

	// Update pagination URLs
	c.updatePagination(body)

	return areas, nil
}

// GetPreviousLocationAreas fetches previous page of location areas
func (c *Client) GetPreviousLocationAreas() ([]location.LocationArea, error) {
	if c.prev == "" {
		return nil, fmt.Errorf("no previous locations available")
	}

	// Temporarily swap next and prev
	originalNext := c.next
	c.next = c.prev

	areas, err := c.GetLocationAreas()
	if err != nil {
		c.next = originalNext
		return nil, err
	}

	return areas, nil
}

// ExploreLocation fetches Pokemon in a specific location
func (c *Client) ExploreLocation(areaName string) (location.LocationArea, error) {
	url := fmt.Sprintf("%s/location-area/%s", c.baseURL, areaName)

	// Check cache first
	if data, exists := c.cache.Get(url); exists {
		return c.parseLocationAreaResponse(data)
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return location.LocationArea{}, fmt.Errorf("failed to explore location: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return location.LocationArea{}, fmt.Errorf("location %s not found", areaName)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return location.LocationArea{}, fmt.Errorf("failed to read response: %w", err)
	}

	// Cache the response
	c.cache.Set(url, body)

	return c.parseLocationAreaResponse(body)
}

// HasNext checks if there are more locations to fetch
func (c *Client) HasNext() bool {
	return c.next != ""
}

// HasPrev checks if there are previous locations to fetch
func (c *Client) HasPrev() bool {
	return c.prev != ""
}

// parsePokemonResponse converts API response to Pokemon domain model
func (c *Client) parsePokemonResponse(data []byte) (pokemon.Pokemon, error) {
	var apiResp struct {
		Name           string `json:"name"`
		Height         int    `json:"height"`
		Weight         int    `json:"weight"`
		BaseExperience int    `json:"base_experience"`
		Types          []struct {
			Type struct {
				Name string `json:"name"`
			} `json:"type"`
		} `json:"types"`
		Stats []struct {
			BaseStat int `json:"base_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		} `json:"stats"`
	}

	if err := json.Unmarshal(data, &apiResp); err != nil {
		return pokemon.Pokemon{}, fmt.Errorf("failed to parse Pokemon response: %w", err)
	}

	// Convert to domain model
	p := pokemon.Pokemon{
		Name:           apiResp.Name,
		Height:         apiResp.Height,
		Weight:         apiResp.Weight,
		BaseExperience: apiResp.BaseExperience,
		Types:          make([]string, len(apiResp.Types)),
		Stats:          make(map[string]int),
	}

	// Convert types
	for i, t := range apiResp.Types {
		p.Types[i] = t.Type.Name
	}

	// Convert stats
	for _, s := range apiResp.Stats {
		p.Stats[s.Stat.Name] = s.BaseStat
	}

	return p, nil
}

// parseLocationResponse converts location areas API response
func (c *Client) parseLocationResponse(data []byte) ([]location.LocationArea, error) {
	var apiResp struct {
		Results []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}

	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse location response: %w", err)
	}

	areas := make([]location.LocationArea, len(apiResp.Results))
	for i, result := range apiResp.Results {
		areas[i] = location.LocationArea{
			Name:    result.Name,
			Pokemon: []string{}, // Will be populated when exploring
		}
	}

	return areas, nil
}

// parseLocationAreaResponse converts specific location area API response
func (c *Client) parseLocationAreaResponse(data []byte) (location.LocationArea, error) {
	var apiResp struct {
		Name              string `json:"name"`
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}

	if err := json.Unmarshal(data, &apiResp); err != nil {
		return location.LocationArea{}, fmt.Errorf("failed to parse location area response: %w", err)
	}

	// Extract unique Pokemon names
	pokemonMap := make(map[string]bool)
	for _, encounter := range apiResp.PokemonEncounters {
		pokemonMap[encounter.Pokemon.Name] = true
	}

	// Convert to slice
	pokemon := make([]string, 0, len(pokemonMap))
	for name := range pokemonMap {
		pokemon = append(pokemon, name)
	}

	return location.LocationArea{
		Name:    apiResp.Name,
		Pokemon: pokemon,
	}, nil
}

// updatePagination updates next and prev URLs from API response
func (c *Client) updatePagination(data []byte) {
	var apiResp struct {
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
	}

	if err := json.Unmarshal(data, &apiResp); err != nil {
		return
	}

	if apiResp.Next != nil {
		c.next = *apiResp.Next
	} else {
		c.next = ""
	}

	if apiResp.Previous != nil {
		c.prev = *apiResp.Previous
	} else {
		c.prev = ""
	}
}
