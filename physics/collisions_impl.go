package physics

import (
	"fmt"
	"math"

	"github.com/quasilyte/gmath"
)

// TODO: a static flag for a Body, so we can avoid layer recalculations.
// TODO: add layer/sector separation, so we don't have O(n^2) complexity.
// TODO: add benchmarks.

type bodyKind int

const (
	bodyCircle bodyKind = iota
	bodyRotatedRect
)

type collisionResolver struct {
	config CollisionConfig

	engine *CollisionEngine

	collisions []Collision
}

func (resolver *collisionResolver) needCollisionNormal() bool {
	return !resolver.config.Velocity.IsZero()
}

func (resolver *collisionResolver) collectCollisionsWith(b, translated *Body, limit int, layerMask uint16, bodies []*Body) {
	for _, b2 := range bodies {
		if len(resolver.collisions) >= limit {
			break
		}
		// This is awkward, but we need to avoid checking the collision with
		// the body itself.
		if b2 == b {
			continue
		}
		intersectedLayers := layerMask & b2.LayerMask
		if intersectedLayers == 0 {
			continue
		}
		collision, ok := resolver.checkCollision(translated, b2)
		if ok {
			collision.Body = b2
			collision.LayerMask = intersectedLayers
			resolver.collisions = append(resolver.collisions, collision)
		}
	}
}

func (resolver *collisionResolver) findCollisions(b, translated *Body, layerMask uint16) []Collision {
	limit := resolver.config.Limit
	if limit == 0 {
		limit = math.MaxInt
	}
	resolver.collectCollisionsWith(b, translated, limit, layerMask, resolver.engine.bodies)
	if !b.static {
		resolver.collectCollisionsWith(b, translated, limit, layerMask, resolver.engine.staticBodies)
	}
	return resolver.collisions
}

func (resolver *collisionResolver) inverseNormal(c Collision, ok bool) (Collision, bool) {
	if ok {
		c.Normal = c.Normal.Neg()
		return c, true
	}
	return c, false
}

func (resolver *collisionResolver) checkCollision(b1, b2 *Body) (Collision, bool) {
	switch b1.kind {
	case bodyCircle:
		switch b2.kind {
		case bodyCircle:
			return resolver.checkCirclesCollision(b1, b2)
		case bodyRotatedRect:
			return resolver.checkCircleRotatedRectCollision(b1, b2)
		}
	case bodyRotatedRect:
		switch b2.kind {
		case bodyCircle:
			return resolver.inverseNormal(resolver.checkCircleRotatedRectCollision(b2, b1))
		case bodyRotatedRect:
			return resolver.checkRotatedRectsCollision(b1, b2)
		}
	}

	panic(fmt.Sprintf("unexpected body kinds combination: %s and %s", b1, b2))
}

func (resolver *collisionResolver) checkCirclesCollision(b1, b2 *Body) (Collision, bool) {
	var result Collision
	if !b1.BoundsRect().Overlaps(b2.BoundsRect()) {
		return result, false
	}
	r1 := b1.CircleRadius()
	r2 := b2.CircleRadius()
	distSqr := b1.Pos.DistanceSquaredTo(b2.Pos)
	if distSqr >= (r1+r2)*(r1+r2) {
		return result, false
	}
	if resolver.needCollisionNormal() {
		// No matter which part of the circle you hit, its normal is always pointing at you.
		result.Normal = b1.Pos.DirectionTo(b2.Pos)
		result.Depth = r1 + r2 - math.Sqrt(distSqr)
	}
	return result, true
}

func (resolver *collisionResolver) checkCircleRotatedRectCollision(circle, rr *Body) (Collision, bool) {
	var result Collision
	if !circle.BoundsRect().Overlaps(rr.BoundsRect()) {
		return result, false
	}

	rectMin := gmath.Vec{
		X: rr.Pos.X - rr.RotatedRectWidth()/2,
		Y: rr.Pos.Y - rr.RotatedRectHeight()/2,
	}

	cosA := math.Cos(float64(-rr.Rotation))
	sinA := math.Sin(float64(-rr.Rotation))
	alignedCirclePos := gmath.Vec{
		X: cosA*(circle.Pos.X-rr.Pos.X) - sinA*(circle.Pos.Y-rr.Pos.Y) + rr.Pos.X,
		Y: sinA*(circle.Pos.X-rr.Pos.X) + cosA*(circle.Pos.Y-rr.Pos.Y) + rr.Pos.Y,
	}

	var closestPos gmath.Vec
	xside := xsideNone
	yside := ysideNone
	if alignedCirclePos.X < rectMin.X {
		closestPos.X = rectMin.X
		xside = xsideLeft
	} else if alignedCirclePos.X > rectMin.X+rr.RotatedRectWidth() {
		closestPos.X = rectMin.X + rr.RotatedRectWidth()
		xside = xsideRight
	} else {
		closestPos.X = alignedCirclePos.X
	}
	if alignedCirclePos.Y < rectMin.Y {
		closestPos.Y = rectMin.Y
		yside = ysideUpper
	} else if alignedCirclePos.Y > rectMin.Y+rr.RotatedRectHeight() {
		closestPos.Y = rectMin.Y + rr.RotatedRectHeight()
		yside = ysideLower
	} else {
		closestPos.Y = alignedCirclePos.Y
	}

	distance := alignedCirclePos.DistanceTo(closestPos)
	if distance >= circle.CircleRadius() {
		return result, false
	}

	if xside+yside == 0 {
		if resolver.needCollisionNormal() {
			result.Normal = resolver.config.Velocity.Normalized().Neg()
			result.Depth = circle.CircleRadius() - distance
		}
		return result, true
	}

	if resolver.needCollisionNormal() {
		result.Normal = gmath.RadToVec(rr.Rotation + rectRotationDelta[xside][yside])
		result.Depth = circle.CircleRadius() - distance
	}

	return result, true
}

func (resolver *collisionResolver) getAxisOverlap(poly1, poly2 []gmath.Vec) (gmath.Vec, float64) {
	var minAxis gmath.Vec
	minOverlap := math.MaxFloat64
	for i := 0; i < len(poly1); i++ {
		axis := getAxisNormal(poly1, i)
		p1 := getPolyProjection(axis, poly1)
		p2 := getPolyProjection(axis, poly2)
		if !p1.HasOverlap(p2) {
			return gmath.Vec{}, 0
		}
		overlap := p1.max - p2.min
		if overlap < minOverlap {
			minOverlap = overlap
			minAxis = axis
		}
	}
	return minAxis, minOverlap
}

func (resolver *collisionResolver) checkRotatedRectsCollision(b1, b2 *Body) (Collision, bool) {
	var result Collision
	if !b1.BoundsRect().Overlaps(b2.BoundsRect()) {
		return result, false
	}

	a := unpackRotatedRect(b1)
	b := unpackRotatedRect(b2)

	axisAB, overlapAB := resolver.getAxisOverlap(a[:], b[:])
	if overlapAB == 0 {
		return result, false
	}
	axisBA, overlapBA := resolver.getAxisOverlap(b[:], a[:])
	if overlapBA == 0 {
		return result, false
	}

	if resolver.needCollisionNormal() {
		if overlapAB < overlapBA {
			result.Depth = overlapAB
			result.Normal = axisAB.Neg()
		} else {
			result.Depth = overlapBA
			result.Normal = axisBA
		}
	}

	return result, true
}

func getAxisNormal(poly []gmath.Vec, i int) gmath.Vec {
	pt1 := poly[i]
	pt2 := poly[0]
	if i+1 < len(poly) {
		pt2 = poly[i+1]
	}
	axis := gmath.Vec{X: pt2.X - pt1.X, Y: pt2.Y - pt1.Y}
	normal := gmath.Vec{X: -axis.Y, Y: axis.X}.Normalized()
	return normal
}

func getPolyProjection(axis gmath.Vec, poly []gmath.Vec) projection {
	if len(poly) == 0 {
		return projection{}
	}
	pmin := axis.Dot(poly[0])
	pmax := pmin
	for i := 1; i < len(poly); i++ {
		dot := axis.Dot(poly[i])
		pmin = fastMin(pmin, dot)
		pmax = fastMax(pmax, dot)
	}
	return projection{min: pmin, max: pmax}
}

type projection struct {
	min float64
	max float64
}

func (p *projection) HasOverlap(other projection) bool {
	return other.min <= p.max && other.max >= p.min
}

var rectRotationDelta = [3][3]gmath.Rad{
	xsideNone: {
		ysideUpper: -(math.Pi / 2),
		ysideLower: math.Pi / 2,
	},
	xsideLeft: {
		ysideNone:  math.Pi,
		ysideUpper: math.Pi + (math.Pi / 4),
		ysideLower: math.Pi - (math.Pi / 4),
	},
	xsideRight: {
		ysideUpper: -(math.Pi / 4),
		ysideLower: math.Pi / 4,
	},
}

func unpackRect(b *Body) RectVertices {
	w2 := b.RotatedRectWidth() / 2
	h2 := b.RotatedRectHeight() / 2
	lr := gmath.Vec{X: b.Pos.X + w2, Y: b.Pos.Y + h2}
	ur := gmath.Vec{X: b.Pos.X + w2, Y: b.Pos.Y - h2}
	ul := gmath.Vec{X: b.Pos.X - w2, Y: b.Pos.Y - h2}
	ll := gmath.Vec{X: b.Pos.X - w2, Y: b.Pos.Y + h2}
	return RectVertices{ur, lr, ll, ul}
}

func unpackRotatedRect(b *Body) RectVertices {
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
	ll := gmath.Vec{X: cx - w2cos - h2sin, Y: cy - w2sin + h2cos}
	lr := gmath.Vec{X: cx + w2cos - h2sin, Y: cy + w2sin + h2cos}
	ul := gmath.Vec{X: cx - w2cos + h2sin, Y: cy - w2sin - h2cos}
	ur := gmath.Vec{X: cx + w2cos + h2sin, Y: cy + w2sin - h2cos}
	return RectVertices{ur, lr, ll, ul}
}

const (
	xsideNone  = 0
	xsideLeft  = 1
	xsideRight = 2
	ysideNone  = 0
	ysideUpper = 1
	ysideLower = 2
)

func fastMin(x, y float64) float64 {
	if x > y {
		return y
	}
	return x
}

func fastMax(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}
