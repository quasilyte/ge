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
	drawLine(screen, *l.BeginPos, *l.EndPos, l.Width, l.ColorScale)
}

func drawLine(dst *ebiten.Image, pos1, pos2 gemath.Vec, width float64, c ColorScale) {
	x1 := pos1.X
	y1 := pos1.Y
	x2 := pos2.X
	y2 := pos2.Y

	length := math.Hypot(x2-x1, y2-y1)

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Scale(length, width)
	drawOptions.GeoM.Rotate(math.Atan2(y2-y1, x2-x1))
	drawOptions.GeoM.Translate(x1, y1)

	applyColorScale(c, &drawOptions)

	dst.DrawImage(primitives.WhitePixel, &drawOptions)
}
