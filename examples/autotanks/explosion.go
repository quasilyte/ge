package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

type explosion struct {
	anim *ge.Animation
	pos  gemath.Vec

	Scale          float64
	AnimationSpeed float64
}

func newExplosion(pos gemath.Vec) *explosion {
	return &explosion{
		pos:            pos,
		Scale:          1,
		AnimationSpeed: 1,
	}
}

func (e *explosion) Init(scene *ge.Scene) {
	sprite := scene.NewSprite(ImageExplosion)
	sprite.Pos.Base = &e.pos
	sprite.FrameWidth = 64
	sprite.Scale = e.Scale
	scene.AddGraphics(sprite)

	e.anim = ge.NewAnimation(sprite, -1)
	e.anim.SetSecondsPerFrame(0.08)
}

func (e *explosion) Dispose() { e.anim.Sprite().Dispose() }

func (e *explosion) IsDisposed() bool { return e.anim.IsDisposed() }

func (e *explosion) Update(delta float64) {
	if e.anim.Tick(delta * e.AnimationSpeed) {
		e.Dispose()
	}
}
