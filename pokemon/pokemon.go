package pokemon

import (
	"fmt"
	"strings"
)

// Pokemon represents our core domain entity
type Pokemon struct {
	Name           string
	Height         int
	Weight         int
	BaseExperience int
	Types          []string
	Stats          map[string]int // hp, attack, defense, etc.
}

// Business rules for Pokemon
func (p Pokemon) String() string {
	return fmt.Sprintf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n%s\nType: %s",
		p.Name, p.Height, p.Weight, p.formatStats(), p.formatTypes())
}

func (p Pokemon) formatStats() string {
	var result strings.Builder
	for stat, value := range p.Stats {
		result.WriteString(fmt.Sprintf("  -%s: %d\n", stat, value))
	}
	return result.String()
}

func (p Pokemon) formatTypes() string {
	return strings.Join(p.Types, ", ")
}

// IsCatchable determines if a Pokemon can be caught based on business rules
func (p Pokemon) IsCatchable() bool {
	return p.BaseExperience > 0 // Simple business rule
}

// CatchDifficulty returns how hard this Pokemon is to catch (0.0 to 1.0)
func (p Pokemon) CatchDifficulty() float64 {
	if p.BaseExperience <= 0 {
		return 0.0
	}
	difficulty := float64(p.BaseExperience) / 255.0
	if difficulty > 1.0 {
		return 1.0
	}
	return difficulty
}

// IsLegendary checks if this Pokemon is likely legendary based on base experience
func (p Pokemon) IsLegendary() bool {
	return p.BaseExperience > 200
}
