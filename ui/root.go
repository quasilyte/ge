package ui

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/xslices"
)

type Root struct {
	ActivationAction input.Action
	PrevInputAction  input.Action
	NextInputAction  input.Action

	elems []uiElement

	inputElems  []inputElement
	focused     inputElement
	lastFocused inputElement

	ctx   *ge.Context
	input *input.Handler

	disposed bool
}

func NewRoot(ctx *ge.Context, h *input.Handler) *Root {
	return &Root{
		ctx:              ctx,
		input:            h,
		ActivationAction: actionUnset,
		PrevInputAction:  actionUnset,
		NextInputAction:  actionUnset,
	}
}

func (r *Root) Init(scene *ge.Scene) {}

func (r *Root) Update(delta float64) {
	if r.PrevInputAction != actionUnset {
		if r.input.ActionIsJustPressed(r.PrevInputAction) {
			r.FocusPrevInput()
		}
	}
	if r.NextInputAction != actionUnset {
		if r.input.ActionIsJustPressed(r.NextInputAction) {
			r.FocusNextInput()
		}
	}
}

func (r *Root) IsDisposed() bool {
	return r.disposed
}

func (r *Root) Dispose() {
	r.disposed = true
	for _, e := range r.inputElems {
		e.Dispose()
	}
	for _, e := range r.elems {
		e.Dispose()
	}
}

func (r *Root) ConnectInputs(from, to inputElement) {
	from.setNextInput(to)
	to.setPrevInput(from)
}

func (r *Root) FocusInputBefore(e inputElement) {
	current := e
	for {
		current = current.prevInput()
		if current == e {
			break // Loop?
		}
		if current == nil {
			break
		}
		if current.IsDisabled() || current.IsDisposed() {
			continue
		}
		r.setFocus(current, true)
		break
	}
}

func (r *Root) FocusInputAfter(e inputElement) {
	current := e
	for {
		current = current.nextInput()
		if current == e {
			break // Loop?
		}
		if current == nil {
			break
		}
		if current.IsDisabled() || current.IsDisposed() {
			continue
		}
		r.setFocus(current, true)
		break
	}
}

func (r *Root) FocusNextInput() {
	var toFocus inputElement
	if r.focused != nil {
		toFocus = r.focused
	} else if r.lastFocused != nil {
		toFocus = r.lastFocused
	}
	if toFocus == nil {
		return
	}
	r.FocusInputAfter(toFocus)
}

func (r *Root) FocusPrevInput() {
	var toFocus inputElement
	if r.focused != nil {
		toFocus = r.focused
	} else if r.lastFocused != nil {
		toFocus = r.lastFocused
	}
	if toFocus == nil {
		return
	}
	r.FocusInputBefore(toFocus)
}

func (r *Root) NewImageButton(style ImageButtonStyle) *ImageButton {
	e := &ImageButton{
		style:   style,
		Visible: true,
		root:    r,
	}
	e.eventDisposed.Connect(r, func(b *ImageButton) {
		if r.disposed {
			return
		}
		r.inputElems = xslices.RemoveIf(r.inputElems, func(e inputElement) bool {
			return b == e
		})
	})
	r.inputElems = append(r.inputElems, e)
	return e
}

func (r *Root) NewButton(style ButtonStyle) *Button {
	e := &Button{
		style:   style,
		Visible: true,
		root:    r,
	}
	e.eventDisposed.Connect(r, func(b *Button) {
		if r.disposed {
			return
		}
		r.inputElems = xslices.RemoveIf(r.inputElems, func(e inputElement) bool {
			return b == e
		})
	})
	r.inputElems = append(r.inputElems, e)
	return e
}

func (r *Root) setFocus(e inputElement, focused bool) {
	if focused {
		if r.focused == e {
			return // already focused
		}
		prevFocused := r.focused
		r.focused = e
		r.focused.onFocusChanged(true) // focused
		if prevFocused != nil {
			prevFocused.onFocusChanged(false) // unfocused
		}
		return
	}

	if r.focused == e {
		r.lastFocused = r.focused
		r.focused = nil
		e.onFocusChanged(false) // unfocused
	}
}

type inputElement interface {
	uiElement

	IsFocused() bool
	SetFocus(focused bool)
	IsDisabled() bool

	onFocusChanged(focused bool)

	prevInput() inputElement
	nextInput() inputElement
	setPrevInput(inputElement)
	setNextInput(inputElement)
}

type uiElement interface {
	Dispose()
	IsDisposed() bool
}

var actionUnset input.Action = math.MaxUint32
