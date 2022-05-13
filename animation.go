package ge

type Animation struct {
	sprite *Sprite

	SecondsPerFrame float64

	frame       int
	numFrames   int
	frameTicker float64
	frameWidth  float64
}

func NewAnimation(s *Sprite) *Animation {
	return &Animation{
		sprite:          s,
		frameWidth:      s.FrameWidth,
		numFrames:       int(s.ImageWidth() / s.FrameWidth),
		SecondsPerFrame: 0.05,
	}
}

func (a *Animation) Sprite() *Sprite {
	return a.sprite
}

func (a *Animation) IsDisposed() bool {
	return a.sprite.IsDisposed()
}

func (a *Animation) Dispose() {
	a.sprite.Dispose()
}

func (a *Animation) Tick(delta float64) bool {
	finished := false
	a.frameTicker += delta
	if a.frameTicker > a.SecondsPerFrame {
		a.frameTicker = a.frameTicker - a.SecondsPerFrame
		a.frame++
		if a.frame > a.numFrames {
			a.frame = 0
			a.sprite.FrameOffset.X = 0
			finished = true
		} else {
			a.sprite.FrameOffset.X += a.frameWidth
		}
	}
	return finished
}
