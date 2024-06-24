package main

import (
	"fmt"

	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
)

type levelTransitionController struct {
	gameState  *gameState
	scene      *ge.Scene
	numPlayers int
	victory    bool
}

func newLevelTransitionController(s *gameState, numPlayers int, victory bool) *levelTransitionController {
	return &levelTransitionController{gameState: s, numPlayers: numPlayers, victory: victory}
}

func (c *levelTransitionController) Init(scene *ge.Scene) {
	c.scene = scene

	scene.Audio().PlaySound(AudioScreenReset)

	title := scene.NewLabel(FontBig)
	title.Text = "Game Over"
	if c.victory {
		if c.gameState.gameLevel == 0 {
			title.Text = "Stage Cleared!"
		} else if c.gameState.gameLevel == lastLevel {
			title.Text = "You Won!"
		} else {
			title.Text = fmt.Sprintf("Stage %d Clear", c.gameState.gameLevel)
		}
	}
	title.Pos.Offset = scene.Context().WindowRect().Center().Add(gmath.Vec{X: -480, Y: -80})
	title.SetColorScaleRGBA(
		ge.RGB(0xe42cca).R,
		ge.RGB(0xe42cca).G,
		ge.RGB(0xe42cca).B,
		ge.RGB(0xe42cca).A,
	)
	title.Width = 960
	title.AlignHorizontal = ge.AlignHorizontalCenter
	scene.AddGraphics(title)

	h := c.gameState.PlayerInput[0]
	action1 := scene.NewLabel(FontSmall)
	action1.Text = "Press " + formattedActionString(h, ActionConfirm) + " To Restart"
	if c.victory {
		action1.Text = "Press " + formattedActionString(h, ActionConfirm) + " To Continue"
	}
	action1.Pos.Offset = scene.Context().WindowRect().Center().Add(gmath.Vec{X: -320})
	action1.SetColorScaleRGBA(
		ge.RGB(0xe42cca).R,
		ge.RGB(0xe42cca).G,
		ge.RGB(0xe42cca).B,
		ge.RGB(0xe42cca).A,
	)
	action1.Width = 640
	action1.AlignHorizontal = ge.AlignHorizontalCenter
	scene.AddGraphics(action1)

	if !c.victory {
		action2 := scene.NewLabel(FontSmall)
		action2.Text = "Press " + formattedActionString(h, ActionEscape) + " To Exit"
		action2.Pos.Offset = scene.Context().WindowRect().Center().Add(gmath.Vec{X: -320, Y: 40})
		action2.SetColorScaleRGBA(
			ge.RGB(0xe42cca).R,
			ge.RGB(0xe42cca).G,
			ge.RGB(0xe42cca).B,
			ge.RGB(0xe42cca).A,
		)
		action2.Width = 640
		action2.AlignHorizontal = ge.AlignHorizontalCenter
		scene.AddGraphics(action2)
	}
}

func (c *levelTransitionController) Update(delta float64) {
	if c.gameState.PlayerInput[0].ActionIsJustPressed(ActionConfirm) {
		c.onAction1()
		return
	}
	if c.gameState.PlayerInput[0].ActionIsJustPressed(ActionEscape) {
		c.onAction2()
		return
	}
}

func (c *levelTransitionController) onAction1() {
	if c.victory {
		if c.gameState.gameLevel == lastLevel || c.gameState.gameLevel == 0 {
			c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
			return
		}
		c.gameState.gameLevel++
		if c.gameState.gameLevel > c.gameState.saveState.LastClearedLevel {
			c.gameState.saveState.LastClearedLevel = c.gameState.gameLevel
			c.scene.Context().SaveGameData("save", c.gameState.saveState)
		}
		c.scene.Context().ChangeScene(newPrebattleController(c.gameState))
		return
	}
	c.scene.Context().ChangeScene(newBattleController(c.gameState, battleConfig{
		levelData:  getLevelData(c.scene, c.gameState.gameLevel),
		numPlayers: c.numPlayers,
	}))
}

func (c *levelTransitionController) onAction2() {
	if c.gameState.gameLevel == 0 {
		c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
	} else {
		c.scene.Context().ChangeScene(newPrebattleController(c.gameState))
	}
}
