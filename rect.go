package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/internal/primitives"
	"github.com/quasilyte/gmath"
)

var borderBoxIndices = []uint16{0, 2, 4, 2, 4, 6, 1, 3, 5, 3, 5, 7}

type Rect struct {
	Pos Pos

	Width  float64
	Height float64

	Centered bool

	OutlineColorScale ColorScale
	FillColorScale    ColorScale

	OutlineWidth float64

	outlineVertices *[8]ebiten.Vertex

	Visible bool

	disposed bool
}

func NewRect(ctx *Context, width, height float64) *Rect {
	return &Rect{
		Visible:           true,
		Centered:          true,
		FillColorScale:    defaultColorScale,
		OutlineColorScale: transparentColor,
		Width:             width,
		Height:            height,
	}
}

func (rect *Rect) BoundsRect() gmath.Rect {
	pos := rect.Pos.Resolve()
	if rect.Centered {
		offset := gmath.Vec{X: rect.Width * 0.5, Y: rect.Height * 0.5}
		return gmath.Rect{
			Min: pos.Sub(offset),
			Max: pos.Add(offset),
		}
	}
	return gmath.Rect{
		Min: pos,
		Max: pos.Add(gmath.Vec{X: rect.Width, Y: rect.Height}),
	}
}

func (rect *Rect) IsDisposed() bool {
	return rect.disposed
}

func (rect *Rect) Dispose() {
	rect.disposed = true
}

// AnchorPos returns a top-left position.
// When Centered is false, it's identical to Pos, otherwise
// it will apply the computations to get the right anchor for the centered rect.
func (rect *Rect) AnchorPos() Pos {
	if rect.Centered {
		return rect.Pos.WithOffset(-rect.Width/2, -rect.Height/2)
	}
	return rect.Pos
}

func (rect *Rect) calculateFinalOffset(offset gmath.Vec) gmath.Vec {
	var origin gmath.Vec
	if rect.Centered {
		origin = gmath.Vec{X: rect.Width * 0.5, Y: rect.Height * 0.5}
	}
	origin = origin.Sub(rect.Pos.Offset)

	var pos gmath.Vec
	if rect.Pos.Base != nil {
		pos = rect.Pos.Base.Sub(origin)
	} else if !origin.IsZero() {
		pos = gmath.Vec{X: -origin.X, Y: -origin.Y}
	}
	return pos.Add(offset)
}

func (rect *Rect) calculateGeom(w, h float64, pos gmath.Vec) ebiten.GeoM {
	var geom ebiten.GeoM
	geom.Scale(w, h)
	geom.Translate(pos.X, pos.Y)
	return geom
}

func (rect *Rect) Draw(screen *ebiten.Image) {
	rect.DrawWithOffset(screen, gmath.Vec{})
}

func (rect *Rect) DrawWithOffset(screen *ebiten.Image, offset gmath.Vec) {
	if !rect.Visible {
		return
	}
	if rect.OutlineColorScale.A == 0 && rect.FillColorScale.A == 0 {
		return
	}

	// TODO: compare the peformance of this method with vector package.
	// TODO: implement the rotation.
	// TODO: implement the scaling.
	// TODO: maybe add a special case for opaque rectangles.

	finalOffset := rect.calculateFinalOffset(offset)

	if rect.OutlineColorScale.A == 0 || rect.OutlineWidth < 1 {
		// Fill-only mode.
		var drawOptions ebiten.DrawImageOptions
		drawOptions.GeoM = rect.calculateGeom(rect.Width, rect.Height, finalOffset)
		drawOptions.ColorScale = rect.FillColorScale.toEbitenColorScale()
		screen.DrawImage(primitives.WhitePixel, &drawOptions)
		return
	}

	if rect.FillColorScale.A == 0 && rect.OutlineWidth >= 1 {
		// Outline-only mode.
		rect.drawOutline(screen, finalOffset)
		return
	}

	rect.drawOutline(screen, finalOffset)

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Scale(rect.Width-rect.OutlineWidth*2, rect.Height-rect.OutlineWidth*2)
	drawOptions.GeoM.Translate(rect.OutlineWidth+finalOffset.X, rect.OutlineWidth+finalOffset.Y)
	drawOptions.ColorScale = rect.FillColorScale.toEbitenColorScale()
	screen.DrawImage(primitives.WhitePixel, &drawOptions)
}

func (rect *Rect) drawOutline(screen *ebiten.Image, offset gmath.Vec) {
	if rect.outlineVertices == nil {
		// Allocate these vertices lazily when we need them and then re-use them.
		rect.outlineVertices = new([8]ebiten.Vertex)
	}

	borderWidth := float32(rect.OutlineWidth)
	x := float32(offset.X)
	y := float32(offset.Y)
	r := float32(rect.OutlineColorScale.R)
	g := float32(rect.OutlineColorScale.G)
	b := float32(rect.OutlineColorScale.B)
	a := float32(rect.OutlineColorScale.A)
	width := float32(rect.Width)
	height := float32(rect.Height)

	rect.outlineVertices[0] = ebiten.Vertex{
		DstX:   x,
		DstY:   y,
		SrcX:   0,
		SrcY:   0,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	rect.outlineVertices[1] = ebiten.Vertex{
		DstX: x + borderWidth,
		DstY: y + borderWidth,
		SrcX: 0,
		SrcY: 0,
	}
	rect.outlineVertices[2] = ebiten.Vertex{
		DstX:   x + width,
		DstY:   y,
		SrcX:   1,
		SrcY:   0,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	rect.outlineVertices[3] = ebiten.Vertex{
		DstX: x + width - borderWidth,
		DstY: y + borderWidth,
		SrcX: 1,
		SrcY: 0,
	}
	rect.outlineVertices[4] = ebiten.Vertex{
		DstX:   x,
		DstY:   y + height,
		SrcX:   0,
		SrcY:   1,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	rect.outlineVertices[5] = ebiten.Vertex{
		DstX: x + borderWidth,
		DstY: y + height - borderWidth,
		SrcX: 0,
		SrcY: 1,
	}
	rect.outlineVertices[6] = ebiten.Vertex{
		DstX:   x + width,
		DstY:   y + height,
		SrcX:   1,
		SrcY:   1,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	rect.outlineVertices[7] = ebiten.Vertex{
		DstX: x + width - borderWidth,
		DstY: y + height - borderWidth,
		SrcX: 1,
		SrcY: 1,
	}

	options := ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	screen.DrawTriangles(rect.outlineVertices[:], borderBoxIndices, primitives.WhitePixel, &options)
}
