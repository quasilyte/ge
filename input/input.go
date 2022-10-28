package input

import (
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/quasilyte/ge/gemath"
)

type Action uint32

type Keymap map[Action][]Key

type InputDeviceKind uint8

const (
	KeyboardInput InputDeviceKind = 1 << iota
	GamepadInput
	MouseInput
	TouchInput
)

const AnyInput InputDeviceKind = KeyboardInput | GamepadInput | MouseInput | TouchInput

type System struct {
	gamepadIDs  []ebiten.GamepadID
	gamepadInfo []gamepadInfo

	touchesEnabled bool
	touchIDs       []ebiten.TouchID
	touchTapID     ebiten.TouchID
	touchHasTap    bool
	touchTapPos    gemath.Vec

	cursorPos gemath.Vec
}

func (sys *System) Init(touchesEnabled bool) {
	sys.touchesEnabled = touchesEnabled

	sys.gamepadIDs = make([]ebiten.GamepadID, 0, 8)
	sys.gamepadInfo = make([]gamepadInfo, 8)

	if sys.touchesEnabled {
		sys.touchIDs = make([]ebiten.TouchID, 0, 8)
	}
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

	if sys.touchesEnabled {
		sys.touchHasTap = false
		for _, id := range sys.touchIDs {
			if id == sys.touchTapID && inpututil.IsTouchJustReleased(id) {
				sys.touchHasTap = true
				break
			}
		}
		if !sys.touchHasTap {
			sys.touchIDs = inpututil.AppendJustPressedTouchIDs(sys.touchIDs)
			for _, id := range sys.touchIDs {
				x, y := ebiten.TouchPosition(id)
				sys.touchTapPos = gemath.Vec{X: float64(x), Y: float64(y)}
				sys.touchTapID = id
				break
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
	kind keyKind
	Pos  gemath.Vec
}

func (e EventInfo) IsTouch() bool { return e.kind == keyTouch }

func (h *Handler) GamepadConnected() bool {
	for _, id := range h.sys.gamepadIDs {
		if id == ebiten.GamepadID(h.id) {
			return true
		}
	}
	return false
}

func (h *Handler) TouchEventsEnabled() bool {
	return h.sys.touchesEnabled
}

func (h *Handler) TapPos() (gemath.Vec, bool) {
	return h.sys.touchTapPos, h.sys.touchHasTap
}

func (h *Handler) CursorPos() gemath.Vec {
	return h.sys.cursorPos
}

func (h *Handler) DefaultInputMask() InputDeviceKind {
	if h.GamepadConnected() {
		return GamepadInput
	}
	return KeyboardInput | MouseInput
}

func (h *Handler) ActionKeyNames(action Action, mask InputDeviceKind) []string {
	keys, ok := h.keymap[action]
	if !ok {
		return nil
	}
	gamepadConnected := h.GamepadConnected()
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		enabled := true
		switch k.kind {
		case keyKeyboard:
			enabled = mask&KeyboardInput != 0
		case keyMouse:
			enabled = mask&MouseInput != 0
		case keyGamepad, keyGamepadLeftStick, keyGamepadRightStick:
			enabled = gamepadConnected && (mask&GamepadInput != 0)
		case keyTouch:
			enabled = h.sys.touchesEnabled && (mask&TouchInput != 0)
		}
		if enabled {
			result = append(result, k.name)
		}
	}
	return result
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
		info.kind = k.kind
		switch k.kind {
		case keyMouse:
			info.Pos = h.sys.cursorPos
			return info, true
		case keyTouch:
			info.Pos = h.sys.touchTapPos
			return info, true
		}
		break
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
	case keyTouch:
		if k.code == int(touchTap) {
			return h.sys.touchHasTap
		}
		return false
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
	case keyGamepadLeftStick:
		return h.gamepadStickIsPressed(stickCode(k.code), ebiten.StandardGamepadAxisLeftStickHorizontal, ebiten.StandardGamepadAxisLeftStickVertical)
	case keyGamepadRightStick:
		return h.gamepadStickIsPressed(stickCode(k.code), ebiten.StandardGamepadAxisRightStickHorizontal, ebiten.StandardGamepadAxisRightStickVertical)
	case keyMouse:
		return ebiten.IsMouseButtonPressed(ebiten.MouseButton(k.code))
	default:
		return ebiten.IsKeyPressed(ebiten.Key(k.code))
	}
}

func (h *Handler) gamepadStickIsPressed(code stickCode, axis1, axis2 ebiten.StandardGamepadAxis) bool {
	if h.gamepadInfo().model == gamepadStandard {
		switch stickCode(code) {
		case stickUp:
			vec := h.leftStickVec(axis1, axis2)
			if vec.Len() < 0.5 {
				return false
			}
			angle := vec.Angle().Normalized()
			return angle > (math.Pi+math.Pi/4) && angle <= (2*math.Pi-math.Pi/4)
		case stickRight:
			vec := h.leftStickVec(axis1, axis2)
			if vec.Len() < 0.5 {
				return false
			}
			angle := vec.Angle().Normalized()
			return angle <= (math.Pi/4) || angle > (2*math.Pi-math.Pi/4)
		case stickDown:
			vec := h.leftStickVec(axis1, axis2)
			if vec.Len() < 0.5 {
				return false
			}
			angle := vec.Angle().Normalized()
			return angle > (math.Pi/4) && angle <= (math.Pi-math.Pi/4)
		case stickLeft:
			vec := h.leftStickVec(axis1, axis2)
			if vec.Len() < 0.5 {
				return false
			}
			angle := vec.Angle().Normalized()
			return angle > (math.Pi-math.Pi/4) && angle <= (math.Pi+math.Pi/4)
		}
	}
	return false // TODO
}

func (h *Handler) leftStickVec(axis1, axis2 ebiten.StandardGamepadAxis) gemath.Vec {
	x := ebiten.StandardGamepadAxisValue(ebiten.GamepadID(h.id), axis1)
	y := ebiten.StandardGamepadAxisValue(ebiten.GamepadID(h.id), axis2)
	return gemath.Vec{X: x, Y: y}
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
