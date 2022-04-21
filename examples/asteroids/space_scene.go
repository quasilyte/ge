package main

import (
	"fmt"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/gesignal"
)

type spaceSceneController struct {
	ctx           *ge.Context
	player        *playerVessel
	asteroidsLeft int
}

func newSpaceSceneController() *spaceSceneController {
	return &spaceSceneController{
		player: newPlayerVessel(gemath.Vec{X: 400, Y: 300}),
	}
}

func (c *spaceSceneController) Init(scene *ge.Scene) {
	c.ctx = scene.Context()

	// Preload resources.
	// It'll make the gaming experience smoother.
	preloadImages := []string{
		"vessel.png",
		"bullet.png",
		"explosion.png",
		"asteroid.png",
	}
	for _, p := range preloadImages {
		c.ctx.Loader.PreloadImage(p)
	}
	c.ctx.Loader.PreloadWAV("fire.wav")

	// Deploy the initial objects of the scene.
	scene.AddObject(c.player)
	c.player.EventDestroyed.Connect(nil, c.onPlayerDestroyed)
	asteroids := []gemath.Vec{
		{X: 100, Y: 200},
		{X: 123, Y: 560},
		{X: 600, Y: 50},
		{X: 700, Y: 330},
	}
	c.asteroidsLeft = 4
	for _, pos := range asteroids {
		a := newAsteroid(pos, 125, 3)
		c.connectAsteroid(a)
		scene.AddObject(a)
	}
}

func (c *spaceSceneController) Update(delta float64) {}

func (c *spaceSceneController) onPlayerDestroyed(gesignal.Void) {
	fmt.Println("defeat")
}

func (c *spaceSceneController) connectAsteroid(a *asteroid) {
	a.EventDestroyed.Connect(nil, c.onAsteroidDestroyed)
	a.EventShardCreated.Connect(nil, c.onAsteroidShardCreated)
}

func (c *spaceSceneController) onAsteroidShardCreated(shard *asteroid) {
	c.asteroidsLeft++
	c.connectAsteroid(shard)
}

func (c *spaceSceneController) onAsteroidDestroyed(a *asteroid) {
	c.asteroidsLeft--
	if c.asteroidsLeft == 0 {
		fmt.Println("victory")
	}
}
