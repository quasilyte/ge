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
	
	if int(ctx.ScreenWidth) == 0 && int(ctx.ScreenHeight) == 0 {
		ctx.ScreenWidth = ctx.WindowWidth
		ctx.ScreenHeight = ctx.WindowHeight
	}
	
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

	var delta float64
	if g.ctx.fixedDelta {
		delta = 1.0 / 60.0
	} else {
		now := time.Now()
		delta = now.Sub(g.prevTime).Seconds()
		g.prevTime = now
	}

	g.ctx.CurrentScene.update(delta)
	return nil
}

func (g *gameRunner) Draw(screen *ebiten.Image) {
	g.ctx.Draw(screen)
}

func (g *gameRunner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.ctx.ScreenWidth), int(g.ctx.ScreenHeight)
}
