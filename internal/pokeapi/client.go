package pokeapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

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

	fmt.Println("you're on the first page")
	return nil
}
