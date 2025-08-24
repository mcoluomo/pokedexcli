package pokeapi

import (
	"math/rand"
	"time"
)

var Rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func SimpleCatch(captureRate int, ballBonus float64) bool {
	adjusted := float64(captureRate) * ballBonus

	chance := adjusted / 255.0

	return Rng.Float64() < chance
}

func RandomBonus() float64 {
	rolls := []float64{0.8, 1.0, 1.2}
	return rolls[Rng.Intn(len(rolls))]
}
