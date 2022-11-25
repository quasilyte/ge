package ge

import "github.com/quasilyte/gmath"

// Pos represents a position with optional offset relative to its base.
type Pos struct {
	Base   *gmath.Vec
	Offset gmath.Vec
}

func MakePos(base gmath.Vec) Pos {
	return Pos{Base: &base}
}

func (p Pos) Resolve() gmath.Vec {
	if p.Base == nil {
		return p.Offset
	}
	return p.Base.Add(p.Offset)
}

func (p *Pos) SetBase(base gmath.Vec) {
	p.Base = &base
}

func (p *Pos) Set(base *gmath.Vec, offsetX, offsetY float64) {
	p.Base = base
	p.Offset.X = offsetX
	p.Offset.Y = offsetY
}

func (p Pos) WithOffset(offsetX, offsetY float64) Pos {
	return Pos{
		Base:   p.Base,
		Offset: gmath.Vec{X: p.Offset.X + offsetX, Y: p.Offset.Y + offsetY},
	}
}
