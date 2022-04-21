package collision

type Engine struct {
	bodies []*Body

	info Info
}

func (e *Engine) AddBody(b *Body) {
	e.bodies = append(e.bodies, b)
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
