package ge

import "github.com/quasilyte/ge/gemath"

func NewVec(x, y float64) *gemath.Vec {
	return &gemath.Vec{X: x, Y: y}
}

func NewRotation(deg float64) *gemath.Rad {
	rad := gemath.DegToRad(deg)
	return &rad
}
