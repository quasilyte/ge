package primitives

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var emptyImage = ebiten.NewImage(3, 3)

var WhitePixel = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

func init() {
	emptyImage.Fill(color.White)
}
