package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/physics"
)

const brickDefaultWidth float64 = 64
const brickShardWidth float64 = 8
const brickShardHeight float64 = 8

type brick struct {
	scene  *ge.Scene
	body   physics.Body
	sprite *ge.Sprite
	scale  float64
	hp     float64

	shapeHeight float64
}

func newBrick(scale float64, rotation gemath.Rad) *brick {
	b := &brick{hp: 4, scale: scale, shapeHeight: 32}
	b.body.InitRotatedRect(b, brickDefaultWidth*scale, 32*scale)
	b.body.Rotation = rotation
	return b
}

func newCircleBrick(scale float64) *brick {
	b := &brick{hp: 3, scale: scale, shapeHeight: brickDefaultWidth}
	b.body.InitCircle(b, 32*scale)
	return b
}

func (b *brick) Init(scene *ge.Scene) {
	b.scene = scene
	if b.body.IsCircle() {
		b.sprite = scene.LoadSprite("brick_circle.png")
	} else {
		b.sprite = scene.LoadSprite("brick_purple.png")
		b.sprite.Rotation = &b.body.Rotation
	}
	b.sprite.Width = brickDefaultWidth
	b.sprite.Pos = &b.body.Pos
	b.sprite.Scaling = b.scale
	scene.AddGraphics(b.sprite)
	scene.AddBody(&b.body)
}

func (b *brick) Hit(hitPos gemath.Vec) bool {
	b.hp--
	b.sprite.Offset.X += brickDefaultWidth

	if b.hp <= 0 {
		b.scene.Audio().PlaySound(AudioBrickDestroyed)
		b.Destroy()
		return true
	}

	for i := 0; i < 3; i++ {
		shardPos := hitPos.Add(gemath.Vec{X: float64(i*8) - 8})
		shard := newBrickShard(shardPos)
		b.scene.AddObject(shard)
	}

	b.scene.Audio().PlaySound(AudioBrickHit)
	return false
}

func (b *brick) IsDisposed() bool { return b.body.IsDisposed() }

func (b *brick) Update(delta float64) {}

func (b *brick) Dispose() {
	b.body.Dispose()
	b.sprite.Dispose()
}

func (b *brick) Destroy() {
	width := brickDefaultWidth * b.scale
	height := b.shapeHeight * b.scale
	offset := gemath.Vec{
		X: b.body.Pos.X - width/2,
		Y: b.body.Pos.Y - height/2,
	}
	startX := float64(0)
	startY := float64(0)
	if b.body.IsCircle() {
		startX = 8
		startY = 8
		width -= 8
		height -= 8
	}
	for y := startY; y < height; y += 8 {
		for x := startX; x < width; x += 8 {
			unrotatedPos := offset.Add(gemath.Vec{X: x, Y: y})
			center := b.body.Pos
			diff := unrotatedPos.Sub(center)
			rotatedPos := diff.Rotated(b.body.Rotation).Add(center)
			shard := newBrickShard(rotatedPos)
			b.scene.AddObject(shard)
		}
	}

	b.Dispose()
}
