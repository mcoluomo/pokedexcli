package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     string
	Previous string
}

var usableCommands map[string]cliCommand

func init() {
	usableCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},

		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "displays 20 location areas in the Pokemon world.",
			callback:    commandMap,
		},
	}
}

func statRepl() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		for _, word := range words {
			if cmd, ok := usableCommands[word]; ok {
				cmd.callback(&config{Next: "https://pokeapi.co/api/v2/location/", Previous: ""}) // previous field string)
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandExit(c *config) error {
	defer os.Exit(0)
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp(c *config) error {
	helpMsg := "Welcome to the Pokedex!\nUsage:\n\n"
	for cmd := range usableCommands {
		helpMsg += cmd + ": " + usableCommands[cmd].description + "\n"
	}

	fmt.Println(helpMsg)
	return nil
}

func commandMap(c *config) error {
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
