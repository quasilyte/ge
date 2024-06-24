package main

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
)

type explosion struct {
	anim     *ge.Animation
	pos      gmath.Vec
	rotation gmath.Rad

	Image          resource.ImageID
	Scale          float64
	Rotating       bool
	Hue            gmath.Rad
	AnimationSpeed float64
}

func newExplosion(pos gmath.Vec) *explosion {
	return &explosion{
		pos:            pos,
		Scale:          1,
		AnimationSpeed: 1,
		Image:          ImageExplosion1,
	}
}

func (e *explosion) Init(scene *ge.Scene) {
	e.rotation = scene.Rand().Rad()
	sprite := scene.NewSprite(e.Image)
	sprite.Pos.Base = &e.pos
	SetHue(sprite, e.Hue)
	sprite.SetScale(e.Scale, e.Scale)

	sprite.Rotation = &e.rotation
	scene.AddGraphics(sprite)

	e.anim = ge.NewAnimation(sprite, -1)
	e.anim.SetSecondsPerFrame(0.04)
}

func (e *explosion) Dispose() { e.anim.Sprite().Dispose() }

func (e *explosion) IsDisposed() bool { return e.anim.IsDisposed() }

func (e *explosion) Update(delta float64) {
	if e.anim.Tick(delta * e.AnimationSpeed) {
		e.Dispose()
	}
	if e.Rotating {
		e.rotation += gmath.Rad(delta * 4)
	}
}
