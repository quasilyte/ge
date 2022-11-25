package main

import (
	"fmt"
	"strconv"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/ui"
	"github.com/quasilyte/gmath"
)

type prebattleController struct {
	gameState *gameState
	scene     *ge.Scene

	uiRoots      []*ui.Root
	readyButtons []*ui.Button
	numPlayers   int
	playersReady int

	totalCredits      *ui.ValueLabel[int]
	totalCreditsValue int
	creditsAvailable  int

	allWeapons  []string
	allSpecials []string
	allBoosts   []string
}

func newPrebattleController(s *gameState) *prebattleController {
	return &prebattleController{
		gameState: s,
	}
}

func (c *prebattleController) calculateTotalCost() int {
	armorPrices := [...]int{0, 5, 11, 18, 26, 35}
	rotationPrices := [...]int{0, 3, 8, 15, 24}
	speedPrices := [...]int{0, 4, 10, 18, 28}

	numPlayers := len(c.uiRoots)
	totalCost := 0
	for _, config := range c.gameState.playerConfig[:numPlayers] {
		totalCost += weaponDesignByName(allWeapons[config.weaponID]).cost
		if config.specialID != 0 {
			totalCost += specialWeaponDesignByName(allSpecials[config.specialID]).extra.cost
		}
		totalCost += armorPrices[config.armorLevel]
		totalCost += rotationPrices[config.rotationLevel]
		totalCost += speedPrices[config.speedLevel]
		if config.boostID != 0 {
			totalCost += boostDesigns[allBoosts[config.boostID]].cost
		}
	}

	return totalCost
}

func (c *prebattleController) updateTotalCost() {
	c.totalCreditsValue = c.calculateTotalCost()
	if c.totalCreditsValue > c.creditsAvailable {
		c.playersReady = 0
		for _, b := range c.readyButtons {
			b.SetDisabled(false)
		}
	}
}

func (c *prebattleController) Init(scene *ge.Scene) {
	c.scene = scene

	c.numPlayers = gmath.ClampMin(c.gameState.NumGamepadsConnected(), 1)
	c.uiRoots = make([]*ui.Root, c.numPlayers)

	for i := range c.uiRoots {
		root := ui.NewRoot(scene.Context(), c.gameState.PlayerInput[i])
		root.ActivationAction = ActionPressButton
		root.NextInputAction = ActionMoveDown
		root.PrevInputAction = ActionMoveUp
		c.uiRoots[i] = root
		scene.AddObject(root)
	}

	uiRoot := ui.NewRoot(scene.Context(), c.gameState.PlayerInput[0])

	c.creditsAvailable = 10 + (c.gameState.gameLevel * 10)

	uiWidth := float64((320+256)*c.numPlayers + (64 * (c.numPlayers - 1)))
	offset := gmath.Vec{X: (scene.Context().WindowWidth - uiWidth) / 2}

	levelLabel := scene.NewLabel(FontBig)
	levelLabel.Text = fmt.Sprintf("LEVEL %d", c.gameState.gameLevel)
	levelLabel.AlignHorizontal = ge.AlignHorizontalCenter
	levelLabel.AlignVertical = ge.AlignVerticalCenter
	levelLabel.Width = scene.Context().WindowWidth
	levelLabel.Height = 64
	levelLabel.Pos.Offset = gmath.Vec{Y: 32}
	scene.AddGraphics(levelLabel)

	totalCreditsStyle := ui.DefaultValueLabelStyle()
	totalCreditsStyle.Width = scene.Context().WindowWidth
	totalCreditsStyle.Height = 64
	totalCreditsStyle.Font = FontMedium
	c.totalCredits = ui.NewValueLabel[int](uiRoot, totalCreditsStyle)
	c.totalCredits.Pos.Offset = gmath.Vec{Y: 928 + 32}
	c.totalCredits.SetText("CREDITS USED: %d/" + strconv.Itoa(c.creditsAvailable))
	c.totalCredits.BindValue(&c.totalCreditsValue)
	scene.AddObject(c.totalCredits)

	for i := 0; i < c.numPlayers; i++ {
		offset.Y = 160
		if i != 0 {
			offset.X += 64 + 576
		}

		title := scene.NewLabel(FontMedium)
		title.Text = fmt.Sprintf("PLAYER %d", i+1)
		title.AlignHorizontal = ge.AlignHorizontalCenter
		title.AlignVertical = ge.AlignVerticalCenter
		title.Width = 576
		title.Height = 64
		title.Pos.Offset = offset
		scene.AddGraphics(title)
		offset.Y += 96

		playerConfig := c.gameState.playerConfig[i]

		button1 := c.createOption(i, playerConfig.weaponID, len(allWeapons)-1, offset, "primary", func(b *ui.Button, v int) {
			playerConfig.weaponID = v
			b.Text = allWeapons[playerConfig.weaponID]
		})
		offset.Y += 96

		button2 := c.createOption(i, playerConfig.specialID, len(allSpecials)-1, offset, "special", func(b *ui.Button, v int) {
			playerConfig.specialID = v
			b.Text = allSpecials[playerConfig.specialID]
		})
		offset.Y += 96

		button3 := c.createOption(i, playerConfig.armorLevel, 5, offset, "armor", func(b *ui.Button, v int) {
			playerConfig.armorLevel = v
			b.Text = fmt.Sprintf("LEVEL %d", playerConfig.armorLevel+1)
		})
		offset.Y += 96

		button4 := c.createOption(i, playerConfig.speedLevel, 4, offset, "speed", func(b *ui.Button, v int) {
			playerConfig.speedLevel = v
			b.Text = fmt.Sprintf("LEVEL %d", playerConfig.speedLevel+1)
		})
		offset.Y += 96

		button5 := c.createOption(i, playerConfig.rotationLevel, 4, offset, "rotation", func(b *ui.Button, v int) {
			playerConfig.rotationLevel = v
			b.Text = fmt.Sprintf("LEVEL %d", playerConfig.rotationLevel+1)
		})
		offset.Y += 96

		button6 := c.createOption(i, playerConfig.boostID, len(allBoosts)-1, offset, "boost", func(b *ui.Button, v int) {
			playerConfig.boostID = v
			b.Text = allBoosts[playerConfig.boostID]
		})
		offset.Y += 128

		readyButton := c.uiRoots[i].NewButton(c.gameState.uiTheme.prebattleButton.Resized(576, 64))
		readyButton.Pos.Offset = offset
		readyButton.Text = "READY"
		readyButton.EventActivated.Connect(nil, func(b *ui.Button) {
			if c.totalCreditsValue > c.creditsAvailable {
				return
			}
			b.SetDisabled(true)
			c.playersReady++
			if c.playersReady != c.numPlayers {
				return
			}
			scene.Context().ChangeScene(newBattleController(c.gameState, battleConfig{
				levelData:  getLevelData(c.scene, c.gameState.gameLevel),
				numPlayers: c.numPlayers,
			}))
		})
		scene.AddObject(readyButton)
		c.readyButtons = append(c.readyButtons, readyButton)

		readyButton.SetFocus(true)
		c.uiRoots[i].ConnectInputs(button1, button2)
		c.uiRoots[i].ConnectInputs(button2, button3)
		c.uiRoots[i].ConnectInputs(button3, button4)
		c.uiRoots[i].ConnectInputs(button4, button5)
		c.uiRoots[i].ConnectInputs(button5, button6)
		c.uiRoots[i].ConnectInputs(button6, readyButton)
		c.uiRoots[i].ConnectInputs(readyButton, button1)
	}

	c.updateTotalCost()
}

func (c *prebattleController) Update(delta float64) {
	if c.gameState.PlayerInput[0].ActionIsJustPressed(ActionEscape) {
		c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
	}
}

func (c *prebattleController) createOption(playerID, value, maxValue int, offset gmath.Vec, label string, update func(b *ui.Button, v int)) *ui.Button {
	scene := c.scene
	uiRoot := c.uiRoots[playerID]

	l := scene.NewLabel(FontSmall)
	l.Pos.Offset = offset
	l.Text = label
	l.Width = 256
	l.Height = 64
	l.AlignHorizontal = ge.AlignHorizontalLeft
	l.AlignVertical = ge.AlignVerticalCenter
	scene.AddGraphics(l)

	b := uiRoot.NewButton(c.gameState.uiTheme.prebattleButton)
	b.Pos.Offset = offset.Add(gmath.Vec{X: l.Width})
	var slider gmath.Slider
	slider.SetBounds(0, maxValue)
	slider.TrySetValue(value)
	update(b, value)
	b.EventActivated.Connect(nil, func(b *ui.Button) {
		slider.Inc()
		update(b, slider.Value())
		c.updateTotalCost()
	})
	scene.AddObject(b)

	return b
}
