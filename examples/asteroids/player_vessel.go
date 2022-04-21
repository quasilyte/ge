package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/collision"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/gesignal"
)

type playerVessel struct {
	ctx    *ge.Context
	body   collision.Body
	sprite *ge.Sprite
	reload float64

	EventDestroyed gesignal.Event[gesignal.Void]
}

func newPlayerVessel(pos gemath.Vec) *playerVessel {
	p := &playerVessel{}
	p.body.InitCircle(p, 32)
	p.body.CollisionHandler = p
	p.body.Pos = pos
	return p
}

func (p *playerVessel) Init(scene *ge.Scene) {
	p.ctx = scene.Context()
	p.sprite = p.ctx.LoadSprite("vessel.png")
	p.sprite.Pos = &p.body.Pos
	p.sprite.Rotation = &p.body.Rotation
	scene.AddGraphics(p.sprite)
	scene.AddBody(&p.body)
}

func (p *playerVessel) OnCollision(info *collision.Info) {
	a, ok := info.Object.(*asteroid)
	if !ok {
		return
	}
	a.Dispose()
	p.Destroy()
}

func (p *playerVessel) IsDisposed() bool { return p.body.IsDisposed() }

func (p *playerVessel) Dispose() {
	p.body.Dispose()
	p.sprite.Dispose()
}

func (p *playerVessel) Update(delta float64) {
	p.reload = gemath.ClampMin(p.reload-delta, 0)

	if p.ctx.Input.ActionIsPressed(ActionLeft) {
		p.body.Rotation -= 0.05
	}
	if p.ctx.Input.ActionIsPressed(ActionRight) {
		p.body.Rotation += 0.05
	}
	if p.ctx.Input.ActionIsPressed(ActionForward) {
		p.body.Pos.MoveInDirection(100*delta, p.body.Rotation)
	}

	if p.reload == 0 && p.ctx.Input.ActionIsPressed(ActionFire) {
		pos := p.body.Pos.MoveInDirectionResult(64, p.body.Rotation)
		obj := newBullet(pos, p.body.Rotation)
		p.ctx.CurrentScene.AddObject(obj)
		p.reload += 0.4
	}
}

func (p *playerVessel) Destroy() {
	p.EventDestroyed.Emit(gesignal.Void{})
	p.ctx.CurrentScene.AddObject(newVesselExplosion(p.body.Pos))
	p.Dispose()
}
