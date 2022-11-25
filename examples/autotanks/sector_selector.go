package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type sectorSelector struct {
	playerID int
	sprite   *ge.Sprite
	sector   *sector
	pos      gmath.Vec
}

func newSectorSelector(playerID int) *sectorSelector {
	return &sectorSelector{playerID: playerID}
}

func (selector *sectorSelector) Init(scene *ge.Scene) {
	selector.sprite = scene.NewSprite(ImageSectorSelector)
	selector.sprite.Pos.Base = &selector.pos
	rotation := scene.Rand().Rad()
	selector.sprite.Rotation = &rotation
	applyPlayerColor(selector.playerID, selector.sprite)
	scene.AddGraphics(selector.sprite)
}

func (selector *sectorSelector) IsDisposed() bool { return false }

func (selector *sectorSelector) Update(delta float64) {
	*selector.sprite.Rotation += gmath.Rad(delta) * 0.3
}

func (selector *sectorSelector) SetSector(s *sector) {
	selector.sector = s
	selector.pos = s.Center()
}

func (selector *sectorSelector) Sector() *sector {
	return selector.sector
}
