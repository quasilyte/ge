package ge

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
)

type Shader struct {
	compiled *ebiten.Shader

	shaderData map[string]any

	Enabled bool

	Texture1 resource.Image
	Texture2 resource.Image
	Texture3 resource.Image
}

func (s *Shader) IsNil() bool { return s.compiled == nil }

func (s *Shader) SetIntValue(key string, v int) {
	s.setFloat32Value(key, float32(v))
}

func (s *Shader) SetFloatValue(key string, v float64) {
	s.setFloat32Value(key, float32(v))
}

func (s *Shader) setFloat32Value(key string, v float32) {
	if oldValue, ok := s.shaderData[key].(float32); ok && oldValue == v {
		return
	}
	if s.shaderData == nil {
		s.shaderData = make(map[string]any, 2)
	}
	s.shaderData[key] = v
}

type Sprite struct {
	image *ebiten.Image
	id    resource.ImageID

	Pos      Pos
	Rotation *gmath.Rad

	colorsChanged bool
	colorScale    ColorScale
	hue           gmath.Rad
	colorM        ebiten.ColorM

	Scale float64

	FlipHorizontal bool
	FlipVertical   bool
	Visible        bool
	Centered       bool

	FrameOffset     gmath.Vec
	FrameWidth      float64
	FrameHeight     float64
	FrameTrimTop    float64
	FrameTrimBottom float64

	imageWidth  float64
	imageHeight float64

	Shader Shader

	imageCache *imageCache

	disposed bool
}

type ColorScale struct {
	R float32
	G float32
	B float32
	A float32
}

func (c *ColorScale) SetColor(rgba color.RGBA) {
	c.SetRGBA(rgba.R, rgba.G, rgba.B, rgba.A)
}

func (c *ColorScale) SetRGBA(r, g, b, a uint8) {
	c.R = float32(r) / 255
	c.G = float32(g) / 255
	c.B = float32(b) / 255
	c.A = float32(a) / 255
}

var defaultColorScale = ColorScale{1, 1, 1, 1}
var transparentColor = ColorScale{0, 0, 0, 0}

func NewSprite(ctx *Context) *Sprite {
	s := &Sprite{
		Visible:    true,
		Centered:   true,
		Scale:      1,
		colorScale: defaultColorScale,
		imageCache: &ctx.imageCache,
	}
	return s
}

func (s *Sprite) SetColorScaleRGBA(r, g, b, a uint8) {
	var scale ColorScale
	scale.SetRGBA(r, g, b, a)
	s.SetColorScale(scale)
}

func (s *Sprite) GetColorScale() ColorScale {
	return s.colorScale
}

func (s *Sprite) GetAlpha() float32 {
	return s.colorScale.A
}

func (s *Sprite) SetAlpha(a float32) {
	if s.colorScale.A == a {
		return
	}
	s.colorScale.A = a
	s.colorsChanged = true
}

func (s *Sprite) SetColorScale(colorScale ColorScale) {
	if s.colorScale == colorScale {
		return
	}
	s.colorScale = colorScale
	s.colorsChanged = true
}

func (s *Sprite) SetHue(hue gmath.Rad) {
	if s.hue == hue {
		return
	}
	s.hue = hue
	s.colorsChanged = true
}

func (s *Sprite) recalculateColorM() {
	var colorM ebiten.ColorM
	if s.colorScale != defaultColorScale {
		colorM.Scale(float64(s.colorScale.R), float64(s.colorScale.G), float64(s.colorScale.B), float64(s.colorScale.A))
	}
	if s.hue != 0 {
		colorM.RotateHue(float64(s.hue))
	}
	s.colorM = colorM
}

func (s *Sprite) ImageID() resource.ImageID {
	return s.id
}

func (s *Sprite) SetImage(img resource.Image) {
	s.id = img.ID
	w, h := img.Data.Size()
	s.imageWidth = float64(w)
	s.imageHeight = float64(h)
	s.image = img.Data
	s.FrameWidth = img.DefaultFrameWidth
	if s.FrameWidth == 0 {
		s.FrameWidth = s.imageWidth
	}
	s.FrameHeight = img.DefaultFrameHeight
	if s.FrameHeight == 0 {
		s.FrameHeight = s.imageHeight
	}
}

func (s *Sprite) SetRepeatedImage(img resource.Image, width, height float64) {
	s.id = img.ID
	w, h := img.Data.Size()
	s.imageWidth = float64(w)
	s.imageHeight = float64(h)
	repeated := ebiten.NewImage(int(width), int(height))
	var op ebiten.DrawImageOptions
	for y := float64(0); y < height; y += s.imageHeight {
		for x := float64(0); x < width; x += s.imageWidth {
			op.GeoM.Reset()
			op.GeoM.Translate(x, y)
			repeated.DrawImage(img.Data, &op)
		}
	}
	s.image = repeated
	s.FrameWidth = img.DefaultFrameWidth
	if s.FrameWidth == 0 {
		s.FrameWidth = width
	}
	s.FrameHeight = img.DefaultFrameHeight
	if s.FrameHeight == 0 {
		s.FrameHeight = height
	}
}

// AnchorPos returns a top-left position.
// When Centered is false, it's identical to Pos, otherwise
// it will apply the computations to get the right anchor for the centered sprite.
func (s *Sprite) AnchorPos() Pos {
	if s.Centered {
		return s.Pos.WithOffset(-s.FrameWidth/2, -s.FrameHeight/2)
	}
	return s.Pos
}

func (s *Sprite) ImageWidth() float64 {
	w, _ := s.image.Size()
	return float64(w)
}

func (s *Sprite) ImageHeight() float64 {
	_, h := s.image.Size()
	return float64(h)
}

func (s *Sprite) BoundsRect() gmath.Rect {
	pos := s.Pos.Resolve()
	if s.Centered {
		offset := gmath.Vec{X: s.FrameWidth * 0.5, Y: s.FrameHeight * 0.5}
		return gmath.Rect{
			Min: pos.Sub(offset),
			Max: pos.Add(offset),
		}
	}
	return gmath.Rect{
		Min: pos,
		Max: pos.Add(gmath.Vec{X: s.FrameWidth, Y: s.FrameHeight}),
	}
}

func (s *Sprite) IsDisposed() bool {
	return s.disposed
}

func (s *Sprite) Dispose() {
	s.disposed = true
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	if !s.Visible || s.image == nil {
		return
	}

	var drawOptions ebiten.DrawImageOptions

	var origin gmath.Vec
	if s.Centered {
		origin = gmath.Vec{X: s.FrameWidth / 2, Y: s.FrameHeight / 2}
	}
	origin = origin.Sub(s.Pos.Offset)

	if s.FlipHorizontal {
		drawOptions.GeoM.Scale(-1, 1)
		drawOptions.GeoM.Translate(s.FrameWidth, 0)
	}
	if s.FlipVertical {
		drawOptions.GeoM.Scale(1, -1)
		drawOptions.GeoM.Translate(0, s.FrameHeight)
	}

	drawOptions.GeoM.Translate(-origin.X, -origin.Y)
	if s.Rotation != nil {
		drawOptions.GeoM.Rotate(float64(*s.Rotation))
	}
	if s.Scale != 1 {
		drawOptions.GeoM.Scale(s.Scale, s.Scale)
	}
	drawOptions.GeoM.Translate(origin.X, origin.Y)

	if s.Pos.Base != nil {
		drawOptions.GeoM.Translate(s.Pos.Base.X-origin.X, s.Pos.Base.Y-origin.Y)
	} else if !origin.IsZero() {
		drawOptions.GeoM.Translate(0-origin.X, 0-origin.Y)
	}

	if s.colorsChanged {
		s.colorsChanged = false
		s.recalculateColorM()
	}
	drawOptions.ColorM = s.colorM

	var srcImage *ebiten.Image
	var srcBounds image.Rectangle
	needSubImage := (s.FrameOffset != gmath.Vec{}) ||
		s.FrameTrimTop != 0 ||
		s.FrameTrimBottom != 0 ||
		s.FrameWidth != s.imageWidth ||
		s.FrameHeight != s.imageHeight
	if needSubImage {
		srcBounds = image.Rectangle{
			Min: image.Point{
				X: int(s.FrameOffset.X),
				Y: int(s.FrameOffset.Y + s.FrameTrimTop),
			},
			Max: image.Point{
				X: int(s.FrameOffset.X + s.FrameWidth),
				Y: int(s.FrameOffset.Y + s.FrameHeight - s.FrameTrimBottom),
			},
		}
		// Basically, we're doing this:
		// > srcImage = s.image.SubImage(srcBounds)
		// But without redundant allocation.
		unsafeSrc := toUnsafeImage(s.image)
		unsafeSubImage := s.imageCache.UnsafeImageForSubImage()
		unsafeSubImage.original = unsafeSrc
		unsafeSubImage.bounds = srcBounds
		unsafeSubImage.image = unsafeSrc.image
		srcImage = toEbitenImage(unsafeSubImage)
	} else {
		srcImage = s.image
		srcBounds = s.image.Bounds()
	}

	shaderEnabled := s.Shader.Enabled && !s.Shader.IsNil()
	if !shaderEnabled {
		screen.DrawImage(srcImage, &drawOptions)
	} else {
		var drawDest *ebiten.Image
		var options ebiten.DrawRectShaderOptions
		usesColor := s.colorScale != defaultColorScale || s.hue != 0
		if usesColor {
			drawDest = s.imageCache.NewTempImage(srcBounds.Dx(), srcBounds.Dy())
		} else {
			drawDest = screen
			options.GeoM = drawOptions.GeoM
		}
		options.CompositeMode = drawOptions.CompositeMode
		options.Images[0] = srcImage
		options.Images[1] = s.Shader.Texture1.Data
		options.Images[2] = s.Shader.Texture2.Data
		options.Images[3] = s.Shader.Texture3.Data
		options.Uniforms = s.Shader.shaderData
		drawDest.DrawRectShader(srcBounds.Dx(), srcBounds.Dy(), s.Shader.compiled, &options)
		if usesColor {
			screen.DrawImage(drawDest, &drawOptions)
		}
	}
}

var tmpImage = ebiten.NewImage(64, 64)
