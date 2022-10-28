package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type keyKind uint8

const (
	keyKeyboard keyKind = iota
	keyGamepad
	keyGamepadLeftStick
	keyGamepadRightStick
	keyMouse
	keyTouch
)

type touchCode int

const (
	touchUnknown touchCode = iota
	touchTap
)

type stickCode int

const (
	stickUnknown stickCode = iota
	stickUp
	stickRight
	stickDown
	stickLeft
)

type Key struct {
	code int
	kind keyKind
	name string
}

// Mouse keys.
var (
	KeyMouseLeft   = Key{code: int(ebiten.MouseButtonLeft), kind: keyMouse, name: "mouse_left_button"}
	KeyMouseRight  = Key{code: int(ebiten.MouseButtonRight), kind: keyMouse, name: "mouse_right_button"}
	KeyMouseMiddle = Key{code: int(ebiten.MouseButtonMiddle), kind: keyMouse, name: "mouse_middle_button"}
)

// Touch keys.
var (
	KeyTouchTap = Key{code: int(touchTap), kind: keyTouch, name: "screen_tap"}
)

// Keyboard keys.
var (
	KeyLeft  Key = Key{code: int(ebiten.KeyLeft), name: "left"}
	KeyRight Key = Key{code: int(ebiten.KeyRight), name: "right"}
	KeyUp    Key = Key{code: int(ebiten.KeyUp), name: "up"}
	KeyDown  Key = Key{code: int(ebiten.KeyDown), name: "down"}

	KeyTab Key = Key{code: int(ebiten.KeyTab), name: "tab"}

	Key1 Key = Key{code: int(ebiten.Key1), name: "1"}
	Key2 Key = Key{code: int(ebiten.Key2), name: "2"}
	Key3 Key = Key{code: int(ebiten.Key3), name: "3"}
	Key4 Key = Key{code: int(ebiten.Key4), name: "4"}
	Key5 Key = Key{code: int(ebiten.Key5), name: "5"}

	KeyA Key = Key{code: int(ebiten.KeyA), name: "a"}
	KeyW Key = Key{code: int(ebiten.KeyW), name: "w"}
	KeyS Key = Key{code: int(ebiten.KeyS), name: "s"}
	KeyD Key = Key{code: int(ebiten.KeyD), name: "d"}
	KeyE Key = Key{code: int(ebiten.KeyE), name: "e"}
	KeyR Key = Key{code: int(ebiten.KeyR), name: "r"}
	KeyT Key = Key{code: int(ebiten.KeyT), name: "t"}
	KeyY Key = Key{code: int(ebiten.KeyY), name: "y"}
	KeyQ Key = Key{code: int(ebiten.KeyQ), name: "q"}

	KeyEscape Key = Key{code: int(ebiten.KeyEscape), name: "escape"}
	KeyEnter  Key = Key{code: int(ebiten.KeyEnter), name: "enter"}

	KeySpace Key = Key{code: int(ebiten.KeySpace), name: "space"}
)

// Gamepad keys.
var (
	KeyGamepadStart  Key = Key{code: int(ebiten.StandardGamepadButtonCenterRight), kind: keyGamepad, name: "gamepad_start"}
	KeyGamepadSelect Key = Key{code: int(ebiten.StandardGamepadButtonCenterLeft), kind: keyGamepad, name: "gamepad_select"}

	KeyGamepadUp    Key = Key{code: int(ebiten.StandardGamepadButtonLeftTop), kind: keyGamepad, name: "gamepad_up"}
	KeyGamepadRight Key = Key{code: int(ebiten.StandardGamepadButtonLeftRight), kind: keyGamepad, name: "gamepad_right"}
	KeyGamepadDown  Key = Key{code: int(ebiten.StandardGamepadButtonLeftBottom), kind: keyGamepad, name: "gamepad_down"}
	KeyGamepadLeft  Key = Key{code: int(ebiten.StandardGamepadButtonLeftLeft), kind: keyGamepad, name: "gamepad_left"}

	KeyGamepadLStickUp    = Key{code: int(stickUp), kind: keyGamepadLeftStick, name: "gamepad_lstick_up"}
	KeyGamepadLStickRight = Key{code: int(stickRight), kind: keyGamepadLeftStick, name: "gamepad_lstick_right"}
	KeyGamepadLStickDown  = Key{code: int(stickDown), kind: keyGamepadLeftStick, name: "gamepad_lstick_down"}
	KeyGamepadLStickLeft  = Key{code: int(stickLeft), kind: keyGamepadLeftStick, name: "gamepad_lstick_left"}
	KeyGamepadRStickUp    = Key{code: int(stickUp), kind: keyGamepadRightStick, name: "gamepad_rstick_up"}
	KeyGamepadRStickRight = Key{code: int(stickRight), kind: keyGamepadRightStick, name: "gamepad_rstick_right"}
	KeyGamepadRStickDown  = Key{code: int(stickDown), kind: keyGamepadRightStick, name: "gamepad_rstick_down"}
	KeyGamepadRStickLeft  = Key{code: int(stickLeft), kind: keyGamepadRightStick, name: "gamepad_rstick_left"}

	KeyGamepadA Key = Key{code: int(ebiten.StandardGamepadButtonRightBottom), kind: keyGamepad, name: "gamepad_a"}
	KeyGamepadB Key = Key{code: int(ebiten.StandardGamepadButtonRightRight), kind: keyGamepad, name: "gamepad_b"}
	KeyGamepadX Key = Key{code: int(ebiten.StandardGamepadButtonRightLeft), kind: keyGamepad, name: "gamepad_x"}
	KeyGamepadY Key = Key{code: int(ebiten.StandardGamepadButtonRightTop), kind: keyGamepad, name: "gamepad_y"}

	KeyGamepadL1 Key = Key{code: int(ebiten.StandardGamepadButtonFrontTopLeft), kind: keyGamepad, name: "gamepad_l1"}
	KeyGamepadR1 Key = Key{code: int(ebiten.StandardGamepadButtonFrontTopRight), kind: keyGamepad, name: "gamepad_r1"}
)
