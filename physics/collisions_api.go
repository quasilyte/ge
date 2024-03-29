package physics

import (
	"fmt"
	"math"

	"github.com/quasilyte/gmath"
)

type CollisionEngine struct {
	bodies       []*Body
	staticBodies []*Body

	translatedBody Body

	collisionPool []Collision
}

type CollisionConfig struct {
	Offset gmath.Vec

	// If velocity magnitude is not 0, collisions are calculated
	// in the dynamics of the movement.
	Velocity gmath.Vec

	// If not 0, this mask will be used instead of the object own mask.
	// It can be used to both include or exclude some other objects from
	// the collisions calculation.
	LayerMask uint16

	Limit int
}

// CalculateFrame re-calculates layers and other things for the upcoming frame.
// This function is called from the framework itself in the beginning of each frame.
func (e *CollisionEngine) CalculateFrame() {
	live := e.bodies[:0]
	liveStatic := e.staticBodies[:0]
	for _, b := range e.bodies {
		if b.IsDisposed() {
			continue
		}
		if b.static {
			liveStatic = append(liveStatic, b)
		} else {
			live = append(live, b)
		}
	}
	for _, b := range e.staticBodies {
		if b.IsDisposed() {
			continue
		}
		if b.static {
			liveStatic = append(liveStatic, b)
		} else {
			live = append(live, b)
		}
	}
	e.bodies = live
	e.staticBodies = liveStatic
}

// AddBody includes the given body into the collision space.
//
// Note that it will not be ready to register or cause any collisions until
// the next frame.
func (e *CollisionEngine) AddBody(b *Body) {
	if b.static {
		e.staticBodies = append(e.staticBodies, b)
	} else {
		e.bodies = append(e.bodies, b)
	}
}

// GetCollisions returns all colliders for the specified body.
// A config can affect the rules of this collision computation.
func (e *CollisionEngine) GetCollisions(b *Body, config CollisionConfig) []Collision {
	translated := b
	if !config.Offset.IsZero() {
		translated = &e.translatedBody
		*translated = *b
		translated.Pos = b.Pos.Add(config.Offset)
	}
	layerMask := b.LayerMask
	if config.LayerMask != 0 {
		layerMask = config.LayerMask
	}
	resolver := collisionResolver{
		config:     config,
		engine:     e,
		collisions: e.collisionPool[:0],
	}
	e.collisionPool = resolver.findCollisions(b, translated, layerMask)
	return e.collisionPool
}

type Collision struct {
	// Body is a body being collided.
	// To get a collision object, access `Collision.Body.Object`.
	Body *Body

	// LayerMask represents the layer masks intersection of the colliding objects.
	LayerMask uint16

	// Normal is a contacted surface collision normal vector.
	// Collision normal vector has unit length (it's normalized).
	//
	// Note: a normal is computed only when resolving with non-zero velocity.
	Normal gmath.Vec

	Depth float64
}

type Body struct {
	Object interface{}

	Rotation gmath.Rad

	Pos gmath.Vec

	LayerMask uint16

	kind     bodyKind
	disposed bool
	static   bool

	value1 float64
	value2 float64
}

func (b *Body) IsDisposed() bool {
	return b.disposed
}

func (b *Body) Dispose() { b.disposed = true }

func (b *Body) InitStaticCircle(o interface{}, radius float64) {
	b.InitCircle(o, radius)
	b.static = true
}

func (b *Body) InitCircle(o interface{}, radius float64) {
	*b = Body{
		Pos:       b.Pos,
		Rotation:  b.Rotation,
		Object:    o,
		LayerMask: 1,
		kind:      bodyCircle,
		value1:    radius,
	}
}

func (b *Body) InitStaticRotatedRect(o interface{}, width, height float64) {
	b.InitRotatedRect(o, width, height)
	b.static = true
}

func (b *Body) InitRotatedRect(o interface{}, width, height float64) {
	*b = Body{
		Pos:       b.Pos,
		Rotation:  b.Rotation,
		Object:    o,
		LayerMask: 1,
		kind:      bodyRotatedRect,
		value1:    width,
		value2:    height,
	}
}

func (b *Body) IsCircle() bool { return b.kind == bodyCircle }

func (b *Body) CircleRadius() float64 { return b.value1 }

func (b *Body) IsRotatedRect() bool { return b.kind == bodyRotatedRect }

func (b *Body) RotatedRectWidth() float64 { return b.value1 }

func (b *Body) RotatedRectHeight() float64 { return b.value2 }

func (b *Body) RotatedRectVertices() RectVertices {
	return unpackRotatedRect(b)
}

func (b *Body) BoundsRect() gmath.Rect {
	// TODO: could precompute dx and dy for both min/max points
	// to make bounds computation in relation to `b.Pos` faster.
	switch b.kind {
	case bodyCircle:
		min := gmath.Vec{
			X: b.Pos.X - b.CircleRadius(),
			Y: b.Pos.Y - b.CircleRadius(),
		}
		max := gmath.Vec{
			X: b.Pos.X + b.CircleRadius(),
			Y: b.Pos.Y + b.CircleRadius(),
		}
		return gmath.Rect{Min: min, Max: max}

	case bodyRotatedRect:
		side := math.Max(b.RotatedRectWidth(), b.RotatedRectHeight())
		xy1 := gmath.Vec{
			X: b.Pos.X - side/2,
			Y: b.Pos.Y - side/2,
		}
		xy2 := gmath.Vec{
			X: b.Pos.X + side/2,
			Y: b.Pos.Y + side/2,
		}
		return gmath.Rect{Min: xy1, Max: xy2}

	default:
		return gmath.Rect{}
	}
}

func (b Body) String() string {
	switch b.kind {
	case bodyCircle:
		return fmt.Sprintf("circle{pos:%v, radius:%f}", b.Pos, b.CircleRadius())
	case bodyRotatedRect:
		return fmt.Sprintf("rotatedRect{pos:%v, rotation:%v, width:%f, height: %f}",
			b.Pos, b.Rotation, b.RotatedRectWidth(), b.RotatedRectHeight())
	default:
		return "?"
	}
}

type RectVertices [4]gmath.Vec

func (v *RectVertices) UR() gmath.Vec { return (*v)[0] }
func (v *RectVertices) LR() gmath.Vec { return (*v)[1] }
func (v *RectVertices) LL() gmath.Vec { return (*v)[2] }
func (v *RectVertices) UL() gmath.Vec { return (*v)[3] }
