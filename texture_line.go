package ge

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/gmath"
)

type TextureLine struct {
	BeginPos Pos
	EndPos   Pos

	Shader Shader

	texture    *Texture
	imageCache *imageCache

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

func (l *TextureLine) SetTexture(tex *Texture) {
	l.texture = tex
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
}

func (l *TextureLine) DrawWithOffset(screen *ebiten.Image, offset gmath.Vec) {
	if !l.Visible {
		return
	}

	if l.colorsChanged {
		l.colorsChanged = false
		l.recalculateColorM()
	}

	pos1 := l.BeginPos.Resolve()
	pos2 := l.EndPos.Resolve()

	origin := gmath.Vec{Y: l.texture.height * 0.5}

	angle := pos1.AngleToPoint(pos2)

	// Maybe use sin+cos to compute the length?
	// TODO: compare sin+cos vs sqrt performance.
	length := gmath.ClampMax(math.Round(pos1.DistanceTo(pos2)), l.texture.width)

	if length == 0 {
		return
	}

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Translate(-origin.X, -origin.Y)
	drawOptions.GeoM.Rotate(float64(angle))
	drawOptions.GeoM.Translate(origin.X, origin.Y)
	drawOptions.GeoM.Translate(pos1.X, pos1.Y)
	drawOptions.GeoM.Translate(offset.X, offset.Y)
	drawOptions.ColorM = l.colorM

	bounds := image.Rectangle{
		Max: image.Point{X: int(length), Y: int(l.texture.height)},
	}
	subImage := l.imageCache.UnsafeSubImage(l.texture.image, bounds)

	shaderEnabled := l.Shader.Enabled && !l.Shader.IsNil()
	if !shaderEnabled {
		screen.DrawImage(subImage, &drawOptions)
	} else {
		var drawDest *ebiten.Image
		var options ebiten.DrawRectShaderOptions
		usesColor := l.colorScale != defaultColorScale
		if usesColor {
			drawDest = l.imageCache.NewTempImage(bounds.Dx(), bounds.Dy())
		} else {
			drawDest = screen
			options.GeoM = drawOptions.GeoM
		}
		options.CompositeMode = drawOptions.CompositeMode
		options.Images[0] = subImage
		options.Images[1] = l.Shader.Texture1.Data
		options.Images[2] = l.Shader.Texture2.Data
		options.Images[3] = l.Shader.Texture3.Data
		options.Uniforms = l.Shader.shaderData
		drawDest.DrawRectShader(bounds.Dx(), bounds.Dy(), l.Shader.compiled, &options)
		if usesColor {
			screen.DrawImage(drawDest, &drawOptions)
		}
	}
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
