package ui

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/gmath"
)

type Button struct {
	Visible bool

	Pos  ge.Pos
	Text string

	PrevInput inputElement
	NextInput inputElement

	EventActivated gesignal.Event[*Button]
	eventDisposed  gesignal.Event[*Button]

	checkKeyboardInput bool
	disabled           bool
	style              ButtonStyle

	label *ge.Label
	rect  *ge.Rect

	geom gmath.Rect

	root *Root
}

type ButtonStyle struct {
	Width  float64
	Height float64

	BorderWidth float64

	Font resource.FontID

	BorderColor     ge.ColorScale
	BackgroundColor ge.ColorScale
	TextColor       ge.ColorScale

	FocusedBorderColor     ge.ColorScale
	FocusedBackgroundColor ge.ColorScale
	FocusedTextColor       ge.ColorScale

	DisabledBorderColor     ge.ColorScale
	DisabledBackgroundColor ge.ColorScale
	DisabledTextColor       ge.ColorScale
}

func DefaultButtonStyle() ButtonStyle {
	return ButtonStyle{
		BorderWidth: 1,
		Width:       256,
		Height:      64,

		BorderColor:     whiteColor,
		BackgroundColor: darkGrayColor,
		TextColor:       grayColor,

		FocusedBorderColor:     whiteColor,
		FocusedBackgroundColor: darkGrayColor,
		FocusedTextColor:       whiteColor,

		DisabledBorderColor:     withAlpha(whiteColor, 0.8),
		DisabledBackgroundColor: withAlpha(darkGrayColor, 0.8),
		DisabledTextColor:       withAlpha(grayColor, 0.8),
	}
}

func (style ButtonStyle) Resized(w, h float64) ButtonStyle {
	style.Width = w
	style.Height = h
	return style
}

func (b *Button) Init(scene *ge.Scene) {
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
	b.rect.Visible = b.Visible
	scene.AddGraphics(b.rect)

	b.label = scene.NewLabel(b.style.Font)
	b.label.Width = b.style.Width
	b.label.Height = b.style.Height
	b.label.AlignHorizontal = ge.AlignHorizontalCenter
	b.label.AlignVertical = ge.AlignVerticalCenter
	b.label.Pos = b.Pos
	b.label.SetColorScale(b.style.TextColor)
	b.label.Visible = b.Visible
	scene.AddGraphics(b.label)
}

func (b *Button) prevInput() inputElement     { return b.PrevInput }
func (b *Button) nextInput() inputElement     { return b.NextInput }
func (b *Button) setPrevInput(e inputElement) { b.PrevInput = e }
func (b *Button) setNextInput(e inputElement) { b.NextInput = e }

func (b *Button) IsDisabled() bool      { return b.disabled }
func (b *Button) IsFocused() bool       { return b.root.focused == b }
func (b *Button) SetFocus(focused bool) { b.root.setFocus(b, focused) }

func (b *Button) SetDisabled(disabled bool) {
	if b.disabled == disabled {
		return
	}
	b.disabled = disabled
	if b.disabled {
		b.SetFocus(false)
		b.label.SetColorScale(b.style.DisabledTextColor)
		b.rect.FillColorScale = b.style.DisabledBackgroundColor
		b.rect.OutlineColorScale = b.style.DisabledBorderColor
	} else {
		b.label.SetColorScale(b.style.TextColor)
		b.rect.FillColorScale = b.style.BackgroundColor
		b.rect.OutlineColorScale = b.style.BorderColor
	}
}

func (b *Button) onFocusChanged(focused bool) {
	if focused {
		b.label.SetColorScale(b.style.FocusedTextColor)
		b.rect.FillColorScale = b.style.FocusedBackgroundColor
		b.rect.OutlineColorScale = b.style.FocusedBorderColor
	} else {
		b.label.SetColorScale(b.style.TextColor)
		b.rect.FillColorScale = b.style.BackgroundColor
		b.rect.OutlineColorScale = b.style.BorderColor
	}
}

func (b *Button) Update(delta float64) {
	if !b.disabled {
		if b.geom.Contains(b.root.input.CursorPos()) {
			b.root.setFocus(b, true)
		}
		b.checkInput()
	}

	b.label.Text = b.Text
	b.label.Visible = b.Visible
	b.rect.Visible = b.Visible

	b.checkKeyboardInput = b.IsFocused()
}

func (b *Button) IsDisposed() bool {
	return b.rect.IsDisposed()
}

func (b *Button) Dispose() {
	b.eventDisposed.Emit(b)
	b.label.Dispose()
	b.rect.Dispose()
}

func (b *Button) Activate() {
	if !b.disabled {
		b.EventActivated.Emit(b)
	}
}

func (b *Button) checkInput() {
	if b.root.ActivationAction == actionUnset {
		return
	}
	if info, ok := b.root.input.JustPressedActionInfo(b.root.ActivationAction); ok {
		if !b.checkKeyboardInput && !info.IsTouchEvent() {
			return
		}
		if !info.HasPos() || b.geom.Contains(info.Pos) {
			b.Activate()
		}
	}
}
