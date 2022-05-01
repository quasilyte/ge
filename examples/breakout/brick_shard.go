package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

type brickShard struct {
	pos      gemath.Vec
	velocity gemath.Vec
	sprite   *ge.Sprite
}

func newBrickShard(pos gemath.Vec) *brickShard {
	return &brickShard{pos: pos}
}

func (shard *brickShard) Init(scene *ge.Scene) {
	angle := scene.Rand().Rad()
	shard.velocity = gemath.RadToVec(angle).Mulf(100)
	shard.sprite = scene.LoadSprite(ImageBrickShard)
	shard.sprite.Pos = &shard.pos
	scene.AddGraphics(shard.sprite)
}

func (shard *brickShard) IsDisposed() bool { return shard.sprite.IsDisposed() }

func (shard *brickShard) Dispose() {
	shard.sprite.Dispose()
}

func (shard *brickShard) Update(delta float64) {
	shard.pos = shard.pos.Add(shard.velocity.Mulf(delta))

	shard.sprite.ColorModulation.A -= float32(delta)
	if shard.sprite.ColorModulation.A < 0.2 {
		shard.Dispose()
	}
}
