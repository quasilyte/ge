package main

import (
	"github.com/quasilyte/ge"
)

type playerUnitLabel struct {
	label *ge.Label
	pos   ge.Pos
	text  string
}

func newPlayerUnitLayer(pos ge.Pos, text string) *playerUnitLabel {
	return &playerUnitLabel{pos: pos, text: text}
}

func (l *playerUnitLabel) Init(scene *ge.Scene) {
	l.label = scene.NewLabel(FontSmall)
	l.label.Pos = l.pos
	l.label.Width = 64
	l.label.Height = 48
	l.label.AlignHorizontal = ge.AlignHorizontalCenter
	l.label.AlignVertical = ge.AlignVerticalCenter
	l.label.Text = l.text
	scene.AddGraphicsAbove(l.label, 1)
}

func (l *playerUnitLabel) IsDisposed() bool {
	return l.label.IsDisposed()
}

func (l *playerUnitLabel) Dispose() {
	l.label.Dispose()
}

func (l *playerUnitLabel) Update(delta float64) {}
