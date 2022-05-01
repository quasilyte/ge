package physics

import (
	"math"

	"github.com/quasilyte/ge/gemath"
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

func (resolver *collisionResolver) findCollisions(b, translated *Body, layerMask uint16) []Collision {
	limit := resolver.config.Limit
	if limit == 0 {
		limit = math.MaxInt
	}
	for _, b2 := range resolver.engine.bodies {
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
	return resolver.collisions
}

func (resolver *collisionResolver) checkCollision(b1, b2 *Body) (Collision, bool) {
	switch b1.kind {
	case bodyCircle:
		switch b2.kind {
		case bodyCircle:
			return resolver.checkCircleCircleCollision(b1, b2)
		case bodyRotatedRect:
			return resolver.checkCircleRotatedRectCollision(b1, b2)
		}
	}

	panic("unexpected body kinds combination")
}

func (resolver *collisionResolver) checkCircleCircleCollision(b1, b2 *Body) (Collision, bool) {
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

	rectMin := gemath.Vec{
		X: rr.Pos.X - rr.RotatedRectWidth()/2,
		Y: rr.Pos.Y - rr.RotatedRectHeight()/2,
	}

	cosA := math.Cos(float64(-rr.Rotation))
	sinA := math.Sin(float64(-rr.Rotation))
	alignedCirclePos := gemath.Vec{
		X: cosA*(circle.Pos.X-rr.Pos.X) - sinA*(circle.Pos.Y-rr.Pos.Y) + rr.Pos.X,
		Y: sinA*(circle.Pos.X-rr.Pos.X) + cosA*(circle.Pos.Y-rr.Pos.Y) + rr.Pos.Y,
	}

	var closestPos gemath.Vec
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
			result.Normal = resolver.config.Velocity.Neg().Normalized()
			result.Depth = circle.CircleRadius() - distance
		}
		return result, true
	}

	if resolver.needCollisionNormal() {
		result.Normal = gemath.RadToVec(rr.Rotation + rectRotationDelta[xside][yside])
		result.Depth = circle.CircleRadius() - distance
	}

	return result, true
}

var rectRotationDelta = [3][3]gemath.Rad{
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
	ll := gemath.Vec{X: cx - w2cos - h2sin, Y: cy - w2sin + h2cos}
	lr := gemath.Vec{X: cx + w2cos - h2sin, Y: cy + w2sin + h2cos}
	ul := gemath.Vec{X: cx - w2cos + h2sin, Y: cy - w2sin - h2cos}
	ur := gemath.Vec{X: cx + w2cos + h2sin, Y: cy + w2sin - h2cos}
	return RectVertices{ul, ur, lr, ll}
}

const (
	xsideNone  = 0
	xsideLeft  = 1
	xsideRight = 2
	ysideNone  = 0
	ysideUpper = 1
	ysideLower = 2
)

type rotatedRect struct {
	ul gemath.Vec // upper-left
	ur gemath.Vec // upper-right
	ll gemath.Vec // lower-left
	lr gemath.Vec // lower-right
}

func getAxis(pt1, pt2 gemath.Vec) gemath.Vec {
	return gemath.Vec{X: pt1.X - pt2.X, Y: pt1.Y - pt2.Y}
}
