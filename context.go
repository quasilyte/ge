package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
)

type Context struct {
	Loader   *Loader
	Renderer *Renderer

	Input Input
	Audio Audio

	Rand gemath.Rand

	CurrentScene *Scene

	OnCriticalError func(err error)

	WindowTitle  string
	WindowWidth  float64
	WindowHeight float64
}

func NewContext() *Context {
	ctx := &Context{
		WindowTitle: "GE Game",
	}
	ctx.Loader = NewLoader(ctx)
	ctx.Renderer = NewRenderer()
	ctx.Rand.SetSeed(0)
	ctx.Input.init()
	ctx.Audio.init(ctx)
	ctx.OnCriticalError = func(err error) {
		panic(err)
	}
	ctx.Loader.audio = &ctx.Audio
	return ctx
}

func (ctx *Context) NewScene(name string, controller SceneController) *Scene {
	scene := newScene()
	scene.Name = name
	scene.context = ctx

	scene.controller = controller
	controller.Init(scene)
	scene.addQueuedObjects()

	return scene
}

func (ctx *Context) LoadSprite(path string) *Sprite {
	return NewSprite(ctx.Loader.LoadImage(path))
}

func (ctx *Context) Draw(screen *ebiten.Image) {
	ctx.CurrentScene.graphics = ctx.Renderer.Draw(screen, ctx.CurrentScene.graphics)
}
