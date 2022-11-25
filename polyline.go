package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/gmath"
)

type PolyLine struct {
	Points []gmath.Vec

	ColorScale ColorScale

	Width float64

	Visible bool

	disposed bool
}

func NewPolyLine() *PolyLine {
	return &PolyLine{
		Visible:    true,
		ColorScale: defaultColorScale,
		Width:      1,
	}
}

func (l *PolyLine) ResetPoints() {
	l.Points = l.Points[:0]
}

func (l *PolyLine) PushPoint(pt gmath.Vec) {
	l.Points = append(l.Points, pt)
}

func (l *PolyLine) PopPoint() gmath.Vec {
	pt := l.Points[len(l.Points)-1]
	l.Points = l.Points[:len(l.Points)-1]
	return pt
}

func (l *PolyLine) IsDisposed() bool {
	return l.disposed
}

func (l *PolyLine) Dispose() {
	l.disposed = true
}

func (l *PolyLine) Draw(screen *ebiten.Image) {
	if !l.Visible {
		return
	}
	if len(l.Points) < 2 {
		return
	}

	points := l.Points
	for i := 0; i < len(points)-1; i++ {
		pt1 := points[i]
		pt2 := points[i+1]
		drawLine(screen, pt1, pt2, l.Width, l.ColorScale)
	}
}
