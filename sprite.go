package ge

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
)

type Sprite struct {
	image *ebiten.Image

	Pos      *gemath.Vec
	Rotation *gemath.Rad

	Scale float64

	ColorScale ColorScale

	Hue gemath.Rad

	Visible  bool
	Centered bool
	Origin   gemath.Vec

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
	origin = origin.Sub(s.Origin)

	drawOptions.GeoM.Translate(-origin.X, -origin.Y)
	if s.Rotation != nil {
		drawOptions.GeoM.Rotate(float64(*s.Rotation))
	}
	if s.Scale != 1 {
		drawOptions.GeoM.Scale(s.Scale, s.Scale)
	}
	drawOptions.GeoM.Translate(origin.X, origin.Y)

	if s.Pos != nil {
		drawOptions.GeoM.Translate(s.Pos.X-origin.X, s.Pos.Y-origin.Y)
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
