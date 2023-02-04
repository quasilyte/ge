package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/gmath"
)

type PolyLine struct {
	Points []LinePoint

	Width float64

	Visible bool

	disposed bool
}

type LinePoint struct {
	Pos        gmath.Vec
	ColorScale ColorScale
}

func NewPolyLine() *PolyLine {
	return &PolyLine{
		Visible: true,
		Width:   1,
	}
}

func (l *PolyLine) ResetPoints() {
	l.Points = l.Points[:0]
}

func (l *PolyLine) PushColorPoint(pt gmath.Vec, c ColorScale) {
	l.Points = append(l.Points, LinePoint{
		ColorScale: c,
		Pos:        pt,
	})
}

func (l *PolyLine) PushPoint(pt gmath.Vec) {
	l.Points = append(l.Points, LinePoint{
		ColorScale: defaultColorScale,
		Pos:        pt,
	})
}

func (l *PolyLine) PopPoint() gmath.Vec {
	pt := l.Points[len(l.Points)-1]
	l.Points = l.Points[:len(l.Points)-1]
	return pt.Pos
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
	var colorM ebiten.ColorM
	for i := 0; i < len(points)-1; i++ {
		pt1 := points[i]
		pt2 := points[i+1]
		applyColorScale(pt2.ColorScale, &colorM)
		drawLine(screen, pt1.Pos, pt2.Pos, l.Width, colorM)
	}
}
