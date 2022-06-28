package main

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
)

type resultsController struct {
	scene     *ge.Scene
	gameState *gameState
	config    battleConfig
	input     *input.MultiHandler
	result    battleResult
}

type battleResult struct {
	alliance int
	players  []playerResult
}

type playerResult struct {
	ID       int
	Alliance int

	Iron          int
	Gold          int
	Oil           int
	UnitsProduced int
}

func newResultsController(state *gameState, config battleConfig, result battleResult) *resultsController {
	return &resultsController{
		gameState: state,
		config:    config,
		input:     state.MenuInput,
		result:    result,
	}
}

func (c *resultsController) Init(scene *ge.Scene) {
	c.scene = scene

	window := scene.Context().WindowRect()

	bg := scene.NewRepeatedSprite(ImageMenuBackground, window.Width(), window.Height())
	bg.Centered = false
	scene.AddGraphics(bg)

	if len(c.result.players) == 0 {
		c.leave()
		return
	}

	offsetY := 384.0

	headerColumns := []string{
		"PLAYER",
		"IRON COLLECTED",
		"GOLD COLLECTED",
		"OIL COLLECTED",
		"UNITS PRODUCED",
		"TEAM",
	}

	bgHeight := float64(len(c.result.players))*48 + 64 + 32
	bgWidth := float64(len(headerColumns)) * 256

	baseOffsetX := 196.0

	descriptionBg := ge.NewRect(bgWidth, bgHeight)
	descriptionBg.Centered = false
	descriptionBg.ColorScale.SetRGBA(0x26, 0x2b, 0x44, 255)
	descriptionBg.Pos.Offset.X = baseOffsetX
	descriptionBg.Pos.Offset.Y = offsetY - 32.0
	scene.AddGraphics(descriptionBg)

	newLabel := func(text string, c color.RGBA, x, y float64) *ge.Label {
		l := scene.NewLabel(FontBig)
		l.Text = text
		l.Pos.Offset.X = x
		l.Pos.Offset.Y = y
		l.AlignHorizontal = ge.AlignHorizontalCenter
		l.Width = 256
		l.ColorScale.SetColor(c)
		return l
	}

	{
		offsetX := baseOffsetX
		for _, col := range headerColumns {
			headerLabel := newLabel(col, ge.RGB(0xffffff), offsetX, offsetY)
			scene.AddGraphics(headerLabel)
			offsetX += 256
		}

		offsetY += 64
	}

	for _, p := range c.result.players {
		textColor := getPlayerTextColor(p.ID)

		teamValue := fmt.Sprintf("%d", p.Alliance+1)
		if p.Alliance == c.result.alliance {
			teamValue = fmt.Sprintf("*%d*", p.Alliance+1)
		}
		columnValues := []string{
			fmt.Sprintf("Player %d", p.ID+1),
			strconv.Itoa(p.Iron),
			strconv.Itoa(p.Gold),
			strconv.Itoa(p.Oil),
			strconv.Itoa(p.UnitsProduced),
			teamValue,
		}

		offsetX := baseOffsetX
		for _, col := range columnValues {
			scene.AddGraphics(newLabel(col, textColor, offsetX, offsetY))
			offsetX += 256
		}

		offsetY += 48
	}

	{
		buttonPos := ge.Pos{}
		buttonPos.Offset.Y += offsetY + 64
		buttonPos.Offset.X = scene.Context().WindowRect().Center().X
		b := newButton("EXIT", buttonPos)
		b.Focused = true
		scene.AddObject(b)
	}
}

func (c *resultsController) Update(delta float64) {
	if c.input.ActionIsJustPressed(ActionExit) || c.input.ActionIsJustPressed(ActionConfirm) || c.input.ActionIsJustPressed(ActionOpenMenu) {
		c.leave()
	}
}

func (c *resultsController) leave() {
	c.scene.Context().ChangeScene("game", newGameController(c.gameState, c.config))
}
