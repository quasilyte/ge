package main

import (
	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
)

type battleWall struct {
	alliance       int
	gridX          int
	gridY          int
	sprite         *ge.Sprite
	Body           physics.Body
	scene          *ge.Scene
	hp             float64
	hpComponent    *healthComponent
	EventDestroyed gesignal.Event[*battleWall]
}

func newBattleWall(pos gmath.Vec, alliance int) *battleWall {
	b := &battleWall{alliance: alliance}
	b.Body.Pos = pos
	b.gridX = int(pos.X) / 64
	b.gridY = int(pos.Y) / 64
	return b
}

func (b *battleWall) Pos() gmath.Vec { return b.Body.Pos }

func (b *battleWall) Init(scene *ge.Scene) {
	b.scene = scene

	b.hp = 30
	b.hpComponent = newHealthComponent(&b.hp, &b.Body)
	b.hpComponent.isBuilding = true

	b.Body.InitStaticCircle(b, 32)
	b.Body.LayerMask = unitLayerMask(b.alliance) | wallLayerMask
	scene.AddBody(&b.Body)

	b.sprite = scene.NewSprite(ImageWall)
	SetHue(b.sprite, spriteHue(b.alliance))
	b.sprite.Pos.Base = &b.Body.Pos
	b.sprite.Shader = scene.NewShader(ShaderBuildingDamage)
	b.sprite.Shader.SetFloatValue("HP", 1.0)
	b.sprite.Shader.Texture1 = scene.LoadImage(gmath.RandElem(scene.Rand(), damageMaskImages))
	scene.AddGraphics(b.sprite)
}

func (b *battleWall) IsDisposed() bool {
	return b.Body.IsDisposed()
}

func (b *battleWall) Dispose() {
	b.Body.Dispose()
	b.sprite.Dispose()
}

func (b *battleWall) Destroy() {
	b.EventDestroyed.Emit(b)
	createExplosions(b.scene, b.Body.Pos, 3, 4)
	b.Dispose()
}

func (b *battleWall) Update(delta float64) {
	if !b.hpComponent.CheckProjectileCollisions(b.scene) {
		b.Destroy()
		return
	}
	b.sprite.Shader.Enabled = b.hpComponent.hpPercentage != 1
	b.sprite.Shader.SetFloatValue("HP", b.hpComponent.hpPercentage)
}

func (b *battleWall) SetFrame(mask uint8) {
	index := int(mask)
	b.sprite.FrameOffset.X = float64(index * 64)
}
