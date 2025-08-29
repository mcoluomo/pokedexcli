package pokemon

import (
	"math/rand"
	"time"
)

// CatchService handles the catching business logic
type CatchService struct {
	rng *rand.Rand
}

// NewCatchService creates a new catch service
func NewCatchService() *CatchService {
	return &CatchService{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// AttemptCatch tries to catch a Pokemon using domain business rules
// Returns (caught, catchRate)
func (cs *CatchService) AttemptCatch(p Pokemon) (bool, float64) {
	if !p.IsCatchable() {
		return false, 0.0
	}

	// Base catch rate calculation
	difficulty := p.CatchDifficulty()

	// Apply bonuses
	bonus := cs.calculateBonus(p)

	// Final catch rate (higher is better for player)
	catchRate := (1.0 - difficulty) * bonus

	// Ensure catch rate is between 0 and 1
	if catchRate < 0.0 {
		catchRate = 0.0
	}
	if catchRate > 1.0 {
		catchRate = 1.0
	}

	// Roll for catch
	roll := cs.rng.Float64()
	caught := roll < catchRate

	return caught, catchRate
}

// calculateBonus applies various bonuses to catch rate
func (cs *CatchService) calculateBonus(p Pokemon) float64 {
	// Base bonus between 0.5 and 0.8
	bonus := 0.5 + cs.rng.Float64()*0.3

	// Legendary Pokemon are harder to catch
	if p.IsLegendary() {
		bonus *= 0.5
	}

	// Small chance of critical catch (extra bonus)
	if cs.rng.Float64() < 0.05 { // 5% chance
		bonus *= 1.5
	}

	return bonus
}

// GetCatchPreview calculates catch probability without attempting
func (cs *CatchService) GetCatchPreview(p Pokemon) float64 {
	if !p.IsCatchable() {
		return 0.0
	}

	difficulty := p.CatchDifficulty()
	avgBonus := 0.65 // Average bonus

	if p.IsLegendary() {
		avgBonus *= 0.5
	}

	catchRate := (1.0 - difficulty) * avgBonus

	if catchRate < 0.0 {
		return 0.0
	}
	if catchRate > 1.0 {
		return 1.0
	}

	return catchRate
}
