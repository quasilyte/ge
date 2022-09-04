package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
)

type targetedUnit interface {
	IsDisposed() bool
	Velocity() gemath.Vec
	Pos() gemath.Vec
	OnDamage(damage float64, kind damageKind)
}

type battleTank struct {
	scene *ge.Scene

	Body physics.Body

	design *hullDesign

	Turret *battleTurret

	inCombat         bool
	leaveCombatDelay float64

	Waypoint gemath.Vec

	Player *playerData

	speed         float64
	rotationSpeed gemath.Rad

	hp float64

	hullSprite    *ge.Sprite
	selectionAura *ge.Sprite
	Selected      bool

	personalRangePreference float64

	EventDestroyed       gesignal.Event[*battleTank]
	EventWaypointReached gesignal.Event[*battleTank]
}

func newBattleTank(p *playerData, design tankDesign, mudTerrain bool) *battleTank {
	bt := &battleTank{
		design: design.Hull,
		Player: p,
	}
	bt.hp = bt.design.HP + design.Turret.HPBonus
	bt.speed = bt.design.Speed
	bt.rotationSpeed = bt.design.RotationSpeed
	bt.speed -= design.Turret.SpeedPenalty
	if mudTerrain {
		bt.speed *= 0.6
		bt.rotationSpeed *= 0.8

	}
	bt.Turret = newBattleTurret(p, &bt.Body.Pos, design.Turret)
	bt.Body.Rotation = math.Pi / 2
	bt.Turret.Rotation = math.Pi / 2
	bt.personalRangePreference = design.Turret.FireRange
	return bt
}

func (bt *battleTank) Init(scene *ge.Scene) {
	bt.scene = scene

	bt.personalRangePreference *= scene.Rand().FloatRange(0.85, 0.95)

	switch bt.design.Size {
	case hullSmall:
		bt.Body.InitRotatedRect(bt, 32, 20)
	case hullMedium:
		bt.Body.InitRotatedRect(bt, 36, 22)
	case hullLarge:
		bt.Body.InitRotatedRect(bt, 40, 24)
	}

	bt.selectionAura = scene.NewSprite(ImageUnitSelector)
	bt.selectionAura.Pos.Base = &bt.Body.Pos
	bt.selectionAura.Visible = false
	applyPlayerColor(bt.Player.ID, bt.selectionAura)
	scene.AddGraphics(bt.selectionAura)

	bt.hullSprite = scene.NewSprite(bt.design.Image)
	bt.hullSprite.Pos.Base = &bt.Body.Pos
	bt.hullSprite.Pos.Offset.X = bt.design.OriginX
	bt.hullSprite.Rotation = &bt.Body.Rotation
	applyPlayerColor(bt.Player.ID, bt.hullSprite)
	scene.AddGraphics(bt.hullSprite)

	bt.Turret.Rotation = bt.Body.Rotation

	scene.AddObject(bt.Turret)
	scene.AddBody(&bt.Body)
}

func (bt *battleTank) IsDisposed() bool { return bt.Body.IsDisposed() }

func (bt *battleTank) Dispose() {
	bt.Body.Dispose()
	bt.selectionAura.Dispose()
	bt.hullSprite.Dispose()
	bt.Turret.Dispose()
}

func (bt *battleTank) Update(delta float64) {
	bt.selectionAura.Visible = bt.Selected

	bt.leaveCombatDelay = gemath.ClampMin(bt.leaveCombatDelay-delta, 0)
	if bt.leaveCombatDelay == 0 {
		bt.inCombat = bt.enterCombat()
		bt.leaveCombatDelay = bt.scene.Rand().FloatRange(0.1, 0.25)
	}
	if bt.inCombat {
		return
	}

	if !bt.Waypoint.IsZero() {
		bt.processMovement(delta)
	}
}

func (bt *battleTank) enterCombat() bool {
	_, dist := closestTarget(bt.Player, bt.Body.Pos, false)
	return dist < bt.personalRangePreference
}

func (bt *battleTank) Pos() gemath.Vec { return bt.Body.Pos }

func (bt *battleTank) Velocity() gemath.Vec {
	if bt.Waypoint.IsZero() || bt.inCombat {
		return gemath.Vec{}
	}
	return gemath.RadToVec(bt.Body.Rotation).Mulf(bt.speed)
}

func (bt *battleTank) OnDamage(damage float64, kind damageKind) {
	if !bt.inCombat && bt.Waypoint.IsZero() {
		bt.relocate()
	}

	if kind == damageThermal {
		damage *= 0.8
	}
	bt.hp -= damage
	if bt.hp <= 0 {
		bt.Destroy()
	}
}

func (bt *battleTank) Destroy() {
	bt.EventDestroyed.Emit(bt)
	if !bt.Waypoint.IsZero() {
		bt.EventWaypointReached.Emit(bt)
	}
	e := newExplosion(bt.Body.Pos)
	bt.scene.AddObject(e)
	bt.Dispose()
}

func (bt *battleTank) processMovement(delta float64) {
	dstAngle := bt.Body.Pos.AngleToPoint(bt.Waypoint)
	rotationAmount := bt.rotationSpeed * gemath.Rad(delta)
	newRotation := bt.Body.Rotation.RotatedTowards(dstAngle, rotationAmount)
	bt.Turret.Rotation -= bt.Body.Rotation - newRotation
	bt.Body.Rotation = newRotation

	if newRotation != dstAngle {
		return
	}

	movementAmount := bt.speed * delta
	if bt.Body.Pos.DistanceTo(bt.Waypoint) > movementAmount {
		bt.Body.Pos = bt.Body.Pos.MoveInDirection(movementAmount, bt.Body.Rotation)
		return
	}

	if bt.scene.HasCollisionsAt(&bt.Body, gemath.Vec{}) {
		bt.relocate()
		return
	}

	bt.Body.Pos = bt.Waypoint
	bt.Waypoint = gemath.Vec{}
	bt.EventWaypointReached.Emit(bt)
}

func (bt *battleTank) relocate() {
	locationProbes := [...]gemath.Vec{
		{X: -64, Y: -64},
		{X: 0, Y: -64},
		{X: 64, Y: -64},
		{X: -64, Y: 0},
		{X: 64, Y: 0},
		{X: -64, Y: 64},
		{X: 0, Y: 64},
		{X: 64, Y: 64},
	}
	initialProbe := bt.scene.Rand().IntRange(0, len(locationProbes)-1)
	probe := initialProbe
	for i := 0; i < len(locationProbes); i++ {
		if !bt.scene.HasCollisionsAt(&bt.Body, locationProbes[probe]) {
			bt.Waypoint = bt.Body.Pos.Add(locationProbes[probe])
			return
		}
		probe++
		if probe >= len(locationProbes) {
			probe = 0
		}
	}
	bt.Waypoint = bt.Body.Pos.Add(locationProbes[initialProbe])
}
