package main

import (
	"image/color"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

func getPlayerTextColor(playerID int) color.RGBA {
	switch playerID {
	case 0:
		return ge.RGB(0x629692)
	case 1:
		return ge.RGB(0x7585b0)
	case 2:
		return ge.RGB(0x9a696d)
	default:
		return ge.RGB(0xa6719e)
	}
}

func applyPlayerColor(playerID int, s *ge.Sprite) {
	switch playerID {
	case 1:
		s.Hue = gemath.DegToRad(50)
	case 2:
		s.Hue = -gemath.DegToRad(180)
	case 3:
		s.Hue = gemath.DegToRad(135)
	}
}
