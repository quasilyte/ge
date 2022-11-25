package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type botProgramKind uint8

const (
	botDelayedAttack botProgramKind = iota
	botGuard
	botBaseAttack
	botHQAttack
	botTankHunt
	botRoam
)

type targetable interface {
	Pos() gmath.Vec
	IsDisposed() bool
}

type localBot struct {
	state *battleState
	unit  *battleUnit

	target     targetable
	attackAura physics.Body
	fireDelay  float64
	fireDist   float64

	program botProgramKind

	delayedAttackCountdown float64

	roamRounds  int
	nextProgram botProgramKind

	prevPos              gmath.Vec
	movementQueueStorage []gmath.Rad
	movementDest         gmath.Vec
	movementQueue        []gmath.Rad
	movementDelay        float64
}

type localBotConfig struct {
	state       *battleState
	unit        *battleUnit
	program     botProgramKind
	attackDelay float64
}

func newLocalBot(config localBotConfig) *localBot {
	return &localBot{
		state:                  config.state,
		unit:                   config.unit,
		program:                config.program,
		delayedAttackCountdown: config.attackDelay,
		movementQueueStorage:   make([]gmath.Rad, 0, 6),
	}
}

func (b *localBot) IsDisposed() bool { return b.unit.IsDisposed() }

func (b *localBot) Init(scene *ge.Scene) {
	b.attackAura.InitCircle(nil, b.unit.config.weapon.maxRange*1.2)

	b.fireDist = b.unit.config.weapon.maxRange - 64

	if b.program == botRoam {
		b.roamRounds = scene.Rand().IntRange(1, 3)
		if scene.Rand().Bool() {
			b.nextProgram = botTankHunt
		} else {
			b.nextProgram = botBaseAttack
		}
	}

	if b.program == botDelayedAttack && b.delayedAttackCountdown == 0 {
		b.delayedAttackCountdown = scene.Rand().FloatRange(5, 10)
	}
}

func (b *localBot) SetDestination(pos gmath.Vec) {
	b.movementDest = pos
}

func (b *localBot) shouldFire(delta float64) bool {
	rand := b.state.rand

	if b.fireDelay != 0 {
		return false
	}

	// On average, bot will fire once in ~8 secs even if there
	// is no need in doing so. Just to make things less predictable and more fun.
	if rand.Chance(delta / 8) {
		return true
	}

	b.attackAura.Pos = b.unit.Body.Pos
	b.attackAura.LayerMask = projectileLayerMask(b.unit.config.alliance)
	for _, collision := range b.state.scene.GetCollisions(&b.attackAura) {
		var aggroMultiplier float64
		switch collision.Body.Object.(type) {
		case *battleUnit, *battleFort:
			aggroMultiplier = 1.1
		case *battleFactory:
			aggroMultiplier = 0.9
		case *battleBase:
			aggroMultiplier = 1
		case *battleWall:
			aggroMultiplier = 0.6
		default:
			continue
		}
		dist := collision.Body.Pos.DistanceTo(b.unit.Body.Pos)
		fireChance := (1 - (dist / (b.unit.config.weapon.maxRange * aggroMultiplier)))
		if b.program != botRoam {
			fireChance += 0.2
		}
		if fireChance > 0 && rand.Chance(fireChance) {
			return true
		}
	}
	b.fireDelay = rand.FloatRange(b.unit.config.weapon.reload/5, b.unit.config.weapon.reload)

	return false
}

func (b *localBot) resetMovementQueue() {
	b.movementQueue = b.movementQueueStorage[:0]
}

func (b *localBot) buildMovementQueue(numSteps int) bool {
	b.resetMovementQueue()
	pos := b.unit.Body.Pos
	if b.prevPos == pos {
		b.prevPos = gmath.Vec{}
		return false
	}
	facing := b.unit.Body.Rotation
	for i := 0; i < 4 && numSteps > 0; i++ {
		side := b.state.rand.IntRange(0, 3)
		switch side {
		case 0:
			for numSteps > 0 && pos.X < b.movementDest.X {
				b.movementQueue = append(b.movementQueue, facingRight)
				if facing == facingRight {
					numSteps--
					pos.X += 64
				} else {
					facing = facingRight
				}
			}
		case 1:
			for numSteps > 0 && pos.X > b.movementDest.X {
				b.movementQueue = append(b.movementQueue, facingLeft)
				if facing == facingLeft {
					numSteps--
					pos.X -= 64
				} else {
					facing = facingLeft
				}
			}
		case 2:
			for numSteps > 0 && pos.Y < b.movementDest.Y {
				b.movementQueue = append(b.movementQueue, facingDown)
				if facing == facingDown {
					numSteps--
					pos.Y += 64
				} else {
					facing = facingDown
				}
			}
		default:
			for numSteps > 0 && pos.Y > b.movementDest.Y {
				b.movementQueue = append(b.movementQueue, facingUp)
				if facing == facingUp {
					numSteps--
					pos.Y -= 64
				} else {
					facing = facingUp
				}
			}
		}
		if pos == b.movementDest {
			break
		}
	}
	if len(b.movementQueue) != 0 && pos != b.unit.Body.Pos {
		b.prevPos = b.unit.Body.Pos
	}
	return true
}

func (b *localBot) maybeTurnTowards(pos gmath.Vec) bool {
	facing, ok := facingTowards(b.unit.Body.Pos, b.target.Pos())
	if !ok {
		return false
	}
	b.resetMovementQueue()
	if facing != b.unit.Body.Rotation {
		b.movementQueue = append(b.movementQueue, facing)
	}
	b.movementDelay = b.state.rand.FloatRange(0.6, 1)
	b.movementDest = gmath.Vec{}
	return true
}

func (b *localBot) tryUnstuck() {
	pos := b.unit.Body.Pos
	b.prevPos = gmath.Vec{}
	roll := b.state.rand.Float()
	if roll < 0.05 {
		b.movementDest = gmath.Vec{}
		return
	}
	if roll < 0.4 {
		b.nextProgram = b.program
		b.program = botRoam
		b.roamRounds = 1
		return
	}
	facing := gmath.RandElem(b.state.rand, facingAngles)
	dist := b.state.rand.IntRange(2, 5)
	for i := 0; i < dist; i++ {
		pos = fixedPos(pos.Add(cellDelta(facing)))
		if b.state.getWallAt(pos) == nil {
			break
		}
	}
	b.movementDest = pos
}

func (b *localBot) scheduleMovement() {
	if b.movementDelay != 0 {
		return
	}

	switch b.program {
	case botDelayedAttack:
		if b.unit.hp != b.unit.config.maxHP {
			b.delayedAttackCountdown = 0
			return
		}
		b.movementDelay = b.state.rand.FloatRange(0.2, 1.2)

	case botGuard:
		if b.unit.hp != b.unit.config.maxHP {
			b.program = botTankHunt
			return
		}
		b.target = b.state.findTargetUnit(b.unit.Body.Pos, b.fireDist, b.unit.config.alliance)
		if b.target != nil {
			b.program = botTankHunt
			return
		}
		b.movementDelay = b.state.rand.FloatRange(0.4, 2)

	case botHQAttack:
		if b.delayedAttackCountdown > 0 {
			b.movementDelay = b.state.rand.FloatRange(0.2, 0.3)
			return
		}
		if b.target == nil {
			b.target = b.state.findClosestBase(b.unit.Body.Pos, b.unit.config.alliance)
		}
		if b.target == nil {
			b.movementDelay = b.state.rand.FloatRange(0.6, 1)
			return
		}
		if b.unit.Body.Pos.DistanceTo(b.target.Pos()) <= b.fireDist && b.maybeTurnTowards(b.target.Pos()) {
			return
		}
		if b.movementDest.IsZero() || b.unit.Body.Pos == b.movementDest {
			b.movementDest = b.target.Pos()
			return
		}
		if !b.buildMovementQueue(b.state.rand.IntRange(2, 5)) {
			b.tryUnstuck()
		}

	case botBaseAttack:
		if b.target == nil {
			b.target = b.state.findTargetBuilding(b.unit.config.alliance)
		}
		if b.target == nil {
			b.program = botTankHunt
			b.movementDelay = b.state.rand.FloatRange(0.6, 2)
			return
		}
		if b.unit.Body.Pos.DistanceTo(b.target.Pos()) <= b.fireDist && b.maybeTurnTowards(b.target.Pos()) {
			return
		}
		if b.movementDest.IsZero() || b.unit.Body.Pos == b.movementDest {
			b.movementDest = b.target.Pos()
			return
		}
		if b.state.rand.Chance(0.4) {
			b.movementDelay = b.state.rand.FloatRange(0.3, 0.7)
			return
		}
		if !b.buildMovementQueue(b.state.rand.IntRange(3, 7)) {
			b.tryUnstuck()
		}

	case botTankHunt:
		if b.target == nil {
			b.target = b.state.findTargetUnit(gmath.Vec{}, 0, b.unit.config.alliance)
		}
		if b.target == nil {
			b.program = botBaseAttack
			b.movementDelay = b.state.rand.FloatRange(0.6, 2)
			return
		}
		if b.unit.Body.Pos.DistanceTo(b.target.Pos()) <= 192 && b.maybeTurnTowards(b.target.Pos()) {
			return
		}
		if b.movementDest.IsZero() || b.unit.Body.Pos == b.movementDest {
			b.movementDest = roundedPos(b.target.Pos())
			for i := 0; i < b.state.rand.IntRange(0, 2); i++ {
				pos := fixedPos(b.movementDest.Add(cellDelta(gmath.RandElem(b.state.rand, facingAngles))))
				if b.state.getWallAt(pos) == nil {
					b.movementDest = pos
				}
			}
			return
		}
		if b.state.rand.Chance(0.2) {
			b.movementDelay = b.state.rand.FloatRange(0.25, 0.9)
			return
		}
		if !b.buildMovementQueue(b.state.rand.IntRange(4, 9)) {
			b.tryUnstuck()
		}

	case botRoam:
		if b.movementDest.IsZero() || b.unit.Body.Pos == b.movementDest {
			if b.roamRounds == 0 {
				b.program = b.nextProgram
				return
			}
			// Choose next movement destination.
			var dest gmath.Vec
			for i := 0; i < 3; i++ {
				destAttempt := fixedPos(roundedPos(b.unit.Body.Pos.Add(b.state.rand.Offset(-256, 256))))
				if b.state.getWallAt(destAttempt) != nil {
					continue
				}
				if b.state.getCellInfoAt(destAttempt)&cellRoamStop != 0 {
					continue
				}
				dest = destAttempt
			}
			if dest.IsZero() {
				b.movementDelay = b.state.rand.FloatRange(0.05, 0.2)
				return
			}
			b.movementDest = dest
			b.roamRounds--
			return
		}
		// Continue following the destination.
		if b.state.rand.Chance(0.2) {
			b.movementDelay = b.state.rand.FloatRange(0.5, 1.5)
			return
		}
		if !b.buildMovementQueue(b.state.rand.IntRange(1, 7)) {
			b.movementDest = gmath.Vec{}
		}
	}
}

func (b *localBot) Update(delta float64) {
	if b.unit.FireOrderAvailable() {
		if b.shouldFire(delta) {
			b.unit.FireOrder()
		}
	}

	if b.target != nil && b.target.IsDisposed() {
		b.target = nil
	}

	b.delayedAttackCountdown -= delta
	if b.program == botDelayedAttack {
		if b.delayedAttackCountdown <= 0 {
			if b.state.rand.Bool() {
				b.program = botBaseAttack
			} else {
				b.program = botTankHunt
			}
		}
	}

	b.fireDelay = gmath.ClampMin(b.fireDelay-delta, 0)
	b.movementDelay = gmath.ClampMin(b.movementDelay-delta, 0)

	if b.unit.MoveOrderAvailable() {
		if len(b.movementQueue) == 0 {
			b.scheduleMovement()
		} else {
			b.unit.MoveOrder(b.movementQueue[0])
			b.movementQueue = b.movementQueue[1:]
		}
	}
}
