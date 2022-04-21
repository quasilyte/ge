package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeymapAction uint32

type Keymap struct {
	mapping map[KeymapAction]ebiten.Key
}

func (m *Keymap) init() {
	m.mapping = make(map[KeymapAction]ebiten.Key)
}

func (m *Keymap) Set(action KeymapAction, key ebiten.Key) {
	m.mapping[action] = key
}

type Input struct {
	Keymap Keymap
}

func (input *Input) init() {
	input.Keymap.init()
}

func (input *Input) ActionIsJustPressed(action KeymapAction) bool {
	key, ok := input.Keymap.mapping[action]
	if !ok {
		return false
	}
	return inpututil.IsKeyJustPressed(key)
}

func (input *Input) ActionIsPressed(action KeymapAction) bool {
	key, ok := input.Keymap.mapping[action]
	if !ok {
		return false
	}
	return ebiten.IsKeyPressed(key)
}
