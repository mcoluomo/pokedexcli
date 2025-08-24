package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mcoluomo/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config, string) error
}

var UsableCommands map[string]cliCommand

func init() {
	UsableCommands = map[string]cliCommand{
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
			description: "displays next 20 location areas in the Pokemon world.",
			callback:    pokeapi.CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays previous 20 location areas in the Pokemon world.",
			callback:    pokeapi.CommandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Lists all the pokemon of an area",
			callback:    pokeapi.CommandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Lets us catch a pokemon",
			callback:    pokeapi.CommandCatch,
		},
	}
}

func statRepl() {
	c := &pokeapi.Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "https://pokeapi.co/api/v2/location-area/",
	}
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		cmd, exists := UsableCommands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		if commandName == "explore" || commandName == "catch" {
			if len(words) < 2 {
				fmt.Printf("Argument missing. Please provide: %s <argument_name>\n", commandName)
				continue
			}
			err := cmd.callback(c, words[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}

		} else {
			err := cmd.callback(c, "")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)

	if len(words) == 0 {
		return []string{""}
	}

	return words
}

func commandExit(c *pokeapi.Config, areaName string) error {
	defer os.Exit(0)
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func commandHelp(c *pokeapi.Config, areaName string) error {
	helpMsg := "Welcome to the Pokedex!\nUsage:\n\n"
	for cmd := range UsableCommands {
		helpMsg += cmd + ": " + UsableCommands[cmd].description + "\n"
	}
	fmt.Println(helpMsg)
	return nil
}
