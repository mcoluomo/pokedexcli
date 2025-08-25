package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mcoluomo/pokedexcli/internal/pokecache"
)

var cache = pokecache.NewCache(15)

var Dex = NewPokedex()

type LocationAreas struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next     string
	Previous string
}

type PokemonLocations struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func CommandMap(c *Config, areaName string) error {
	if c.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	if cachedResBody, isCached := cache.Get(c.Next); isCached {
		fmt.Println("Using cached data...")
		DecodeResBody(c, cachedResBody)
		return nil
	}

	res, err := http.Get(c.Next)
	if err != nil {
		log.Fatal(err)
	}

	responseBody, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, responseBody)
	}

	cache.Add(c.Next, responseBody)

	DecodeResBody(c, responseBody)

	return nil
}

func CommandMapBack(c *Config, areaName string) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	if cachedResBody, inCache := cache.Get(c.Previous); inCache {
		fmt.Println("Using cached data...")
		DecodeResBody(c, cachedResBody)
		return nil
	}
	res, err := http.Get(c.Previous)
	if err != nil {
		log.Fatal(err)
	}

	responseBody, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, responseBody)
	}

	cache.Add(c.Previous, responseBody)

	DecodeResBody(c, responseBody)

	return nil
}

func CommandExplore(c *Config, areaName string) error {
	locUrl := "https://pokeapi.co/api/v2/location-area/" + areaName
	if cachedResBody, inCache := cache.Get(locUrl); inCache {
		fmt.Println("Using cached data...")
		DecodeResBody(c, cachedResBody)
		return nil
	}
	res, err := http.Get(locUrl)
	if err != nil {
		log.Fatal(err)
	}

	responseBody, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, responseBody)
	}

	cache.Add(locUrl, responseBody)

	var pokeLocation PokemonLocations
	if err := json.Unmarshal(responseBody, &pokeLocation); err != nil {
		log.Fatalf("Json decode failure: %v", err)
	}

	fmt.Printf("Exploring %s...\nFound Pokemon:\n", areaName)
	for _, pstruct := range pokeLocation.PokemonEncounters {
		fmt.Println("-", pstruct.Pokemon.Name)
	}

	return nil
}

func DecodeResBody(c *Config, body []byte) {
	var locationAreas LocationAreas
	if err := json.Unmarshal(body, &locationAreas); err != nil {
		log.Fatalf("Json decode failure: %v", err)
	}
	if locationAreas.Next != nil {
		c.Next = *locationAreas.Next
	} else {
		c.Next = ""
	}

	if locationAreas.Previous != nil {
		c.Previous = *locationAreas.Previous
	} else {
		c.Previous = ""
	}

	for _, pokeArea := range locationAreas.Results {
		fmt.Println(pokeArea.Name)
	}
}
