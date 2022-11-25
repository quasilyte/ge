package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/ui"
)

type uiTheme struct {
	mainMenuButton  ui.ButtonStyle
	prebattleButton ui.ButtonStyle
}

func newUITheme() *uiTheme {
	theme := &uiTheme{}

	mainMenuButton := ui.DefaultButtonStyle()
	mainMenuButton.Font = FontMedium
	mainMenuButton.BorderWidth = 4
	mainMenuButton.BackgroundColor.A = 0
	mainMenuButton.BorderColor.G = 0.5
	mainMenuButton.FocusedBackgroundColor = mainMenuButton.BackgroundColor
	mainMenuButton.FocusedBorderColor = mainMenuButton.BorderColor
	mainMenuButton.FocusedTextColor.G = 0.5
	mainMenuButton.DisabledBackgroundColor = mainMenuButton.BackgroundColor
	mainMenuButton.DisabledTextColor = mainMenuButton.TextColor
	mainMenuButton.DisabledTextColor.A = 0.7
	mainMenuButton.DisabledBorderColor = ge.ColorScale{R: 0.6, G: 1, B: 0.8, A: 1}
	mainMenuButton.Width = 480
	mainMenuButton.Height = 80
	theme.mainMenuButton = mainMenuButton

	prebattleButton := mainMenuButton
	prebattleButton.Font = FontSmall
	prebattleButton.Width = 320
	prebattleButton.Height = 64
	theme.prebattleButton = prebattleButton

	return theme
}
