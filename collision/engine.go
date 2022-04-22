package collision

import "github.com/quasilyte/ge/gemath"

type Engine struct {
	bodies []*Body

	info Info
}

func (e *Engine) AddBody(b *Body) {
	e.bodies = append(e.bodies, b)
}

func (e *Engine) HasCollisionsAt(b *Body, pos gemath.Vec) bool {
	b1 := *b
	b1.Pos = pos
	for _, b2 := range e.bodies {
		if b == b2 || b2.IsDisposed() {
			continue
		}
		if b1.LayerMask&b2.LayerMask == 0 {
			continue
		}
		if !Intersection(&b1, b2).IsEmpty() {
			return true
		}
	}
	return false
}

func (e *Engine) Calculate() {
	liveBodies := e.bodies[:0]
	for _, b := range e.bodies {
		if b.IsDisposed() {
			continue
		}
		liveBodies = append(liveBodies, b)
	}
	e.bodies = liveBodies

	info := &e.info
	for i, b1 := range e.bodies {
		for j := i + 1; j < len(e.bodies); j++ {
			b2 := e.bodies[j]
			// Even if these two bodies collide, they won't notice it.
			if b1.CollisionHandler == nil && b2.CollisionHandler == nil {
				continue
			}
			if b1.LayerMask&b2.LayerMask == 0 {
				continue
			}
			area := Intersection(b1, b2)
			if !area.IsEmpty() {
				info.Rect = area
				if b1.CollisionHandler != nil {
					info.Object = b2.Object
					b1.CollisionHandler.OnCollision(info)
				}
				if b2.CollisionHandler != nil {
					info.Object = b1.Object
					b2.CollisionHandler.OnCollision(info)
				}
			}
		}
	}
}
