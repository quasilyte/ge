package gedebug

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/internal/primitives"
)

type VecLine struct {
	Vec   *gemath.Vec
	Pos   *gemath.Vec
	Color color.RGBA
}

func (l *VecLine) Draw(screen *ebiten.Image) {
	c := l.Color
	if c == (color.RGBA{}) {
		c = color.RGBA{G: 100, B: 200, A: 100}
	}

	angle := l.Vec.Angle()
	arrowPoint := l.Pos.MoveInDirection(48, angle)
	left := arrowPoint.MoveInDirection(8, angle-2.2)
	right := arrowPoint.MoveInDirection(8, angle+2.2)
	primitives.DrawLine(screen, l.Pos.X, l.Pos.Y, arrowPoint.X, arrowPoint.Y, c)
	primitives.DrawLine(screen, arrowPoint.X, arrowPoint.Y, left.X, left.Y, c)
	primitives.DrawLine(screen, arrowPoint.X, arrowPoint.Y, right.X, right.Y, c)
}

func (l *VecLine) Dispose() {
	l.Vec = nil
}

func (l *VecLine) IsDisposed() bool { return l.Vec == nil }
