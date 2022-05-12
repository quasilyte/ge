package gemath

import (
	"math"
)

// Rad represents a radian value.
// It's not capped in [0, 2*Pi] range.
//
// In terms of the orientations, Pi rotation points the object down (South).
// Zero radians point towards the right side (East).
type Rad float64

func DegToRad(deg float64) Rad {
	return Rad(deg * (math.Pi / 180))
}

// Positive returns the equivalent radian value expressed as a positive value.
func (r Rad) Positive() Rad {
	if r >= 0 {
		return r
	}
	return r + 2*math.Pi
}

// Normalized returns the equivalent radians value in [0, 2*Pi] range.
// For example, 3*Pi becomes just Pi.
func (r Rad) Normalized() Rad {
	angle := float64(r)
	angle -= math.Floor(angle/(2*math.Pi)) * 2 * math.Pi
	return Rad(angle)
}

// EqualApprox compares two radian values using EqualApprox function.
// Note that you may want to normalize the operands in some way before doing this.
func (r Rad) EqualApprox(other Rad) bool {
	return EqualApprox(float64(r), float64(other))
}

// AngleDelta returns an angle delta between two radian values.
// The sign is preserved.
//
// It doesn't need the angles to be normalized,
// r=0 and r=2*Pi are considered to have no delta.
func (r Rad) AngleDelta(r2 Rad) Rad {
	angle1 := math.Mod(float64(r-r2), 2*math.Pi)
	angle2 := math.Mod(float64(r2-r), 2*math.Pi)
	if angle1 < angle2 {
		return Rad(-angle1)
	}
	return Rad(angle2)
}

func (r Rad) LerpAngle(toAngle Rad, weight float64) Rad {
	difference := math.Mod(float64(toAngle)-float64(r), 2*math.Pi)
	dist := math.Mod(2.0*difference, 2*math.Pi) - difference
	return Rad(float64(r) + dist*weight)
}

func (r Rad) RotatedTowards(toAngle, amount Rad) Rad {
	difference := math.Mod(float64(toAngle)-float64(r), 2*math.Pi)
	dist := math.Mod(2.0*difference, 2*math.Pi) - difference
	if EqualApprox(dist, 0) {
		return toAngle
	}
	lerpa1 := Rad(float64(r) + dist)
	if min := r - amount; lerpa1 < min {
		return min
	}
	if max := r + amount; lerpa1 > max {
		return max
	}
	return lerpa1
}

func (r Rad) Cos() float64 {
	return math.Cos(float64(r))
}

func (r Rad) Sin() float64 {
	return math.Sin(float64(r))
}
