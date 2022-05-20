package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/resource"
)

type Context struct {
	Loader   *resource.Loader
	Renderer *Renderer

	Input input.System
	Audio resource.AudioSystem

	Rand gemath.Rand

	CurrentScene *RootScene

	OnCriticalError func(err error)

	FullScreen   bool
	WindowTitle  string
	WindowWidth  float64
	WindowHeight float64
}

func NewContext() *Context {
	ctx := &Context{
		WindowTitle: "GE Game",
	}
	ctx.Loader = resource.NewLoader(&ctx.Audio, &ctx.Audio)
	ctx.Renderer = NewRenderer()
	ctx.Rand.SetSeed(0)
	ctx.Audio.Init(ctx.Loader)
	ctx.Input.Init()
	ctx.OnCriticalError = func(err error) {
		panic(err)
	}
	return ctx
}

func (ctx *Context) ChangeScene(name string, controller SceneController) {
	ctx.CurrentScene = ctx.NewRootScene(name, controller)
}

func (ctx *Context) NewRootScene(name string, controller SceneController) *RootScene {
	rootScene := newRootScene()
	rootScene.Name = name
	rootScene.context = ctx
	rootScene.controller = controller

	scene0 := &rootScene.subSceneArray[1]

	controller.Init(scene0)

	return rootScene
}

func (ctx *Context) Draw(screen *ebiten.Image) {
	ctx.Renderer.Draw(screen, &ctx.CurrentScene.graphics)
}

func (ctx *Context) WindowRect() gemath.Rect {
	return gemath.Rect{
		Max: gemath.Vec{X: ctx.WindowWidth, Y: ctx.WindowHeight},
	}
}
