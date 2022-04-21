package gemath

import (
	"fmt"
	"math"
)

type Vec struct {
	X float64
	Y float64
}

func (v Vec) String() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
}

func (v Vec) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vec) DistanceTo(v2 Vec) float64 {
	return math.Sqrt((v.X-v2.X)*(v.X-v2.X) + (v.Y-v2.Y)*(v.Y-v2.Y))
}

func (v Vec) Dot(v2 Vec) float64 {
	return (v.X * v2.X) + (v.Y * v2.Y)
}

func (v *Vec) MoveInDirection(dist float64, dir Rad) {
	v.X += dist * dir.Cos()
	v.Y += dist * dir.Sin()
}

func (v *Vec) MoveInDirectionResult(dist float64, dir Rad) Vec {
	cloned := *v
	cloned.MoveInDirection(dist, dir)
	return cloned
}

func (v *Vec) Add(other Vec) {
	v.X += other.X
	v.Y += other.Y
}
