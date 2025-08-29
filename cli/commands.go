package cli

import (
	"fmt"
	"os"

	"github.com/mcoluomo/pokedexcli/api"
	"github.com/mcoluomo/pokedexcli/location"
	"github.com/mcoluomo/pokedexcli/pokemon"
)

// App contains all the application services
type App struct {
	client          *api.Client
	pokedex         *pokemon.Pokedex
	catchService    *pokemon.CatchService
	locationService *location.LocationService
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{
		client:          api.NewClient(),
		pokedex:         pokemon.NewPokedex(),
		catchService:    pokemon.NewCatchService(),
		locationService: location.NewLocationService(),
	}
}

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	RequiresArg bool
	Callback    func(*App, string) error
}

// GetCommands returns all available commands
func GetCommands() map[string]Command {
	return map[string]Command{
		"help": {
			Name:        "help",
			Description: "Display this help message",
			RequiresArg: false,
			Callback:    (*App).HelpCommand,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			RequiresArg: false,
			Callback:    (*App).ExitCommand,
		},
		"map": {
			Name:        "map",
			Description: "Display the next 20 location areas",
			RequiresArg: false,
			Callback:    (*App).MapCommand,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Display the previous 20 location areas",
			RequiresArg: false,
			Callback:    (*App).MapBackCommand,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore a location area to see Pokemon",
			RequiresArg: true,
			Callback:    (*App).ExploreCommand,
		},
		"catch": {
			Name:        "catch",
			Description: "Attempt to catch a Pokemon",
			RequiresArg: true,
			Callback:    (*App).CatchCommand,
		},
		"inspect": {
			Name:        "inspect",
			Description: "View details of a caught Pokemon",
			RequiresArg: true,
			Callback:    (*App).InspectCommand,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "List all caught Pokemon",
			RequiresArg: false,
			Callback:    (*App).PokedexCommand,
		},
	}
}

// HelpCommand displays help information
func (app *App) HelpCommand(args string) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Printf("Usage:\n")

	commands := GetCommands()
	for _, cmd := range commands {
		if cmd.RequiresArg {
			fmt.Printf("%s <arg>: %s\n", cmd.Name, cmd.Description)
		} else {
			fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
		}
	}
	fmt.Println()

	return nil
}

// ExitCommand exits the application
func (app *App) ExitCommand(args string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// MapCommand displays location areas
func (app *App) MapCommand(args string) error {
	if !app.client.HasNext() {
		fmt.Println("You're on the last page!")
		return nil
	}

	areas, err := app.client.GetLocationAreas()
	if err != nil {
		return fmt.Errorf("failed to get location areas: %w", err)
	}

	app.locationService.DisplayAreas(areas)
	return nil
}

// MapBackCommand displays previous location areas
func (app *App) MapBackCommand(args string) error {
	if !app.client.HasPrev() {
		fmt.Println("You're on the first page!")
		return nil
	}

	areas, err := app.client.GetPreviousLocationAreas()
	if err != nil {
		return fmt.Errorf("failed to get previous location areas: %w", err)
	}

	app.locationService.DisplayAreas(areas)
	return nil
}

// ExploreCommand explores a specific location
func (app *App) ExploreCommand(areaName string) error {
	if areaName == "" {
		return fmt.Errorf("please provide an area name to explore")
	}

	area, err := app.client.ExploreLocation(areaName)
	if err != nil {
		return fmt.Errorf("failed to explore %s: %w", areaName, err)
	}

	app.locationService.ExploreArea(area)
	return nil
}

// CatchCommand attempts to catch a Pokemon
func (app *App) CatchCommand(pokemonName string) error {
	if pokemonName == "" {
		return fmt.Errorf("please provide a Pokemon name to catch")
	}

	// Check if already caught
	if app.pokedex.HasCaught(pokemonName) {
		fmt.Printf("You have already caught %s! Use 'inspect %s' to view it.\n", pokemonName, pokemonName)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// Get Pokemon from API
	p, err := app.client.GetPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("couldn't find Pokemon %s: %w", pokemonName, err)
	}

	// Use domain logic to attempt catch
	caught, rate := app.catchService.AttemptCatch(p)

	if caught {
		app.pokedex.Catch(p)
		fmt.Printf("%s was caught! (catch rate: %.2f)\n", pokemonName, rate)
		fmt.Printf("You may now inspect it with the 'inspect %s' command.\n", pokemonName)
	} else {
		fmt.Printf("%s escaped! (catch rate: %.2f)\n", pokemonName, rate)
		fmt.Printf("Try again with the 'catch %s' command.\n", pokemonName)
	}

	return nil
}

// InspectCommand displays details of a caught Pokemon
func (app *App) InspectCommand(pokemonName string) error {
	if pokemonName == "" {
		return fmt.Errorf("please provide a Pokemon name to inspect")
	}

	p, exists := app.pokedex.Get(pokemonName)
	if !exists {
		return fmt.Errorf("you haven't caught %s yet! Use 'catch %s' to catch it first", pokemonName, pokemonName)
	}

	fmt.Printf("\n=== %s ===\n", p.Name)
	fmt.Print(p.String())
	fmt.Println()

	return nil
}

// PokedexCommand lists all caught Pokemon
func (app *App) PokedexCommand(args string) error {
	fmt.Printf("\n=== Your Pokedex (%d Pokemon) ===\n", app.pokedex.Count())
	app.pokedex.ListAll()
	fmt.Println()

	return nil
}
