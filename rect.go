package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/internal/primitives"
	"github.com/quasilyte/gmath"
)

type Rect struct {
	Pos      Pos
	Rotation *gmath.Rad

	Width  float64
	Height float64

	Scale float64

	Centered bool

	OutlineColorScale ColorScale
	FillColorScale    ColorScale

	OutlineWidth float64

	Visible bool

	imageCache *imageCache

	disposed bool
}

func NewRect(ctx *Context, width, height float64) *Rect {
	return &Rect{
		Visible:           true,
		Centered:          true,
		FillColorScale:    defaultColorScale,
		OutlineColorScale: transparentColor,
		Scale:             1,
		Width:             width,
		Height:            height,
		imageCache:        &ctx.imageCache,
	}
}

func (rect *Rect) BoundsRect() gmath.Rect {
	pos := rect.Pos.Resolve()
	if rect.Centered {
		offset := gmath.Vec{X: rect.Width * 0.5, Y: rect.Height * 0.5}
		return gmath.Rect{
			Min: pos.Sub(offset),
			Max: pos.Add(offset),
		}
	}
	return gmath.Rect{
		Min: pos,
		Max: pos.Add(gmath.Vec{X: rect.Width, Y: rect.Height}),
	}
}

func (rect *Rect) IsDisposed() bool {
	return rect.disposed
}

func (rect *Rect) Dispose() {
	rect.disposed = true
}

// AnchorPos returns a top-left position.
// When Centered is false, it's identical to Pos, otherwise
// it will apply the computations to get the right anchor for the centered rect.
func (rect *Rect) AnchorPos() Pos {
	if rect.Centered {
		return rect.Pos.WithOffset(-rect.Width/2, -rect.Height/2)
	}
	return rect.Pos
}

func (rect *Rect) calculateGeom(w, h float64, offset gmath.Vec) ebiten.GeoM {
	var origin gmath.Vec
	if rect.Centered {
		origin = gmath.Vec{X: rect.Width / 2, Y: rect.Height / 2}
	}
	origin = origin.Sub(rect.Pos.Offset)

	var geom ebiten.GeoM
	geom.Scale(w, h)

	geom.Translate(-origin.X, -origin.Y)
	if rect.Rotation != nil {
		geom.Rotate(float64(*rect.Rotation))
	}
	if rect.Scale != 1 {
		geom.Scale(rect.Scale, rect.Scale)
	}
	geom.Translate(origin.X, origin.Y)

	if rect.Pos.Base != nil {
		geom.Translate(rect.Pos.Base.X-origin.X, rect.Pos.Base.Y-origin.Y)
	} else if !origin.IsZero() {
		geom.Translate(0-origin.X, 0-origin.Y)
	}

	geom.Translate(offset.X, offset.Y)

	return geom
}

func (rect *Rect) Draw(screen *ebiten.Image) {
	rect.DrawWithOffset(screen, gmath.Vec{})
}

func (rect *Rect) DrawWithOffset(screen *ebiten.Image, offset gmath.Vec) {
	if !rect.Visible {
		return
	}
	if rect.OutlineColorScale.A == 0 && rect.FillColorScale.A == 0 {
		return
	}

	if rect.OutlineColorScale.A == 0 || rect.OutlineWidth < 1 {
		// Fill-only mode.
		var drawOptions ebiten.DrawImageOptions
		drawOptions.GeoM = rect.calculateGeom(rect.Width, rect.Height, offset)
		applyColorScale(rect.FillColorScale, &drawOptions.ColorM)
		screen.DrawImage(primitives.WhitePixel, &drawOptions)
		return
	}

	if rect.FillColorScale.A == 0 && rect.OutlineWidth >= 1 {
		// Outline-only mode.
		var tmpDrawOptions ebiten.DrawImageOptions
		dst := rect.imageCache.NewTempImage(int(rect.Width), int(rect.Height))
		applyColorScale(rect.OutlineColorScale, &tmpDrawOptions.ColorM)
		tmpDrawOptions.GeoM.Scale(rect.Width, rect.Height)
		dst.DrawImage(primitives.WhitePixel, &tmpDrawOptions)

		var geom ebiten.GeoM
		geom.Scale(rect.Width-rect.OutlineWidth*2, rect.Height-rect.OutlineWidth*2)
		geom.Translate(rect.OutlineWidth, rect.OutlineWidth)
		tmpDrawOptions.GeoM = geom
		tmpDrawOptions.CompositeMode = ebiten.CompositeModeClear
		dst.DrawImage(primitives.WhitePixel, &tmpDrawOptions)

		var drawOptions ebiten.DrawImageOptions
		drawOptions.GeoM = rect.calculateGeom(1, 1, offset)
		screen.DrawImage(dst, &drawOptions)
		return
	}

	// TODO: it doesn't work with a fill color with alpha not equal to 1.
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM = rect.calculateGeom(rect.Width, rect.Height, offset)
	applyColorScale(rect.OutlineColorScale, &drawOptions.ColorM)
	screen.DrawImage(primitives.WhitePixel, &drawOptions)
	outlineDrawOptions := drawOptions
	outlineDrawOptions.ColorM.Reset()
	applyColorScale(rect.FillColorScale, &outlineDrawOptions.ColorM)
	outlineDrawOptions.GeoM = rect.calculateGeom(rect.Width-rect.OutlineWidth*2, rect.Height-rect.OutlineWidth*2, offset)
	outlineDrawOptions.GeoM.Translate(rect.OutlineWidth, rect.OutlineWidth)
	outlineDrawOptions.CompositeMode = ebiten.CompositeModeCopy
	screen.DrawImage(primitives.WhitePixel, &outlineDrawOptions)
}
