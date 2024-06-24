package main

import (
	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
)

type battleBonusGenerator struct {
	alliance    int
	sprite      *ge.Sprite
	Body        physics.Body
	scene       *ge.Scene
	hp          float64
	hpComponent *healthComponent

	bonusAlive         bool
	hasProductionOrder bool
	productionDelay    float64
	productionOrder    pickupBonusKind

	EventDestroyed gesignal.Event[*battleBonusGenerator]
}

type battleBonusGeneratorConfig struct {
	alliance int
	pos      gmath.Vec
	rotation gmath.Rad
}

func newBattleBonusGenerator(config battleBonusGeneratorConfig) *battleBonusGenerator {
	b := &battleBonusGenerator{alliance: config.alliance}
	b.Body.Pos = config.pos
	b.Body.Rotation = config.rotation
	return b
}

func (b *battleBonusGenerator) Pos() gmath.Vec { return b.Body.Pos }

func (b *battleBonusGenerator) Init(scene *ge.Scene) {
	b.scene = scene

	b.hp = 35
	b.hpComponent = newHealthComponent(&b.hp, &b.Body)
	b.hpComponent.isBuilding = true

	b.Body.InitStaticCircle(b, 32)
	b.Body.LayerMask = unitLayerMask(b.alliance) | buildingLayerMask
	scene.AddBody(&b.Body)

	b.sprite = scene.NewSprite(ImageBonusGenerator)
	SetHue(b.sprite, spriteHue(b.alliance))
	b.sprite.Pos.Base = &b.Body.Pos
	b.sprite.Rotation = &b.Body.Rotation
	b.sprite.Shader = scene.NewShader(ShaderBuildingDamage)
	b.sprite.Shader.SetFloatValue("HP", 1.0)
	b.sprite.Shader.Texture1 = scene.LoadImage(gmath.RandElem(scene.Rand(), damageMaskImages))
	scene.AddGraphics(b.sprite)

	b.maybeOrderProduction()
}

func (b *battleBonusGenerator) IsDisposed() bool {
	return b.Body.IsDisposed()
}

func (b *battleBonusGenerator) Dispose() {
	b.Body.Dispose()
	b.sprite.Dispose()
}

func (b *battleBonusGenerator) Destroy() {
	b.EventDestroyed.Emit(b)
	createExplosions(b.scene, b.Body.Pos, 3, 4)
	b.Dispose()
}

func (b *battleBonusGenerator) Update(delta float64) {
	if !b.hpComponent.CheckProjectileCollisions(b.scene) {
		b.Destroy()
		return
	}
	b.sprite.Shader.Enabled = b.hpComponent.hpPercentage != 1
	b.sprite.Shader.SetFloatValue("HP", b.hpComponent.hpPercentage)

	b.productionDelay = gmath.ClampMin(b.productionDelay-delta, 0)
	if b.hasProductionOrder && b.productionDelay == 0 {
		spawnPos := b.Body.Pos.Add(cellDelta(b.Body.Rotation))
		bonus := newPickupBonus(spawnPos, b.productionOrder)
		b.scene.AddObject(bonus)
		b.bonusAlive = true
		b.hasProductionOrder = false
		bonus.EventDestroyed.Connect(nil, func(u *pickupBonus) {
			b.bonusAlive = false
			b.maybeOrderProduction()
		})
	}
}

func (b *battleBonusGenerator) maybeOrderProduction() {
	if b.bonusAlive {
		return
	}

	kind := gmath.RandElem(b.scene.Rand(), pickupKindList)
	b.productionOrder = kind
	b.productionDelay = b.scene.Rand().FloatRange(20, 30)
	b.hasProductionOrder = true
}
