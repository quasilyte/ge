package ge

import "github.com/hajimehoshi/ebiten/v2"

type Renderer struct {
	op ebiten.DrawImageOptions
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) Draw(screen *ebiten.Image, graphics []SceneGraphics) []SceneGraphics {
	liveGraphics := graphics[:0]

	for _, g := range graphics {
		if g.IsDisposed() {
			continue
		}

		g.Draw(screen)
		liveGraphics = append(liveGraphics, g)
	}

	return liveGraphics
}
