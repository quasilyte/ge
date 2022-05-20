package ge

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func RunGame(ctx *Context) error {
	if ctx.CurrentScene == nil {
		panic("running game without a scene: Context.CurrentScene is nil")
	}
	g := &gameRunner{
		ctx:      ctx,
		prevTime: time.Now(),
	}
	if ctx.FullScreen {
		ebiten.SetFullscreen(true)
	}
	ebiten.SetWindowTitle(ctx.WindowTitle)
	ebiten.SetWindowSize(int(ctx.WindowWidth), int(ctx.WindowHeight))
	return ebiten.RunGame(g)
}

type gameRunner struct {
	ctx      *Context
	prevTime time.Time
}

func (g *gameRunner) Update() error {
	now := time.Now()
	timeDelta := now.Sub(g.prevTime).Seconds()
	g.prevTime = now

	g.ctx.Input.Update()
	g.ctx.CurrentScene.update(timeDelta)
	return nil
}

func (g *gameRunner) Draw(screen *ebiten.Image) {
	g.ctx.Draw(screen)
}

func (g *gameRunner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.ctx.WindowWidth), int(g.ctx.WindowHeight)
}
