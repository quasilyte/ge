package input

import "github.com/hajimehoshi/ebiten/v2"

type Key struct {
	code      int
	isGamepad bool
}

// Keyboard keys.
var (
	KeyLeft  Key = Key{code: int(ebiten.KeyLeft)}
	KeyRight Key = Key{code: int(ebiten.KeyRight)}

	KeyA Key = Key{code: int(ebiten.KeyA)}
	KeyW Key = Key{code: int(ebiten.KeyW)}
	KeyS Key = Key{code: int(ebiten.KeyS)}
	KeyD Key = Key{code: int(ebiten.KeyD)}
	KeyE Key = Key{code: int(ebiten.KeyE)}
	KeyQ Key = Key{code: int(ebiten.KeyQ)}

	KeyEscape Key = Key{code: int(ebiten.KeyEscape)}
	KeyEnter  Key = Key{code: int(ebiten.KeyEnter)}

	KeySpace Key = Key{code: int(ebiten.KeySpace)}
)

// Gamepad keys.
var (
	KeyGamepadStart Key = Key{code: int(ebiten.GamepadButton7), isGamepad: true}

	KeyGamepadUp    Key = Key{code: int(ebiten.GamepadButton11), isGamepad: true}
	KeyGamepadRight Key = Key{code: int(ebiten.GamepadButton12), isGamepad: true}
	KeyGamepadDown  Key = Key{code: int(ebiten.GamepadButton13), isGamepad: true}
	KeyGamepadLeft  Key = Key{code: int(ebiten.GamepadButton14), isGamepad: true}

	KeyGamepadA Key = Key{code: int(ebiten.GamepadButton0), isGamepad: true}
	KeyGamepadB Key = Key{code: int(ebiten.GamepadButton1), isGamepad: true}
	KeyGamepadX Key = Key{code: int(ebiten.GamepadButton2), isGamepad: true}
	KeyGamepadY Key = Key{code: int(ebiten.GamepadButton3), isGamepad: true}

	KeyGamepadL1 Key = Key{code: int(ebiten.GamepadButton4), isGamepad: true}
	KeyGamepadR1 Key = Key{code: int(ebiten.GamepadButton5), isGamepad: true}
)
