package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/quasilyte/ge/gemath"
	"golang.org/x/image/font"
)

type AlignHorizontal int

const (
	AlignLeft AlignHorizontal = iota
	AlignCenterHorizontal
	AlignRight
)

type AlignVertical int

const (
	AlignTop AlignVertical = iota
	AlignCenter
	AlignBottom
)

type Label struct {
	Text string

	ColorScale ColorScale

	Pos *gemath.Vec

	Origin gemath.Vec

	Visible bool
	HAlign  AlignHorizontal
	VAlign  AlignVertical

	face font.Face
}

func NewLabel(ff font.Face) *Label {
	label := &Label{
		face:       ff,
		ColorScale: defaultColorScale,
		HAlign:     AlignLeft,
		VAlign:     AlignTop,
		Visible:    true,
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
	if !l.Visible {
		return
	}

	var drawOptions ebiten.DrawImageOptions

	var origin gemath.Vec

	bounds := text.BoundString(l.face, l.Text)
	boundsWidth := float64(bounds.Dx())
	switch l.VAlign {
	case AlignTop:
		origin.Y = float64(l.face.Metrics().CapHeight.Round())
	}
	switch l.HAlign {
	case AlignLeft:
		// Do nothing.
	case AlignCenterHorizontal:
		origin.X = boundsWidth / 2
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
