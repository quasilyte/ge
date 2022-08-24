package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SimpleLayer struct {
	Visible bool

	graphics []SceneGraphics
}

func NewSimpleLayer() *SimpleLayer {
	return &SimpleLayer{
		Visible:  true,
		graphics: make([]SceneGraphics, 0, 8),
	}
}

func (l *SimpleLayer) AddGraphics(g SceneGraphics) {
	l.graphics = append(l.graphics, g)
}

func (l *SimpleLayer) IsDisposed() bool {
	return false
}

func (l *SimpleLayer) Draw(screen *ebiten.Image) {
	if !l.Visible {
		return
	}

	list := l.graphics[:0]
	for _, g := range l.graphics {
		if g.IsDisposed() {
			continue
		}
		g.Draw(screen)
		list = append(list, g)
	}
	l.graphics = list
}
