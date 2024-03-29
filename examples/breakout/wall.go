package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type wall struct {
	body physics.Body
}

func newWall(width, height float64, rotation gmath.Rad) *wall {
	w := &wall{}
	w.body.InitRotatedRect(w, width, height)
	w.body.Rotation = rotation
	return w
}

func (w *wall) Init(scene *ge.Scene) {
	scene.AddBody(&w.body)
}

func (w *wall) IsDisposed() bool { return w.body.IsDisposed() }

func (w *wall) Update(delta float64) {}
