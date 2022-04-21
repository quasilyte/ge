package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

type vesselExplosion struct {
	anim *ge.Animation
	pos  gemath.Vec
}

func newVesselExplosion(pos gemath.Vec) *vesselExplosion {
	return &vesselExplosion{
		pos: pos,
	}
}

func (e *vesselExplosion) Init(scene *ge.Scene) {
	ctx := scene.Context()

	sprite := ctx.LoadSprite("explosion.png")
	sprite.Pos = &e.pos
	sprite.Rotation = ge.NewRotation(0)
	sprite.Width = 40
	sprite.Scaling = 1.7
	scene.AddGraphics(sprite)

	e.anim = ge.NewAnimation(sprite)
	e.anim.SecondsPerFrame = 0.08
}

func (e *vesselExplosion) Dispose() { e.anim.Dispose() }

func (e *vesselExplosion) IsDisposed() bool { return e.anim.IsDisposed() }

func (e *vesselExplosion) Update(delta float64) {
	*e.anim.Sprite().Rotation += gemath.Rad(delta * 2)
	if e.anim.Tick(delta) {
		e.Dispose()
	}
}
