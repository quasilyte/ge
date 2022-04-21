package gedraw

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quasilyte/ge/gemath"
)

var (
	emptyImage = ebiten.NewImage(3, 3)
)

func init() {
	emptyImage.Fill(color.White)
}

func DrawPath(dst *ebiten.Image, points []gemath.Vec, c color.RGBA) {
	if len(points) == 0 {
		return
	}

	src := emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	var drawOptions ebiten.DrawTrianglesOptions

	var p vector.Path
	p.MoveTo(float32(points[0].X), float32(points[0].Y))
	for _, pt := range points[1:] {
		p.LineTo(float32(pt.X), float32(pt.Y))
	}

	var vertices []ebiten.Vertex
	var indices []uint16
	vertices, indices = p.AppendVerticesAndIndicesForFilling(vertices, indices)
	assignColors(vertices, c)
	dst.DrawTriangles(vertices, indices, src, &drawOptions)
}

func DrawCircle(dst *ebiten.Image, pos gemath.Vec, radius float64, c color.RGBA) {
	src := emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	var drawOptions ebiten.DrawTrianglesOptions

	var p vector.Path
	p.Arc(float32(pos.X), float32(pos.Y), float32(radius), 0, 2*math.Pi, vector.Clockwise)
	var vertices []ebiten.Vertex
	var indices []uint16
	vertices, indices = p.AppendVerticesAndIndicesForFilling(vertices, indices)
	assignColors(vertices, c)
	dst.DrawTriangles(vertices, indices, src, &drawOptions)
}

func assignColors(vertices []ebiten.Vertex, c color.RGBA) {
	colorR := float32(c.R) / 0xff
	colorG := float32(c.G) / 0xff
	colorB := float32(c.B) / 0xff
	colorA := float32(c.A) / 0xff
	for i := range vertices {
		v := &vertices[i]
		v.ColorR = colorR
		v.ColorG = colorG
		v.ColorB = colorB
		v.ColorA = colorA
	}
}
