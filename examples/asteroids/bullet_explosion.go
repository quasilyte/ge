package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

type bulletExplosion struct {
	anim *ge.Animation
	pos  gemath.Vec
}

func newBulletExplosion(pos gemath.Vec) *bulletExplosion {
	return &bulletExplosion{
		pos: pos,
	}
}

func (b *bulletExplosion) Init(scene *ge.Scene) {
	ctx := scene.Context()

	sprite := ctx.LoadSprite("explosion.png")
	sprite.Pos = &b.pos
	sprite.Width = 40
	sprite.Scaling = 1.3
	scene.AddGraphics(sprite)

	b.anim = ge.NewAnimation(sprite)
	b.anim.SecondsPerFrame = 0.03
}

func (b *bulletExplosion) Dispose() { b.anim.Dispose() }

func (b *bulletExplosion) IsDisposed() bool { return b.anim.IsDisposed() }

func (b *bulletExplosion) Update(delta float64) {
	if b.anim.Tick(delta) {
		b.Dispose()
	}
}
