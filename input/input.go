package input

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/quasilyte/ge/gemath"
)

type Action uint32

type Keymap map[Action][]Key

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

func (sys *System) NewMultiHandler() *MultiHandler {
	return &MultiHandler{}
}

type MultiHandler struct {
	list []*Handler
}

func (h *MultiHandler) AddHandler(handler *Handler) {
	h.list = append(h.list, handler)
}

func (h *MultiHandler) ActionIsJustPressed(action Action) bool {
	for i := range h.list {
		if h.list[i].ActionIsJustPressed(action) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) ActionIsPressed(action Action) bool {
	for i := range h.list {
		if h.list[i].ActionIsPressed(action) {
			return true
		}
	}
	return false
}

type Handler struct {
	id     int
	keymap Keymap
	sys    *System
}

type EventInfo struct {
	Pos gemath.Vec
}

func (h *Handler) GamepadConnected() bool {
	for _, id := range h.sys.gamepadIDs {
		if id == ebiten.GamepadID(h.id) {
			return true
		}
	}
	return false
}

func (h *Handler) CursorPos() gemath.Vec {
	return h.sys.cursorPos
}

func (h *Handler) JustPressedActionInfo(action Action) (EventInfo, bool) {
	var info EventInfo
	keys, ok := h.keymap[action]
	if !ok {
		return info, false
	}
	for _, k := range keys {
		if !h.keyIsJustPressed(k) {
			continue
		}
		switch k.kind {
		case keyMouse:
			info.Pos = h.sys.cursorPos
		}
		return info, true
	}
	return info, false
}

func (h *Handler) ActionIsJustPressed(action Action) bool {
	keys, ok := h.keymap[action]
	if !ok {
		return false
	}
	for _, k := range keys {
		if h.keyIsJustPressed(k) {
			return true
		}
	}
	return false
}

func (h *Handler) keyIsJustPressed(k Key) bool {
	switch k.kind {
	case keyGamepad:
		if h.gamepadInfo().model == gamepadStandard {
			return inpututil.IsStandardGamepadButtonJustPressed(ebiten.GamepadID(h.id), ebiten.StandardGamepadButton(k.code))
		}
		return inpututil.IsGamepadButtonJustPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(k.code))
	case keyMouse:
		return inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(k.code))
	default:
		return inpututil.IsKeyJustPressed(ebiten.Key(k.code))
	}
}

func (h *Handler) ActionIsPressed(action Action) bool {
	keys, ok := h.keymap[action]
	if !ok {
		return false
	}
	for _, k := range keys {
		if h.keyIsPressed(k) {
			return true
		}
	}
	return false
}

func (h *Handler) keyIsPressed(k Key) bool {
	switch k.kind {
	case keyGamepad:
		if h.gamepadInfo().model == gamepadStandard {
			return ebiten.IsStandardGamepadButtonPressed(ebiten.GamepadID(h.id), ebiten.StandardGamepadButton(k.code))
		}
		return ebiten.IsGamepadButtonPressed(ebiten.GamepadID(h.id), h.mappedGamepadKey(k.code))
	case keyMouse:
		return ebiten.IsMouseButtonPressed(ebiten.MouseButton(k.code))
	default:
		return ebiten.IsKeyPressed(ebiten.Key(k.code))
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
