package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// REPL represents the Read-Eval-Print Loop
type REPL struct {
	app      *App
	commands map[string]Command
	running  bool
}

// NewREPL creates a new REPL instance
func NewREPL() *REPL {
	return &REPL{
		app:      NewApp(),
		commands: GetCommands(),
		running:  true,
	}
}

// Start begins the REPL loop
func (r *REPL) Start() {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Type 'help' for available commands.")

	scanner := bufio.NewScanner(os.Stdin)

	for r.running {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		r.processInput(input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}

// processInput handles user input
func (r *REPL) processInput(input string) {
	words := r.CleanInput(input)
	if len(words) == 0 {
		return
	}

	commandName := words[0]
	cmd, exists := r.commands[commandName]
	if !exists {
		fmt.Printf("Unknown command '%s'. Type 'help' for available commands.\n", commandName)
		return
	}

	// Check if command requires an argument
	var arg string
	if cmd.RequiresArg {
		if len(words) < 2 {
			fmt.Printf("Error: %s requires an argument.\n", commandName)
			fmt.Printf("Usage: %s <argument>\n", commandName)
			return
		}
		arg = words[1]
	}

	// Execute command
	if err := cmd.Callback(r.app, arg); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// cleanInput normalizes user input
func (r *REPL) CleanInput(text string) []string {
	// Convert to lowercase and split into words
	normalized := strings.ToLower(strings.TrimSpace(text))
	words := strings.Fields(normalized)

	// Filter out empty strings
	result := make([]string, 0, len(words))
	for _, word := range words {
		if word != "" {
			result = append(result, word)
		}
	}

	return result
}

// Stop terminates the REPL loop
func (r *REPL) Stop() {
	r.running = false
}
