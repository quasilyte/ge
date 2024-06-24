package main

import (
	"fmt"

	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
)

type battleBase struct {
	alliance       int
	level          int
	sprite         *ge.Sprite
	Body           physics.Body
	scene          *ge.Scene
	hp             float64
	hpComponent    *healthComponent
	EventDestroyed gesignal.Event[*battleBase]
}

func newBattleBase(pos gmath.Vec, alliance, level int) *battleBase {
	b := &battleBase{alliance: alliance, level: level}
	b.Body.Pos = pos
	return b
}

func (b *battleBase) Pos() gmath.Vec { return b.Body.Pos }

func (b *battleBase) Init(scene *ge.Scene) {
	b.scene = scene

	switch b.level {
	case 1:
		b.hp = 50
	case 2:
		b.hp = 90
	case 3:
		b.hp = 120
	default:
		panic(fmt.Sprintf("unexpected base level: %d", b.level))
	}
	b.hpComponent = newHealthComponent(&b.hp, &b.Body)
	b.hpComponent.isBuilding = true

	b.Body.InitStaticCircle(b, 32)
	b.Body.LayerMask = unitLayerMask(b.alliance) | buildingLayerMask
	scene.AddBody(&b.Body)

	b.sprite = scene.NewSprite(ImageBase)
	SetHue(b.sprite, spriteHue(b.alliance))
	b.sprite.Pos.Base = &b.Body.Pos
	b.sprite.Shader = scene.NewShader(ShaderBuildingDamage)
	b.sprite.Shader.SetFloatValue("HP", 1.0)
	b.sprite.Shader.Texture1 = scene.LoadImage(gmath.RandElem(scene.Rand(), damageMaskImages))
	scene.AddGraphics(b.sprite)
}

func (b *battleBase) IsDisposed() bool {
	return b.Body.IsDisposed()
}

func (b *battleBase) Dispose() {
	b.Body.Dispose()
	b.sprite.Dispose()
}

func (b *battleBase) Destroy() {
	b.EventDestroyed.Emit(b)
	createExplosions(b.scene, b.Body.Pos, 4, 5)
	b.Dispose()
}

func (b *battleBase) Update(delta float64) {
	if !b.hpComponent.CheckProjectileCollisions(b.scene) {
		b.Destroy()
		return
	}
	b.sprite.Shader.Enabled = b.hpComponent.hpPercentage != 1
	b.sprite.Shader.SetFloatValue("HP", b.hpComponent.hpPercentage)
}
