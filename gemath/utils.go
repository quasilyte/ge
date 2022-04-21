package gemath

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
