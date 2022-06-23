package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func applyColorScale(c ColorScale, opts *ebiten.DrawImageOptions) {
	if c == defaultColorScale {
		return
	}
	r := float64(c.R)
	g := float64(c.G)
	b := float64(c.B)
	a := float64(c.A)
	opts.ColorM.Scale(r, g, b, a)
}

func assignColors(vertices []ebiten.Vertex, c ColorScale) {
	colorR := float32(c.R)
	colorG := float32(c.G)
	colorB := float32(c.B)
	colorA := float32(c.A)
	for i := range vertices {
		v := &vertices[i]
		v.ColorR = colorR
		v.ColorG = colorG
		v.ColorB = colorB
		v.ColorA = colorA
	}
}
