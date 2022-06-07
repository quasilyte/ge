package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

type railgunRay struct {
	from     gemath.Vec
	to       gemath.Vec
	line     *ge.Line
	alpha    float64
	disposed bool
}

func newRailgunRay(from, to gemath.Vec) *railgunRay {
	return &railgunRay{from: from, to: to}
}

func (ray *railgunRay) IsDisposed() bool {
	return ray.line.IsDisposed()
}

func (ray *railgunRay) Dispose() {
	ray.line.Dispose()
}

func (ray *railgunRay) Init(scene *ge.Scene) {
	ray.line = ge.NewLine(&ray.from, &ray.to)
	ray.line.ColorScale.SetRGBA(255, 100, 180, 255)
	ray.line.Width = 3
	scene.AddGraphics(ray.line)
}

func (ray *railgunRay) Update(delta float64) {
	ray.line.ColorScale.A -= float32(delta * 4)
	if ray.line.ColorScale.A < 0.2 {
		ray.Dispose()
	}
}
