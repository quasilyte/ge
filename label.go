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

	Pos Pos

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
	case AlignCenter:
		origin.Y = float64(l.face.Metrics().CapHeight.Round() / 2)
	}
	switch l.HAlign {
	case AlignLeft:
		// Do nothing.
	case AlignCenterHorizontal:
		origin.X = boundsWidth / 2
	}
	origin = origin.Sub(l.Pos.Offset)

	if l.Pos.Base != nil {
		drawOptions.GeoM.Translate(l.Pos.Base.X-origin.X, l.Pos.Base.Y-origin.Y)
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
