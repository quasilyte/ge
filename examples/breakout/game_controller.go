package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/resource"
)

type gameController struct {
	scene              *ge.Scene
	platform           *platform
	lifeSpheres        []*ge.Sprite
	rotatingBricks     brickGroup
	slidingLeftBricks  brickGroup
	slidingRightBricks brickGroup
	numBricks          int
	wave               int
}

func newGameController() *gameController {
	return &gameController{wave: 1}
}

func (c *gameController) Init(scene *ge.Scene) {
	ctx := scene.Context()

	c.scene = scene

	preloadImages := []resource.ID{
		ImageBackground,
		ImageBrickCircle,
		ImageBrickRect,
		ImageBrickShard,
		ImageBall,
		ImagePlatform,
	}
	for _, p := range preloadImages {
		ctx.Loader.PreloadImage(p)
	}
	preloadAudio := []resource.ID{
		AudioBrickDestroyed,
		AudioBrickHit,
		AudioMusic,
	}
	for _, p := range preloadAudio {
		ctx.Loader.PreloadAudio(p)
	}

	ctx.Audio.PlayMusic(AudioMusic)

	{
		bg := scene.LoadSprite(ImageBackground)
		bg.Centered = false
		scene.AddGraphics(bg)
	}

	// Deploy the initial objects of the scene.
	c.platform = newPlatform()
	c.platform.EventBallLost.Connect(nil, c.onBallLost)
	c.platform.body.Pos = gemath.Vec{X: 400, Y: 600}
	scene.AddObject(c.platform)

	for i := 0; i < c.platform.numLives; i++ {
		pos := gemath.Vec{X: 32, Y: float64(i*32) + 64}
		lifeSphere := scene.LoadSprite(ImageBall)
		lifeSphere.Pos = &pos
		c.lifeSpheres = append(c.lifeSpheres, lifeSphere)
		scene.AddGraphics(lifeSphere)
	}

	leftWall := newWall(800, 16, gemath.DegToRad(100))
	leftWall.body.Pos = gemath.Vec{X: 40, Y: 320}
	scene.AddObject(leftWall)
	rightWall := newWall(800, 16, gemath.DegToRad(-100))
	rightWall.body.Pos = gemath.Vec{X: 760, Y: 320}
	scene.AddObject(rightWall)
	topWall := newWall(800, 16, 0)
	topWall.body.Pos = gemath.Vec{X: 400, Y: 6}
	scene.AddObject(topWall)

	c.initLevel()
}

func (c *gameController) Update(delta float64) {
	c.rotatingBricks.Update(delta)
	c.slidingLeftBricks.Update(delta)
	c.slidingRightBricks.Update(delta)
}

func (c *gameController) newCircleBrick(scale float64, pos gemath.Vec) *brick {
	c.numBricks++
	b := newCircleBrick(scale, pos)
	b.EventDestroyed.Connect(b, c.onBrickDestroyed)
	return b
}

func (c *gameController) newBrick(scale float64, rotation gemath.Rad, pos gemath.Vec) *brick {
	c.numBricks++
	b := newBrick(scale, rotation, pos)
	b.EventDestroyed.Connect(b, c.onBrickDestroyed)
	return b
}

func (c *gameController) onBallLost(gesignal.Void) {
	c.lifeSpheres[len(c.lifeSpheres)-1].Dispose()
	c.lifeSpheres = c.lifeSpheres[:len(c.lifeSpheres)-1]
	if c.platform.numLives == 0 {
		c.platform.Dispose()
		c.scene.DelayedCall(2, func() {
			ctx := c.scene.Context()
			ctx.CurrentScene = ctx.NewScene("game", newGameController())
		})
	}
}

func (c *gameController) onBrickDestroyed(gesignal.Void) {
	c.numBricks--
	if c.numBricks != 0 {
		return
	}
	c.wave++
	switch c.wave {
	case 2:
		c.scene.DelayedCall(1, c.initWave2)
	case 3:
		// Victory.
		c.platform.Dispose()
	}
}

func (c *gameController) initLevel() {
	c.rotatingBricks.rotate = 1
	c.slidingLeftBricks.dx = -16
	c.slidingLeftBricks.slideTime = 4
	c.slidingRightBricks.dx = 16
	c.slidingRightBricks.slideTime = 4

	{
		pos := gemath.Vec{X: 544, Y: 64}
		rotation := gemath.Rad(0)
		for i := 0; i < 14; i++ {
			b := c.newBrick(1, rotation, pos)
			rotation += 0.45
			pos = pos.MoveInDirection(80, rotation)
			c.rotatingBricks.bricks = append(c.rotatingBricks.bricks, b)
			c.scene.AddObject(b)
		}
	}

	for i := 0; i < 5; i++ {
		x := float64(144)
		b := c.newBrick(1.2, 0, gemath.Vec{X: x, Y: float64(i*96 + 64)})
		if i%2 == 0 {
			b.body.Pos.X += 96
			c.slidingLeftBricks.bricks = append(c.slidingLeftBricks.bricks, b)
		} else {
			c.slidingRightBricks.bricks = append(c.slidingRightBricks.bricks, b)
		}
		c.scene.AddObject(b)
	}

	c.scene.AddObject(c.newCircleBrick(1, gemath.Vec{X: 512, Y: 240}))
	c.scene.AddObject(c.newCircleBrick(0.8, gemath.Vec{X: 512 - 64, Y: 240 - 64}))
	c.scene.AddObject(c.newCircleBrick(0.8, gemath.Vec{X: 512 + 64, Y: 240 - 64}))
	c.scene.AddObject(c.newCircleBrick(0.8, gemath.Vec{X: 512 - 64, Y: 240 + 64}))
	c.scene.AddObject(c.newCircleBrick(0.8, gemath.Vec{X: 512 + 64, Y: 240 + 64}))
}

func (c *gameController) initWave2() {
	c.rotatingBricks.Reset()
	c.slidingLeftBricks.Reset()
	c.slidingRightBricks.Reset()
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			angle := gemath.Rad(0.2)
			if i%2 == 0 {
				angle = -angle
			}
			pos := gemath.Vec{X: float64(i*96) + 160, Y: float64(j*96) + 64}
			b := c.newBrick(1, angle, pos)
			c.scene.AddObject(b)
			if j%2 != 0 {
				c.rotatingBricks.bricks = append(c.rotatingBricks.bricks, b)
			}
		}
	}
}
