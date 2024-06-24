package ui

import (
	"fmt"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/xslices"
)

type ValueLabel[T comparable] struct {
	Visible bool

	Pos ge.Pos

	text string

	style ValueLabelStyle

	eventDisposed gesignal.Event[*ValueLabel[T]]

	value        *T
	currentValue T

	label *ge.Label
	rect  *ge.Rect

	root *Root
}

type ValueLabelStyle struct {
	Width  float64
	Height float64

	BorderWidth float64

	Font resource.FontID

	BorderColor     ge.ColorScale
	BackgroundColor ge.ColorScale
	TextColor       ge.ColorScale
}

func DefaultValueLabelStyle() ValueLabelStyle {
	return ValueLabelStyle{
		BorderWidth: 1,
		Width:       256,
		Height:      64,

		BorderColor:     transparentColor,
		BackgroundColor: transparentColor,
		TextColor:       grayColor,
	}
}

func NewValueLabel[T comparable](r *Root, style ValueLabelStyle) *ValueLabel[T] {
	e := &ValueLabel[T]{
		style:   style,
		Visible: true,
		root:    r,
	}
	e.eventDisposed.Connect(r, func(b *ValueLabel[T]) {
		if r.disposed {
			return
		}
		r.elems = xslices.RemoveIf(r.elems, func(e uiElement) bool {
			return b == e
		})
	})
	r.elems = append(r.elems, e)
	return e
}

func (l *ValueLabel[T]) Init(scene *ge.Scene) {
	if l.style.BorderWidth != 0 && l.style.BorderColor.A != 0 {
		l.rect = ge.NewRect(l.root.ctx, l.style.Width, l.style.Height)
		l.rect.FillColorScale = l.style.BackgroundColor
		l.rect.OutlineColorScale = l.style.BorderColor
		l.rect.OutlineWidth = l.style.BorderWidth
		l.rect.Pos = l.Pos
		l.rect.Centered = false
		scene.AddGraphics(l.rect)
	}

	l.label = scene.NewLabel(l.style.Font)
	l.label.Width = l.style.Width
	l.label.Height = l.style.Height
	l.label.AlignHorizontal = ge.AlignHorizontalCenter
	l.label.AlignVertical = ge.AlignVerticalCenter
	l.label.Pos = l.Pos
	l.label.SetColorScale(l.style.TextColor)
	scene.AddGraphics(l.label)

	if l.value != nil {
		l.currentValue = *l.value
	}
	l.updateText()
}

func (l *ValueLabel[T]) Dispose() {
	l.eventDisposed.Emit(l)
	l.label.Dispose()
	if l.rect != nil {
		l.rect.Dispose()
	}
}

func (l *ValueLabel[T]) IsDisposed() bool {
	return l.label.IsDisposed()
}

func (l *ValueLabel[T]) Update(delta float64) {
	if l.currentValue != *l.value {
		l.currentValue = *l.value
		l.updateText()
	}
}

func (l *ValueLabel[T]) BindValue(value *T) {
	l.value = value
}

func (l *ValueLabel[T]) SetText(s string) {
	if l.text == s {
		return
	}
	l.text = s
	if l.label != nil {
		l.updateText()
	}
}

func (l *ValueLabel[T]) updateText() {
	l.label.Text = fmt.Sprintf(l.text, l.currentValue)
}
