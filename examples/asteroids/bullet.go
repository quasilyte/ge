package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/collision"
	"github.com/quasilyte/ge/gemath"
)

type bullet struct {
	ctx    *ge.Context
	body   collision.Body
	sprite *ge.Sprite
	hp     float64
}

func newBullet(pos gemath.Vec, rotation gemath.Rad) *bullet {
	b := &bullet{
		hp: 400,
	}
	b.body.InitRotatedRect(b, 4, 8)
	b.body.Pos = pos
	b.body.Rotation = rotation
	return b
}

func (b *bullet) Init(scene *ge.Scene) {
	ctx := scene.Context()

	b.ctx = ctx
	b.sprite = ctx.LoadSprite("bullet.png")
	b.sprite.Pos = &b.body.Pos
	b.sprite.Rotation = &b.body.Rotation
	scene.AddGraphics(b.sprite)
	scene.AddBody(&b.body)

	ctx.Audio.PlayWAV(ctx.Loader.LoadWAV("fire.wav"), true)
}

func (b *bullet) IsDisposed() bool { return b.body.IsDisposed() }

func (b *bullet) Dispose() {
	b.body.Dispose()
	b.sprite.Dispose()
}

func (b *bullet) Destroy() {
	b.Dispose()
	b.ctx.CurrentScene.AddObject(newBulletExplosion(b.body.Pos))
}

func (b *bullet) Update(delta float64) {
	if b.hp <= 0 {
		b.Dispose()
		return
	}
	const bulletSpeed float64 = 400
	travelled := bulletSpeed * delta
	b.hp -= travelled
	b.body.Pos.MoveInDirection(travelled, b.body.Rotation)
}
