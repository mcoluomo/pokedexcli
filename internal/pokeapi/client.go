package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LocationAreaList struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
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
	res, err := http.Get(c.Next)
	if err != nil {
		log.Fatal(err)
	}
	bodyDate, err := io.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, bodyDate)
	}

	var locationAreas LocationAreaList
	if err := json.Unmarshal(bodyDate, &locationAreas); err != nil {
		log.Fatalf("Json decode failure: %v", err)
	}

	// case when when next is nil

	if locationAreas.Previous == nil {
		fmt.Println("you're on the first page")
	}
	for _, pokeLocation := range locationAreas.Results {
		fmt.Println(pokeLocation.Name)
	}

	return nil
}
