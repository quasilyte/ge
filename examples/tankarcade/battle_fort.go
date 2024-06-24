package main

import (
	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
)

type battleFort struct {
	alliance       int
	sprite         *ge.Sprite
	turretSprite   *ge.Sprite
	Body           physics.Body
	scene          *ge.Scene
	hp             float64
	hpComponent    *healthComponent
	battleState    *battleState
	reload         float64
	flamer         bool
	EventDestroyed gesignal.Event[*battleFort]
}

type battleFortConfig struct {
	alliance    int
	pos         gmath.Vec
	rotation    gmath.Rad
	battleState *battleState
	flamer      bool
}

func newBattleFort(config battleFortConfig) *battleFort {
	b := &battleFort{alliance: config.alliance}
	b.Body.Pos = config.pos
	b.Body.Rotation = config.rotation
	b.battleState = config.battleState
	b.flamer = config.flamer
	return b
}

func (b *battleFort) Pos() gmath.Vec { return b.Body.Pos }

func (b *battleFort) Init(scene *ge.Scene) {
	b.scene = scene

	b.hp = 50
	b.hpComponent = newHealthComponent(&b.hp, &b.Body)
	b.hpComponent.isBuilding = true

	b.Body.InitStaticCircle(b, 32)
	b.Body.LayerMask = unitLayerMask(b.alliance) | buildingLayerMask
	scene.AddBody(&b.Body)

	b.sprite = scene.NewSprite(ImageWall)
	SetHue(b.sprite, spriteHue(b.alliance))
	b.sprite.Pos.Base = &b.Body.Pos
	scene.AddGraphics(b.sprite)

	turretImage := ImageTurret1
	if b.flamer {
		turretImage = ImageTurret2
	}
	b.turretSprite = scene.NewSprite(turretImage)
	SetHue(b.turretSprite, spriteHue(b.alliance))
	b.turretSprite.Pos.Base = &b.Body.Pos
	b.turretSprite.Rotation = &b.Body.Rotation
	b.turretSprite.Shader = scene.NewShader(ShaderBuildingDamage)
	b.turretSprite.Shader.SetFloatValue("HP", 1.0)
	b.turretSprite.Shader.Texture1 = scene.LoadImage(gmath.RandElem(scene.Rand(), damageMaskImages))
	scene.AddGraphics(b.turretSprite)
}

func (b *battleFort) IsDisposed() bool {
	return b.Body.IsDisposed()
}

func (b *battleFort) Dispose() {
	b.Body.Dispose()
	b.sprite.Dispose()
	b.turretSprite.Dispose()
}

func (b *battleFort) Destroy() {
	b.EventDestroyed.Emit(b)
	createExplosions(b.scene, b.Body.Pos, 3, 4)
	b.Dispose()
}

func (b *battleFort) Update(delta float64) {
	if !b.hpComponent.CheckProjectileCollisions(b.scene) {
		b.Destroy()
		return
	}
	b.turretSprite.Shader.Enabled = b.hpComponent.hpPercentage != 1
	b.turretSprite.Shader.SetFloatValue("HP", b.hpComponent.hpPercentage)

	b.reload = gmath.ClampMin(b.reload-delta, 0)
	if b.reload == 0 {
		b.maybeShoot()
	}
}

func (b *battleFort) maybeShootFlamer() {
	target := b.battleState.findTargetUnit(b.Body.Pos, flameTurretWeaponDesign.maxRange+32, b.alliance)
	if target == nil {
		b.reload = b.battleState.rand.FloatRange(0.1, 0.2)
		return
	}
	b.reload = flameTurretWeaponDesign.reload
	p := newProjectile(projectileConfig{
		alliance: b.alliance,
		pos:      b.Body.Pos,
		target:   roundedPos(target.Pos()),
		design:   flameTurretWeaponDesign,
	})
	b.scene.AddObject(p)
	b.scene.Audio().PlaySound(flameTurretWeaponDesign.fireSound)
}

func (b *battleFort) maybeShootCannon() {
	target := b.battleState.findTargetUnit(b.Body.Pos, turretWeaponDesign.maxRange+32, b.alliance)
	if target == nil {
		b.reload = b.battleState.rand.FloatRange(0.3, 0.4)
		return
	}
	facing, ok := facingTowards(b.Body.Pos, target.Pos())
	if !ok {
		b.reload = b.battleState.rand.FloatRange(0.2, 0.35)
		return
	}
	if turretWeaponDesign.maxRange < target.Pos().DistanceTo(b.Body.Pos) {
		b.reload = b.battleState.rand.FloatRange(0.05, 0.2)
		return
	}
	b.reload = turretWeaponDesign.reload
	b.Body.Rotation = facing
	p := newProjectile(projectileConfig{
		alliance: b.alliance,
		pos:      b.Body.Pos,
		target:   fireTargetPos(b.Body.Pos, facing, turretWeaponDesign.maxRange),
		design:   turretWeaponDesign,
	})
	b.scene.AddObject(p)
	b.scene.Audio().PlaySound(turretWeaponDesign.fireSound)
}

func (b *battleFort) maybeShoot() {
	if b.flamer {
		b.maybeShootFlamer()
	} else {
		b.maybeShootCannon()
	}
}
