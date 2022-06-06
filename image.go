package ge

import "github.com/hajimehoshi/ebiten/v2"

type Image struct {
	Texture *ebiten.Image

	DefaultFrameWidth  float64
	DefaultFrameHeight float64
}
