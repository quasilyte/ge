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
	ImageBackground resource.ID = iota
	ImageBall
	ImageBrickCircle
	ImageBrickRect
	ImageBrickShard
	ImagePlatform
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

	// Bind controls.
	keyBindings := map[ge.KeymapAction]ebiten.Key{
		ActionLeft:  ebiten.KeyA,
		ActionRight: ebiten.KeyD,
		ActionFire:  ebiten.KeySpace,
	}
	for id, key := range keyBindings {
		ctx.Input.Keymap.Set(id, key)
	}

	// Associate audio resources.
	audioResources := map[resource.ID]resource.Audio{
		AudioBrickHit:       {Path: "brick_hit.wav", Volume: -0.3},
		AudioBrickDestroyed: {Path: "brick_destroyed.wav", Volume: -0.1},
		AudioMusic:          {Path: "music.ogg"},
	}
	for id, res := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, res)
	}

	// Associate image resources.
	imageResources := map[resource.ID]resource.Image{
		ImageBackground:  {Path: "background.png"},
		ImageBall:        {Path: "ball.png"},
		ImageBrickCircle: {Path: "brick_circle.png"},
		ImageBrickRect:   {Path: "brick_rect.png"},
		ImageBrickShard:  {Path: "brick_shard.png"},
		ImagePlatform:    {Path: "platform.png"},
	}
	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
	}

	ctx.CurrentScene = ctx.NewScene("game", newGameController())

	if err := ge.RunGame(ctx); err != nil {
		panic(err)
	}
}
