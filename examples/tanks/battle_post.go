package main

import (
	"strconv"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
)

type battlePost struct {
	scene *ge.Scene

	Turret            *battleTurret
	turretSpritesheet *ge.Animation

	HQ bool

	body physics.Body

	Player *playerData

	spritesheet *ge.Animation

	UnderAttack float64

	turretHp float64

	HP          float64
	maxHP       float64
	frameOffset float64

	startingTurret *turretDesign

	production    float64
	product       tankDesign
	progressLabel *ge.Label

	EventDestroyed           gesignal.Event[*battlePost]
	EventProductionCompleted gesignal.Event[tankDesign]
}

func newBattlePost(p *playerData, pos gemath.Vec, turret *turretDesign) *battlePost {
	bp := &battlePost{
		Player:         p,
		maxHP:          750,
		startingTurret: turret,
	}
	bp.HP = bp.maxHP
	bp.body.Pos = pos
	return bp
}

func (bp *battlePost) Init(scene *ge.Scene) {
	bp.scene = scene

	// TODO: InitRect.
	bp.body.InitRotatedRect(bp, 50, 50)

	sprite := scene.NewSprite(ImageBattlePost)
	sprite.FrameWidth = 64
	sprite.Pos.Base = &bp.body.Pos
	applyPlayerColor(bp.Player.ID, sprite)
	bp.spritesheet = ge.NewAnimation(sprite, 10)
	bp.spritesheet.SetAnimationSpan(bp.maxHP)
	bp.spritesheet.EventFrameChanged.Connect(bp, bp.onFrameChanged)
	scene.AddGraphics(sprite)

	bp.progressLabel = scene.NewLabel(FontSmall)
	bp.progressLabel.Visible = false
	bp.progressLabel.Pos.Set(&bp.body.Pos, -120, -122)
	scene.AddGraphics(bp.progressLabel)

	scene.AddBody(&bp.body)

	if bp.startingTurret != nil {
		bp.InstallTurret(bp.startingTurret)
		bp.startingTurret = nil
	}
}

func (bp *battlePost) IsDisposed() bool {
	return bp.body.IsDisposed()
}

func (bp *battlePost) Dispose() {
	bp.body.Dispose()
	bp.spritesheet.Sprite().Dispose()
	bp.progressLabel.Dispose()
}

func (bp *battlePost) Pos() gemath.Vec {
	return bp.body.Pos
}

func (bp *battlePost) Velocity() gemath.Vec { return gemath.Vec{} }

func (bp *battlePost) Update(delta float64) {
	bp.UnderAttack = gemath.ClampMin(bp.UnderAttack-delta, 0)
	if !bp.product.IsEmpty() {
		bp.handleProcution(delta)
	}
}

func (bp *battlePost) InstallTurret(design *turretDesign) {
	bp.turretHp = design.HP
	bp.Turret = newBattlePostTurret(bp.Player, &bp.body.Pos, design)
	bp.Turret.Rotation = bp.body.Pos.AngleToPoint(gemath.Vec{X: 1920 / 2, Y: 1080 / 2})
	bp.scene.AddObject(bp.Turret)

	sprite := bp.Turret.Sprite()
	sprite.FrameWidth = 80
	bp.turretSpritesheet = ge.NewAnimation(sprite, 5)
	bp.turretSpritesheet.SetAnimationSpan(design.HP)
}

func (bp *battlePost) OnDamage(damage float64, kind damageKind) {
	bp.UnderAttack = 1

	switch kind {
	case damageEnergy:
		damage *= 0.4
	case damageThermal:
		damage *= 1.2
	}

	if bp.Turret != nil {
		bp.turretHp -= damage
		bp.turretSpritesheet.Tick(damage)
		if bp.turretHp < 0 {
			bp.turretHp = 0
			bp.destroyTurret()
		}
		return
	}

	bp.HP -= damage
	bp.spritesheet.Tick(damage)
	if bp.HP <= 0 {
		bp.Destroy()
	}
}

func (bp *battlePost) onFrameChanged(frameDelta int) {
	for i := 0; i < frameDelta; i++ {
		offset := bp.scene.Rand().Offset(-24, 24)
		e := newExplosion(bp.body.Pos.Add(offset))
		e.Scale = 0.6
		e.AnimationSpeed = 2
		bp.scene.AddObject(e)
	}
}

func (bp *battlePost) destroyTurret() {
	e := newExplosion(bp.body.Pos)
	bp.scene.AddObject(e)
	bp.Turret.Dispose()
	bp.Turret = nil
	bp.turretSpritesheet = nil
}

func (bp *battlePost) Destroy() {
	if bp.Turret != nil {
		bp.destroyTurret()
	}
	bp.Dispose()
	bp.EventDestroyed.Emit(bp)
	e := newExplosion(bp.body.Pos)
	bp.scene.AddObject(e)
}

func (bp *battlePost) IsBusy() bool {
	return !bp.product.IsEmpty()
}

func (bp *battlePost) StartProduction(design tankDesign) {
	bp.product = design
	bp.progressLabel.Visible = true
}

func (bp *battlePost) handleProcution(delta float64) {
	bp.production += delta

	if bp.production >= bp.product.ProductionTime() {
		bp.EventProductionCompleted.Emit(bp.product)
		bp.production = 0
		bp.product = tankDesign{}
		bp.progressLabel.Visible = false
	} else {
		percent := gemath.Percentage(bp.production, bp.product.ProductionTime())
		bp.progressLabel.Text = strconv.Itoa(int(percent)) + "%"
	}
}
