package ge

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/internal/primitives"
	"github.com/quasilyte/gmath"
)

type Line struct {
	BeginPos Pos
	EndPos   Pos

	Width float64

	colorScale    ColorScale
	colorM        ebiten.ColorM
	colorsChanged bool

	Visible bool

	disposed bool
}

func NewLine(begin, end Pos) *Line {
	return &Line{
		Visible:    true,
		colorScale: defaultColorScale,
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
	if l.colorsChanged {
		l.colorsChanged = false
		l.recalculateColorM()
	}
	pos1 := l.BeginPos.Resolve()
	pos2 := l.EndPos.Resolve()
	drawLine(screen, pos1, pos2, l.Width, l.colorM)
}

func (l *Line) SetColorScaleRGBA(r, g, b, a uint8) {
	var scale ColorScale
	scale.SetRGBA(r, g, b, a)
	l.SetColorScale(scale)
}

func (l *Line) GetAlpha() float32 {
	return l.colorScale.A
}

func (l *Line) SetAlpha(a float32) {
	if l.colorScale.A == a {
		return
	}
	l.colorScale.A = a
	l.colorsChanged = true
}

func (l *Line) SetColorScale(colorScale ColorScale) {
	if l.colorScale == colorScale {
		return
	}
	l.colorScale = colorScale
	l.colorsChanged = true
}

func (l *Line) recalculateColorM() {
	var colorM ebiten.ColorM
	if l.colorScale != defaultColorScale {
		colorM.Scale(float64(l.colorScale.R), float64(l.colorScale.G), float64(l.colorScale.B), float64(l.colorScale.A))
	}
	l.colorM = colorM
}

func drawLine(dst *ebiten.Image, pos1, pos2 gmath.Vec, width float64, colorM ebiten.ColorM) {
	x1 := pos1.X
	y1 := pos1.Y
	x2 := pos2.X
	y2 := pos2.Y

	length := math.Hypot(x2-x1, y2-y1)

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Scale(length, width)
	drawOptions.GeoM.Rotate(math.Atan2(y2-y1, x2-x1))
	drawOptions.GeoM.Translate(x1, y1)

	drawOptions.ColorM = colorM

	dst.DrawImage(primitives.WhitePixel, &drawOptions)
}
