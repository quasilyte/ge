package ui

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/gmath"
)

type Image struct {
	Resource       resource.Image
	FlipHorizontal bool
	FlipVertical   bool
}

type ImageButton struct {
	Visible bool

	Pos ge.Pos

	PrevInput inputElement
	NextInput inputElement

	EventActivated gesignal.Event[*ImageButton]
	eventDisposed  gesignal.Event[*ImageButton]

	checkKeyboardInput bool
	disabled           bool
	style              ImageButtonStyle

	sprite *ge.Sprite
	rect   *ge.Rect

	geom gmath.Rect

	root *Root
}

type ImageButtonStyle struct {
	Width  float64
	Height float64

	BorderWidth float64

	BorderColor     ge.ColorScale
	BackgroundColor ge.ColorScale
	ImageColor      ge.ColorScale

	FocusedBorderColor     ge.ColorScale
	FocusedBackgroundColor ge.ColorScale
	FocusedImageColor      ge.ColorScale
}

func DefaultImageButtonStyle() ImageButtonStyle {
	return ImageButtonStyle{
		BorderWidth: 1,
		Width:       256,
		Height:      64,

		BorderColor:     whiteColor,
		BackgroundColor: darkGrayColor,
		ImageColor:      grayColor,

		FocusedBorderColor:     whiteColor,
		FocusedBackgroundColor: darkGrayColor,
		FocusedImageColor:      whiteColor,
	}
}

func (style ImageButtonStyle) Resized(w, h float64) ImageButtonStyle {
	style.Width = w
	style.Height = h
	return style
}

func (b *ImageButton) SetImage(image Image) {
	b.sprite.SetImage(image.Resource)
	dx := (b.style.Width - b.sprite.FrameWidth) / 2
	dy := (b.style.Height - b.sprite.FrameHeight) / 2
	b.sprite.Pos = b.Pos.WithOffset(dx, dy)
	b.sprite.FlipHorizontal = image.FlipHorizontal
	b.sprite.FlipVertical = image.FlipVertical
}

func (b *ImageButton) Init(scene *ge.Scene) {
	pos := b.Pos.Resolve()
	b.geom = gmath.Rect{
		Min: pos,
		Max: pos.Add(gmath.Vec{X: b.style.Width, Y: b.style.Height}),
	}

	b.rect = ge.NewRect(b.root.ctx, b.style.Width, b.style.Height)
	b.rect.FillColorScale = b.style.BackgroundColor
	b.rect.OutlineColorScale = b.style.BorderColor
	b.rect.OutlineWidth = b.style.BorderWidth
	b.rect.Pos = b.Pos
	b.rect.Centered = false
	scene.AddGraphics(b.rect)

	b.sprite = ge.NewSprite(scene.Context())
	b.sprite.Centered = false
	b.onFocusChanged(b.IsFocused())
	scene.AddGraphics(b.sprite)
}

func (b *ImageButton) prevInput() inputElement     { return b.PrevInput }
func (b *ImageButton) nextInput() inputElement     { return b.NextInput }
func (b *ImageButton) setPrevInput(e inputElement) { b.PrevInput = e }
func (b *ImageButton) setNextInput(e inputElement) { b.NextInput = e }

func (b *ImageButton) IsDisabled() bool      { return b.disabled }
func (b *ImageButton) IsFocused() bool       { return b.root.focused == b }
func (b *ImageButton) SetFocus(focused bool) { b.root.setFocus(b, focused) }

func (b *ImageButton) onFocusChanged(focused bool) {
	if focused {
		b.sprite.SetColorScale(b.style.FocusedImageColor)
		b.rect.FillColorScale = b.style.FocusedBackgroundColor
		b.rect.OutlineColorScale = b.style.FocusedBorderColor
	} else {
		b.sprite.SetColorScale(b.style.ImageColor)
		b.rect.FillColorScale = b.style.BackgroundColor
		b.rect.OutlineColorScale = b.style.BorderColor
	}
}

func (b *ImageButton) Update(delta float64) {
	if !b.disabled {
		if b.geom.Contains(b.root.input.CursorPos()) {
			b.root.setFocus(b, true)
		}
		if b.IsFocused() {
			b.checkInput()
		}
	}

	b.sprite.Visible = b.Visible
	b.rect.Visible = b.Visible

	b.checkKeyboardInput = b.IsFocused()
}

func (b *ImageButton) IsDisposed() bool {
	return b.rect.IsDisposed()
}

func (b *ImageButton) Dispose() {
	b.eventDisposed.Emit(b)
	b.sprite.Dispose()
	b.rect.Dispose()
}

func (b *ImageButton) Activate() {
	if !b.disabled {
		b.EventActivated.Emit(b)
	}
}

func (b *ImageButton) checkInput() {
	if b.root.ActivationAction != actionUnset {
		if info, ok := b.root.input.JustPressedActionInfo(b.root.ActivationAction); ok {
			if !b.checkKeyboardInput && !info.IsTouchEvent() {
				return
			}
			if !info.HasPos() || b.geom.Contains(info.Pos) {
				b.Activate()
			}
		}
	}
}
