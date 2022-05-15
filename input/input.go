package input

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type System struct {
	gamepadIDs  []ebiten.GamepadID
	gamepadInfo []gamepadInfo
}

func (sys *System) Init() {
	sys.gamepadInfo = make([]gamepadInfo, 8)
}

func (sys *System) Update() {
	sys.gamepadIDs = ebiten.AppendGamepadIDs(sys.gamepadIDs[:0])
	if len(sys.gamepadIDs) != 0 {
		for i, id := range sys.gamepadIDs {
			info := &sys.gamepadInfo[i]
			modelName := ebiten.GamepadName(id)
			if info.modelName != modelName {
				info.modelName = modelName
				info.model = guessGamepadModel(modelName)
			}
		}
	}
}

func (sys *System) NewHandler(playerID int, keymap Keymap) *Handler {
	return &Handler{
		id:     playerID,
		keymap: keymap,
		sys:    sys,
	}
}

type Handler struct {
	id     int
	keymap Keymap
	sys    *System
}

func (h *Handler) ActionIsJustPressed(action Action) bool {
	key, ok := h.keymap.mapping[action]
	if !ok {
		return false
	}
	if key.isGamepad {
		return inpututil.IsGamepadButtonJustPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(ebiten.GamepadButton(key.code)))
	}
	return inpututil.IsKeyJustPressed(ebiten.Key(key.code))
}

func (h *Handler) ActionIsPressed(action Action) bool {
	key, ok := h.keymap.mapping[action]
	if !ok {
		return false
	}
	if key.isGamepad {
		return ebiten.IsGamepadButtonPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(ebiten.GamepadButton(key.code)))
	}
	return ebiten.IsKeyPressed(ebiten.Key(key.code))
}

func (h *Handler) mappedGamepadKey(b ebiten.GamepadButton) ebiten.GamepadButton {
	model := h.sys.gamepadInfo[h.id].model
	switch model {
	case gamepadXbox:
		return b
	case gamepadMicront:
		return microntToXbox(b)
	default:
		return b
	}
}

type gamepadModel int

const (
	gamepadUnknown gamepadModel = iota
	gamepadXbox
	gamepadMicront
)

func guessGamepadModel(s string) gamepadModel {
	s = strings.ToLower(s)
	if strings.Contains(s, "xinput") {
		return gamepadXbox
	}
	if s == "micront" {
		return gamepadMicront
	}
	return gamepadUnknown
}

type gamepadInfo struct {
	model     gamepadModel
	modelName string
}

func microntToXbox(b ebiten.GamepadButton) ebiten.GamepadButton {
	switch b {
	case ebiten.GamepadButton7:
		return ebiten.GamepadButton9

	case ebiten.GamepadButton11:
		return ebiten.GamepadButton12
	case ebiten.GamepadButton12:
		return ebiten.GamepadButton13
	case ebiten.GamepadButton13:
		return ebiten.GamepadButton14
	case ebiten.GamepadButton14:
		return ebiten.GamepadButton15

	case ebiten.GamepadButton0:
		return ebiten.GamepadButton2
	case ebiten.GamepadButton2:
		return ebiten.GamepadButton3
	case ebiten.GamepadButton3:
		return ebiten.GamepadButton0

	default:
		return b
	}
}
