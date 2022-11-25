package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type railgunRayEffect struct {
	sprite *ge.Sprite
	pos    gmath.Vec
	angle  gmath.Rad
}

func newRailgunRayEffect(pos gmath.Vec, angle gmath.Rad) *railgunRayEffect {
	return &railgunRayEffect{pos: pos, angle: angle}
}

func (ray *railgunRayEffect) Init(scene *ge.Scene) {
	ray.sprite = scene.NewSprite(ImageRailgunRay)
	ray.sprite.Pos.Base = &ray.pos
	ray.sprite.Pos.Offset = gmath.Vec{X: -32}
	ray.sprite.Rotation = &ray.angle
	scene.AddGraphics(ray.sprite)
}

func (ray *railgunRayEffect) Dispose() {
	ray.sprite.Dispose()
}

func (ray *railgunRayEffect) IsDisposed() bool {
	return ray.sprite.IsDisposed()
}

func (ray *railgunRayEffect) Update(delta float64) {
	alpha := ray.sprite.GetAlpha() - float32(delta*6)
	if alpha < 0.1 {
		ray.Dispose()
		return
	}
	ray.sprite.SetAlpha(alpha)
}
