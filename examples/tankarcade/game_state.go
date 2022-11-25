package main

import (
	"github.com/quasilyte/ge/input"
)

type gameSaveState struct {
	LastClearedLevel int
}

type playerConfig struct {
	weaponID      int
	specialID     int
	armorLevel    int
	speedLevel    int
	rotationLevel int
	boostID       int
}

func newPlayerConfig() *playerConfig {
	return &playerConfig{
		armorLevel:    0,
		speedLevel:    0,
		rotationLevel: 0,
	}
}

type gameState struct {
	PlayerInput [3]*input.Handler

	playerConfig [3]*playerConfig

	uiTheme *uiTheme

	gameLevel int

	saveState gameSaveState
}

func newGameState() *gameState {
	state := &gameState{
		saveState: gameSaveState{
			LastClearedLevel: 1,
		},
	}
	state.resetGame()
	return state
}

func (state *gameState) resetGame() {
	for i := range state.playerConfig {
		state.playerConfig[i] = newPlayerConfig()
	}
	state.gameLevel = 1
	state.uiTheme = newUITheme()
}

func (state *gameState) NumGamepadsConnected() int {
	num := 0
	for _, h := range state.PlayerInput {
		if h.GamepadConnected() {
			num++
		}
	}
	return num
}
