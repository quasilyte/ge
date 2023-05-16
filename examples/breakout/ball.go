package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type ball struct {
	scene    *ge.Scene
	body     physics.Body
	sprite   *ge.Sprite
	velocity gmath.Vec

	EventDestroyed gesignal.Event[gesignal.Void]
}

func newBall() *ball {
	b := &ball{}
	b.body.InitCircle(b, 12)
	return b
}

func (b *ball) Init(scene *ge.Scene) {
	b.scene = scene
	b.sprite = scene.NewSprite(ImageBall)
	b.sprite.Pos.Base = &b.body.Pos
	scene.AddGraphics(b.sprite)
	scene.AddBody(&b.body)
}

func (b *ball) IsDisposed() bool { return b.body.IsDisposed() }

func (b *ball) Dispose() {
	b.body.Dispose()
	b.sprite.Dispose()
}

func (b *ball) Update(delta float64) {
	b.handleMovement(delta)

	if !b.scene.Context().WindowRect().Contains(b.body.Pos) {
		b.EventDestroyed.Emit(gesignal.Void{})
		b.Dispose()
	}
}

func (b *ball) handleMovement(delta float64) {
	if collision := b.scene.GetMovementCollision(&b.body, b.velocity); collision != nil {
		bounced := false
		extraRotation := float64(0)
		switch o := collision.Body.Object.(type) {
		case *brick:
			// If brick is destroyed, we'll go through it without a bounce.
			bounced = !o.Hit(b.body.Pos)
		case *platform:
			// Since platform can rotate, it can result in a tricky corner cases
			// where default reflection calculations can fail.
			// We'll try to fix these situations here by adding an extra rotation.
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
			// Collided with wall, etc.
			bounced = true
		}
		if bounced {
			b.body.Pos = b.body.Pos.Add(collision.Normal.Mulf(collision.Depth + 6))
			b.velocity = b.reflect(collision.Normal.Rotated(gmath.Rad(extraRotation)))
			return
		}
	}

	b.body.Pos = b.body.Pos.Add(b.velocity.Mulf(delta))
}

func (b *ball) reflect(n gmath.Vec) gmath.Vec {
	v := b.velocity
	return v.Sub(n.Mulf(2 * v.Dot(n)))
}
