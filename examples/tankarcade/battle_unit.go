package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/gmath"
)

var facingAngles = []gmath.Rad{
	facingRight,
	facingDown,
	facingLeft,
	facingUp,
}

const facingRight = 0 * (math.Pi / 2)
const facingDown = 1 * (math.Pi / 2)
const facingLeft = 2 * (math.Pi / 2)
const facingUp = 3 * (math.Pi / 2)

type battleUnit struct {
	scene *ge.Scene

	bot    *localBot
	config battleUnitConfig

	Body physics.Body

	SpecialAmmo     int
	maxSpecialAmmo  int
	specialCooldown float64

	lerpProgress   float64
	lerpTarget     float64
	lerpStart      gmath.Rad
	rotationTarget gmath.Rad

	movementVelocity gmath.Vec
	movementTarget   gmath.Vec

	weaponCooldown float64

	hp          float64
	hpComponent *healthComponent

	shield     *unitShield
	shieldTime float64

	sprite *ge.Sprite

	EventDestroyed gesignal.Event[*battleUnit]
}

type battleUnitConfig struct {
	group string

	alliance int

	maxHP float64

	playerID int

	rotation gmath.Rad
	pos      gmath.Vec
	image    resource.ImageID

	weapon  *weaponDesign
	special *specialWeaponDesign

	improvedSpecial bool

	battleState *battleState

	weaponReloadMultiplier float64

	speed        float64
	rotationTime float64
}

func newBattleUnit(config battleUnitConfig) *battleUnit {
	return &battleUnit{config: config}
}

func (u *battleUnit) Pos() gmath.Vec { return u.Body.Pos }

func (u *battleUnit) Init(scene *ge.Scene) {
	u.scene = scene

	u.hp = u.config.maxHP
	u.hpComponent = newHealthComponent(&u.hp, &u.Body)

	u.weaponCooldown = 0.5

	if u.config.special != nil {
		u.maxSpecialAmmo = u.config.special.ammo
		if u.config.improvedSpecial && u.config.special == rocketLauncherDesign {
			u.maxSpecialAmmo++
		}
		u.SpecialAmmo = u.maxSpecialAmmo
	}

	u.Body.InitCircle(u, 32)
	u.Body.Pos = u.config.pos
	u.Body.Rotation = u.config.rotation
	u.Body.LayerMask = unitLayerMask(u.config.alliance)
	scene.AddBody(&u.Body)

	u.sprite = scene.NewSprite(u.config.image)
	u.sprite.SetHue(spriteHue(u.config.alliance))
	u.sprite.Pos.Base = &u.Body.Pos
	u.sprite.Rotation = &u.Body.Rotation
	scene.AddGraphics(u.sprite)
}

func (u *battleUnit) IsDisposed() bool {
	return u.sprite.IsDisposed()
}

func (u *battleUnit) Dispose() {
	u.sprite.Dispose()
	u.Body.Dispose()
}

func (u *battleUnit) Destroy() {
	u.EventDestroyed.Emit(u)
	if u.shield != nil {
		u.shield.Dispose()
		u.shield = nil
	}
	createExplosions(u.scene, u.Body.Pos, 2, 3)
	u.Dispose()
}

func (u *battleUnit) Update(delta float64) {
	if !u.hpComponent.CheckProjectileCollisions(u.scene) {
		u.Destroy()
		return
	}

	u.shieldTime = gmath.ClampMin(u.shieldTime-delta, 0)
	if u.shieldTime == 0 && u.shield != nil {
		u.hpComponent.shieldLevel = 0
		u.shield.Dispose()
		u.shield = nil
	}

	for _, collision := range u.scene.GetCollisions(&u.Body) {
		switch obj := collision.Body.Object.(type) {
		case *pickupBonus:
			switch obj.kind {
			case pickupHP:
				u.hp = u.config.maxHP
			case pickupAmmo:
				if u.config.special != nil {
					u.SpecialAmmo = u.maxSpecialAmmo
				}
			}
			u.scene.Audio().PlaySound(AudioBonus)
			obj.Destroy()
		}
	}

	u.weaponCooldown = gmath.ClampMin(u.weaponCooldown-delta, 0)
	u.specialCooldown = gmath.ClampMin(u.specialCooldown-delta, 0)

	if !u.movementVelocity.IsZero() {
		if u.Body.Pos.DistanceTo(u.movementTarget) <= u.config.speed*delta {
			u.Body.Pos = u.movementTarget
			u.movementVelocity = gmath.Vec{}
		} else {
			u.Body.Pos = u.Body.Pos.Add(u.movementVelocity.Mulf(delta))
		}
	}
	if u.lerpTarget != 0 {
		u.lerpProgress += delta
		u.Body.Rotation = u.lerpStart.LerpAngle(u.rotationTarget, u.lerpProgress/u.lerpTarget).Normalized()
		if u.Body.Rotation.EqualApprox(u.rotationTarget) {
			u.Body.Rotation = u.rotationTarget
			u.lerpTarget = 0
			u.lerpProgress = 0
		}
	}

	if !u.MoveOrderAvailable() {
		u.sprite.FrameOffset.X = u.sprite.FrameWidth
	} else {
		u.sprite.FrameOffset.X = 0
	}
}

func (u *battleUnit) SpecialOrderAvailable() bool {
	if u.config.special == nil {
		return false // No special weapon
	}
	if u.lerpProgress != 0 {
		return false // Rotation is in progress
	}
	if u.specialCooldown != 0 {
		return false
	}
	if u.SpecialAmmo <= 0 {
		return false
	}
	return true
}

func (u *battleUnit) FireOrderAvailable() bool {
	if u.lerpProgress != 0 {
		return false // Rotation is in progress
	}
	if u.weaponCooldown != 0 {
		return false
	}
	return true
}

func (u *battleUnit) SpecialOrder() {
	if !u.SpecialOrderAvailable() {
		return
	}
	u.SpecialAmmo--
	special := u.config.special

	if special == shieldDesign {
		u.shieldTime = 1.3
		u.hpComponent.shieldLevel = 1
		if u.config.improvedSpecial {
			u.hpComponent.shieldLevel = 2
		}
		u.shield = newUnitShield(ge.Pos{Base: &u.Body.Pos})
		u.scene.AddObject(u.shield)
		u.specialCooldown = special.extra.reload
		u.scene.Audio().PlaySound(special.extra.fireSound)
		return
	}

	if special == mineLayerDesign {
		mineLifespan := 30.0
		if u.config.improvedSpecial {
			mineLifespan = -1
		}
		u.scene.AddObject(newBattleMine(roundedPos(u.Body.Pos), u.config.alliance, mineLifespan))
		u.specialCooldown = special.extra.reload
		u.scene.Audio().PlaySound(special.extra.fireSound)
		return
	}

	if special == rocketLauncherDesign {
		target := fireTargetPos(u.Body.Pos, u.Body.Rotation, special.extra.maxRange)
		p := newProjectile(projectileConfig{
			alliance: u.config.alliance,
			pos:      u.Body.Pos,
			target:   target,
			design:   special.extra,
		})
		u.scene.AddObject(p)
		u.specialCooldown = special.extra.reload
		u.scene.Audio().PlaySound(special.extra.fireSound)
		return
	}

	if special == flamethrowerDesign {
		target0 := u.Body.Pos.Add(cellDelta(u.Body.Rotation))
		targets := [3]gmath.Vec{
			target0,
			target0.Add(cellDelta((u.Body.Rotation - math.Pi/2).Normalized())),
			target0.Add(cellDelta((u.Body.Rotation + math.Pi/2).Normalized())),
		}
		for _, target := range targets {
			p := newProjectile(projectileConfig{
				alliance: u.config.alliance,
				pos:      u.Body.Pos,
				target:   target,
				design:   special.extra,
			})
			u.scene.AddObject(p)
		}
		reloadMultiplier := 1.0
		if u.config.improvedSpecial {
			reloadMultiplier = 0.35
		}
		u.specialCooldown = special.extra.reload * reloadMultiplier
		u.scene.Audio().PlaySound(special.extra.fireSound)
		return
	}

	if special == regeneratorDesign {
		firePos := u.Body.Pos.Add(cellDelta(u.Body.Rotation))
		target := fireTargetPos(u.Body.Pos, u.Body.Rotation, special.extra.maxRange)
		config := projectileConfig{
			alliance: u.config.alliance,
			pos:      firePos,
			target:   target,
			design:   special.extra,
		}
		if u.config.improvedSpecial {
			config.design = improvedRegeneratorDesign
		}
		p := newProjectile(config)
		u.scene.AddObject(p)
		u.specialCooldown = special.extra.reload
		u.scene.Audio().PlaySound(special.extra.fireSound)
		return
	}

	panic("unexpected special weapon")
}

func (u *battleUnit) FireOrder() {
	if !u.FireOrderAvailable() {
		return
	}

	if u.config.weapon == railgunDesign {
		firePos := roundedPos(u.Body.Pos)
		for i := 0; i < 5; i++ {
			if firePos != fixedPos(firePos) {
				break
			}
			if u.config.battleState.getCellInfoAt(firePos)&cellDarkTile != 0 {
				break
			}
			p := newProjectile(projectileConfig{
				alliance: u.config.alliance,
				pos:      firePos,
				target:   firePos,
				design:   u.config.weapon,
			})
			firePos = firePos.Add(cellDelta(u.Body.Rotation))
			u.scene.AddObject(p)
			if i != 4 {
				ray := newRailgunRayEffect(firePos, u.Body.Rotation)
				u.scene.AddObject(ray)
			}
		}
	} else {
		target := fireTargetPos(u.Body.Pos, u.Body.Rotation, u.config.weapon.maxRange)
		p := newProjectile(projectileConfig{
			alliance: u.config.alliance,
			pos:      u.Body.Pos,
			target:   target,
			design:   u.config.weapon,
		})
		u.scene.AddObject(p)
	}

	u.weaponCooldown = u.config.weapon.reload * u.config.weaponReloadMultiplier
	u.scene.Audio().PlaySound(u.config.weapon.fireSound)
}

func (u *battleUnit) MoveOrderAvailable() bool {
	if u.lerpProgress != 0 {
		return false // Rotation is in progress
	}
	if !u.movementVelocity.IsZero() {
		return false // Movement is in progress
	}
	return true
}

func (u *battleUnit) MoveOrder(angle gmath.Rad) {
	if !u.MoveOrderAvailable() {
		return
	}
	if u.Body.Rotation == angle {
		var movementStep gmath.Vec
		var movementVelocity gmath.Vec
		switch angle {
		case facingRight:
			movementStep.X = 64
			movementVelocity = gmath.Vec{X: u.config.speed}
		case facingDown:
			movementStep.Y = 64
			movementVelocity = gmath.Vec{Y: u.config.speed}
		case facingLeft:
			movementStep.X = -64
			movementVelocity = gmath.Vec{X: -u.config.speed}
		case facingUp:
			movementStep.Y = -64
			movementVelocity = gmath.Vec{Y: -u.config.speed}
		}
		if u.scene.HasCollisionsAtLayer(&u.Body, movementStep, blockingLayerMask) {
			return
		}
		movementTarget := u.Body.Pos.Add(movementStep)
		if fixedPos(movementTarget) != movementTarget {
			return
		}
		u.movementTarget = movementTarget
		u.movementVelocity = movementVelocity
		return
	}
	u.lerpStart = u.Body.Rotation
	u.lerpTarget = u.config.rotationTime
	u.rotationTarget = angle
}
