package gemath

import (
	"math"
	"math/rand"
)

type Rand struct {
	rng *rand.Rand
}

func (r *Rand) SetSeed(seed int64) {
	src := rand.NewSource(seed)
	r.rng = rand.New(src)
}

func (r *Rand) Chance(probability float64) bool {
	return r.rng.Float64() <= probability
}

func (r *Rand) IntRange(min, max int) int {
	return min + r.rng.Intn(max-min+1)
}

func (r *Rand) FloatRange(min, max float64) float64 {
	return min + r.rng.Float64()*(max-min+1)
}

func (r *Rand) Rad() Rad {
	return Rad(r.FloatRange(0, 2*math.Pi))
}
