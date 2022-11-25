package ge

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/gmath"
)

type disposeState uint8

const (
	disposeUnset disposeState = iota
	disposeNow
	disposeDetached
)

type Particle struct {
	Velocity gmath.Vec
	Rotation gmath.Rad

	Lifetime float64
}

type ParticleConfig struct {
	Lifetime float64

	Amount int

	InitFunc func(p *Particle)

	RotationFunc func(p *Particle, delta float64) gmath.Rad
	VelocityFunc func(p *Particle, delta float64) gmath.Vec
	ColorFunc    func(p *Particle) ColorScale
}

type ParticleEmitter struct {
	Visible bool

	Centered bool

	Pos Pos

	FrameOffset gmath.Vec
	FrameWidth  float64
	FrameHeight float64

	Hue gmath.Rad

	image *ebiten.Image

	disposed disposeState

	emitReload       float64
	emitCooldown     float64
	particleLifetime float64
	particleIndex    gmath.Slider

	tmpParticle  Particle
	initFunc     func(p *Particle)
	rotationFunc func(p *Particle, delta float64) gmath.Rad
	velocityFunc func(p *Particle, delta float64) gmath.Vec
	colorFunc    func(p *Particle) ColorScale

	particles []particle
}

func NewParticleEmitter() *ParticleEmitter {
	return &ParticleEmitter{
		Visible:  true,
		Centered: true,
	}
}

func (e *ParticleEmitter) Init(scene *Scene) {}

func (e *ParticleEmitter) SetImage(img resource.Image) {
	w, h := img.Data.Size()
	e.image = img.Data
	e.FrameWidth = img.DefaultFrameWidth
	if e.FrameWidth == 0 {
		e.FrameWidth = float64(w)
	}
	e.FrameHeight = img.DefaultFrameHeight
	if e.FrameHeight == 0 {
		e.FrameHeight = float64(h)
	}
}

func (e *ParticleEmitter) SetConfig(config ParticleConfig) {
	e.particleLifetime = config.Lifetime

	if len(e.particles) < config.Amount {
		e.particles = make([]particle, config.Amount)
	}

	e.particleIndex.SetBounds(0, config.Amount-1)
	e.emitReload = config.Lifetime / float64(config.Amount)
	e.emitCooldown = 0

	e.initFunc = config.InitFunc
	if e.initFunc == nil {
		e.initFunc = func(p *Particle) { *p = Particle{} }
	}
	e.rotationFunc = config.RotationFunc
	e.velocityFunc = config.VelocityFunc
	e.colorFunc = config.ColorFunc
}

func (e *ParticleEmitter) IsDisposed() bool {
	return e.disposed == disposeNow
}

func (e *ParticleEmitter) Dispose() {
	e.disposed = disposeNow
}

func (e *ParticleEmitter) DisposeDetached() {
	e.disposed = disposeDetached
}

func (e *ParticleEmitter) Update(delta float64) {
	if len(e.particles) == 0 {
		return
	}

	e.emitCooldown = gmath.ClampMin(e.emitCooldown-delta, 0)
	if e.emitCooldown == 0 {
		// Create new particles only if this object is not being disposed.
		if e.disposed == disposeUnset {
			e.emitCooldown = e.emitReload
			p := &e.particles[e.particleIndex.Value()]
			e.particleReset(p)
			e.particleIndex.Inc()
		}
	}

	hasActiveParticles := false
	for i := range e.particles {
		p := &e.particles[i]
		if p.hp == 0 {
			continue
		}
		hasActiveParticles = true
		e.particleTick(delta, p)
	}
	// If all particles are gone and emitter is being disposed in detached mode,
	// it's time for it to be deleted.
	if !hasActiveParticles && e.disposed == disposeDetached {
		e.Dispose()
	}
}

func (e *ParticleEmitter) setTmpParticle(p *particle) {
	e.tmpParticle.Velocity = p.velocity
	e.tmpParticle.Rotation = p.rotation
	e.tmpParticle.Lifetime = p.hp
}

func (e *ParticleEmitter) particleTick(delta float64, p *particle) {
	e.setTmpParticle(p)

	if e.velocityFunc != nil {
		p.velocity = e.velocityFunc(&e.tmpParticle, delta)
	}
	if e.rotationFunc != nil {
		p.rotation = e.rotationFunc(&e.tmpParticle, delta)
	}

	p.hp = gmath.ClampMin(p.hp-delta, 0)
	p.offset = p.offset.Add(p.velocity.Mulf(delta))
}

func (e *ParticleEmitter) particleReset(p *particle) {
	p.offset = gmath.Vec{}
	p.hp = e.particleLifetime

	e.initFunc(&e.tmpParticle)
	p.rotation = e.tmpParticle.Rotation
	p.velocity = e.tmpParticle.Velocity
}

func (e *ParticleEmitter) Draw(screen *ebiten.Image) {
	if !e.Visible || len(e.particles) == 0 {
		return
	}

	var origin gmath.Vec
	if e.Centered {
		origin = gmath.Vec{X: e.FrameWidth / 2, Y: e.FrameHeight / 2}
	}
	origin = origin.Sub(e.Pos.Offset)

	subImage := e.image.SubImage(image.Rectangle{
		Min: image.Point{
			X: int(e.FrameOffset.X),
			Y: int(e.FrameOffset.Y),
		},
		Max: image.Point{
			X: int(e.FrameOffset.X + e.FrameWidth),
			Y: int(e.FrameOffset.Y + e.FrameHeight),
		},
	}).(*ebiten.Image)

	var posX float64
	var posY float64
	if e.Pos.Base != nil {
		posX = e.Pos.Base.X - origin.X
		posY = e.Pos.Base.Y - origin.Y
	} else {
		posX = 0 - origin.X
		posY = 0 - origin.Y
	}
	for i := range e.particles {
		p := &e.particles[i]
		if p.hp == 0 {
			continue
		}

		e.setTmpParticle(p)

		var drawOptions ebiten.DrawImageOptions

		drawOptions.GeoM.Translate(-origin.X, -origin.Y)
		drawOptions.GeoM.Rotate(float64(p.rotation))
		drawOptions.GeoM.Translate(origin.X, origin.Y)
		drawOptions.GeoM.Translate(posX, posY)

		if e.colorFunc != nil {
			colorScale := e.colorFunc(&e.tmpParticle)
			drawOptions.ColorM.Scale(float64(colorScale.R), float64(colorScale.G), float64(colorScale.B), float64(colorScale.A))
		}
		if e.Hue != 0 {
			drawOptions.ColorM.RotateHue(float64(e.Hue))
		}

		drawOptions.GeoM.Translate(p.offset.X, p.offset.Y)
		screen.DrawImage(subImage, &drawOptions)
	}
}

type particle struct {
	hp       float64
	rotation gmath.Rad
	velocity gmath.Vec
	offset   gmath.Vec
}
