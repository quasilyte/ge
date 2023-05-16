package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type platform struct {
	scene    *ge.Scene
	body     physics.Body
	sprite   *ge.Sprite
	ball     *ball
	numLives int

	input *input.Handler

	EventBallLost gesignal.Event[gesignal.Void]
}

func newPlatform(h *input.Handler) *platform {
	p := &platform{numLives: 4, input: h}
	p.body.InitRotatedRect(p, 100, 22)
	return p
}

func (p *platform) Init(scene *ge.Scene) {
	p.scene = scene
	p.sprite = p.scene.NewSprite(ImagePlatform)
	p.sprite.Pos.Base = &p.body.Pos
	p.sprite.Rotation = &p.body.Rotation
	scene.AddGraphics(p.sprite)
	scene.AddBody(&p.body)
}

func (p *platform) IsDisposed() bool { return p.body.IsDisposed() }

func (p *platform) Dispose() {
	p.body.Dispose()
	p.sprite.Dispose()
	if p.ball != nil {
		p.ball.Dispose()
	}
}

func (p *platform) Update(delta float64) {
	moving := false
	if p.input.ActionIsPressed(ActionLeft) {
		moving = true
		p.body.Pos.X -= 250 * delta
		p.body.Rotation = gmath.Rad(gmath.ClampMin(float64(p.body.Rotation)-1.5*delta, -0.3))
	} else if p.input.ActionIsPressed(ActionRight) {
		moving = true
		p.body.Pos.X += 250 * delta
		p.body.Rotation = gmath.Rad(gmath.ClampMax(float64(p.body.Rotation)+1.5*delta, 0.3))
	}
	if !moving {
		if p.body.Rotation < 0 {
			p.body.Rotation = gmath.Rad(gmath.ClampMax(float64(p.body.Rotation)+1.1*delta, 0))
		} else if p.body.Rotation > 0 {
			p.body.Rotation = gmath.Rad(gmath.ClampMin(float64(p.body.Rotation)-1.1*delta, 0))
		}
	}

	if p.ball != nil && p.ball.IsDisposed() {
		p.ball = nil
	}
	if p.ball == nil && p.numLives > 0 && p.input.ActionIsPressed(ActionFire) {
		b := newBall()
		b.EventDestroyed.Connect(p, p.onBallDestroyed)
		p.ball = b
		b.velocity = gmath.Vec{X: 0, Y: -300}
		b.body.Pos = gmath.Vec{X: p.body.Pos.X, Y: p.body.Pos.Y - 40}
		p.scene.AddObject(b)
	}
}

func (p *platform) onBallDestroyed(gesignal.Void) {
	p.numLives--
	p.EventBallLost.Emit(gesignal.Void{})
}
