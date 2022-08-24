package ge

import "github.com/hajimehoshi/ebiten/v2"

type MultiLayer struct {
	List []SceneGraphicsLayer
}

func NewMultiLayer(layers ...SceneGraphicsLayer) *MultiLayer {
	if len(layers) == 0 {
		panic("numLayers should be greated than 0")
	}
	return &MultiLayer{
		List: layers,
	}
}

func NewMultiSimpleLayer(numLayers int) *MultiLayer {
	layers := make([]SceneGraphicsLayer, numLayers)
	for i := range layers {
		layers[i] = NewSimpleLayer()
	}
	return NewMultiLayer(layers...)
}

func (l *MultiLayer) AddGraphics(g SceneGraphics) {
	l.List[0].AddGraphics(g)
}

func (l *MultiLayer) IsDisposed() bool {
	return false
}

func (l *MultiLayer) Draw(screen *ebiten.Image) {
	for i := range l.List {
		l.List[i].Draw(screen)
	}
}
