package input

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/quasilyte/ge/gemath"
)

type System struct {
	gamepadIDs  []ebiten.GamepadID
	gamepadInfo []gamepadInfo

	cursorPos gemath.Vec
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

	{
		x, y := ebiten.CursorPosition()
		sys.cursorPos = gemath.Vec{X: float64(x), Y: float64(y)}
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

type EventInfo struct {
	Pos gemath.Vec
}

func (h *Handler) CursorPos() gemath.Vec {
	return h.sys.cursorPos
}

func (h *Handler) EventInfo(action Action) EventInfo {
	var info EventInfo
	key, ok := h.keymap.mapping[action]
	if !ok {
		return info
	}
	switch key.kind {
	case keyMouse:
		info.Pos = h.sys.cursorPos
	}
	return info
}

func (h *Handler) ActionIsJustPressed(action Action) bool {
	key, ok := h.keymap.mapping[action]
	if !ok {
		return false
	}
	switch key.kind {
	case keyGamepad:
		return inpututil.IsGamepadButtonJustPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(ebiten.GamepadButton(key.code)))
	case keyMouse:
		return inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(key.code))
	default:
		return inpututil.IsKeyJustPressed(ebiten.Key(key.code))
	}
}

func (h *Handler) ActionIsPressed(action Action) bool {
	key, ok := h.keymap.mapping[action]
	if !ok {
		return false
	}
	switch key.kind {
	case keyGamepad:
		return ebiten.IsGamepadButtonPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(ebiten.GamepadButton(key.code)))
	case keyMouse:
		return ebiten.IsMouseButtonPressed(ebiten.MouseButton(key.code))
	default:
		return ebiten.IsKeyPressed(ebiten.Key(key.code))
	}
}

func (h *Handler) mappedGamepadKey(b ebiten.GamepadButton) ebiten.GamepadButton {
	model := h.sys.gamepadInfo[h.id].model
	switch model {
	case gamepadXbox:
		return b
	case gamepadDualSense:
		return dualsenseToXbox(b)
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
	gamepadDualSense
	gamepadMicront
)

func guessGamepadModel(s string) gamepadModel {
	s = strings.ToLower(s)
	if strings.Contains(s, "xinput") {
		return gamepadXbox
	}
	if strings.Contains(s, "dualsense") {
		return gamepadDualSense
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

func dualsenseToXbox(b ebiten.GamepadButton) ebiten.GamepadButton {
	switch b {
	case ebiten.GamepadButton11:
		return ebiten.GamepadButton13
	case ebiten.GamepadButton12:
		return ebiten.GamepadButton14
	case ebiten.GamepadButton13:
		return ebiten.GamepadButton15
	case ebiten.GamepadButton14:
		return ebiten.GamepadButton16

	case ebiten.GamepadButton7:
		return ebiten.GamepadButton10

	default:
		return b
	}
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
