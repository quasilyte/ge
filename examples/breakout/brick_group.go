package main

import (
	"github.com/quasilyte/ge/gemath"
)

type brickGroup struct {
	bricks    []*brick
	rotate    gemath.Rad
	dx        float64
	slideTime float64
	slide     float64
}

func (g *brickGroup) Update(delta float64) {
	for _, b := range g.bricks {
		if b.IsDisposed() {
			continue
		}
		b.body.Rotation += g.rotate * gemath.Rad(delta)
		b.body.Pos.X += g.dx * delta
		if g.slide >= g.slideTime {
			g.dx = -g.dx
			g.slide = 0
		}
	}
	g.slide += delta
}

func (g *brickGroup) Reset() {
	g.bricks = g.bricks[:0]
}
