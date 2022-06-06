package gemath

import (
	"math"
)

const Epsilon = 1e-9

func EqualApprox(a, b float64) bool {
	return math.Abs(a-b) <= Epsilon
}

func ArcContains(angle, measure Rad, r float64, pos, point Vec) bool {
	if pos.DistanceSquaredTo(point) > r*r {
		return false
	}
	startAngle := (angle - measure/2)
	endAngle := (angle + measure/2)
	if endAngle < startAngle {
		endAngle += 2 * math.Pi
	}
	half := (endAngle - startAngle) / 2
	mid := (endAngle + startAngle) / 2
	coshalf := math.Cos(float64(half))
	angleToPoint := pos.AngleToPoint(point).Normalized()
	return math.Cos(float64(angleToPoint-mid)) >= coshalf
}

func Clamp[T numeric](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func ClampMin[T numeric](v, min T) T {
	if v < min {
		return min
	}
	return v
}

func ClampMax[T numeric](v, max T) T {
	if v > max {
		return max
	}
	return v
}

func Percentage[T numeric](value, max T) T {
	if max == 0 && value == 0 {
		return 0
	}
	return T(100 * (float64(value) / float64(max)))
}

type numeric interface {
	int | float64
}
