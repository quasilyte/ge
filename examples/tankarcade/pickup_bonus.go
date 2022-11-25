package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/gmath"
)

type pickupBonusKind uint8

const (
	pickupHP pickupBonusKind = iota
	pickupAmmo
)

var pickupKindList = []pickupBonusKind{
	pickupHP,
	pickupAmmo,
}

type pickupBonus struct {
	scene  *ge.Scene
	kind   pickupBonusKind
	sprite *ge.Sprite
	body   physics.Body

	EventDestroyed gesignal.Event[*pickupBonus]
}

func newPickupBonus(pos gmath.Vec, kind pickupBonusKind) *pickupBonus {
	b := &pickupBonus{kind: kind}
	b.body.Pos = pos
	return b
}

func (b *pickupBonus) Init(scene *ge.Scene) {
	b.scene = scene

	var spriteID resource.ImageID
	switch b.kind {
	case pickupHP:
		spriteID = ImagePickupHP
	case pickupAmmo:
		spriteID = ImagePickupAmmo
	default:
		panic("unexpected pickup bonus kind")
	}
	b.sprite = scene.NewSprite(spriteID)
	b.sprite.Pos.Base = &b.body.Pos
	scene.AddGraphics(b.sprite)

	b.body.InitStaticCircle(b, 16)
	b.body.LayerMask = 0b11
	scene.AddBody(&b.body)
}

func (b *pickupBonus) IsDisposed() bool {
	return b.sprite.IsDisposed()
}

func (b *pickupBonus) Destroy() {
	e := newExplosion(b.body.Pos)
	e.Image = ImageExplosion2
	e.Hue = 2
	b.scene.AddObject(e)

	b.EventDestroyed.Emit(b)
	b.Dispose()
}

func (b *pickupBonus) Dispose() {
	b.sprite.Dispose()
	b.body.Dispose()
}

func (b *pickupBonus) Update(delta float64) {}
