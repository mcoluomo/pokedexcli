package pokeapi

import (
	"testing"
)

func TestDecodeAndOutputRequstData(t *testing.T) {
	// Sample JSON response similar to what the API would return
	jsonResponse := `{
		"count": 1302,
		"next": "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
		"previous": null,
		"results": [
			{
				"name": "canalave-city-area",
				"url": "https://pokeapi.co/api/v2/location-area/1/"
			},
			{
				"name": "eterna-city-area",
				"url": "https://pokeapi.co/api/v2/location-area/2/"
			}
		]
	}`

	config := &Config{
		Next:     "",
		Previous: "",
	}

	// Test the decode function
	DecodeAndOutputRequstData(config, []byte(jsonResponse))

	// Verify that config was updated correctly
	expectedNext := "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20"
	if config.Next != expectedNext {
		t.Errorf("Expected Next to be %s, got %s", expectedNext, config.Next)
	}

	if config.Previous != "" {
		t.Errorf("Expected Previous to be empty, got %s", config.Previous)
	}
}

func TestDecodeAndOutputRequstDataWithPrevious(t *testing.T) {
	// JSON with both next and previous
	jsonResponse := `{
		"count": 1302,
		"next": "https://pokeapi.co/api/v2/location-area/?offset=40&limit=20",
		"previous": "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		"results": [
			{
				"name": "test-area",
				"url": "https://pokeapi.co/api/v2/location-area/1/"
			}
		]
	}`

	config := &Config{}
	DecodeAndOutputRequstData(config, []byte(jsonResponse))

	expectedNext := "https://pokeapi.co/api/v2/location-area/?offset=40&limit=20"
	expectedPrevious := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

	if config.Next != expectedNext {
		t.Errorf("Expected Next to be %s, got %s", expectedNext, config.Next)
	}

	if config.Previous != expectedPrevious {
		t.Errorf("Expected Previous to be %s, got %s", expectedPrevious, config.Previous)
	}
}
