package ge

import (
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/quasilyte/ge/gemath"
	"golang.org/x/image/font"
)

type AlignVertical uint8

const (
	AlignVerticalTop AlignVertical = iota
	AlignVerticalCenter
	AlignVerticalBottom
)

type AlignHorizontal uint8

const (
	AlignHorizontalLeft AlignHorizontal = iota
	AlignHorizontalCenter
	AlignHorizontalRight
)

type GrowVertical uint8

const (
	GrowVerticalDown GrowVertical = iota
	GrowVerticalUp
	GrowVerticalBoth
	GrowVerticalNone
)

type GrowHorizontal uint8

const (
	GrowHorizontalRight GrowHorizontal = iota
	GrowHorizontalLeft
	GrowHorizontalBoth
	GrowHorizontalNone
)

type Label struct {
	Text string

	ColorScale ColorScale

	Hue gemath.Rad

	Pos    Pos
	Width  float64
	Height float64

	Visible bool

	AlignVertical   AlignVertical
	AlignHorizontal AlignHorizontal
	GrowVertical    GrowVertical
	GrowHorizontal  GrowHorizontal

	face       font.Face
	capHeight  float64
	lineHeight float64
}

func NewLabel(ff font.Face) *Label {
	m := ff.Metrics()
	capHeight := math.Abs(float64(m.CapHeight.Floor()))
	lineHeight := float64(m.Height.Floor())
	label := &Label{
		face:            ff,
		capHeight:       capHeight,
		lineHeight:      lineHeight,
		ColorScale:      defaultColorScale,
		Visible:         true,
		AlignHorizontal: AlignHorizontalLeft,
		AlignVertical:   AlignVerticalTop,
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
	if !l.Visible || l.Text == "" {
		return
	}

	pos := l.Pos.Resolve()

	// Adjust the pos, since "dot position" (baseline) is not a top-left corner.
	pos.Y += l.capHeight

	numLines := strings.Count(l.Text, "\n") + 1

	var containerRect gemath.Rect
	bounds := text.BoundString(l.face, l.Text)
	boundsWidth := float64(bounds.Dx())
	boundsHeight := float64(bounds.Dy())
	if l.Width == 0 && l.Height == 0 {
		// Auto-sized container.
		containerRect = gemath.Rect{
			Min: pos,
			Max: pos.Add(gemath.Vec{X: float64(bounds.Dx()), Y: float64(bounds.Dy())}),
		}
	} else {
		containerRect = gemath.Rect{
			Min: pos,
			Max: pos.Add(gemath.Vec{X: l.Width, Y: l.Height}),
		}
		if delta := boundsWidth - l.Width; delta > 0 {
			switch l.GrowHorizontal {
			case GrowHorizontalRight:
				containerRect.Max.X += delta
			case GrowHorizontalLeft:
				containerRect.Min.X -= delta
			case GrowHorizontalBoth:
				containerRect.Min.X -= delta / 2
				containerRect.Max.X += delta / 2
			case GrowHorizontalNone:
				// Do nothing.
			}
		}
		if delta := boundsHeight - l.Height; delta > 0 {
			switch l.GrowVertical {
			case GrowVerticalDown:
				containerRect.Min.Y += delta
			case GrowVerticalUp:
				containerRect.Min.Y -= delta
				pos.Y -= delta
			case GrowVerticalBoth:
				containerRect.Min.Y -= delta / 2
				containerRect.Max.Y += delta / 2
				pos.Y -= delta / 2
			case GrowVerticalNone:
				// Do nothing.
			}
		}
	}

	switch l.AlignVertical {
	case AlignVerticalTop:
		// Do nothing.
	case AlignVerticalCenter:
		pos.Y += (containerRect.Height() - l.estimateHeight(numLines)) / 2
	case AlignVerticalBottom:
		pos.Y += containerRect.Height() - l.estimateHeight(numLines)
	}

	var drawOptions ebiten.DrawImageOptions
	applyColorScale(l.ColorScale, &drawOptions)
	if l.Hue != 0 {
		drawOptions.ColorM.RotateHue(float64(l.Hue))
	}
	drawOptions.Filter = ebiten.FilterLinear

	if l.AlignHorizontal == AlignHorizontalLeft {
		drawOptions.GeoM.Translate(pos.X, pos.Y)
		text.DrawWithOptions(screen, l.Text, l.face, &drawOptions)
		return
	}

	textRemaining := l.Text
	offsetY := 0.0
	for {
		nextLine := strings.IndexByte(textRemaining, '\n')
		lineText := textRemaining
		if nextLine != -1 {
			lineText = textRemaining[:nextLine]
			textRemaining = textRemaining[nextLine+len("\n"):]
		}
		lineBounds := text.BoundString(l.face, lineText)
		lineBoundsWidth := float64(lineBounds.Dx())
		offsetX := 0.0
		switch l.AlignHorizontal {
		case AlignHorizontalCenter:
			offsetX = (containerRect.Width() - lineBoundsWidth) / 2
		case AlignHorizontalRight:
			offsetX = containerRect.Width() - lineBoundsWidth
		}
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(pos.X+offsetX, pos.Y+offsetY)
		text.DrawWithOptions(screen, lineText, l.face, &drawOptions)
		if nextLine == -1 {
			break
		}
		offsetY += l.lineHeight
	}
}

func (l *Label) estimateHeight(numLines int) float64 {
	estimatedHeight := l.capHeight
	if numLines >= 2 {
		estimatedHeight += (float64(numLines) - 1) * l.lineHeight
	}
	return estimatedHeight
}
