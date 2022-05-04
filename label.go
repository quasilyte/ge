package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/quasilyte/ge/gemath"
	"golang.org/x/image/font"
)

type Label struct {
	Text string

	ColorScale ColorScale

	Pos *gemath.Vec

	Origin gemath.Vec

	Centered bool

	face font.Face
}

func NewLabel(ff font.Face) *Label {
	label := &Label{
		face:       ff,
		ColorScale: defaultColorScale,
		Centered:   true,
	}
	return label
}

func (l *Label) IsDisposed() bool {
	return l.face == nil
}

func (l *Label) Dispose() {
	l.face = nil
}

func (l *Label) SetColor(r, g, b, a uint8) {
	l.ColorScale = ColorScale{
		R: float32(r) / 255,
		G: float32(g) / 255,
		B: float32(b) / 255,
		A: float32(a) / 255,
	}
}

func (l *Label) Draw(screen *ebiten.Image) {
	var drawOptions ebiten.DrawImageOptions

	var origin gemath.Vec
	bounds := text.BoundString(l.face, l.Text)
	boundsHeight := float64(bounds.Dy())
	if l.Centered {
		boundsWidth := float64(bounds.Dx())
		origin = gemath.Vec{X: boundsWidth / 2, Y: boundsHeight / 2}
	} else {
		origin = gemath.Vec{Y: boundsHeight}
	}
	origin = origin.Sub(l.Origin)

	if l.Pos != nil {
		drawOptions.GeoM.Translate(l.Pos.X-origin.X, l.Pos.Y-origin.Y)
	} else {
		drawOptions.GeoM.Translate(origin.X, origin.Y)
	}

	if l.ColorScale != defaultColorScale {
		r := float64(l.ColorScale.R)
		g := float64(l.ColorScale.G)
		b := float64(l.ColorScale.B)
		a := float64(l.ColorScale.A)
		drawOptions.ColorM.Scale(r, g, b, a)
	}

	drawOptions.Filter = ebiten.FilterLinear
	text.DrawWithOptions(screen, l.Text, l.face, &drawOptions)
}
