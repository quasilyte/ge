package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/internal/primitives"
)

type Rect struct {
	Pos      Pos
	Rotation *gemath.Rad

	Width  float64
	Height float64

	Scale float64

	Centered bool

	ColorScale ColorScale

	Visible bool

	disposed bool
}

func NewRect(width, height float64) *Rect {
	return &Rect{
		Visible:    true,
		Centered:   true,
		ColorScale: defaultColorScale,
		Scale:      1,
		Width:      width,
		Height:     height,
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

func (rect *Rect) Draw(screen *ebiten.Image) {
	if !rect.Visible {
		return
	}

	var origin gemath.Vec
	if rect.Centered {
		origin = gemath.Vec{X: rect.Width / 2, Y: rect.Height / 2}
	}
	origin = origin.Sub(rect.Pos.Offset)

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Scale(rect.Width, rect.Height)

	drawOptions.GeoM.Translate(-origin.X, -origin.Y)
	if rect.Rotation != nil {
		drawOptions.GeoM.Rotate(float64(*rect.Rotation))
	}
	if rect.Scale != 1 {
		drawOptions.GeoM.Scale(rect.Scale, rect.Scale)
	}
	drawOptions.GeoM.Translate(origin.X, origin.Y)

	if rect.Pos.Base != nil {
		drawOptions.GeoM.Translate(rect.Pos.Base.X-origin.X, rect.Pos.Base.Y-origin.Y)
	} else if !origin.IsZero() {
		drawOptions.GeoM.Translate(0-origin.X, 0-origin.Y)
	}

	if rect.ColorScale != defaultColorScale {
		r := float64(rect.ColorScale.R)
		g := float64(rect.ColorScale.G)
		b := float64(rect.ColorScale.B)
		a := float64(rect.ColorScale.A)
		drawOptions.ColorM.Scale(r, g, b, a)
	}

	screen.DrawImage(primitives.WhitePixel, &drawOptions)
}
