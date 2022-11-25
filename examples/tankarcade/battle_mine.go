package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type battleMine struct {
	scene    *ge.Scene
	sprite   *ge.Sprite
	body     physics.Body
	alliance int
	lifespan float64
}

func newBattleMine(pos gmath.Vec, alliance int, lifespan float64) *battleMine {
	b := &battleMine{alliance: alliance, lifespan: lifespan}
	b.body.Pos = pos
	return b
}

func (b *battleMine) Init(scene *ge.Scene) {
	b.scene = scene

	b.sprite = scene.NewSprite(ImageMine)
	b.sprite.Pos.Base = &b.body.Pos
	b.sprite.SetHue(spriteHue(b.alliance))
	scene.AddGraphicsBelow(b.sprite, 1)

	b.body.InitStaticCircle(b, 8)
	b.body.LayerMask = projectileLayerMask(b.alliance) | mineLayer
	scene.AddBody(&b.body)

	// Mines can't be stacked.
	collided := false
	for _, collision := range scene.GetCollisionsAtLayer(&b.body, gmath.Vec{}, mineLayer) {
		m, ok := collision.Body.Object.(*battleMine)
		if ok {
			m.Destroy()
			collided = true
		}
	}
	if collided {
		scene.Audio().PlaySound(mineLayerDesign.extra.hitSound)
		b.Destroy()
	}
}

func (b *battleMine) IsDisposed() bool {
	return b.sprite.IsDisposed()
}

func (b *battleMine) Destroy() {
	e := newExplosion(b.body.Pos)
	b.scene.AddObject(e)

	b.Dispose()
}

func (b *battleMine) Dispose() {
	b.sprite.Dispose()
	b.body.Dispose()
}

func (b *battleMine) Update(delta float64) {
	for _, collision := range b.scene.GetCollisionsAtLayer(&b.body, gmath.Vec{}, mineLayer) {
		p, ok := collision.Body.Object.(*projectile)
		if ok && p.config.alliance != b.alliance {
			p.Destroy()
			b.Destroy()
			return
		}
	}

	if b.lifespan < 0 {
		return
	}

	b.lifespan = gmath.ClampMin(b.lifespan-delta, 0)
	if b.lifespan == 0 {
		b.Dispose()
	}
}
