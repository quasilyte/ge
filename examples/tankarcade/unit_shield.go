package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type unitShield struct {
	pos      ge.Pos
	sprite   *ge.Sprite
	rotation gmath.Rad
}

func newUnitShield(pos ge.Pos) *unitShield {
	return &unitShield{pos: pos}
}

func (s *unitShield) Init(scene *ge.Scene) {
	s.sprite = scene.NewSprite(ImageShield)
	s.sprite.Pos = s.pos
	s.sprite.Rotation = &s.rotation
	scene.AddGraphics(s.sprite)
}

func (s *unitShield) IsDisposed() bool {
	return s.sprite.IsDisposed()
}

func (s *unitShield) Dispose() {
	s.sprite.Dispose()
}

func (s *unitShield) Update(delta float64) {
	s.rotation += gmath.Rad(delta)
}
