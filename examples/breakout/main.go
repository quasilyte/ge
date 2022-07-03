package main

import (
	"embed"
	"io"
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/resource"

	_ "image/png"
)

//go:embed assets/*
var gameAssets embed.FS

const (
	ActionLeft input.Action = iota
	ActionRight
	ActionFire
)

const (
	ImageBackground resource.ImageID = iota
	ImageBall
	ImageBrickCircle
	ImageBrickRect
	ImageBrickShard
	ImagePlatform
)

const (
	AudioBrickHit resource.AudioID = iota
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
		f, err := gameAssets.Open("assets/" + path)
		if err != nil {
			ctx.OnCriticalError(err)
		}
		return f
	}

	// Bind controls.
	var keymap input.Keymap
	keymap.Set(ActionLeft, input.KeyA)
	keymap.Set(ActionRight, input.KeyD)
	keymap.Set(ActionFire, input.KeySpace)

	// Associate audio resources.
	audioResources := map[resource.AudioID]resource.Audio{
		AudioBrickHit:       {Path: "brick_hit.wav", Volume: -0.3},
		AudioBrickDestroyed: {Path: "brick_destroyed.wav", Volume: -0.1},
		AudioMusic:          {Path: "music.ogg"},
	}
	for id, res := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, res)
	}

	// Associate image resources.
	imageResources := map[resource.ImageID]resource.ImageInfo{
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

	if err := ge.RunGame(ctx, newGameController(ctx.Input.NewHandler(0, keymap))); err != nil {
		panic(err)
	}
}
