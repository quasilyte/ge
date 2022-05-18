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

func (r *Rand) Offset(min, max float64) Vec {
	return Vec{X: r.FloatRange(min, max), Y: r.FloatRange(min, max)}
}

func (r *Rand) Chance(probability float64) bool {
	return r.rng.Float64() <= probability
}

func (r *Rand) Bool() bool {
	return r.rng.Float64() < 0.5
}

func (r *Rand) IntRange(min, max int) int {
	return min + r.rng.Intn(max-min+1)
}

func (r *Rand) Float() float64 {
	return r.rng.Float64()
}

func (r *Rand) FloatRange(min, max float64) float64 {
	return min + r.rng.Float64()*(max-min)
}

func (r *Rand) Rad() Rad {
	return Rad(r.FloatRange(0, 2*math.Pi))
}

func RandIndex[T any](r *Rand, slice []T) int {
	if len(slice) == 0 {
		return -1
	}
	return r.IntRange(0, len(slice)-1)
}

func RandElem[T any](r *Rand, slice []T) (elem T) {
	if len(slice) == 0 {
		return elem // Zero value
	}
	if len(slice) == 1 {
		return slice[0]
	}
	return slice[r.rng.Intn(len(slice))]
}
