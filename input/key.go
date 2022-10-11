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
	KeyUp    Key = Key{code: int(ebiten.KeyUp)}
	KeyDown  Key = Key{code: int(ebiten.KeyDown)}

	KeyTab Key = Key{code: int(ebiten.KeyTab)}

	KeyA Key = Key{code: int(ebiten.KeyA)}
	KeyW Key = Key{code: int(ebiten.KeyW)}
	KeyS Key = Key{code: int(ebiten.KeyS)}
	KeyD Key = Key{code: int(ebiten.KeyD)}
	KeyE Key = Key{code: int(ebiten.KeyE)}
	KeyR Key = Key{code: int(ebiten.KeyR)}
	KeyT Key = Key{code: int(ebiten.KeyT)}
	KeyY Key = Key{code: int(ebiten.KeyY)}
	KeyQ Key = Key{code: int(ebiten.KeyQ)}

	KeyEscape Key = Key{code: int(ebiten.KeyEscape)}
	KeyEnter  Key = Key{code: int(ebiten.KeyEnter)}

	KeySpace Key = Key{code: int(ebiten.KeySpace)}
)

// Gamepad keys.
var (
	KeyGamepadStart  Key = Key{code: int(ebiten.StandardGamepadButtonCenterRight), kind: keyGamepad}
	KeyGamepadSelect Key = Key{code: int(ebiten.StandardGamepadButtonCenterLeft), kind: keyGamepad}

	KeyGamepadUp    Key = Key{code: int(ebiten.StandardGamepadButtonLeftTop), kind: keyGamepad}
	KeyGamepadRight Key = Key{code: int(ebiten.StandardGamepadButtonLeftRight), kind: keyGamepad}
	KeyGamepadDown  Key = Key{code: int(ebiten.StandardGamepadButtonLeftBottom), kind: keyGamepad}
	KeyGamepadLeft  Key = Key{code: int(ebiten.StandardGamepadButtonLeftLeft), kind: keyGamepad}

	KeyGamepadA Key = Key{code: int(ebiten.StandardGamepadButtonRightBottom), kind: keyGamepad}
	KeyGamepadB Key = Key{code: int(ebiten.StandardGamepadButtonRightRight), kind: keyGamepad}
	KeyGamepadX Key = Key{code: int(ebiten.StandardGamepadButtonRightLeft), kind: keyGamepad}
	KeyGamepadY Key = Key{code: int(ebiten.StandardGamepadButtonRightTop), kind: keyGamepad}

	KeyGamepadL1 Key = Key{code: int(ebiten.StandardGamepadButtonFrontTopLeft), kind: keyGamepad}
	KeyGamepadR1 Key = Key{code: int(ebiten.StandardGamepadButtonFrontTopRight), kind: keyGamepad}
)
