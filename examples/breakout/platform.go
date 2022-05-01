package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/physics"
)

type platform struct {
	ctx    *ge.Context
	body   physics.Body
	sprite *ge.Sprite
	ball   *ball
}

func newPlatform() *platform {
	p := &platform{}
	p.body.InitRotatedRect(p, 100, 22)
	return p
}

func (p *platform) Init(scene *ge.Scene) {
	p.ctx = scene.Context()
	p.sprite = p.ctx.LoadSprite("platform.png")
	p.sprite.Pos = &p.body.Pos
	p.sprite.Rotation = &p.body.Rotation
	scene.AddGraphics(p.sprite)
	scene.AddBody(&p.body)
}

func (p *platform) IsDisposed() bool { return p.body.IsDisposed() }

func (p *platform) Dispose() {
	p.body.Dispose()
	p.sprite.Dispose()
}

func (p *platform) Update(delta float64) {
	moving := false
	if p.ctx.Input.ActionIsPressed(ActionLeft) {
		moving = true
		p.body.Pos.X -= 250 * delta
		p.body.Rotation = gemath.Rad(gemath.ClampMin(float64(p.body.Rotation)-1.5*delta, -0.3))
	} else if p.ctx.Input.ActionIsPressed(ActionRight) {
		moving = true
		p.body.Pos.X += 250 * delta
		p.body.Rotation = gemath.Rad(gemath.ClampMax(float64(p.body.Rotation)+1.5*delta, 0.3))
	}
	if !moving {
		if p.body.Rotation < 0 {
			p.body.Rotation = gemath.Rad(gemath.ClampMax(float64(p.body.Rotation)+1.1*delta, 0))
		} else if p.body.Rotation > 0 {
			p.body.Rotation = gemath.Rad(gemath.ClampMin(float64(p.body.Rotation)-1.1*delta, 0))
		}
	}

	if p.ball != nil && p.ball.IsDisposed() {
		p.ball = nil
	}
	canFire := p.ball == nil
	if canFire && p.ctx.Input.ActionIsPressed(ActionFire) {
		x, y := ebiten.CursorPosition()
		b := newBall()
		b.velocity = p.body.Pos.VecTowards(350, gemath.Vec{X: float64(x), Y: float64(y)})
		b.body.Pos = gemath.Vec{X: p.body.Pos.X, Y: p.body.Pos.Y - 40}
		p.ctx.CurrentScene.AddObject(b)
		p.ball = b
	}
}
