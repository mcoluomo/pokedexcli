package pokemon

import "fmt"

// Pokedex is our aggregate root for managing caught Pokemon
type Pokedex struct {
	caught map[string]Pokemon
}

// NewPokedex creates a new Pokedex instance
func NewPokedex() *Pokedex {
	return &Pokedex{
		caught: make(map[string]Pokemon),
	}
}

// Catch adds a Pokemon to the Pokedex
func (pd *Pokedex) Catch(p Pokemon) {
	pd.caught[p.Name] = p
}

// HasCaught checks if a Pokemon has been caught
func (pd *Pokedex) HasCaught(name string) bool {
	_, exists := pd.caught[name]
	return exists
}

// Get retrieves a caught Pokemon by name
func (pd *Pokedex) Get(name string) (Pokemon, bool) {
	p, exists := pd.caught[name]
	return p, exists
}

// ListAll displays all caught Pokemon
func (pd *Pokedex) ListAll() {
	if len(pd.caught) == 0 {
		fmt.Println("You haven't caught any Pokemon yet!")
		return
	}

	fmt.Println("Your Pokemon:")
	for name := range pd.caught {
		fmt.Printf("  - %s\n", name)
	}
}

// Count returns the number of caught Pokemon
func (pd *Pokedex) Count() int {
	return len(pd.caught)
}

// Release removes a Pokemon from the Pokedex
func (pd *Pokedex) Release(name string) bool {
	if !pd.HasCaught(name) {
		return false
	}
	delete(pd.caught, name)
	return true
}
