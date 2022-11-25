package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/quasilyte/ge"
)

type battlePanel struct {
	battleState *battleState

	hpLabels          []*ge.Label
	lastHPValues      []float64
	specialLabels     []*ge.Label
	lastSpecialValues []int
}

func newBattlePanel(s *battleState) *battlePanel {
	return &battlePanel{battleState: s}
}

func (p *battlePanel) Init(scene *ge.Scene) {
	panel := scene.NewSprite(ImageInterfacePanel)
	panel.Pos.Offset.Y = 896
	panel.Centered = false
	scene.AddGraphicsAbove(panel, 1)

	{
		p.lastHPValues = make([]float64, 3)
		p.lastSpecialValues = []int{-1, -1, -1}

		offset := 0.0
		for i, player := range p.battleState.players {
			if !player.active {
				break
			}
			l := scene.NewLabel(FontSmall)
			l.Text = fmt.Sprintf("Player %d |", i+1)
			l.ColorScale.SetColor(ge.RGB(0xe42cca))
			l.Pos.Offset.Y = 896 + 32 + offset
			l.Pos.Offset.X = 64
			scene.AddGraphicsAbove(l, 1)

			value := scene.NewLabel(FontSmall)
			value.Text = "HP: ?"
			value.ColorScale.SetColor(ge.RGB(0xe42cca))
			value.Pos.Offset.Y = 896 + 32 + offset
			value.Pos.Offset.X = 64 + 288
			scene.AddGraphicsAbove(value, 1)
			p.hpLabels = append(p.hpLabels, value)

			specialLabel := scene.NewLabel(FontSmall)
			specialLabel.Text = "Special: ?"
			specialLabel.ColorScale.SetColor(ge.RGB(0xe42cca))
			specialLabel.Pos.Offset.Y = 896 + 32 + offset
			specialLabel.Pos.Offset.X = 64 + 288 + 224
			scene.AddGraphicsAbove(specialLabel, 1)
			p.specialLabels = append(p.specialLabels, specialLabel)

			offset += 48
		}
	}
}

func (p *battlePanel) IsDisposed() bool {
	return false
}

func (p *battlePanel) Update(delta float64) {
	for i, player := range p.battleState.players {
		if !player.active {
			break
		}
		hpLabel := p.hpLabels[i]
		specialLabel := p.specialLabels[i]
		if player.unit == nil {
			hpLabel.Text = "HP: ?"
			specialLabel.Text = "Special: ?"
			continue
		}
		if player.unit.hp != p.lastHPValues[i] {
			p.lastHPValues[i] = player.unit.hp
			hpLabel.Text = "HP: " + strconv.Itoa(int(math.Round(player.unit.hp)))
		}
		if player.unit.SpecialAmmo != p.lastSpecialValues[i] {
			p.lastSpecialValues[i] = player.unit.SpecialAmmo
			specialLabel.Text = "Special: " + strconv.Itoa(player.unit.SpecialAmmo)
		}
	}
}
