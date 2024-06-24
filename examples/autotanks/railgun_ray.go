package main

import (
	"github.com/quasilyte/ge"

	"github.com/quasilyte/gmath"
)

type railgunRay struct {
	from       gmath.Vec
	to         gmath.Vec
	line       *ge.Line
	colorScale ge.ColorScale
	alpha      float64
	disposed   bool
}

func newRailgunRay(from, to gmath.Vec) *railgunRay {
	return &railgunRay{from: from, to: to}
}

func (ray *railgunRay) IsDisposed() bool {
	return ray.line.IsDisposed()
}

func (ray *railgunRay) Dispose() {
	ray.line.Dispose()
}

func (ray *railgunRay) Init(scene *ge.Scene) {
	ray.line = ge.NewLine(ge.MakePos(ray.from), ge.MakePos(ray.to))
	ray.colorScale = ge.ColorScale{R: 255, G: 100, B: 180, A: 255}
	ray.line.SetColorScale(ray.colorScale)
	ray.line.Width = 3
	scene.AddGraphics(ray.line)
}

func (ray *railgunRay) Update(delta float64) {
	ray.colorScale.A -= float32(delta * 4)
	ray.line.SetColorScale(ray.colorScale)
	if ray.colorScale.A < 0.2 {
		ray.Dispose()
	}
}
