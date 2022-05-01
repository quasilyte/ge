package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/physics"
)

type ball struct {
	ctx      *ge.Context
	body     physics.Body
	sprite   *ge.Sprite
	velocity gemath.Vec
}

func newBall() *ball {
	b := &ball{}
	b.body.InitCircle(b, 12)
	return b
}

func (b *ball) Init(scene *ge.Scene) {
	b.ctx = scene.Context()
	b.sprite = b.ctx.LoadSprite("ball.png")
	b.sprite.Pos = &b.body.Pos
	scene.AddGraphics(b.sprite)
	scene.AddBody(&b.body)
}

func (b *ball) IsDisposed() bool { return b.body.IsDisposed() }

func (b *ball) Dispose() {
	b.body.Dispose()
	b.sprite.Dispose()
}

func (b *ball) Update(delta float64) {
	bounced := false
	if collision := b.ctx.CurrentScene.GetMovementCollision(&b.body, b.velocity); collision != nil {
		extraRotation := float64(0)
		switch o := collision.Body.Object.(type) {
		case *brick:
			bounced = !o.Hit()
			for i := 0; i < 3; i++ {
				shardPos := b.body.Pos.Add(gemath.Vec{X: float64(i*8) - 8})
				shard := newBrickShard(shardPos)
				b.ctx.CurrentScene.AddObject(shard)
			}
		case *platform:
			platformRotation := o.body.Rotation
			if platformRotation < 0 {
				platformRotation += math.Pi
			}
			delta := b.velocity.Normalized().Angle().AngleDelta(platformRotation)
			if math.Abs(float64(delta)) < 0.4 {
				if delta > 0 {
					extraRotation = 0.45 - float64(delta)
				} else {
					extraRotation = -0.45 - float64(delta)
				}
			}
			bounced = true
		default:
			bounced = true
		}
		if bounced {
			b.body.Pos = b.body.Pos.Add(collision.Normal.Mulf(collision.Depth + 6))
			b.velocity = b.calculateReflectedVec(collision.Normal.Rotated(gemath.Rad(extraRotation)))
		}
	}
	if !bounced {
		b.body.Pos = b.body.Pos.Add(b.velocity.Mulf(delta))
	}

	if !b.ctx.WindowRect().Contains(b.body.Pos) {
		b.Dispose()
	}
}

func (b *ball) calculateReflectedVec(n gemath.Vec) gemath.Vec {
	v := b.velocity
	return v.Sub(n.Mulf(2 * v.Dot(n)))
}
