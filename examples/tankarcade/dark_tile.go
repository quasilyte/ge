package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type darkTile struct {
	sprite *ge.Sprite
	Body   physics.Body
	scene  *ge.Scene
}

func newDarkTile(pos gmath.Vec) *darkTile {
	t := &darkTile{}
	t.Body.Pos = pos
	return t
}

func (t *darkTile) Pos() gmath.Vec { return t.Body.Pos }

func (t *darkTile) Init(scene *ge.Scene) {
	t.scene = scene

	t.Body.InitStaticCircle(t, 32)
	t.Body.LayerMask = 0b1111
	scene.AddBody(&t.Body)

	t.sprite = scene.NewSprite(ImageDarkTile)
	t.sprite.Pos.Base = &t.Body.Pos
	scene.AddGraphics(t.sprite)
}

func (t *darkTile) IsDisposed() bool {
	return t.Body.IsDisposed()
}

func (t *darkTile) Dispose() {
	t.Body.Dispose()
	t.sprite.Dispose()
}

func (t *darkTile) Update(delta float64) {
	for _, collision := range t.scene.GetCollisions(&t.Body) {
		p, ok := collision.Body.Object.(*projectile)
		if ok {
			p.ForceDestroy()
		}
	}
}
