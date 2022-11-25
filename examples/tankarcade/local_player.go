package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
)

type localPlayer struct {
	input *input.Handler
	unit  *battleUnit
}

func newLocalPlayer(h *input.Handler, u *battleUnit) *localPlayer {
	return &localPlayer{input: h, unit: u}
}

func (p *localPlayer) IsDisposed() bool { return false }

func (p *localPlayer) Init(scene *ge.Scene) {}

func (p *localPlayer) Update(delta float64) {
	if p.input.ActionIsPressed(ActionMoveRight) {
		p.unit.MoveOrder(facingRight)
	} else if p.input.ActionIsPressed(ActionMoveDown) {
		p.unit.MoveOrder(facingDown)
	} else if p.input.ActionIsPressed(ActionMoveLeft) {
		p.unit.MoveOrder(facingLeft)
	} else if p.input.ActionIsPressed(ActionMoveUp) {
		p.unit.MoveOrder(facingUp)
	}
	if p.input.ActionIsPressed(ActionFire) {
		p.unit.FireOrder()
	}
	if p.input.ActionIsPressed(ActionSpecial) {
		p.unit.SpecialOrder()
	}
}
