package main

import (
	"image/color"
	"math"

	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
)

func rgbaToHSL(c ge.ColorScale) (h, s, l gmath.Rad) {
	r := float64(c.R)
	g := float64(c.G)
	b := float64(c.B)
	maxVal := math.Max(math.Max(r, g), b)
	minVal := math.Min(math.Min(r, g), b)
	l = gmath.Rad((maxVal + minVal) / 2)

	if maxVal == minVal {
		h = 0
		s = 0
		return
	}

	d := maxVal - minVal
	if l > 0.5 {
		s = gmath.Rad(d / (2 - maxVal - minVal))
	} else {
		s = gmath.Rad(d / (maxVal + minVal))
	}
	switch maxVal {
	case r:
		h = gmath.Rad((g - b) / d)
		if g < b {
			h += 6
		}
	case g:
		h = gmath.Rad((b-r)/d + 2)
	case b:
		h = gmath.Rad((r-g)/d + 4)
	}
	h /= 6

	return
}

func hslToRGBA(h, s, l gmath.Rad) color.RGBA {
	var r, g, b gmath.Rad
	if s == 0 {
		r, g, b = l, l, l // achromatic
		return color.RGBA{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255), A: 255}
	}

	var q gmath.Rad
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	r = hueToRGB(p, q, h+1/3)
	g = hueToRGB(p, q, h)
	b = hueToRGB(p, q, h-1/3)

	return color.RGBA{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255), A: 255}
}

func hueToRGB(p, q, t gmath.Rad) gmath.Rad {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1/6 {
		return p + (q-p)*6*t
	}
	if t < 1/2 {
		return q
	}
	if t < 2/3 {
		return p + (q-p)*(2/3-t)*6
	}
	return p
}

func SetHue(s *ge.Sprite, nh gmath.Rad) {
	rgba := s.GetColorScale()
	_, oldSat, oldLum := rgbaToHSL(rgba)
	nr, ng, nb, _ := hslToRGBA(nh, oldSat, oldLum).RGBA()
	s.SetColorScaleRGBA(uint8(nr), uint8(ng), uint8(nb), 255)
}
