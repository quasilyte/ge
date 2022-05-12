package ge

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/internal/primitives"
)

type Line struct {
	BeginPos *gemath.Vec
	EndPos   *gemath.Vec

	ColorScale ColorScale

	Width float64

	Visible bool

	disposed bool
}

func NewLine(begin, end *gemath.Vec) *Line {
	return &Line{
		Visible:    true,
		ColorScale: defaultColorScale,
		BeginPos:   begin,
		EndPos:     end,
		Width:      1,
	}
}

func (l *Line) IsDisposed() bool {
	return l.disposed
}

func (l *Line) Dispose() {
	l.disposed = true
}

func (l *Line) Draw(screen *ebiten.Image) {
	if !l.Visible {
		return
	}

	x1 := l.BeginPos.X
	y1 := l.BeginPos.Y
	x2 := l.EndPos.X
	y2 := l.EndPos.Y

	length := math.Hypot(x2-x1, y2-y1)

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Scale(length, l.Width)
	drawOptions.GeoM.Rotate(math.Atan2(y2-y1, x2-x1))
	drawOptions.GeoM.Translate(x1, y1)

	if l.ColorScale != defaultColorScale {
		r := float64(l.ColorScale.R)
		g := float64(l.ColorScale.G)
		b := float64(l.ColorScale.B)
		a := float64(l.ColorScale.A)
		drawOptions.ColorM.Scale(r, g, b, a)
	}

	screen.DrawImage(primitives.WhitePixel, &drawOptions)
}
