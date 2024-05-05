package ge

import (
	"encoding/json"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/internal/gamedata"
	"github.com/quasilyte/ge/langs"
)

type Context struct {
	Loader   *resource.Loader
	Renderer *Renderer

	Input input.System
	Audio AudioSystem

	Dict *langs.Dictionary

	Rand gmath.Rand

	CurrentScene *RootScene
	nextScene    *RootScene

	// If non-nil, this function is used to create a scene controller that will handle the panic.
	// The single arguments holds the occurred panic information.
	// When game panics for whatever reason, instead of crashing, you can assign a
	// recovery controller constructor here.
	// You can just show the error to the user and crash or you may want to recover the game somehow
	// (e.g. run the main menu controller again).
	NewPanicController func(panicInfo *PanicInfo) SceneController

	// updateFn - user defined function is called in every update cycle if not null
	updateFn        func(delta float64)
	OnCriticalError func(err error)

	GameName string

	FullScreen   bool
	WindowTitle  string
	WindowWidth  float64
	WindowHeight float64
	ScreenWidth  float64
	ScreenHeight float64

	layoutWidth  int
	layoutHeight int

	firstController SceneController

	fixedDelta float64

	imageCache imageCache
}

type PanicInfo struct {
	// A controller that was active during the panic.
	Controller SceneController

	// The error trace.
	Trace string

	// A value retrieved from recover().
	Value any
}

type TimeDeltaMode int

const (
	TimeDeltaComputed60 TimeDeltaMode = iota
	TimeDeltaComputed120
	TimeDeltaFixed60
	TimeDeltaFixed120
)

type ContextConfig struct {
	Mute bool

	TimeDeltaMode TimeDeltaMode
}

func NewContext(config ContextConfig) *Context {
	ctx := &Context{
		WindowTitle: "GE Game",
	}
	switch config.TimeDeltaMode {
	case TimeDeltaComputed60:
		// Nothing to do.
	case TimeDeltaComputed120:
		ebiten.SetTPS(120)
	case TimeDeltaFixed60:
		ctx.fixedDelta = 1.0 / 60.0
	case TimeDeltaFixed120:
		ctx.fixedDelta = 1.0 / 120.0
		ebiten.SetTPS(120)
	}
	audioContext := audio.NewContext(44100)
	ctx.Loader = resource.NewLoader(audioContext)
	if config.Mute {
		ctx.Audio.muted = true
	} else {
		ctx.Audio.init(audioContext, ctx.Loader)
	}
	ctx.Renderer = NewRenderer()
	ctx.Rand.SetSeed(0)
	// TODO: some platforms don't need touches
	ctx.Input.Init(input.SystemConfig{DevicesEnabled: input.AnyDevice})
	ctx.OnCriticalError = func(err error) {
		panic(err)
	}
	ctx.imageCache.Init()
	return ctx
}

func (ctx *Context) ChangeScene(controller SceneController) {
	ctx.nextScene = ctx.newRootScene(controller)

}

func (ctx *Context) newRootScene(controller SceneController) *RootScene {
	rootScene := newRootScene()
	rootScene.context = ctx
	rootScene.controller = controller

	return rootScene
}

func (ctx *Context) Draw(screen *ebiten.Image) {
	ctx.Renderer.Draw(screen, &ctx.CurrentScene.graphics)
}

func (ctx *Context) WindowRect() gmath.Rect {
	return gmath.Rect{
		Max: gmath.Vec{X: ctx.WindowWidth, Y: ctx.WindowHeight},
	}
}

func (ctx *Context) LocateGameData(key string) string {
	return gamedata.Locate(ctx.GameName, key)
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

func (ctx *Context) CheckGameData(key string) bool {
	exists, err := gamedata.Exists(ctx.GameName, key)
	return exists && err == nil
}

func (ctx *Context) ReadGameData(key string) ([]byte, error) {
	if ctx.GameName == "" {
		panic("can't load game data with empty Context.GameName")
	}
	exists, err := gamedata.Exists(ctx.GameName, key)
	if err != nil {
		return nil, fmt.Errorf("touch data with key %q: %w", key, err)
	}
	if !exists {
		return nil, nil
	}
	data, err := gamedata.Load(ctx.GameName, key)
	if err != nil {
		return nil, fmt.Errorf("load game data with key %q: %w", key, err)
	}
	return data, nil
}

func (ctx *Context) LoadGameData(key string, dst any) error {
	data, err := ctx.ReadGameData(key)
	if err != nil {
		return err
	}
	if data == nil {
		ctx.SaveGameData(key, dst)
		return nil
	}
	return json.Unmarshal(data, dst)
}

func (ctx *Context) InferDisplayRatio() (int, int) {
	return inferDisplayRatio(ctx.layoutWidth, ctx.layoutHeight)
}

func (ctx *Context) LayoutSize() (int, int) {
	return ctx.layoutWidth, ctx.layoutHeight
}

func (ctx *Context) WithUpdateFn(fn func(delta float64)) *Context {
	ctx.updateFn = fn
	return ctx
}
