package main

import (
	"embed"
	"io"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/resource"

	_ "image/png"
)

//go:embed assets/*
var gameAssets embed.FS

const (
	ActionLeft ge.KeymapAction = iota
	ActionRight
	ActionFire
)

const (
	AudioBrickHit resource.ID = iota
	AudioBrickDestroyed
	AudioMusic
)

func main() {
	ctx := ge.NewContext()
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.WindowTitle = "Breakout"
	ctx.WindowWidth = 800
	ctx.WindowHeight = 640

	ctx.Loader.OpenAssetFunc = func(path string) io.ReadCloser {
		f, err := gameAssets.Open(filepath.Join("assets", path))
		if err != nil {
			ctx.OnCriticalError(err)
		}
		return f
	}

	ctx.Input.Keymap.Set(ActionLeft, ebiten.KeyA)
	ctx.Input.Keymap.Set(ActionRight, ebiten.KeyD)
	ctx.Input.Keymap.Set(ActionFire, ebiten.KeySpace)

	audioResources := map[resource.ID]resource.Audio{
		AudioBrickHit:       {Path: "brick_hit.wav", Volume: -0.3},
		AudioBrickDestroyed: {Path: "brick_destroyed.wav", Volume: -0.1},
		AudioMusic:          {Path: "music.ogg"},
	}
	for id, res := range audioResources {
		ctx.Audio.Registry.Set(id, res)
	}

	ctx.CurrentScene = ctx.NewScene("game", newGameController())

	if err := ge.RunGame(ctx); err != nil {
		panic(err)
	}
}
