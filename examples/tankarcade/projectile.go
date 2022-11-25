package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type projectile struct {
	config projectileConfig

	scene *ge.Scene

	target    gmath.Vec
	velocity  gmath.Vec
	rotation  gmath.Rad
	dist      float64
	hadImpact bool

	sprite *ge.Sprite
	body   physics.Body
}

type projectileConfig struct {
	alliance int
	pos      gmath.Vec
	target   gmath.Vec
	design   *weaponDesign
}

func newProjectile(config projectileConfig) *projectile {
	return &projectile{config: config}
}

func (p *projectile) CanHit() bool {
	return p.dist >= p.config.design.minRange && !p.hadImpact
}

func (p *projectile) Init(scene *ge.Scene) {
	p.scene = scene

	p.body.Pos = p.config.pos
	p.body.InitCircle(p, 6)
	if p.config.design == sonicCannonDesign {
		p.body.LayerMask = projectileLayerMask(p.config.alliance) | wallLayerMask
	} else if p.config.design == improvedRegeneratorDesign {
		p.body.LayerMask = unitLayerMask(p.config.alliance) | mineLayer
	} else if p.config.design == regeneratorDesign.extra {
		p.body.LayerMask = unitLayerMask(p.config.alliance)
	} else {
		p.body.LayerMask = projectileLayerMask(p.config.alliance)
	}
	scene.AddBody(&p.body)

	p.updateTarget(p.config.target)

	if p.config.design.projectileImage != ImageNone {
		p.sprite = scene.NewSprite(p.config.design.projectileImage)
		p.sprite.Pos.Base = &p.body.Pos
		p.sprite.Rotation = &p.rotation
		scene.AddGraphics(p.sprite)
	}
}

func (p *projectile) Update(delta float64) {
	dist := p.config.design.projectileSpeed * delta
	if p.body.Pos.DistanceTo(p.target) <= dist {
		p.Dispose()
		return
	}
	if p.config.design.projectileRotationSpeed != 0 {
		p.rotation += gmath.Rad(delta) * p.config.design.projectileRotationSpeed
	}
	if p.sprite != nil {
		if p.CanHit() {
			p.sprite.SetAlpha(1)
		} else {
			p.sprite.SetAlpha(0.5)
		}
	}
	p.dist += dist
	p.body.Pos = p.body.Pos.Add(p.velocity.Mulf(delta))
}

func (p *projectile) ForceDestroy() {
	p.destroy(true)
}

func (p *projectile) Destroy() {
	p.destroy(false)
}

func (p *projectile) Reflect() {
	p.updateTarget(p.config.pos)
	p.body.LayerMask ^= 0b11
}

func (p *projectile) updateTarget(target gmath.Vec) {
	p.target = target
	targetDir := target.DirectionTo(p.body.Pos)
	p.rotation = targetDir.Angle()
	p.velocity = targetDir.Mulf(p.config.design.projectileSpeed)
}

func (p *projectile) destroy(forced bool) {
	e := newExplosion(p.body.Pos)
	e.Image = p.config.design.projectileExplosion
	e.Scale = p.config.design.projectileExplosionScale
	e.Hue = p.config.design.projectileExplosionHue
	e.Rotating = p.config.design.projectileExplosionRotates
	p.scene.AddObject(e)

	if p.config.design.hitSound != AudioNone {
		p.scene.Audio().PlaySound(p.config.design.hitSound)
	}

	if p.config.design == hurricaneDesign && !forced {
		p.hadImpact = true
	} else {
		p.dispose(forced)
	}
}

func (p *projectile) IsDisposed() bool {
	return p.body.IsDisposed()
}

func (p *projectile) Dispose() {
	p.dispose(false)
}

func (p *projectile) dispose(forced bool) {
	if !forced && p.config.design == hurricaneDesign {
		targetDir := p.target.DirectionTo(p.body.Pos)
		p2 := newProjectile(projectileConfig{
			alliance: p.config.alliance,
			pos:      p.target,
			target:   fireTargetPos(p.target, targetDir.Mulf(-1).Angle().Normalized(), 64*2),
			design:   hurricaneBouncedDesign,
		})
		p.scene.AddObject(p2)
	}

	p.body.Dispose()
	if p.sprite != nil {
		p.sprite.Dispose()
	}
}
