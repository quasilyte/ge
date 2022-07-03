package main

import (
	"runtime"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
)

type tutorialController struct {
	input *input.MultiHandler

	scene *ge.Scene

	gameState *gameState
}

func newTutorialController(state *gameState) *tutorialController {
	return &tutorialController{gameState: state, input: state.MenuInput}
}

func (c *tutorialController) Init(scene *ge.Scene) {
	c.scene = scene

	window := scene.Context().WindowRect()

	bg := scene.NewRepeatedSprite(ImageMenuBackground, window.Width(), window.Height())
	bg.Centered = false
	scene.AddGraphics(bg)

	titleLabel := scene.NewLabel(FontBig)
	titleLabel.Pos.Offset.X = window.Center().X
	titleLabel.Pos.Offset.Y = 128
	titleLabel.Text = "AUTOTANKS"
	titleLabel.GrowHorizontal = ge.GrowHorizontalBoth
	titleLabel.AlignHorizontal = ge.AlignHorizontalCenter
	scene.AddGraphics(titleLabel)

	// When running in browser (wasm build), we want the game iframe to become focused.
	// The current workaround is to ask player to make a mouse click somewhere inside the iframe.
	continueText := scene.Dict().Get("tutorial.continue_keyboard")
	if runtime.GOARCH == "wasm" {
		continueText = scene.Dict().Get("tutorial.continue_wasm")
	} else if c.gameState.Player1gamepad.GamepadConnected() {
		continueText = scene.Dict().Get("tutorial.continue_gamepad")
	}
	continueLabel := scene.NewLabel(FontBig)
	continueLabel.Pos.Offset.X = window.Center().X
	continueLabel.Pos.Offset.Y = 196
	continueLabel.Text = continueText
	continueLabel.GrowHorizontal = ge.GrowHorizontalBoth
	continueLabel.AlignHorizontal = ge.AlignHorizontalCenter
	continueLabel.ColorScale = ge.ColorScale{R: 0.8, G: 1, B: 0.8, A: 1}
	scene.AddGraphics(continueLabel)

	contorolsText := scene.Dict().Get("tutorial.keymap.keyboard")
	if c.gameState.Player1gamepad.GamepadConnected() {
		contorolsText = scene.Dict().Get("tutorial.keymap.gamepad")
	}
	controlsLabel := scene.NewLabel(FontDescription)
	controlsLabel.Pos.Offset.X = 196
	controlsLabel.Pos.Offset.Y = 256 + 32
	controlsLabel.Text = contorolsText
	scene.AddGraphics(controlsLabel)

	hintsLabel := scene.NewLabel(FontDescription)
	hintsLabel.Text = scene.Dict().Get("tutorial.hints")
	hintsLabel.Pos.Offset.X = 1024
	hintsLabel.Pos.Offset.Y = 256 + 32
	scene.AddGraphics(hintsLabel)

	if runtime.GOARCH == "wasm" {
		warningLabel := scene.NewLabel(FontDescription)
		warningLabel.Text = scene.Dict().Get("tutorial.wasm_notice")
		warningLabel.Pos.Offset.X = 1024
		warningLabel.Pos.Offset.Y = 256 + 32 + 450
		warningLabel.ColorScale.SetRGBA(255, 180, 180, 255)
		scene.AddGraphics(warningLabel)
	}
}

func (c *tutorialController) Update(delta float64) {
	skipPressed := c.input.ActionIsJustPressed(ActionConfirm) ||
		c.input.ActionIsJustPressed(ActionOpenMenu) ||
		c.input.ActionIsJustPressed(ActionLeftClick)
	if skipPressed {
		c.scene.Context().ChangeScene("menu", newMenuController(c.gameState))
	}
}
