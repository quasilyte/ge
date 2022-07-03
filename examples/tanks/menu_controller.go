package main

import (
	"os"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/input"
)

type menuController struct {
	input     *input.MultiHandler
	scene     *ge.Scene
	gameState *gameState

	buttons        []*button
	selectedButton gemath.Slider
}

func newMenuController(state *gameState) *menuController {
	return &menuController{gameState: state, input: state.MenuInput}
}

func (c *menuController) Init(scene *ge.Scene) {
	c.scene = scene

	window := scene.Context().WindowRect()

	bg := scene.NewRepeatedSprite(ImageMenuBackground, window.Width(), window.Height())
	bg.Centered = false
	scene.AddGraphics(bg)

	versionLabel := scene.NewLabel(FontBig)
	versionLabel.Text = gameBuildVersion
	versionLabel.Pos = ge.MakePos(gemath.Vec{X: 16, Y: 1080 - 32})
	scene.AddGraphics(versionLabel)

	buttons := []string{
		"menu.new_game",
		"menu.unit_stats",
		"menu.exit",
	}
	c.selectedButton.SetBounds(0, len(buttons)-1)
	buttonHeight := 80.0
	buttonPos := ge.MakePos(window.Center())
	buttonPos.Base.Y -= (buttonHeight * float64(len(buttons)-1)) / 2
	for _, text := range buttons {
		b := newButton(text, buttonPos)
		c.buttons = append(c.buttons, b)
		scene.AddObject(b)
		buttonPos.Offset.Y += buttonHeight
	}
	c.buttons[0].Focused = true
}

func (c *menuController) Update(delta float64) {
	prevSelected := c.selectedButton.Value()
	if c.input.ActionIsJustPressed(ActionPrevCategory) {
		c.selectedButton.Dec()
	}
	if c.input.ActionIsJustPressed(ActionNextCategory) {
		c.selectedButton.Inc()
	}
	if prevSelected != c.selectedButton.Value() {
		c.buttons[prevSelected].Focused = false
		c.buttons[c.selectedButton.Value()].Focused = true
	}

	if c.input.ActionIsJustPressed(ActionConfirm) || c.input.ActionIsJustPressed(ActionOpenMenu) {
		c.onButtonPressed(c.buttons[c.selectedButton.Value()].Text)
	}
}

func (c *menuController) onButtonPressed(op string) {
	switch op {
	case "menu.new_game":
		defaultPlayer1 := pkLocalPlayer1keyboard
		if c.gameState.Player1gamepad.GamepadConnected() {
			defaultPlayer1 = pkLocalPlayer1
		}
		config := battleConfig{
			teamsMode: teams2vs2,
			players: [4]playerKind{
				defaultPlayer1,
				pkEmpty,
				pkBot,
				pkEmpty,
			},
			rules: make(map[string]bool),
		}
		c.scene.Context().ChangeScene("game", newGameController(c.gameState, config))
	case "menu.unit_stats":
		c.scene.Context().ChangeScene("unit stats", newUnitStatsController(c.gameState))
	case "menu.exit":
		os.Exit(0)
	}
}
