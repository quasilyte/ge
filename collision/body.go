package collision

import (
	"fmt"
	"math"

	"github.com/quasilyte/ge/gemath"
)

type Handler interface {
	OnCollision(*Info)
}

type Info struct {
	// Object is an object being collided.
	Object interface{}

	// Rect is a bounding rect of the intersected zone.
	Rect gemath.Rect
}

type Body struct {
	Object interface{}

	CollisionHandler Handler

	LayerMask uint16

	Rotation gemath.Rad
	Pos      gemath.Vec

	value1 float64
	value2 float64

	kind     shapeKind
	disposed bool
}

type shapeKind uint16

const (
	bodyNone shapeKind = iota
	bodyRect
	bodyRotatedRect
	bodyCircle
)

func (b *Body) InitCircle(o interface{}, radius float64) {
	*b = Body{
		Object:    o,
		LayerMask: 1,
		kind:      bodyCircle,
		value1:    radius,
	}
}

func (b *Body) InitRect(o interface{}, width, height float64) {
	*b = Body{
		Object:    o,
		LayerMask: 1,
		kind:      bodyRect,
		value1:    width,
		value2:    height,
	}
}

func (b *Body) InitRotatedRect(o interface{}, width, height float64) {
	*b = Body{
		Object:    o,
		LayerMask: 1,
		kind:      bodyRotatedRect,
		value1:    width,
		value2:    height,
	}
}

func (b *Body) IsDisposed() bool {
	return b.disposed
}

func (b *Body) Dispose() { b.disposed = true }

func (b *Body) IsRect() bool {
	return b.kind == bodyRect
}

func (b *Body) IsRotatedRect() bool {
	return b.kind == bodyRotatedRect
}

func (b *Body) IsCircle() bool {
	return b.kind == bodyCircle
}

func (b *Body) BoundsRect() gemath.Rect {
	switch b.kind {
	case bodyCircle:
		xy1 := gemath.Vec{
			X: b.Pos.X - b.CircleRadius(),
			Y: b.Pos.Y - b.CircleRadius(),
		}
		xy2 := gemath.Vec{
			X: b.Pos.X + b.CircleRadius(),
			Y: b.Pos.Y + b.CircleRadius(),
		}
		return gemath.Rect{Min: xy1, Max: xy2}

	case bodyRotatedRect:
		side := math.Max(b.RotatedRectWidth(), b.RotatedRectHeight())
		xy1 := gemath.Vec{
			X: b.Pos.X - side/2,
			Y: b.Pos.Y - side/2,
		}
		xy2 := gemath.Vec{
			X: b.Pos.X + side/2,
			Y: b.Pos.Y + side/2,
		}
		return gemath.Rect{Min: xy1, Max: xy2}

	case bodyRect:
		return gemath.Rect{Min: b.RectMin(), Max: b.RectMax()}

	default:
		return gemath.Rect{}
	}
}

func (b *Body) RectWidth() float64 {
	return b.value1
}

func (b *Body) RectHeight() float64 {
	return b.value2
}

func (b *Body) RectMin() gemath.Vec {
	return gemath.Vec{X: b.Pos.X - b.RectWidth()/2, Y: b.Pos.Y - b.RectHeight()/2}
}

func (b *Body) RectMax() gemath.Vec {
	return gemath.Vec{X: b.Pos.X + b.RectWidth()/2, Y: b.Pos.Y + b.RectHeight()/2}
}

func (b *Body) RotatedRectWidth() float64 {
	return b.value1
}

func (b *Body) RotatedRectHeight() float64 {
	return b.value2
}

func (b *Body) RotatedRectVertices() [4]gemath.Vec {
	r := unpackRotatedRect(b)
	return [4]gemath.Vec{r.ul, r.ur, r.lr, r.ll}
}

func (b *Body) CircleRadius() float64 {
	return b.value1
}

func (b *Body) String() string {
	if b.IsCircle() {
		return fmt.Sprintf("circle{radius:%f}", b.CircleRadius())
	}
	if b.IsRect() {
		return fmt.Sprintf("rect{width:%f,height:%f}", b.RectWidth(), b.RectHeight())
	}
	return "?"
}

func Intersection(b1, b2 *Body) gemath.Rect {
	switch b1.kind {
	case bodyRect:
		switch b2.kind {
		case bodyRect:
			return rectsIntersection(b1.BoundsRect(), b2.BoundsRect())
		case bodyRotatedRect:
			// TODO: something more clever?
			return rotatedRectsIntersection(b1, b2)
		case bodyCircle:
			return rectCircleIntersection(b1, b2)
		}
	case bodyRotatedRect:
		switch b2.kind {
		case bodyRect:
			return Intersection(b2, b1)
		case bodyRotatedRect:
			return rotatedRectsIntersection(b1, b2)
		case bodyCircle:
			return rotatedRectCircleIntersection(b1, b2)
		}
	case bodyCircle:
		switch b2.kind {
		case bodyRect:
			return Intersection(b2, b1)
		case bodyRotatedRect:
			return Intersection(b2, b1)
		case bodyCircle:
			return circlesIntersection(b1, b2)
		}
	}
	panic(fmt.Sprintf("unimplemented collision: %d vs %d", b1.kind, b2.kind))
}

func rectsIntersection(rect1, rect2 gemath.Rect) gemath.Rect {
	if rect1.Min.X < rect2.Min.X {
		rect1.Min.X = rect2.Min.X
	}
	if rect1.Min.Y < rect2.Min.Y {
		rect1.Min.Y = rect2.Min.Y
	}
	if rect1.Max.X > rect2.Max.X {
		rect1.Max.X = rect2.Max.X
	}
	if rect1.Max.Y > rect2.Max.Y {
		rect1.Max.Y = rect2.Max.Y
	}
	if rect1.IsEmpty() {
		return gemath.Rect{}
	}
	return rect1
}

func circlesIntersection(b1, b2 *Body) gemath.Rect {
	rect := rectsIntersection(b1.BoundsRect(), b2.BoundsRect())
	if rect.IsEmpty() {
		return gemath.Rect{}
	}
	x1 := b1.Pos.X
	y1 := b1.Pos.Y
	r1 := b1.CircleRadius()
	x2 := b2.Pos.X
	y2 := b2.Pos.Y
	r2 := b2.CircleRadius()
	collide := math.Abs((x1-x2)*(x1-x2)+(y1-y2)*(y1-y2)) < (r1+r2)*(r1+r2)
	if collide {
		return rect
	}
	return gemath.Rect{}
}

func rectCircleIntersection(rect, circle *Body) gemath.Rect {
	bounds := rectsIntersection(rect.BoundsRect(), circle.BoundsRect())
	if bounds.IsEmpty() {
		return gemath.Rect{}
	}

	rectMin := rect.RectMin()
	var closestPos gemath.Vec
	if circle.Pos.X < rectMin.X {
		closestPos.X = rectMin.X
	} else if circle.Pos.X > rectMin.X+rect.RotatedRectWidth() {
		closestPos.X = rectMin.X + rect.RotatedRectWidth()
	} else {
		closestPos.X = circle.Pos.X
	}
	if circle.Pos.Y < rectMin.Y {
		closestPos.Y = rectMin.Y
	} else if circle.Pos.Y > rectMin.Y+rect.RotatedRectHeight() {
		closestPos.Y = rectMin.Y + rect.RotatedRectHeight()
	} else {
		closestPos.Y = circle.Pos.Y
	}

	if circle.Pos.DistanceTo(closestPos) >= circle.CircleRadius() {
		return gemath.Rect{}
	}
	return bounds
}

func rotatedRectCircleIntersection(rect, circle *Body) gemath.Rect {
	bounds := rectsIntersection(rect.BoundsRect(), circle.BoundsRect())
	if bounds.IsEmpty() {
		return gemath.Rect{}
	}

	rectMin := rect.RectMin()

	cosA := math.Cos(float64(-rect.Rotation))
	sinA := math.Sin(float64(-rect.Rotation))
	alignedCirclePos := gemath.Vec{
		X: cosA*(circle.Pos.X-rect.Pos.X) - sinA*(circle.Pos.Y-rect.Pos.Y) + rect.Pos.X,
		Y: sinA*(circle.Pos.X-rect.Pos.X) + cosA*(circle.Pos.Y-rect.Pos.Y) + rect.Pos.Y,
	}

	var closestPos gemath.Vec
	if alignedCirclePos.X < rectMin.X {
		closestPos.X = rectMin.X
	} else if alignedCirclePos.X > rectMin.X+rect.RotatedRectWidth() {
		closestPos.X = rectMin.X + rect.RotatedRectWidth()
	} else {
		closestPos.X = alignedCirclePos.X
	}
	if alignedCirclePos.Y < rectMin.Y {
		closestPos.Y = rectMin.Y
	} else if alignedCirclePos.Y > rectMin.Y+rect.RotatedRectHeight() {
		closestPos.Y = rectMin.Y + rect.RotatedRectHeight()
	} else {
		closestPos.Y = alignedCirclePos.Y
	}

	if alignedCirclePos.DistanceTo(closestPos) >= circle.CircleRadius() {
		return gemath.Rect{}
	}
	return bounds
}

func rotatedRectsIntersection(b1, b2 *Body) gemath.Rect {
	rect := rectsIntersection(b1.BoundsRect(), b2.BoundsRect())
	if rect.IsEmpty() {
		return gemath.Rect{}
	}

	a := unpackRotatedRect(b1)
	b := unpackRotatedRect(b2)

	hasCollision := checkAxisOverlap(a.Axis1(), a.ul, a.ur, &b) &&
		checkAxisOverlap(a.Axis2(), a.ur, a.lr, &b) &&
		checkAxisOverlap(b.Axis1(), b.ul, b.ur, &a) &&
		checkAxisOverlap(b.Axis2(), b.ur, b.lr, &a)
	if !hasCollision {
		return gemath.Rect{}
	}

	return rect
}

func unpackRotatedRect(b *Body) rotatedRect {
	cosA := b.Rotation.Cos()
	sinA := b.Rotation.Sin()
	w2 := b.RotatedRectWidth() / 2
	h2 := b.RotatedRectHeight() / 2
	w2cos := w2 * cosA
	w2sin := w2 * sinA
	h2sin := h2 * sinA
	h2cos := h2 * cosA
	cx := b.Pos.X
	cy := b.Pos.Y
	ul := gemath.Vec{X: cx - w2cos - h2sin, Y: cy - w2sin + h2cos}
	ur := gemath.Vec{X: cx + w2cos - h2sin, Y: cy + w2sin + h2cos}
	ll := gemath.Vec{X: cx - w2cos + h2sin, Y: cy - w2sin - h2cos}
	lr := gemath.Vec{X: cx + w2cos + h2sin, Y: cy + w2sin - h2cos}
	return rotatedRect{ul: ul, ur: ur, ll: ll, lr: lr}
}

type rotatedRect struct {
	ul gemath.Vec // upper-left
	ur gemath.Vec // upper-right
	ll gemath.Vec // lower-left
	lr gemath.Vec // lower-right
}

func (rr *rotatedRect) Axis1() gemath.Vec {
	return gemath.Vec{X: rr.ur.X - rr.ul.X, Y: rr.ur.Y - rr.ul.Y}
}

func (rr *rotatedRect) Axis2() gemath.Vec {
	return gemath.Vec{X: rr.ur.X - rr.lr.X, Y: rr.ur.Y - rr.lr.Y}
}

func projectVector(axis, v gemath.Vec) gemath.Vec {
	a := (v.X * axis.X) + (v.Y * axis.Y)
	b := (axis.X * axis.X) + (axis.Y * axis.Y)
	lhs := a / b
	return gemath.Vec{X: lhs * axis.X, Y: lhs * axis.Y}
}

func checkAxisOverlap(axis gemath.Vec, av1, av2 gemath.Vec, b *rotatedRect) bool {
	pav1 := projectVector(axis, av1)
	pav2 := projectVector(axis, av2)
	aMin := axis.Dot(pav1)
	aMax := axis.Dot(pav2)
	if aMax < aMin {
		aMax, aMin = aMin, aMax
	}
	pbv1 := projectVector(axis, b.ll)
	pbv2 := projectVector(axis, b.lr)
	pbv3 := projectVector(axis, b.ul)
	pbv4 := projectVector(axis, b.ur)
	bDot := [4]float64{
		axis.Dot(pbv1),
		axis.Dot(pbv2),
		axis.Dot(pbv3),
		axis.Dot(pbv4),
	}
	bMin := bDot[0]
	bMax := bDot[3]
	for _, bValue := range bDot {
		if bValue < bMin {
			bMin = bValue
		}
		if bValue > bMax {
			bMax = bValue
		}
	}
	return bMin <= aMax && bMax >= aMin
}
