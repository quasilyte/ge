package gedebug

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gedraw"
	"github.com/quasilyte/ge/physics"
)

type BodyAura struct {
	Body *physics.Body
}

func (a *BodyAura) Draw(screen *ebiten.Image) {
	c := color.RGBA{G: 100, B: 200, A: 100}

	// if a.Body.IsRect() {
	// 	rect := a.Body.BoundsRect()
	// 	ebitenutil.DrawRect(screen, float64(rect.X1()), float64(rect.Y1()), float64(rect.Width()), float64(rect.Height()), c)
	// 	return
	// }

	if a.Body.IsRotatedRect() {
		vertices := a.Body.RotatedRectVertices()
		gedraw.DrawPath(screen, vertices[:], c)
		return
	}

	if a.Body.IsCircle() {
		gedraw.DrawCircle(screen, a.Body.Pos, a.Body.CircleRadius(), c)
		return
	}

	panic("unsupported body kind")
}

func (a *BodyAura) IsDisposed() bool { return a.Body.IsDisposed() }
