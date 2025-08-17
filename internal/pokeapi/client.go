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

	res, err := http.Get(c.Next)
	if err != nil {
		log.Fatal(err)
	}

	bodyData, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, bodyData)
	}

	var locationAreas LocationAreaList
	if err := json.Unmarshal(bodyData, &locationAreas); err != nil {
		log.Fatalf("Json decode failure: %v", err)
	}
	// Update config safely
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

	return nil
}

func CommandMapBack(c *Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := http.Get(c.Previous)
	if err != nil {
		log.Fatal(err)
	}

	bodyData, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, bodyData)
	}

	var locationAreas LocationAreaList
	if err := json.Unmarshal(bodyData, &locationAreas); err != nil {
		log.Fatalf("Json decode failure: %v", err)
	}
	// Update config safely
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

	return nil
}
