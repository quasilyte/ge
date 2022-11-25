package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type triggerNode struct {
	name  string
	body  physics.Body
	scene *ge.Scene

	EventActivated gesignal.Event[*triggerNode]
}

func newTriggerNode(pos gmath.Vec, name string) *triggerNode {
	t := &triggerNode{name: name}
	t.body.Pos = pos
	return t
}

func (t *triggerNode) Init(scene *ge.Scene) {
	t.scene = scene
	t.body.InitStaticCircle(t, 16)
	t.body.LayerMask = 0b1
	scene.AddBody(&t.body)
}

func (t *triggerNode) IsDisposed() bool {
	return t.body.IsDisposed()
}

func (t *triggerNode) Activate() {
	t.EventActivated.Emit(t)
	t.body.Dispose()
}

func (t *triggerNode) Update(delta float64) {
	for _, c := range t.scene.GetCollisions(&t.body) {
		if _, ok := c.Body.Object.(*battleUnit); ok {
			t.Activate()
			break
		}
	}
}
