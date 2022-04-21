package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/collision"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/gesignal"
)

type asteroid struct {
	ctx    *ge.Context
	body   collision.Body
	sprite *ge.Sprite
	shards int
	speed  float64
	hp     float64

	EventDestroyed    gesignal.Event[*asteroid]
	EventShardCreated gesignal.Event[*asteroid]
}

func newAsteroid(pos gemath.Vec, speed float64, shards int) *asteroid {
	a := &asteroid{
		shards: shards,
		hp:     100,
		speed:  speed,
	}
	a.body.InitCircle(a, 24)
	a.body.CollisionHandler = a
	a.body.Pos = pos
	return a
}

func (a *asteroid) Init(scene *ge.Scene) {
	a.ctx = scene.Context()
	a.body.Rotation = a.ctx.Rand.Rad()
	a.sprite = scene.Context().LoadSprite("asteroid.png")
	switch a.shards {
	case 3:
		a.sprite.Scaling = 1.6
	case 2:
		a.sprite.Scaling = 1.25
	case 1:
		a.sprite.Scaling = 1
	default:
		a.sprite.Scaling = 0.75
	}
	a.sprite.Pos = &a.body.Pos
	a.sprite.Rotation = ge.NewRotation(0)
	scene.AddGraphics(a.sprite)
	scene.AddBody(&a.body)
}

func (a *asteroid) OnCollision(info *collision.Info) {
	b, ok := info.Object.(*bullet)
	if !ok {
		return
	}
	b.Destroy()

	a.hp -= 50
	if a.hp <= 0 {
		a.Destroy()
	}
}

func (a *asteroid) IsDisposed() bool { return a.body.IsDisposed() }

func (a *asteroid) Dispose() {
	a.body.Dispose()
	a.sprite.Dispose()
}

func (a *asteroid) Update(delta float64) {
	if a.body.Pos.X > a.ctx.WindowWidth {
		a.body.Pos.X = 0
		a.body.Rotation += 0.2
	} else if a.body.Pos.X < 0 {
		a.body.Pos.X = a.ctx.WindowWidth
		a.body.Rotation -= 0.2
	}
	if a.body.Pos.Y > a.ctx.WindowHeight {
		a.body.Pos.Y = 0
		a.body.Rotation += 0.2
	} else if a.body.Pos.Y < 0 {
		a.body.Pos.Y = a.ctx.WindowHeight
		a.body.Rotation -= 0.2
	}

	*a.sprite.Rotation += gemath.Rad(delta * 2)

	travelled := a.speed * delta
	a.body.Pos.MoveInDirection(travelled, a.body.Rotation)
}

func (a *asteroid) CreateShards() {
	for i := 0; i < a.shards; i++ {
		shard := newAsteroid(a.body.Pos, a.speed-25, a.shards-1)
		a.EventShardCreated.Emit(shard)
		a.ctx.CurrentScene.AddObject(shard)
	}
}

func (a *asteroid) Destroy() {
	a.CreateShards()
	a.EventDestroyed.Emit(a)
	a.Dispose()
}
