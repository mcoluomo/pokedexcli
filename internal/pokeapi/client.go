package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/mcoluomo/pokedexcli/internal/pokecache"
)

var jsonHandler = slog.NewJSONHandler(os.Stdout, nil)

var logger = slog.New(jsonHandler)

var cache = pokecache.NewCache(15)

type LocationAreaList struct {
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

func CommandMap(c *Config) error {
	if c.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	if _, ok := cache.Entry[c.Next]; ok {
		cahedresBody, _ := cache.Get(c.Next)
		logger.Info("cached data present in cache")
		DecodeAndOutputRequstData(c, cahedresBody)
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

	DecodeAndOutputRequstData(c, responseBody)

	return nil
}

func CommandMapBack(c *Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	if _, ok := cache.Entry[c.Previous]; ok {
		cahedresBody, _ := cache.Get(c.Previous)
		logger.Info("cached data present in cache")
		DecodeAndOutputRequstData(c, cahedresBody)
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

	DecodeAndOutputRequstData(c, responseBody)

	return nil
}

func DecodeAndOutputRequstData(c *Config, body []byte) {
	var locationAreas LocationAreaList
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

	for _, pokeLocation := range locationAreas.Results {
		fmt.Println(pokeLocation.Name)
	}
}
