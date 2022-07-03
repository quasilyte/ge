package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/langs"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/ge/resource"
)

type SceneController interface {
	Init(*Scene)

	Update(delta float64)
}

type Disposable interface {
	IsDisposed() bool

	Dispose()
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

func (s *Scene) Audio() *resource.AudioSystem {
	return &s.root.context.Audio
}

func (s *Scene) Rand() *gemath.Rand {
	return &s.root.context.Rand
}

func (s *Scene) LoadImage(imageID resource.ImageID) resource.Image {
	return s.root.context.Loader.LoadImage(imageID)
}

func (s *Scene) NewSprite(imageID resource.ImageID) *Sprite {
	sprite := NewSprite()
	sprite.SetImage(s.LoadImage(imageID))
	return sprite
}

func (s *Scene) NewRepeatedSprite(imageID resource.ImageID, width, height float64) *Sprite {
	sprite := NewSprite()
	sprite.SetRepeatedImage(s.LoadImage(imageID), width, height)
	return sprite
}

func (s *Scene) NewLabel(fontID resource.FontID) *Label {
	return NewLabel(s.root.context.Loader.LoadFont(fontID))
}

func (s *Scene) AddBody(b *physics.Body) {
	s.root.collisionEngine.AddBody(b)
}

func (s *Scene) AddGraphics(g SceneGraphics) {
	s.root.graphics[s.zindex] = append(s.root.graphics[s.zindex], g)
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

func (s *Scene) HasCollisionsAt(b *physics.Body, offset gemath.Vec) bool {
	collisions := s.root.collisionEngine.GetCollisions(b, physics.CollisionConfig{
		Offset: offset,
		Limit:  1,
	})
	return len(collisions) != 0
}

func (s *Scene) GetMovementCollision(b *physics.Body, velocity gemath.Vec) *physics.Collision {
	collisions := s.root.collisionEngine.GetCollisions(b, physics.CollisionConfig{
		Velocity: velocity,
		Limit:    1,
	})
	if len(collisions) == 1 {
		return &collisions[0]
	}
	return nil
}
