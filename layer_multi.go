package ge

import "github.com/hajimehoshi/ebiten/v2"

type MultiLayer struct {
	Layers []SceneGraphicsLayer
}

func NewMultiLayer(layers ...SceneGraphicsLayer) *MultiLayer {
	if len(layers) == 0 {
		panic("numLayers should be greated than 0")
	}
	return &MultiLayer{
		Layers: layers,
	}
}

func (l *MultiLayer) AddGraphics(g SceneGraphics) {
	l.Layers[0].AddGraphics(g)
}

func (l *MultiLayer) IsDisposed() bool {
	return false
}

func (l *MultiLayer) Draw(screen *ebiten.Image) {
	for i := range l.Layers {
		l.Layers[i].Draw(screen)
	}
}
