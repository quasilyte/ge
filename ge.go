package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func RunGame(ctx *Context) error {
	if ctx.CurrentScene == nil {
		panic("running game without a scene: Context.CurrentScene is nil")
	}
	g := &gameRunner{ctx: ctx}
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle(ctx.WindowTitle)
	ebiten.SetWindowSize(int(ctx.WindowWidth), int(ctx.WindowHeight))
	return ebiten.RunGame(g)
}

type gameRunner struct {
	ctx *Context
}

func (g *gameRunner) Update() error {
	g.ctx.Input.Update()
	g.ctx.CurrentScene.update(1.0 / 60.0)
	return nil
}

func (g *gameRunner) Draw(screen *ebiten.Image) {
	g.ctx.Draw(screen)
}

func (g *gameRunner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.ctx.WindowWidth), int(g.ctx.WindowHeight)
}
