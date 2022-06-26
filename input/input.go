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
				if ebiten.IsStandardGamepadLayoutAvailable(id) {
					info.model = gamepadStandard
				} else {
					info.model = guessGamepadModel(modelName)
				}
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
		if h.gamepadInfo().model == gamepadStandard {
			return inpututil.IsStandardGamepadButtonJustPressed(ebiten.GamepadID(h.id), ebiten.StandardGamepadButton(key.code))
		}
		return inpututil.IsGamepadButtonJustPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(key.code))
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
		if h.gamepadInfo().model == gamepadStandard {
			return ebiten.IsStandardGamepadButtonPressed(ebiten.GamepadID(h.id), ebiten.StandardGamepadButton(key.code))
		}
		return ebiten.IsGamepadButtonPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(key.code))
	case keyMouse:
		return ebiten.IsMouseButtonPressed(ebiten.MouseButton(key.code))
	default:
		return ebiten.IsKeyPressed(ebiten.Key(key.code))
	}
}

func (h *Handler) gamepadInfo() *gamepadInfo {
	return &h.sys.gamepadInfo[h.id]
}

func (h *Handler) mappedGamepadKey(keyCode int) ebiten.GamepadButton {
	b := ebiten.StandardGamepadButton(keyCode)
	switch h.gamepadInfo().model {
	case gamepadMicront:
		return microntToXbox(b)
	default:
		return ebiten.GamepadButton(keyCode)
	}
}

type gamepadModel int

const (
	gamepadUnknown gamepadModel = iota
	gamepadStandard
	gamepadMicront
)

func guessGamepadModel(s string) gamepadModel {
	s = strings.ToLower(s)
	if s == "micront" {
		return gamepadMicront
	}
	return gamepadUnknown
}

type gamepadInfo struct {
	model     gamepadModel
	modelName string
}

func microntToXbox(b ebiten.StandardGamepadButton) ebiten.GamepadButton {
	switch b {
	case ebiten.StandardGamepadButtonLeftTop:
		return ebiten.GamepadButton12
	case ebiten.StandardGamepadButtonLeftRight:
		return ebiten.GamepadButton13
	case ebiten.StandardGamepadButtonLeftBottom:
		return ebiten.GamepadButton14
	case ebiten.StandardGamepadButtonLeftLeft:
		return ebiten.GamepadButton15

	case ebiten.StandardGamepadButtonRightTop:
		return ebiten.GamepadButton0
	case ebiten.StandardGamepadButtonRightRight:
		return ebiten.GamepadButton1
	case ebiten.StandardGamepadButtonRightBottom:
		return ebiten.GamepadButton2
	case ebiten.StandardGamepadButtonRightLeft:
		return ebiten.GamepadButton3

	default:
		return ebiten.GamepadButton(b)
	}
}
