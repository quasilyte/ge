package main

import "github.com/quasilyte/ge/input"

type MultiHandler struct {
	list []*input.Handler
}

func (h *MultiHandler) AddHandler(handler *input.Handler) {
	h.list = append(h.list, handler)
}

func (h *MultiHandler) ActionIsJustPressed(action input.Action) bool {
	for i := range h.list {
		if h.list[i].ActionIsJustPressed(action) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) ActionIsPressed(action input.Action) bool {
	for i := range h.list {
		if h.list[i].ActionIsPressed(action) {
			return true
		}
	}
	return false
}
