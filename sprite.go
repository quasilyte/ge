package ge

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
)

type Sprite struct {
	image *ebiten.Image

	Pos      Pos
	Rotation *gemath.Rad

	Scale float64

	ColorScale ColorScale

	Hue gemath.Rad

	Visible  bool
	Centered bool

	FrameOffset gemath.Vec
	FrameWidth  float64
	FrameHeight float64

	disposed bool
}

type ColorScale struct {
	R float32
	G float32
	B float32
	A float32
}

func (c *ColorScale) SetRGBA(r, g, b, a uint8) {
	c.R = float32(r) / 255
	c.G = float32(g) / 255
	c.B = float32(b) / 255
	c.A = float32(a) / 255
}

var defaultColorScale = ColorScale{1, 1, 1, 1}

func NewSprite() *Sprite {
	return &Sprite{
		Visible:    true,
		Centered:   true,
		Scale:      1,
		ColorScale: defaultColorScale,
	}
}

func (s *Sprite) SetImage(img *ebiten.Image) {
	w, h := img.Size()
	s.image = img
	s.FrameWidth = float64(w)
	s.FrameHeight = float64(h)
}

func (s *Sprite) SetRepeatedImage(img *ebiten.Image, width, height float64) {
	w, h := img.Size()
	repeated := ebiten.NewImage(int(width), int(height))
	var op ebiten.DrawImageOptions
	for y := float64(0); y < height; y += float64(h) {
		for x := float64(0); x < width; x += float64(w) {
			op.GeoM.Reset()
			op.GeoM.Translate(x, y)
			repeated.DrawImage(img, &op)
		}
	}
	s.image = repeated
	s.FrameWidth = width
	s.FrameHeight = height
}

func (s *Sprite) ImageWidth() float64 {
	w, _ := s.image.Size()
	return float64(w)
}

func (s *Sprite) ImageHeight() float64 {
	_, h := s.image.Size()
	return float64(h)
}

func (s *Sprite) IsDisposed() bool {
	return s.disposed
}

func (s *Sprite) Dispose() {
	s.disposed = true
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	if !s.Visible {
		return
	}

	var drawOptions ebiten.DrawImageOptions

	var origin gemath.Vec
	if s.Centered {
		origin = gemath.Vec{X: s.FrameWidth / 2, Y: s.FrameHeight / 2}
	}
	origin = origin.Sub(s.Pos.Offset)

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

	if s.ColorScale != defaultColorScale {
		r := float64(s.ColorScale.R)
		g := float64(s.ColorScale.G)
		b := float64(s.ColorScale.B)
		a := float64(s.ColorScale.A)
		drawOptions.ColorM.Scale(r, g, b, a)
	}
	if s.Hue != 0 {
		drawOptions.ColorM.RotateHue(float64(s.Hue))
	}

	subImage := s.image.SubImage(image.Rectangle{
		Min: image.Point{
			X: int(s.FrameOffset.X),
			Y: int(s.FrameOffset.Y),
		},
		Max: image.Point{
			X: int(s.FrameOffset.X + s.FrameWidth),
			Y: int(s.FrameOffset.Y + s.FrameHeight),
		},
	}).(*ebiten.Image)
	screen.DrawImage(subImage, &drawOptions)
}
