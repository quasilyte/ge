package ge

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/internal/locales"
)

func InferLanguages() []string {
	return locales.InferLanguages()
}

func RunGame(ctx *Context, controller SceneController) error {
	g := &gameRunner{
		ctx:      ctx,
		prevTime: time.Now(),
	}
	ebiten.SetFullscreen(ctx.FullScreen)
	ctx.firstController = controller
	ebiten.SetWindowTitle(ctx.WindowTitle)
	ebiten.SetWindowSize(int(ctx.WindowWidth), int(ctx.WindowHeight))
	return ebiten.RunGame(g)
}

type gameRunner struct {
	ctx      *Context
	prevTime time.Time
}

func (g *gameRunner) Update() error {
	g.ctx.Input.Update()
	g.ctx.Audio.Update()

	if g.ctx.CurrentScene == nil && g.ctx.firstController != nil {
		g.ctx.ChangeScene(g.ctx.firstController)
		g.ctx.firstController = nil
	}

	now := time.Now()
	timeDelta := now.Sub(g.prevTime).Seconds()
	g.prevTime = now

	g.ctx.CurrentScene.update(timeDelta)
	return nil
}

func (g *gameRunner) Draw(screen *ebiten.Image) {
	g.ctx.Draw(screen)
}

func (g *gameRunner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.ctx.WindowWidth), int(g.ctx.WindowHeight)
}
