package main

import (
	"fmt"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/input"
)

type gameController struct {
	scene *ge.Scene

	input *input.Handler

	gameState *gameState
	config    battleConfig

	focusedRule   gemath.Slider
	checkbox      *checkboxButton
	focusedButton gemath.Slider
	buttons       []focusToggler // All buttons combined
	slotSelectors []*selectButton
	teamsSelector *selectButton
}

func newGameController(state *gameState, config battleConfig) *gameController {
	return &gameController{
		gameState: state,
		input:     state.MainInput,
		config:    config,
	}
}

func (c *gameController) Init(scene *ge.Scene) {
	c.scene = scene

	window := scene.Context().WindowRect()

	bg := scene.NewRepeatedSprite(ImageMenuBackground, window.Width(), window.Height())
	bg.Centered = false
	scene.AddGraphics(bg)

	buttonHeight := 80.0
	numButtons := 8
	c.focusedButton.SetBounds(0, numButtons-1)
	buttonPos := ge.MakePos(window.Center())
	buttonPos.Base.Y -= (buttonHeight * float64(numButtons)) / 2

	{
		c.focusedRule.SetBounds(0, len(battleRules)-1)
		ruleName := battleRules[c.focusedRule.Value()]
		ruleEnabled := c.config.rules[ruleName]
		c.checkbox = newCheckboxButton(ruleName, ruleEnabled, buttonPos)
		scene.AddObject(c.checkbox)
		c.buttons = append(c.buttons, c.checkbox)

		leftArrow := scene.NewSprite(ImageMenuSlideLeft)
		leftArrow.Pos = buttonPos.WithOffset(-190, 0)
		scene.AddGraphics(leftArrow)

		rightArrow := scene.NewSprite(ImageMenuSlideLeft)
		rightArrow.Pos = buttonPos.WithOffset(190, 0)
		rightArrow.FlipHorizontal = true
		scene.AddGraphics(rightArrow)

		buttonPos.Offset.Y += buttonHeight
	}

	{
		for i, pk := range c.config.players {
			l := newLabel(fmt.Sprintf("SLOT %d", i+1), buttonPos.WithOffset(-170, 0))
			scene.AddObject(l)

			b := newSelectButton(playerKindNames, buttonPos.WithOffset(72, 0))
			b.Select(pk.String())
			scene.AddObject(b)
			buttonPos.Offset.Y += buttonHeight

			c.buttons = append(c.buttons, b)
			c.slotSelectors = append(c.slotSelectors, b)
		}
	}
	{
		l := newLabel("TEAMS", buttonPos.WithOffset(-170, 0))
		scene.AddObject(l)

		b := newSelectButton(teamsModeNames, buttonPos.WithOffset(72, 0))
		b.Select(c.config.teamsMode.String())
		scene.AddObject(b)
		buttonPos.Offset.Y += buttonHeight

		c.buttons = append(c.buttons, b)
		c.teamsSelector = b
	}

	{
		b := newButton("START", buttonPos)
		b.Focused = true
		c.focusedButton.TrySetValue(6)
		scene.AddObject(b)
		buttonPos.Offset.Y += buttonHeight
		c.buttons = append(c.buttons, b)
	}
	{
		b := newButton("EXIT", buttonPos)
		scene.AddObject(b)
		buttonPos.Offset.Y += buttonHeight
		c.buttons = append(c.buttons, b)
	}
}

func (c *gameController) Update(delta float64) {
	prevSelected := c.focusedButton.Value()
	if c.input.ActionIsJustPressed(ActionPrevCategory) {
		c.focusedButton.Dec()
	}
	if c.input.ActionIsJustPressed(ActionNextCategory) {
		c.focusedButton.Inc()
	}
	if prevSelected != c.focusedButton.Value() {
		c.buttons[prevSelected].ToggleFocus()
		c.buttons[c.focusedButton.Value()].ToggleFocus()
	}

	if c.input.ActionIsJustPressed(ActionPrevItem) {
		c.onPrevItem()
	} else if c.input.ActionIsJustPressed(ActionNextItem) {
		c.onNextItem()
	}

	if c.input.ActionIsJustPressed(ActionConfirm) || c.input.ActionIsJustPressed(ActionOpenMenu) {
		c.onButtonPressed()
	}

	if c.input.ActionIsJustPressed(ActionExit) {
		c.scene.Context().ChangeScene("menu", newMenuController(c.gameState))
	}
}

func (c *gameController) startBattle() {
	var config battleConfig
	config.rules = c.config.rules

	switch c.teamsSelector.SelectedOption() {
	case "2 VS 2":
		config.teamsMode = teams2vs2
	case "1 VS 3":
		config.teamsMode = teams1vs3
	case "DEATHMATCH":
		config.teamsMode = teamsDeathmatch
	case "VS LEADER":
		config.teamsMode = teamsLeader
	default:
		panic("unexpected option")
	}

	for i, sb := range c.slotSelectors {
		switch sb.SelectedOption() {
		case "EMPTY":
			config.players[i] = pkEmpty
		case "PLAYER 1":
			config.players[i] = pkLocalPlayer1keyboard
		case "PLAYER 1 gamepad":
			config.players[i] = pkLocalPlayer1
		case "PLAYER 2 gamepad":
			config.players[i] = pkLocalPlayer2
		case "PLAYER 3 gamepad":
			config.players[i] = pkLocalPlayer3
		case "PLAYER 4 gamepad":
			config.players[i] = pkLocalPlayer4
		case "EASY BOT":
			config.players[i] = pkEasyBot
		case "BOT":
			config.players[i] = pkBot
		default:
			panic("unexpected option")
		}
	}

	c.scene.Context().ChangeScene("battle", newBattleController(c.gameState, config))
}

func (c *gameController) onButtonPressed() {
	switch b := c.buttons[c.focusedButton.Value()].(type) {
	case *button:
		switch b.Text {
		case "START":
			c.startBattle()
		case "EXIT":
			c.scene.Context().ChangeScene("menu", newMenuController(c.gameState))
		}
	case *checkboxButton:
		c.checkbox.ToggleChecked()
		c.config.rules[c.checkbox.Text] = !c.config.rules[c.checkbox.Text]
	}
}

func (c *gameController) onPrevItem() {
	switch b := c.buttons[c.focusedButton.Value()].(type) {
	case *selectButton:
		b.PrevOption()
	case *checkboxButton:
		c.focusedRule.Dec()
		c.checkbox.Text = battleRules[c.focusedRule.Value()]
		c.checkbox.SetChecked(c.config.rules[c.checkbox.Text])
	}
}

func (c *gameController) onNextItem() {
	switch b := c.buttons[c.focusedButton.Value()].(type) {
	case *selectButton:
		b.NextOption()
	case *checkboxButton:
		c.focusedRule.Inc()
		c.checkbox.Text = battleRules[c.focusedRule.Value()]
		c.checkbox.SetChecked(c.config.rules[c.checkbox.Text])
	}
}
