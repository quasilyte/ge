package main

import (
	"embed"
	"io"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"

	_ "image/png"
)

//go:embed assets/*
var gameAssets embed.FS

const (
	ActionLeft ge.KeymapAction = iota
	ActionRight
	ActionForward
	ActionFire
)

func main() {
	ctx := ge.NewContext()
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.WindowTitle = "Asteroids"
	ctx.WindowWidth = 800
	ctx.WindowHeight = 600
	ctx.Loader.OpenAssetFunc = func(path string) io.ReadCloser {
		f, err := gameAssets.Open(filepath.Join("assets", path))
		if err != nil {
			ctx.OnCriticalError(err)
		}
		return f
	}

	ctx.Input.Keymap.Set(ActionForward, ebiten.KeyW)
	ctx.Input.Keymap.Set(ActionLeft, ebiten.KeyA)
	ctx.Input.Keymap.Set(ActionRight, ebiten.KeyD)
	ctx.Input.Keymap.Set(ActionFire, ebiten.KeySpace)

	ctx.CurrentScene = ctx.NewScene("space_scene", newSpaceSceneController())

	if err := ge.RunGame(ctx); err != nil {
		panic(err)
	}
}
