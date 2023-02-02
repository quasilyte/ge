package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge/langs"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type SceneController interface {
	Init(*Scene)

	Update(delta float64)
}

type SceneObject interface {
	// Init is called once when object is added to the scene.
	//
	// It's a good time to initialize all dependent objects
	// and attach sprites to the scene.
	Init(*Scene)

	// IsDisposed reports whether scene object was disposed.
	//
	// Disposed objects are removed from the scene before their
	// Update method is called for the current frame.
	IsDisposed() bool

	// Update is called for every object during every logical game frame.
	// Delta specifies how many seconds have passed from the previous frame.
	Update(delta float64)
}

type SceneGraphics interface {
	Draw(dst *ebiten.Image)

	IsDisposed() bool
}

type SceneGraphicsLayer interface {
	AddGraphics(g SceneGraphics)

	Draw(dst *ebiten.Image)
}

type Scene struct {
	root *RootScene

	zindex uint8
}

func (s *Scene) Context() *Context {
	return s.root.context
}

func (s *Scene) Dict() *langs.Dictionary {
	return s.root.context.Dict
}

func (s *Scene) Audio() *AudioSystem {
	return &s.root.context.Audio
}

func (s *Scene) Rand() *gmath.Rand {
	return &s.root.context.Rand
}

func (s *Scene) LoadImage(imageID resource.ImageID) resource.Image {
	return s.root.context.Loader.LoadImage(imageID)
}

func (s *Scene) LoadRaw(rawID resource.RawID) resource.Raw {
	return s.root.context.Loader.LoadRaw(rawID)
}

func (s *Scene) NewParticleEmitter(imageID resource.ImageID) *ParticleEmitter {
	emitter := NewParticleEmitter()
	emitter.SetImage(s.LoadImage(imageID))
	return emitter
}

func (s *Scene) NewShader(shaderID resource.ShaderID) Shader {
	compiled := s.Context().Loader.LoadShader(shaderID)
	return Shader{Enabled: true, compiled: compiled.Data}
}

func (s *Scene) NewSprite(imageID resource.ImageID) *Sprite {
	sprite := NewSprite(s.root.context)
	sprite.SetImage(s.LoadImage(imageID))
	return sprite
}

func (s *Scene) NewRepeatedSprite(imageID resource.ImageID, width, height float64) *Sprite {
	sprite := NewSprite(s.root.context)
	sprite.SetRepeatedImage(s.LoadImage(imageID), width, height)
	return sprite
}

func (s *Scene) NewLabel(fontID resource.FontID) *Label {
	return NewLabel(s.root.context.Loader.LoadFont(fontID).Face)
}

func (s *Scene) AddBody(b *physics.Body) {
	s.root.collisionEngine.AddBody(b)
}

func (s *Scene) AddGraphics(g SceneGraphics) {
	s.root.graphics[s.zindex] = append(s.root.graphics[s.zindex], g)
}

func (s *Scene) AddGraphicsAbove(g SceneGraphics, zindex uint8) {
	z := int(s.zindex) + int(zindex)
	if z > zindexMax {
		panic("z index overflow")
	}
	s.root.graphics[z] = append(s.root.graphics[z], g)
}

func (s *Scene) AddGraphicsBelow(g SceneGraphics, zindex uint8) {
	z := int(s.zindex) - int(zindex)
	if z < 0 {
		panic("z index underflow")
	}
	s.root.graphics[z] = append(s.root.graphics[z], g)
}

func (scene *Scene) AddObject(o SceneObject) {
	scene.root.addObject(o, uint(scene.zindex))
}

func (scene *Scene) AddObjectAbove(o SceneObject, zindex uint8) {
	scene.root.addObject(o, uint(scene.zindex+zindex))
}

func (scene *Scene) AddObjectBelow(o SceneObject, zindex uint8) {
	z := int(scene.zindex) - int(zindex)
	if z < 0 {
		panic("z index underflow")
	}
	scene.root.addObject(o, uint(z))
}

func (scene *Scene) DelayedCall(seconds float64, fn func()) {
	scene.root.delayedFuncs = append(scene.root.delayedFuncs, delayedFunc{
		delay:  seconds,
		action: fn,
	})
}

func (s *Scene) GetCollisions(b *physics.Body) []physics.Collision {
	return s.root.collisionEngine.GetCollisions(b, physics.CollisionConfig{})
}

func (s *Scene) HasCollisionsAt(b *physics.Body, offset gmath.Vec) bool {
	collisions := s.root.collisionEngine.GetCollisions(b, physics.CollisionConfig{
		Offset: offset,
		Limit:  1,
	})
	return len(collisions) != 0
}

func (s *Scene) GetCollisionsAtLayer(b *physics.Body, offset gmath.Vec, layerMask uint16) []physics.Collision {
	return s.root.collisionEngine.GetCollisions(b, physics.CollisionConfig{
		Offset:    offset,
		LayerMask: layerMask,
	})
}

func (s *Scene) HasCollisionsAtLayer(b *physics.Body, offset gmath.Vec, layerMask uint16) bool {
	collisions := s.root.collisionEngine.GetCollisions(b, physics.CollisionConfig{
		Offset:    offset,
		Limit:     1,
		LayerMask: layerMask,
	})
	return len(collisions) != 0
}

func (s *Scene) GetMovementCollision(b *physics.Body, velocity gmath.Vec) *physics.Collision {
	collisions := s.root.collisionEngine.GetCollisions(b, physics.CollisionConfig{
		Velocity: velocity,
		Limit:    1,
	})
	if len(collisions) == 1 {
		return &collisions[0]
	}
	return nil
}
