package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

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
