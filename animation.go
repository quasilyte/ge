package ge

import (
	"math"

	"github.com/quasilyte/ge/gesignal"
)

// TODO: change animation speed to FPS instead of seconds per frame.

type Animation struct {
	sprite     *Sprite
	frameWidth float64

	numFrames int
	offsetY   float64

	animationSpan float64
	deltaPerFrame float64

	repeated bool

	frame       int
	frameTicker float64

	EventFrameChanged gesignal.Event[int]
}

func NewRepeatedAnimation(s *Sprite, numFrames int) *Animation {
	a := NewAnimation(s, numFrames)
	a.repeated = true
	return a
}

func NewAnimation(s *Sprite, numFrames int) *Animation {
	a := &Animation{}
	a.SetSprite(s, numFrames)
	a.SetSecondsPerFrame(0.05)
	return a
}

func (a *Animation) SetSprite(s *Sprite, numFrames int) {
	a.sprite = s
	if numFrames < 0 {
		numFrames = int(s.ImageWidth() / s.FrameWidth)
	}
	a.frameWidth = s.FrameWidth
	a.numFrames = numFrames
	a.SetAnimationSpan(a.animationSpan)
}

func (a *Animation) SetOffsetY(offset float64) {
	a.offsetY = offset
}

func (a *Animation) SetAnimationSpan(value float64) {
	a.animationSpan = value
	a.deltaPerFrame = value / float64(a.numFrames)
}

func (a *Animation) SetSecondsPerFrame(seconds float64) {
	a.animationSpan = seconds * float64(a.numFrames)
	a.deltaPerFrame = seconds
}

func (a *Animation) Sprite() *Sprite {
	return a.sprite
}

func (a *Animation) IsDisposed() bool {
	return a.sprite.IsDisposed()
}

func (a *Animation) Rewind() {
	a.frameTicker = 0
	a.frame = 0
}

func (a *Animation) RewindTo(value float64) {
	a.frameTicker = 0
	a.frame = -1
	a.Tick(value)
}

func (a *Animation) Tick(delta float64) bool {
	if !a.repeated {
		if a.frameTicker >= a.animationSpan {
			return true
		}
	}

	a.sprite.FrameOffset.Y = a.offsetY

	finished := false
	a.frameTicker += delta
	var frame int
	if a.frameTicker >= a.animationSpan {
		finished = true
		if a.repeated {
			rem := math.Mod(a.frameTicker, a.animationSpan)
			a.frameTicker = rem
			frame = int(a.frameTicker / a.deltaPerFrame)
		} else {
			a.frameTicker = a.animationSpan
			frame = a.numFrames - 1
		}
	} else {
		frame = int(a.frameTicker / a.deltaPerFrame)
	}

	framesDelta := frame - a.frame
	a.frame = frame
	if framesDelta != 0 {
		// A small optimization: don't call Emit if there are no listeners.
		// This is more useful for repeated animations as they're less likely to have
		// any frame event listeners.
		if !a.EventFrameChanged.IsEmpty() {
			a.EventFrameChanged.Emit(framesDelta)
		}
		a.sprite.FrameOffset.X = a.frameWidth * float64(frame)
	}

	return finished
}
