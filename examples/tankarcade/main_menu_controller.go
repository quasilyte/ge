package main

import (
	"fmt"
	"os"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/ui"
	"github.com/quasilyte/gmath"
)

type mainMenuController struct {
	gameState *gameState
}

func newMainMenuController(s *gameState) *mainMenuController {
	return &mainMenuController{gameState: s}
}

func (c *mainMenuController) Init(scene *ge.Scene) {
	scene.Context().LoadGameData("save", &c.gameState.saveState)

	root := ui.NewRoot(scene.Context(), c.gameState.PlayerInput[0])
	root.ActivationAction = ActionPressButton
	root.PrevInputAction = ActionMoveUp
	root.NextInputAction = ActionMoveDown
	scene.AddObject(root)

	windowRect := scene.Context().WindowRect()
	buttonX := windowRect.Width()/2 - c.gameState.uiTheme.mainMenuButton.Width/2

	offsetY := 320.0

	newGameButton := root.NewButton(c.gameState.uiTheme.mainMenuButton)
	newGameButton.Pos.Offset = gmath.Vec{X: buttonX, Y: offsetY}
	newGameButton.Text = "START"
	newGameButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		scene.Context().ChangeScene(newPrebattleController(c.gameState))
	})
	scene.AddObject(newGameButton)
	newGameButton.SetFocus(true)
	offsetY += 120

	var gameLevelSlider gmath.Slider
	gameLevelSlider.SetBounds(1, c.gameState.saveState.LastClearedLevel)
	gameLevelSlider.TrySetValue(c.gameState.saveState.LastClearedLevel)
	loadGameButton := root.NewButton(c.gameState.uiTheme.mainMenuButton)
	loadGameButton.Pos.Offset = gmath.Vec{X: buttonX, Y: offsetY}
	loadGameButton.Text = fmt.Sprintf("LEVEL %d", c.gameState.saveState.LastClearedLevel)
	c.gameState.gameLevel = c.gameState.saveState.LastClearedLevel
	loadGameButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		gameLevelSlider.Inc()
		loadGameButton.Text = fmt.Sprintf("LEVEL %d", gameLevelSlider.Value())
		c.gameState.gameLevel = gameLevelSlider.Value()
	})
	scene.AddObject(loadGameButton)
	offsetY += 120

	tutorialButton := root.NewButton(c.gameState.uiTheme.mainMenuButton)
	tutorialButton.Pos.Offset = gmath.Vec{X: buttonX, Y: offsetY}
	tutorialButton.Text = "TUTORIAL"
	tutorialButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		c.gameState.gameLevel = 0
		scene.Context().ChangeScene(newBattleController(c.gameState, battleConfig{
			levelData:  getLevelData(scene, 0),
			numPlayers: 1,
		}))
	})
	scene.AddObject(tutorialButton)
	offsetY += 120

	exitButton := root.NewButton(c.gameState.uiTheme.mainMenuButton)
	exitButton.Pos.Offset = gmath.Vec{X: buttonX, Y: offsetY}
	exitButton.Text = "EXIT"
	exitButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		os.Exit(0)
	})
	scene.AddObject(exitButton)

	root.ConnectInputs(newGameButton, loadGameButton)
	root.ConnectInputs(loadGameButton, tutorialButton)
	root.ConnectInputs(tutorialButton, exitButton)
	root.ConnectInputs(exitButton, newGameButton)

	buildLabel := scene.NewLabel(FontSmall)
	buildLabel.Text = fmt.Sprintf("build %d", buildVersion)
	buildLabel.Pos.Offset = gmath.Vec{X: 32, Y: 1080 - 64}
	scene.AddGraphics(buildLabel)
}

func (c *mainMenuController) Update(delta float64) {}
