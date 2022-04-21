package gemath

import "math"

const Epsilon = 1e-9

func EqualApprox(a, b float64) bool {
	return math.Abs(a-b) <= Epsilon
}

func ClampMin(v, min float64) float64 {
	if v < min {
		return min
	}
	return v
}

func ClampMax(v, max float64) float64 {
	if v > max {
		return max
	}
	return v
}
