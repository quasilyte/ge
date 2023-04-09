package ge

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
)

type TextureLine struct {
	BeginPos Pos
	EndPos   Pos

	image      *ebiten.Image
	imageCache *imageCache

	maxLen      float64
	imageHeight float64

	colorScale    ColorScale
	colorM        ebiten.ColorM
	colorsChanged bool

	Visible bool

	disposed bool
}

func NewTextureLine(ctx *Context, begin, end Pos) *TextureLine {
	return &TextureLine{
		Visible:    true,
		BeginPos:   begin,
		EndPos:     end,
		imageCache: &ctx.imageCache,
		colorScale: defaultColorScale,
	}
}

func (l *TextureLine) SetImage(img resource.Image, maxLen float64) {
	w, h := img.Data.Size()
	imageWidth := float64(w)
	l.imageHeight = float64(h)

	texture := ebiten.NewImage(int(math.Round(maxLen)), h)
	x := 0.0
	var drawOptions ebiten.DrawImageOptions
	for x < maxLen {
		segmentWidth := imageWidth
		srcImage := img.Data
		if x+imageWidth > maxLen {
			segmentWidth = maxLen - x
			srcImage = img.Data.SubImage(image.Rectangle{
				Max: image.Point{X: int(math.Round(segmentWidth)), Y: h},
			}).(*ebiten.Image)
		}
		texture.DrawImage(srcImage, &drawOptions)
		drawOptions.GeoM.Translate(segmentWidth, 0)
		x += segmentWidth
	}

	l.image = texture
	l.maxLen = maxLen
}

func (l *TextureLine) IsDisposed() bool {
	return l.disposed
}

func (l *TextureLine) Dispose() {
	l.disposed = true
}

func (l *TextureLine) BoundsRect() gmath.Rect {
	pos1 := l.BeginPos.Resolve()
	pos2 := l.EndPos.Resolve()
	x0 := pos1.X
	x1 := pos2.X
	y0 := pos1.Y
	y1 := pos2.Y
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return gmath.Rect{Min: gmath.Vec{X: x0, Y: y0}, Max: gmath.Vec{X: x1, Y: y1}}
}

func (l *TextureLine) Draw(screen *ebiten.Image) {
	if !l.Visible {
		return
	}

	if l.colorsChanged {
		l.colorsChanged = false
		l.recalculateColorM()
	}

	pos1 := l.BeginPos.Resolve()
	pos2 := l.EndPos.Resolve()

	origin := gmath.Vec{Y: l.imageHeight * 0.5}

	angle := pos1.AngleToPoint(pos2)

	// Maybe use sin+cos to compute the length?
	// TODO: compare sin+cos vs sqrt performance.
	length := gmath.ClampMax(math.Round(pos1.DistanceTo(pos2)), l.maxLen)

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(-origin.X, -origin.Y)
	drawOptions.GeoM.Rotate(float64(angle))
	drawOptions.GeoM.Translate(origin.X, origin.Y)
	drawOptions.GeoM.Translate(pos1.X, pos1.Y)
	drawOptions.ColorM = l.colorM

	bounds := image.Rectangle{
		Max: image.Point{X: int(length), Y: int(l.imageHeight)},
	}
	subImage := l.imageCache.UnsafeSubImage(l.image, bounds)
	screen.DrawImage(subImage, &drawOptions)
}

func (l *TextureLine) GetAlpha() float32 {
	return l.colorScale.A
}

func (l *TextureLine) SetAlpha(a float32) {
	if l.colorScale.A == a {
		return
	}
	l.colorScale.A = a
	l.colorsChanged = true
}

func (l *TextureLine) recalculateColorM() {
	var colorM ebiten.ColorM
	if l.colorScale != defaultColorScale {
		colorM.Scale(float64(l.colorScale.R), float64(l.colorScale.G), float64(l.colorScale.B), float64(l.colorScale.A))
	}
	l.colorM = colorM
}
