package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ShaderLayer struct {
	Visible bool

	Shader Shader

	tmp    *ebiten.Image
	width  int
	height int

	disposed bool

	graphics []SceneGraphics
}

func NewShaderLayer() *ShaderLayer {
	return &ShaderLayer{
		Visible:  true,
		graphics: make([]SceneGraphics, 0, 8),
	}
}

func (l *ShaderLayer) AddGraphics(g SceneGraphics) {
	l.graphics = append(l.graphics, g)
}

func (l *ShaderLayer) IsDisposed() bool {
	return l.disposed
}

func (l *ShaderLayer) Draw(screen *ebiten.Image) {
	if !l.Visible {
		return
	}

	list := l.graphics[:0]
	for _, g := range l.graphics {
		if g.IsDisposed() {
			continue
		}
		list = append(list, g)
	}
	l.graphics = list

	shaderEnabled := l.Shader.Enabled && !l.Shader.IsNil()
	if !shaderEnabled {
		for _, g := range l.graphics {
			g.Draw(screen)
		}
		return
	}

	if l.tmp == nil {
		l.width, l.height = screen.Size()
		l.tmp = ebiten.NewImage(l.width, l.height)
	} else {
		l.tmp.Clear()
	}

	for _, g := range l.graphics {
		g.Draw(l.tmp)
	}

	var options ebiten.DrawRectShaderOptions
	options.Images[0] = l.tmp
	options.Images[1] = l.Shader.Texture1.Data
	options.Images[2] = l.Shader.Texture2.Data
	options.Images[3] = l.Shader.Texture3.Data
	screen.DrawRectShader(l.width, l.height, l.Shader.compiled, &options)
}
