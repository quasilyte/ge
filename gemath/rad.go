package gemath

import (
	"math"
)

type Rad float64

func (r Rad) Cos() float64 {
	return math.Cos(float64(r))
}

func (r Rad) Sin() float64 {
	return math.Sin(float64(r))
}

func DegToRad(deg float64) Rad {
	return Rad(deg * (math.Pi / 180))
}
