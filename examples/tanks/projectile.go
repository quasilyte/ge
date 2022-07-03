package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/physics"
)

type projectile struct {
	scene *ge.Scene

	body     physics.Body
	velocity gemath.Vec

	design *turretDesign

	dist float64

	alliance int

	sprite *ge.Sprite
}

func newProjectile(alliance int, pos gemath.Vec, direction gemath.Rad, design *turretDesign) *projectile {
	p := &projectile{
		design:   design,
		alliance: alliance,
		dist:     design.FireRange,
		velocity: gemath.RadToVec(direction).Mulf(design.ProjectileSpeed),
	}
	p.body.Pos = pos
	p.body.Rotation = direction
	return p
}

func (p *projectile) Init(scene *ge.Scene) {
	p.scene = scene

	p.body.InitCircle(p, 4)

	p.sprite = scene.NewSprite(p.design.AmmoImage)
	p.sprite.Pos.Base = &p.body.Pos
	p.sprite.Rotation = &p.body.Rotation
	scene.AddGraphics(p.sprite)

	scene.AddBody(&p.body)
}

func (p *projectile) IsDisposed() bool { return p.body.IsDisposed() }

func (p *projectile) Dispose() {
	p.body.Dispose()
	p.sprite.Dispose()
}

func (p *projectile) Update(delta float64) {
	for _, collision := range p.scene.GetCollisions(&p.body) {
		impact := false
		switch obj := collision.Body.Object.(type) {
		case *battleTank:
			if obj.Player.Alliance == p.alliance {
				continue
			}
			impact = true
			obj.OnDamage(p.design.Damage, p.design.DamageKind)
		case *battlePost:
			if obj.Player.Alliance == p.alliance {
				continue
			}
			impact = true
			obj.OnDamage(p.design.Damage, p.design.DamageKind)
		}
		if impact {
			p.Destroy()
			return
		}
	}

	if p.dist <= 0 {
		p.Dispose()
		return
	}
	movement := p.velocity.Mulf(delta)
	p.dist -= movement.Len()
	p.body.Pos = p.body.Pos.Add(movement)
}

func (p *projectile) Destroy() {
	switch p.design.Name {
	case "gatling_gun", "gauss", "ion":
		// No explosion.
	case "lancer":
		// Bigger explosion.
		e := newExplosion(p.body.Pos)
		e.Scale = 0.7
		e.AnimationSpeed = 1.75
		p.scene.AddObject(e)
	case "dual_cannon":
		// Two explosions.
		e1 := newExplosion(p.body.Pos.Add(gemath.Vec{X: 3, Y: 3}))
		e1.Scale = 0.4
		e1.AnimationSpeed = 2.5
		p.scene.AddObject(e1)
		e2 := newExplosion(p.body.Pos.Add(gemath.Vec{X: -3, Y: -3}))
		e2.Scale = 0.4
		e2.AnimationSpeed = 2.5
		p.scene.AddObject(e2)
	default:
		e := newExplosion(p.body.Pos)
		e.Scale = 0.5
		e.AnimationSpeed = 2.5
		p.scene.AddObject(e)
	}

	p.Dispose()
}
