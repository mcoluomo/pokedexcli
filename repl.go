package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
				cmd.callback()
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

func commandExit() error {
	defer os.Exit(0)
	return fmt.Errorf("Closing the Pokedex... Goodbye!\n")
}

func commandHelp() error {
	helpMsg := "Welcome to the Pokedex!\nUsage:\n\n"
	for cmd := range usableCommands {
		helpMsg += cmd + ": " + usableCommands[cmd].description + "\n"
	}

	fmt.Printf(helpMsg)
	return nil
}
