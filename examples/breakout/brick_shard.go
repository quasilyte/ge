package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type brickShard struct {
	pos      gmath.Vec
	velocity gmath.Vec
	sprite   *ge.Sprite
}

func newBrickShard(pos gmath.Vec) *brickShard {
	return &brickShard{pos: pos}
}

func (shard *brickShard) Init(scene *ge.Scene) {
	angle := scene.Rand().Rad()
	shard.velocity = gmath.RadToVec(angle).Mulf(100)
	shard.sprite = scene.NewSprite(ImageBrickShard)
	shard.sprite.Pos.Base = &shard.pos
	scene.AddGraphics(shard.sprite)
}

func (shard *brickShard) IsDisposed() bool { return shard.sprite.IsDisposed() }

func (shard *brickShard) Dispose() {
	shard.sprite.Dispose()
}

func (shard *brickShard) Update(delta float64) {
	shard.pos = shard.pos.Add(shard.velocity.Mulf(delta))

	shard.sprite.SetAlpha(shard.sprite.GetAlpha() - float32(delta))
	if shard.sprite.GetAlpha() < 0.2 {
		shard.Dispose()
	}
}
