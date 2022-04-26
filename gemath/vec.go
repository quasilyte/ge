package gemath

import (
	"fmt"
	"math"
)

// Vec is a 2-element structure that is used to represent positions,
// velocities, and other kinds numerical pairs.
//
// Its implementation as well as its API is inspired by Vector2 type
// of the Godot game engine. Where feasible, its adjusted to fit Go
// coding conventions better. Also, earlier versions of Godot used
// 32-bit values for X and Y; our vector uses 64-bit values.
//
// Since Go has no operator overloading, we implement scalar forms of
// operations with "f" suffix. So, Add() is used to add two vectors
// while Addf() is used to add scalar to the vector.
type Vec struct {
	X float64
	Y float64
}

// RadToVec converts a given angle into a normalized vector that encodes that direction.
func RadToVec(angle Rad) Vec {
	return Vec{X: angle.Cos(), Y: angle.Sin()}
}

// String returns a pretty-printed representation of a 2D vector object.
func (v Vec) String() string {
	return fmt.Sprintf("[%f, %f]", v.X, v.Y)
}

// IsZero reports whether v is a zero value vector.
// A zero value vector has X=0 and Y=0, created with Vec{}.
//
// The zero value vector has a property that its length is 0,
// but not all zero length vectors are zero value vectors.
func (v Vec) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

// IsNormalizer reports whether the vector is normalized.
// A vector is considered to be normalized if its length is 1.
func (v Vec) IsNormalized() bool {
	return EqualApprox(v.LenSquared(), 1)
}

// DistanceTo calculates the distance between the two vectors.
func (v Vec) DistanceTo(v2 Vec) float64 {
	return math.Sqrt((v.X-v2.X)*(v.X-v2.X) + (v.Y-v2.Y)*(v.Y-v2.Y))
}

func (v Vec) DistanceSquaredTo(v2 Vec) float64 {
	return ((v.X - v2.X) * (v.X - v2.X)) + ((v.Y - v2.Y) * (v.Y - v2.Y))
}

// Dot returns a dot-product of the two vectors.
func (v Vec) Dot(v2 Vec) float64 {
	return (v.X * v2.X) + (v.Y * v2.Y)
}

// Len reports the length of this vector (also known as magnitude).
func (v Vec) Len() float64 {
	return math.Sqrt(v.LenSquared())
}

// LenSquared returns the squared length of this vector.
//
// This function runs faster than Len(),
// so prefer it if you need to compare vectors
// or need the squared distance for some formula.
func (v Vec) LenSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec) Rotated(angle Rad) Vec {
	sine := angle.Sin()
	cosi := angle.Cos()
	return Vec{
		X: v.X*cosi - v.Y*sine,
		Y: v.X*sine + v.Y*cosi,
	}
}

func (v Vec) Angle() Rad {
	return Rad(math.Atan2(v.Y, v.X))
}

// AngleToPoint returns the angle between the line connecting the two points
// and the X axis, in radians.
func (v Vec) AngleToPoint(pos Vec) Rad {
	return v.Sub(pos).Angle()
}

func (v Vec) DirectionTo(v2 Vec) Vec {
	return v.Sub(v2).Normalized()
}

func (v Vec) VecTowards(length float64, pos Vec) Vec {
	angle := pos.AngleToPoint(v)
	result := Vec{X: angle.Cos(), Y: angle.Sin()}
	return result.Mulf(length)
}

func (v Vec) EqualApprox(other Vec) bool {
	return EqualApprox(v.X, other.X) && EqualApprox(v.Y, other.Y)
}

func (v Vec) MoveInDirection(dist float64, dir Rad) Vec {
	return Vec{
		X: v.X + dist*dir.Cos(),
		Y: v.Y + dist*dir.Sin(),
	}
}

func (v Vec) Mulf(scalar float64) Vec {
	return Vec{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func (v Vec) Mul(other Vec) Vec {
	return Vec{
		X: v.X * other.X,
		Y: v.Y * other.Y,
	}
}

func (v Vec) Divf(scalar float64) Vec {
	return Vec{
		X: v.X / scalar,
		Y: v.Y / scalar,
	}
}

func (v Vec) Div(other Vec) Vec {
	return Vec{
		X: v.X / other.X,
		Y: v.Y / other.Y,
	}
}

func (v Vec) Add(other Vec) Vec {
	return Vec{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vec) Sub(other Vec) Vec {
	return Vec{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

// Normalized returns the vector scaled to unit length.
// Functionally equivalent to `v.Divf(v.Len())`.
//
// Special case: for zero value vectors it returns that unchanged.
func (v Vec) Normalized() Vec {
	l := v.LenSquared()
	if l != 0 {
		return v.Divf(math.Sqrt(l))
	}
	return v
}

// Neg applies unary minus (-) to the vector.
func (v Vec) Neg() Vec {
	return Vec{
		X: -v.X,
		Y: -v.Y,
	}
}
