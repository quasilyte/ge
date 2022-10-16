package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

type battleTurret struct {
	reload float64

	player *playerData

	target            targetedUnit
	snipePos          gemath.Vec
	targetSelectDelay float64
	ignoreBases       bool

	design         *turretDesign
	preferVehicles bool

	ready bool
	delay float64

	scene *ge.Scene

	Rotation gemath.Rad
	pos      *gemath.Vec

	sprite *ge.Sprite
}

func newBattlePostTurret(p *playerData, pos *gemath.Vec, design *turretDesign) *battleTurret {
	turret := newBattleTurret(p, pos, design)
	turret.ignoreBases = true
	return turret
}

func newBattleTurret(p *playerData, pos *gemath.Vec, design *turretDesign) *battleTurret {
	turret := &battleTurret{
		pos:    pos,
		design: design,
		player: p,
		ready:  true,
	}
	return turret
}

func (turret *battleTurret) Init(scene *ge.Scene) {
	turret.preferVehicles = turret.design.Name == "railgun"
	turret.scene = scene
	turret.sprite = scene.NewSprite(turret.design.Image)
	turret.sprite.Pos.Base = turret.pos
	turret.sprite.Rotation = &turret.Rotation
	if turret.design.HP != 0 {
		turret.sprite.FrameWidth = 80
	}
	applyPlayerColor(turret.player.ID, turret.sprite)
	scene.AddGraphics(turret.sprite)
}

func (turret *battleTurret) IsDisposed() bool { return turret.sprite.IsDisposed() }

func (turret *battleTurret) Dispose() {
	turret.sprite.Dispose()
}

func (turret *battleTurret) Sprite() *ge.Sprite {
	return turret.sprite
}

func (turret *battleTurret) seekTarget() {
	if turret.preferVehicles {
		closestTank, closestTankDist := closestTankTarget(turret.player, *turret.pos)
		if closestTank != nil && closestTankDist < turret.design.FireRange {
			turret.target = closestTank
			turret.calculateSnipePos()
			return
		}
	}

	closestEnemy, dist := closestTarget(turret.player, *turret.pos, turret.ignoreBases)
	if closestEnemy != nil && turret.design.FireRange >= dist {
		turret.target = closestEnemy
		turret.calculateSnipePos()
	}
}

func (turret *battleTurret) calculateSnipePos() {
	pos := snipePos(turret.design.ProjectileSpeed, *turret.pos, turret.target.Pos(), turret.target.Velocity())
	r := turret.scene.Rand()
	if r.Chance(0.3) {
		turret.snipePos = pos
	} else {
		turret.snipePos = pos.Add(gemath.Vec{X: r.FloatRange(-8, 8), Y: r.FloatRange(-8, 8)})
	}
}

func (turret *battleTurret) SetDelay(delay float64) {
	if delay != 0 {
		turret.ready = false
		turret.delay = delay
	}
}

func (turret *battleTurret) IsReady() bool {
	return turret.ready
}

func (turret *battleTurret) Update(delta float64) {
	turret.reload = gemath.ClampMin(turret.reload-delta, 0)
	turret.sprite.Pos.Offset.X = gemath.ClampMax(turret.sprite.Pos.Offset.X+delta*20, 0)
	turret.targetSelectDelay = gemath.ClampMin(turret.targetSelectDelay-delta, 0)

	if !turret.ready {
		turret.sprite.SetAlpha(0.5)
		turret.delay = gemath.ClampMin(turret.delay-delta, 0)
		if turret.delay == 0 {
			turret.ready = true
		}
		return
	}

	turret.sprite.SetAlpha(1)

	if turret.target != nil {
		if turret.target.IsDisposed() || (turret.target.Pos().DistanceTo(*turret.pos) > turret.design.FireRange*1.1) {
			turret.snipePos = gemath.Vec{}
			turret.target = nil
		}
	}
	if turret.targetSelectDelay == 0 {
		turret.seekTarget()
		turret.targetSelectDelay = turret.scene.Rand().FloatRange(0.2, 0.6)
	}
	if turret.target == nil {
		return
	}

	rotationTarget := turret.snipePos
	if turret.design.ProjectileSpeed == 0 || rotationTarget.IsZero() {
		rotationTarget = turret.target.Pos()
	}

	dstRotation := turret.pos.AngleToPoint(rotationTarget)
	rotationAmount := turret.design.RotationSpeed * gemath.Rad(delta)
	turret.Rotation = turret.Rotation.RotatedTowards(dstRotation, rotationAmount)
	if turret.Rotation == dstRotation && turret.reload == 0 {
		turret.fire(rotationTarget)
		turret.calculateSnipePos()
	}
}

func (turret *battleTurret) IsLancer() bool {
	return turret.design.Name == "lancer"
}

func (turret *battleTurret) IsBuilder() bool {
	return turret.design.Name == "builder"
}

func (turret *battleTurret) fire(targetPos gemath.Vec) {
	firePos := turret.pos.MoveInDirection(32, turret.Rotation)
	turret.reload += turret.design.Reload
	turret.sprite.Pos.Offset.X = -3
	turret.scene.Audio().PlaySound(turret.design.Sound)

	if turret.design.Name == "railgun" {
		ray := newRailgunRay(firePos, targetPos)
		turret.scene.AddObject(ray)
		turret.target.OnDamage(turret.design.Damage, turret.design.DamageKind)
	} else {
		p := newProjectile(turret.player.Alliance, firePos, turret.Rotation, turret.design)
		turret.scene.AddObject(p)
	}
}

func snipePos(projectileSpeed float64, fireFrom, targetPos, targetVelocity gemath.Vec) gemath.Vec {
	if targetVelocity.IsZero() || projectileSpeed == 0 {
		return targetPos
	}
	dist := targetPos.DistanceTo(fireFrom)
	targetVelocity = targetVelocity.Mulf(1.)
	predictedPos := targetPos.Add(targetVelocity.Mulf(dist / projectileSpeed))
	return predictedPos
}

func closestTarget(p *playerData, pos gemath.Vec, ignoreBases bool) (targetedUnit, float64) {
	if ignoreBases {
		return closestTankTarget(p, pos)
	}
	closestBase, closestBaseDist := closestBaseTarget(p, pos)
	closestTank, closestTankDist := closestTankTarget(p, pos)
	if closestBaseDist < closestTankDist {
		return closestBase, closestBaseDist
	}
	return closestTank, closestTankDist
}

func closestBaseTarget(p *playerData, pos gemath.Vec) (*battlePost, float64) {
	closestBaseDist := math.MaxFloat64
	var closestBase *battlePost
	for _, sector := range p.BattleState.Sectors {
		if sector.Base == nil || sector.Base.Player.Alliance == p.Alliance {
			continue
		}
		dist := sector.Center().DistanceTo(pos)
		if closestBase == nil || dist < closestBaseDist {
			closestBase = sector.Base
			closestBaseDist = dist
		}
	}
	return closestBase, closestBaseDist
}

func closestTankTarget(p *playerData, pos gemath.Vec) (*battleTank, float64) {
	var closestTank *battleTank
	closestTankDist := math.MaxFloat64
	for _, other := range p.BattleState.Players {
		if other.Alliance == p.Alliance {
			continue
		}
		for enemyTank := range p.BattleState.Tanks[other.ID] {
			dist := enemyTank.Body.Pos.DistanceTo(pos)
			if closestTank == nil || dist < closestTankDist {
				closestTank = enemyTank
				closestTankDist = dist
			}
		}
	}
	return closestTank, closestTankDist
}
