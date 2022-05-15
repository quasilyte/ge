package primitives

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var emptyImage = ebiten.NewImage(3, 3)

var WhitePixel = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

func init() {
	emptyImage.Fill(color.White)
}

func DrawLine(dst *ebiten.Image, x1, y1, x2, y2 float64, clr color.Color) {
	length := math.Hypot(x2-x1, y2-y1)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(length, 1)
	op.GeoM.Rotate(math.Atan2(y2-y1, x2-x1))
	op.GeoM.Translate(x1, y1)
	op.ColorM.ScaleWithColor(clr)
	dst.DrawImage(WhitePixel, op)
}
