package ge

import "github.com/quasilyte/ge/gemath"

// Pos represents a position with optional offset relative to its base.
type Pos struct {
	Base   *gemath.Vec
	Offset gemath.Vec
}

func MakePos(base gemath.Vec) Pos {
	return Pos{Base: &base}
}

func (p Pos) Resolve() gemath.Vec {
	return p.Base.Add(p.Offset)
}

func (p *Pos) SetBase(base gemath.Vec) {
	p.Base = &base
}

func (p *Pos) Set(base *gemath.Vec, offsetX, offsetY float64) {
	p.Base = base
	p.Offset.X = offsetX
	p.Offset.Y = offsetY
}

func (p Pos) WithOffset(offsetX, offsetY float64) Pos {
	return Pos{
		Base:   p.Base,
		Offset: gemath.Vec{X: p.Offset.X + offsetX, Y: p.Offset.Y + offsetY},
	}
}
