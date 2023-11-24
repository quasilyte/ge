package ge

import (
	"math"
)

func inferDisplayRatio(layoutWidth, layoutHeight int) (int, int) {
	width := layoutWidth
	height := layoutHeight
	if height > width {
		// Assume a horizontal layout even in this case.
		// Resolutions like 1536x2048 become 2048x1536.
		width, height = height, width
	}

	// Try sizes like [w,h], [w/2,h/2], etc.
	// This loop can infer the exact size that will fit perfectly.
	for i := 1; i <= 6; i++ {
		w := width / i
		h := height / i

		switch [2]int{w, h} {
		case [2]int{512, 384}, [2]int{160, 120}:
			// This is actually a 4:3, but I like to have 9 in there for consistency.
			return 12, 9

		case [2]int{480, 270}:
			return 16, 9

		case [2]int{540, 270}:
			return 18, 9

		case [2]int{570, 270}:
			return 19, 9

		case [2]int{600, 270}:
			return 20, 9

		case [2]int{630, 270}:
			return 21, 9

		case [2]int{640, 400}, [2]int{840, 525}:
			return 16, 10
		}
	}

	switch ratio := float64(width) / float64(height); {
	case math.Abs(float64(ratio-(18.0/9.0))) <= 0.07:
		return 18, 9
	case math.Abs(float64(ratio-(19.0/9.0))) <= 0.07:
		return 19, 9
	case math.Abs(float64(ratio-(20.0/9.0))) <= 0.07:
		return 20, 9
	case math.Abs(float64(ratio-(21.0/9.0))) <= 0.07:
		return 21, 9
	case ratio > (22.0 / 9.0):
		// 21:9 is our limit here.
		return 21, 9
	}

	return 16, 9
}
