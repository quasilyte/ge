package ge

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
)

type Sprite struct {
	image *ebiten.Image

	Pos *gemath.Vec

	Rotation *gemath.Rad

	Scale float64

	ColorScale ColorScale

	Hue gemath.Rad

	Centered bool
	Origin   gemath.Vec

	Offset gemath.Vec
	Width  float64
	Height float64

	disposed bool
}

type ColorScale struct {
	R float32
	G float32
	B float32
	A float32
}

var defaultColorScale = ColorScale{1, 1, 1, 1}

func NewSprite(img *ebiten.Image) *Sprite {
	w, h := img.Size()
	sprite := &Sprite{
		image:      img,
		Width:      float64(w),
		Height:     float64(h),
		Centered:   true,
		Scale:      1,
		ColorScale: defaultColorScale,
	}
	return sprite
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
	var drawOptions ebiten.DrawImageOptions

	var origin gemath.Vec
	if s.Centered {
		origin = gemath.Vec{X: s.Width / 2, Y: s.Height / 2}
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
			X: int(s.Offset.X),
			Y: int(s.Offset.Y),
		},
		Max: image.Point{
			X: int(s.Offset.X + s.Width),
			Y: int(s.Offset.Y + s.Height),
		},
	}).(*ebiten.Image)
	screen.DrawImage(subImage, &drawOptions)
}
