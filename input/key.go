package input

import "github.com/hajimehoshi/ebiten/v2"

type keyKind uint8

const (
	keyKeyboard keyKind = iota
	keyGamepad
	keyMouse
)

type Key struct {
	code int
	kind keyKind
}

// Mouse keys.
var (
	KeyMouseLeft   = Key{code: int(ebiten.MouseButtonLeft), kind: keyMouse}
	KeyMouseRight  = Key{code: int(ebiten.MouseButtonRight), kind: keyMouse}
	KeyMouseMiddle = Key{code: int(ebiten.MouseButtonMiddle), kind: keyMouse}
)

// Keyboard keys.
var (
	KeyLeft  Key = Key{code: int(ebiten.KeyLeft)}
	KeyRight Key = Key{code: int(ebiten.KeyRight)}

	KeyA Key = Key{code: int(ebiten.KeyA)}
	KeyW Key = Key{code: int(ebiten.KeyW)}
	KeyS Key = Key{code: int(ebiten.KeyS)}
	KeyD Key = Key{code: int(ebiten.KeyD)}
	KeyE Key = Key{code: int(ebiten.KeyE)}
	KeyR Key = Key{code: int(ebiten.KeyR)}
	KeyQ Key = Key{code: int(ebiten.KeyQ)}

	KeyEscape Key = Key{code: int(ebiten.KeyEscape)}
	KeyEnter  Key = Key{code: int(ebiten.KeyEnter)}

	KeySpace Key = Key{code: int(ebiten.KeySpace)}
)

// Gamepad keys.
var (
	KeyGamepadStart Key = Key{code: int(ebiten.GamepadButton7), kind: keyGamepad}

	KeyGamepadUp    Key = Key{code: int(ebiten.GamepadButton11), kind: keyGamepad}
	KeyGamepadRight Key = Key{code: int(ebiten.GamepadButton12), kind: keyGamepad}
	KeyGamepadDown  Key = Key{code: int(ebiten.GamepadButton13), kind: keyGamepad}
	KeyGamepadLeft  Key = Key{code: int(ebiten.GamepadButton14), kind: keyGamepad}

	KeyGamepadA Key = Key{code: int(ebiten.GamepadButton0), kind: keyGamepad}
	KeyGamepadB Key = Key{code: int(ebiten.GamepadButton1), kind: keyGamepad}
	KeyGamepadX Key = Key{code: int(ebiten.GamepadButton2), kind: keyGamepad}
	KeyGamepadY Key = Key{code: int(ebiten.GamepadButton3), kind: keyGamepad}

	KeyGamepadL1 Key = Key{code: int(ebiten.GamepadButton4), kind: keyGamepad}
	KeyGamepadR1 Key = Key{code: int(ebiten.GamepadButton5), kind: keyGamepad}
)
