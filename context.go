package ge

import (
	"encoding/json"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/internal/gamedata"
	"github.com/quasilyte/ge/langs"
	"github.com/quasilyte/ge/resource"
)

type Context struct {
	Loader   *resource.Loader
	Renderer *Renderer

	Input input.System
	Audio resource.AudioSystem

	Dict *langs.Dictionary

	Rand gemath.Rand

	CurrentScene *RootScene

	OnCriticalError func(err error)

	GameName string

	FullScreen   bool
	WindowTitle  string
	WindowWidth  float64
	WindowHeight float64

	firstController SceneController

	imageCache imageCache
}

func NewContext() *Context {
	ctx := &Context{
		WindowTitle: "GE Game",
	}
	ctx.Loader = resource.NewLoader(&ctx.Audio, &ctx.Audio)
	ctx.Renderer = NewRenderer()
	ctx.Rand.SetSeed(0)
	ctx.Audio.Init(ctx.Loader)
	// TODO: some platforms don't need touches
	ctx.Input.Init(input.SystemConfig{TouchesEnabled: true})
	ctx.OnCriticalError = func(err error) {
		panic(err)
	}
	ctx.imageCache.Init()
	return ctx
}

func (ctx *Context) ChangeScene(controller SceneController) {
	ctx.CurrentScene = ctx.NewRootScene(controller)
}

func (ctx *Context) NewRootScene(controller SceneController) *RootScene {
	rootScene := newRootScene()
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

func (ctx *Context) SaveGameData(key string, data any) {
	if ctx.GameName == "" {
		panic("can't save game data with empty Context.GameName")
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("can't save game data with key %q: %v", key, err))
	}
	err = gamedata.Save(ctx.GameName, key, jsonData)
	if err != nil {
		panic(fmt.Sprintf("can't save game data with key %q: %v", key, err))
	}
}

func (ctx *Context) LoadGameData(key string, dst any) {
	if ctx.GameName == "" {
		panic("can't load game data with empty Context.GameName")
	}
	exists, err := gamedata.Exists(ctx.GameName, key)
	if err != nil {
		panic(fmt.Sprintf("can't load game data with key %q: %v", key, err))
	}
	if !exists {
		ctx.SaveGameData(key, dst)
		return
	}
	jsonData, err := gamedata.Load(ctx.GameName, key)
	if err != nil {
		panic(fmt.Sprintf("can't load game data with key %q: %v", key, err))
	}
	err = json.Unmarshal(jsonData, dst)
	if err != nil {
		panic(fmt.Sprintf("can't load game data with key %q: %v", key, err))
	}
}
