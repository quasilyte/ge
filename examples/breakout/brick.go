package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/gesignal"
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

	EventDestroyed gesignal.Event[gesignal.Void]
}

func newBrick(scale float64, rotation gemath.Rad, pos gemath.Vec) *brick {
	b := &brick{hp: 4, scale: scale, shapeHeight: 32}
	b.body.InitRotatedRect(b, brickDefaultWidth*scale, 32*scale)
	b.body.Rotation = rotation
	b.body.Pos = pos
	return b
}

func newCircleBrick(scale float64, pos gemath.Vec) *brick {
	b := &brick{
		hp:          3,
		scale:       scale,
		shapeHeight: brickDefaultWidth,
	}
	b.body.InitCircle(b, 32*scale)
	b.body.Pos = pos
	return b
}

func (b *brick) Init(scene *ge.Scene) {
	b.scene = scene
	if b.body.IsCircle() {
		b.sprite = scene.NewSprite(ImageBrickCircle)
	} else {
		b.sprite = scene.NewSprite(ImageBrickRect)
		b.sprite.Rotation = &b.body.Rotation
	}
	b.sprite.FrameWidth = brickDefaultWidth
	b.sprite.Pos.Base = &b.body.Pos
	b.sprite.Scale = b.scale
	scene.AddGraphics(b.sprite)
	scene.AddBody(&b.body)
}

func (b *brick) Hit(hitPos gemath.Vec) bool {
	b.hp--
	b.sprite.FrameOffset.X += brickDefaultWidth

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

	b.EventDestroyed.Emit(gesignal.Void{})
	b.Dispose()
}
