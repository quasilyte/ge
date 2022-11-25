package ge

import (
	"image/color"

	"github.com/quasilyte/gmath"
)

func NewVec(x, y float64) *gmath.Vec {
	return &gmath.Vec{X: x, Y: y}
}

func NewRotation(deg float64) *gmath.Rad {
	rad := gmath.DegToRad(deg)
	return &rad
}

// RGB returns a color.RGBA created from the bits of rgb value.
// RGB(0xAABBCC) is identical to color.RGBA{R: 0xAA, G: 0xBB, B: 0xCC, A: 0xFF}
func RGB(rgb uint64) color.RGBA {
	return color.RGBA{
		R: uint8((rgb & (0xFF << (8 * 2))) >> (8 * 2)),
		G: uint8((rgb & (0xFF << (8 * 1))) >> (8 * 1)),
		B: uint8((rgb & (0xFF << (8 * 0))) >> (8 * 0)),
		A: 0xFF,
	}
}
