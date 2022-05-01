package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

type gameController struct {
	ctx *ge.Context

	rotatingBricks     brickGroup
	slidingLeftBricks  brickGroup
	slidingRightBricks brickGroup
}

func newGameController() *gameController {
	return &gameController{}
}

func (c *gameController) initLevel(scene *ge.Scene) {
	c.rotatingBricks.rotate = 1
	c.slidingLeftBricks.dx = -16
	c.slidingLeftBricks.slideTime = 4
	c.slidingRightBricks.dx = 16
	c.slidingRightBricks.slideTime = 4

	{
		pos := gemath.Vec{X: 544, Y: 64}
		rotation := gemath.Rad(0)
		for i := 0; i < 14; i++ {
			br := newBrick(1, rotation)
			br.body.Pos = pos
			scene.AddObject(br)
			rotation += 0.45
			pos = pos.MoveInDirection(80, rotation)
			c.rotatingBricks.bricks = append(c.rotatingBricks.bricks, br)
		}
	}
	for i := 0; i < 5; i++ {
		br := newBrick(1.2, 0)
		x := 128 + 16
		if i%2 == 0 {
			x += 96
			c.slidingLeftBricks.bricks = append(c.slidingLeftBricks.bricks, br)
		} else {
			c.slidingRightBricks.bricks = append(c.slidingRightBricks.bricks, br)
		}
		br.body.Pos = gemath.Vec{X: float64(x), Y: float64(i*96 + 64)}
		scene.AddObject(br)
	}

	{
		br := newCircleBrick(1)
		br.body.Pos = gemath.Vec{X: 512, Y: 240}
		scene.AddObject(br)

		br = newCircleBrick(0.8)
		br.body.Pos = gemath.Vec{X: 512 - 64, Y: 240 - 64}
		scene.AddObject(br)
		br = newCircleBrick(0.8)
		br.body.Pos = gemath.Vec{X: 512 + 64, Y: 240 - 64}
		scene.AddObject(br)
		br = newCircleBrick(0.8)
		br.body.Pos = gemath.Vec{X: 512 - 64, Y: 240 + 64}
		scene.AddObject(br)
		br = newCircleBrick(0.8)
		br.body.Pos = gemath.Vec{X: 512 + 64, Y: 240 + 64}
		scene.AddObject(br)
	}
}

func (c *gameController) Init(scene *ge.Scene) {
	c.ctx = scene.Context()

	preloadImages := []string{
		"background.png",
		"brick_purple.png",
		"brick_circle.png",
		"ball.png",
		"platform.png",
	}
	for _, p := range preloadImages {
		c.ctx.Loader.PreloadImage(p)
	}
	preloadAudio := []string{
		"brick_destroyed.wav",
		"brick_hit.wav",
	}
	for _, p := range preloadAudio {
		c.ctx.Loader.PreloadWAV(p)
	}
	c.ctx.Loader.PreloadOGG("music.ogg")

	// c.ctx.Audio.PlayMusic(AudioMusic)

	{
		bg := c.ctx.LoadSprite("background.png")
		bg.Pos = ge.NewVec(800/2, 640/2)
		scene.AddGraphics(bg)
	}

	// Deploy the initial objects of the scene.
	p := newPlatform()
	p.body.Pos = gemath.Vec{X: 400, Y: 600}
	scene.AddObject(p)

	leftWall := newWall(800, 16, gemath.DegToRad(100))
	leftWall.body.Pos = gemath.Vec{X: 40, Y: 320}
	scene.AddObject(leftWall)
	rightWall := newWall(800, 16, gemath.DegToRad(-100))
	rightWall.body.Pos = gemath.Vec{X: 760, Y: 320}
	scene.AddObject(rightWall)
	topWall := newWall(800, 16, 0)
	topWall.body.Pos = gemath.Vec{X: 400, Y: 6}
	scene.AddObject(topWall)

	c.initLevel(scene)
}

func (c *gameController) Update(delta float64) {
	c.rotatingBricks.Update(delta)
	c.slidingLeftBricks.Update(delta)
	c.slidingRightBricks.Update(delta)
}
