package main

import (
	"fmt"

	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
)

type battleFactory struct {
	alliance    int
	sprite      *ge.Sprite
	Body        physics.Body
	scene       *ge.Scene
	hp          float64
	hpComponent *healthComponent
	battleState *battleState

	techLevel          int
	maxUnits           int
	numUnits           int
	hasProductionOrder bool
	productionDelay    float64
	productionOrder    battleUnitConfig
	productionOptions  []*standardBodyDesign

	botProgramPicker *gmath.RandPicker[botProgramKind]

	EventDestroyed gesignal.Event[*battleFactory]
}

type battleFactoryConfig struct {
	alliance    int
	techLevel   int
	maxUnits    int
	pos         gmath.Vec
	rotation    gmath.Rad
	battleState *battleState
}

func newBattleFactory(config battleFactoryConfig) *battleFactory {
	b := &battleFactory{alliance: config.alliance}
	b.Body.Pos = config.pos
	b.Body.Rotation = config.rotation
	b.battleState = config.battleState
	b.maxUnits = config.maxUnits
	b.techLevel = config.techLevel

	if b.techLevel < 1 || b.techLevel > 6 {
		panic(fmt.Sprintf("invalid base tech level: %d", b.techLevel))
	}

	b.productionOptions = make([]*standardBodyDesign, 0, len(standardBodyDesignList))
	for _, design := range standardBodyDesignList {
		if b.techLevel >= design.techLevel {
			b.productionOptions = append(b.productionOptions, design)
		}
	}

	return b
}

func (b *battleFactory) Pos() gmath.Vec { return b.Body.Pos }

func (b *battleFactory) Init(scene *ge.Scene) {
	b.scene = scene

	b.hp = 40 + (5 * float64(b.techLevel-1))
	b.hpComponent = newHealthComponent(&b.hp, &b.Body)
	b.hpComponent.isBuilding = true

	b.Body.InitStaticCircle(b, 32)
	b.Body.LayerMask = unitLayerMask(b.alliance) | buildingLayerMask
	scene.AddBody(&b.Body)

	b.sprite = scene.NewSprite(ImageFactory)
	SetHue(b.sprite, spriteHue(b.alliance))
	b.sprite.Pos.Base = &b.Body.Pos
	b.sprite.Rotation = &b.Body.Rotation
	b.sprite.Shader = scene.NewShader(ShaderBuildingDamage)
	b.sprite.Shader.SetFloatValue("HP", 1.0)
	b.sprite.Shader.Texture1 = scene.LoadImage(gmath.RandElem(scene.Rand(), damageMaskImages))
	scene.AddGraphics(b.sprite)

	b.botProgramPicker = gmath.NewRandPicker[botProgramKind](b.battleState.rand)
	b.botProgramPicker.AddOption(botDelayedAttack, 0.15)
	b.botProgramPicker.AddOption(botRoam, 0.2)
	b.botProgramPicker.AddOption(botBaseAttack, 0.35)
	b.botProgramPicker.AddOption(botTankHunt, 0.45)

	b.productionOrder.alliance = b.alliance
	b.productionOrder.pos = b.Body.Pos.Add(cellDelta(b.Body.Rotation))
	b.productionOrder.rotation = b.Body.Rotation
	b.productionOrder.weaponReloadMultiplier = 1
	b.productionOrder.battleState = b.battleState
	b.maybeOrderProduction()
}

func (b *battleFactory) IsDisposed() bool {
	return b.Body.IsDisposed()
}

func (b *battleFactory) Dispose() {
	b.Body.Dispose()
	b.sprite.Dispose()
}

func (b *battleFactory) Destroy() {
	b.EventDestroyed.Emit(b)
	createExplosions(b.scene, b.Body.Pos, 3, 4)
	b.Dispose()
}

func (b *battleFactory) Update(delta float64) {
	if !b.hpComponent.CheckProjectileCollisions(b.scene) {
		b.Destroy()
		return
	}
	b.sprite.Shader.Enabled = b.hpComponent.hpPercentage != 1
	b.sprite.Shader.SetFloatValue("HP", b.hpComponent.hpPercentage)

	b.productionDelay = gmath.ClampMin(b.productionDelay-delta, 0)
	if b.hasProductionOrder && b.productionDelay == 0 {
		u := b.battleState.newBattleUnit(b.productionOrder)
		botConfig := localBotConfig{
			state: b.battleState,
			unit:  u,
		}
		if b.battleState.forceBaseAttack {
			botConfig.program = botBaseAttack
		} else {
			botConfig.program = b.botProgramPicker.Pick()
		}
		bot := newLocalBot(botConfig)
		b.scene.AddObject(bot)
		b.scene.AddObject(u)
		b.hasProductionOrder = false
		b.numUnits++
		b.maybeOrderProduction()
		dist := float64(b.scene.Rand().IntRange(2, 5))
		dest := fixedPos(u.Body.Pos.Add(cellDelta(b.Body.Rotation).Mulf(dist)))
		bot.SetDestination(dest)

		u.EventDestroyed.Connect(nil, func(u *battleUnit) {
			b.numUnits--
			b.maybeOrderProduction()
		})
	}
}

func (b *battleFactory) maybeOrderProduction() {
	if b.numUnits >= b.maxUnits {
		return
	}

	design := gmath.RandElem(b.battleState.rand, b.productionOptions)

	b.productionOrder.maxHP = design.maxHP
	b.productionOrder.image = design.image
	b.productionOrder.speed = design.speed
	b.productionOrder.rotationTime = design.rotationTime
	b.productionOrder.weaponReloadMultiplier = design.reloadModified

	if b.techLevel == 6 && b.battleState.rand.Chance(0.15) {
		b.productionOrder.weapon = railgunDesign
	} else {
		advancedWeaponChance := 0.0
		switch b.techLevel {
		case 3:
			advancedWeaponChance = 0.1
		case 4:
			advancedWeaponChance = 0.3
		case 5, 6:
			advancedWeaponChance = 0.7
		}
		if advancedWeaponChance != 0 && b.battleState.rand.Chance(advancedWeaponChance) {
			b.productionOrder.weapon = gmath.RandElem(b.battleState.rand, advancedWeaponDesignList)
		} else {
			b.productionOrder.weapon = gmath.RandElem(b.battleState.rand, simpleWeaponDesignList)
		}
	}

	productionTime := 1.0 - (float64(b.techLevel-1) * 0.1)
	b.productionDelay = design.productionTime * productionTime
	b.hasProductionOrder = true
}
