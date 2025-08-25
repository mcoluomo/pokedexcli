package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokemonData struct {
	Species PokemonInfo `json:"species"`
	Details Pokemon     `json:"details`
}

type Pokedex struct {
	CaughtPokemon map[string]PokemonInfo
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		CaughtPokemon: make(map[string]PokemonInfo),
	}
}

func (p *Pokedex) AddPokemon(name string, info PokemonInfo, details Pokemon) {
	p.CaughtPokemon[name] = PokemonData{
		Species: species,
		Details: details,
	}
}

func (p *Pokedex) HasPokemon(name string) bool {
	_, exists := p.CaughtPokemon[name]
	return exists
}

func (p *Pokedex) ListCaught() {
	fmt.Println("Your Pokemon:")
	for name := range p.CaughtPokemon {
		fmt.Printf("  - %s\n", name)
	}
}

func CommandCatch(c *Config, pokeName string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokeName)
	pokeSpecies := "https://pokeapi.co/api/v2/pokemon/" + pokeName

	var responseBody []byte

	if cachedResBody, inCache := cache.Get(pokeSpecies); inCache {
		fmt.Println("Using cached data...")
		responseBody = cachedResBody
	} else {
		res, err := http.Get(pokeSpecies)
		if err != nil {
			return fmt.Errorf("failed to fetch response: %w", err)
		}
		defer res.Body.Close()

		responseBody, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if res.StatusCode > 299 {
			return fmt.Errorf("API request failed with status %d", res.StatusCode)
		}

		cache.Add(pokeSpecies, responseBody)

	}

	var pokeInfo PokemonInfo
	if err := json.Unmarshal(responseBody, &pokeInfo); err != nil {
		return fmt.Errorf("failed to decode pokemon data: %w", err)
	}

	Dex.AddPokemon(pokeName, pokeInfo)

	bonus := RandomBonus()
	captureRate, caught := SimpleCatch(pokeInfo.BaseExperience, bonus)

	if caught {
		fmt.Printf("%s was caught! (Capture rate: %f)\n", pokeName, captureRate)
	} else {
		fmt.Printf("%s escaped! (Capture rate: %f)\n", pokeName, captureRate)
	}

	return nil
}

func CommandInspect(c *Config, name string) error {
	if inDex := Dex.HasPokemon(name); inDex {
		fmt.Printf(`
			Height: %d\n
			Weight: %d\n
			Stats:\n
			  -hp: %d\n
			  -attack: %d\n
			  -defense: %d\n
			  -special-attack: %d\n
			  -special-defense: %d\n
			  -speed: %d\n
			Types:
			  - normal
			  - flying
		`, Dex.CaughtPokemon[name].Height, Dex.CaughtPokemon[name].Weight, Dex.CaughtPokemon[name].Abilities)
	} else {
		return fmt.Errorf("%s is not present in the pokedex", name)
	}

	return nil
}
