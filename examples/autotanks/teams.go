package main

import (
	"image/color"

	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
)

func unitLayerMask(alliance int) uint16 {
	return 1 << uint16(alliance)
}

func projectileLayerMask(alliance int) uint16 {
	return 0b1111 ^ unitLayerMask(alliance)
}

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
		SetHue(s, gmath.DegToRad(50))
	case 2:
		SetHue(s, gmath.DegToRad(180))
	case 3:
		SetHue(s, gmath.DegToRad(135))
	}
}
