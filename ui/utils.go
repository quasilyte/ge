package ui

import "github.com/quasilyte/ge"

func withAlpha(c ge.ColorScale, a float32) ge.ColorScale {
	c.A = a
	return c
}
