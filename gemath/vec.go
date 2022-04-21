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

func (v Vec) EqualApprox(other Vec) bool {
	return EqualApprox(v.X, other.X) && EqualApprox(v.Y, other.Y)
}

func (v *Vec) MoveInDirection(dist float64, dir Rad) {
	v.X += dist * dir.Cos()
	v.Y += dist * dir.Sin()
}

func (v Vec) MoveInDirectionResult(dist float64, dir Rad) Vec {
	v.MoveInDirection(dist, dir)
	return v
}

func (v *Vec) MulScalar(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vec) Mul(other Vec) {
	v.X *= other.X
	v.Y *= other.Y
}

func (v *Vec) DivScalar(scalar float64) {
	v.X /= scalar
	v.Y /= scalar
}

func (v *Vec) Div(other Vec) {
	v.X /= other.X
	v.Y /= other.Y
}

func (v *Vec) Add(other Vec) {
	v.X += other.X
	v.Y += other.Y
}

func (v Vec) AddResult(other Vec) Vec {
	v.Add(other)
	return v
}

func (v *Vec) Neg() {
	v.X = -v.X
	v.Y = -v.Y
}

func (v Vec) NegResult() Vec {
	v.Neg()
	return v
}
