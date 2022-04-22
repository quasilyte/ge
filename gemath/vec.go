package gemath

import (
	"fmt"
	"math"
)

type Vec struct {
	X float64
	Y float64
}

// RadToVec converts a given angle into a normalized vector that encodes that direction.
func RadToVec(angle Rad) Vec {
	return Vec{X: angle.Cos(), Y: angle.Sin()}
}

func (v Vec) String() string {
	return fmt.Sprintf("[%f, %f]", v.X, v.Y)
}

func (v Vec) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vec) IsNormalized() bool {
	return EqualApprox(v.LenSquared(), 1)
}

func (v Vec) DistanceTo(v2 Vec) float64 {
	return math.Sqrt((v.X-v2.X)*(v.X-v2.X) + (v.Y-v2.Y)*(v.Y-v2.Y))
}

func (v Vec) Dot(v2 Vec) float64 {
	return (v.X * v2.X) + (v.Y * v2.Y)
}

func (v Vec) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vec) LenSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec) Angle() Rad {
	return Rad(math.Atan2(v.Y, v.X))
}

func (v Vec) AngleTo(pos Vec) Rad {
	return v.SubResult(pos).Angle()
}

func (v Vec) VecTowards(length float64, pos Vec) Vec {
	angle := pos.AngleTo(v)
	result := Vec{X: angle.Cos(), Y: angle.Sin()}
	result.MulScalar(length)
	return result
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

func (v Vec) MulScalarResult(scalar float64) Vec {
	v.MulScalar(scalar)
	return v
}

func (v *Vec) Mul(other Vec) {
	v.X *= other.X
	v.Y *= other.Y
}

func (v *Vec) DivScalar(scalar float64) {
	v.X /= scalar
	v.Y /= scalar
}

func (v Vec) DivScalarResult(scalar float64) Vec {
	v.DivScalar(scalar)
	return v
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

func (v *Vec) Sub(other Vec) {
	v.X -= other.X
	v.Y -= other.Y
}

func (v Vec) SubResult(other Vec) Vec {
	v.Sub(other)
	return v
}

func (v *Vec) Normalize() {
	l := v.LenSquared()
	if l != 0 {
		v.DivScalar(math.Sqrt(l))
	}
}

func (v Vec) NormalizeResult() Vec {
	v.Normalize()
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
