package main

import (
	"fmt"
	"strings"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/input"
)

type unitStatsController struct {
	scene     *ge.Scene
	input     *input.MultiHandler
	gameState *gameState
}

func newUnitStatsController(state *gameState) *unitStatsController {
	return &unitStatsController{gameState: state, input: state.MenuInput}
}

func (c *unitStatsController) Init(scene *ge.Scene) {
	c.scene = scene

	window := scene.Context().WindowRect()

	bg := scene.NewRepeatedSprite(ImageMenuBackground, window.Width(), window.Height())
	bg.Centered = false
	scene.AddGraphics(bg)

	rowHeight := 128.0
	numRows := len(turretDesignListNoBuilder) + 1

	dict := scene.Dict()

	{
		pos := window.Center()
		pos.X -= 464
		pos.Y -= (rowHeight * float64(numRows-1)) / 2
		for _, d := range turretDesignListNoBuilder {
			imageBg := ge.NewRect(112, 112)
			imageBg.FillColorScale.SetRGBA(0x26, 0x2b, 0x44, 255)
			imageBg.Pos.SetBase(pos)
			scene.AddGraphics(imageBg)

			image := scene.NewSprite(d.Image)
			image.Scale = 1.5
			image.Pos.SetBase(pos.Sub(gemath.Vec{X: 8}))
			scene.AddGraphics(image)

			descriptionBg := ge.NewRect(384, 112)
			descriptionBg.Centered = false
			descriptionBg.FillColorScale.SetRGBA(0x26, 0x2b, 0x44, 255)
			descriptionBg.Pos.SetBase(pos.Add(gemath.Vec{X: 72, Y: -112 / 2}))
			scene.AddGraphics(descriptionBg)

			price := newPriceDisplay()
			price.Pos = descriptionBg.Pos.WithOffset(268, 8)
			scene.AddObject(price)
			price.SetPrice(d.Price)

			timeLabel := scene.NewLabel(FontDescription)
			timeLabel.Text = fmt.Sprintf("%.1f %s", d.Production, dict.Get("word.sec"))
			timeLabel.AlignHorizontal = ge.AlignHorizontalCenter
			timeLabel.Pos = descriptionBg.Pos.WithOffset(268, 80)
			timeLabel.Width = 112
			scene.AddGraphics(timeLabel)

			descLabel := scene.NewLabel(FontDescription)
			descLabel.Text = c.turretDescription(d)
			descLabel.Pos = descriptionBg.Pos.WithOffset(12, 12)

			scene.AddGraphics(descLabel)

			pos.Y += rowHeight
		}
	}

	{
		pos := window.Center()
		pos.X += 64
		pos.Y -= (rowHeight * float64(numRows-1)) / 2
		for _, d := range hullDesignList {
			imageBg := ge.NewRect(112, 112)
			imageBg.FillColorScale.SetRGBA(0x26, 0x2b, 0x44, 255)
			imageBg.Pos.SetBase(pos)
			scene.AddGraphics(imageBg)

			image := scene.NewSprite(d.Image)
			image.Scale = 1.5
			image.Pos.SetBase(pos.Sub(gemath.Vec{X: 4}))
			scene.AddGraphics(image)

			descriptionBg := ge.NewRect(384, 112)
			descriptionBg.Centered = false
			descriptionBg.FillColorScale.SetRGBA(0x26, 0x2b, 0x44, 255)
			descriptionBg.Pos.SetBase(pos.Add(gemath.Vec{X: 72, Y: -112 / 2}))
			scene.AddGraphics(descriptionBg)

			price := newPriceDisplay()
			price.Pos = descriptionBg.Pos.WithOffset(268, 8)
			scene.AddObject(price)
			price.SetPrice(d.Price)

			timeLabel := scene.NewLabel(FontDescription)
			timeLabel.Text = fmt.Sprintf("%.1f %s", d.Production, dict.Get("word.sec"))
			timeLabel.AlignHorizontal = ge.AlignHorizontalCenter
			timeLabel.Pos = descriptionBg.Pos.WithOffset(268, 80)
			timeLabel.Width = 112
			scene.AddGraphics(timeLabel)

			descLabel := scene.NewLabel(FontDescription)
			descLabel.Text = c.hullDescription(d)
			descLabel.Pos = descriptionBg.Pos.WithOffset(12, 12)

			scene.AddGraphics(descLabel)

			pos.Y += rowHeight
		}
	}

	{
		buttonPos := ge.MakePos(window.Center())
		buttonPos.Base.Y += (rowHeight * float64(numRows-1)) / 2
		b := newButton("menu.exit", buttonPos)
		b.Focused = true
		scene.AddObject(b)
	}
}

func (c *unitStatsController) hullDescription(d *hullDesign) string {
	var lines []string
	dict := c.scene.Dict()
	lines = append(lines, strings.ToUpper(dict.Get("design.hull."+d.Name)))
	lines = append(lines, fmt.Sprintf("HP: %d", int(d.HP)))
	lines = append(lines, fmt.Sprintf("%s: %d", dict.GetTitleCase("word.speed"), int(d.Speed)))
	lines = append(lines, fmt.Sprintf("%s: %.1f", dict.GetTitleCase("word.rotation"), d.RotationSpeed))
	return strings.Join(lines, "\n")
}

func (c *unitStatsController) turretDescription(d *turretDesign) string {
	var lines []string
	dict := c.scene.Dict()
	lines = append(lines, strings.ToUpper(dict.Get("design.turret."+d.Name)))
	dps := (1 / d.Reload) * d.Damage
	damageKindText := dict.Get("design.damage_" + d.DamageKind.String())
	lines = append(lines, fmt.Sprintf("DPS: %.1f (%s)", dps, damageKindText))
	lines = append(lines, fmt.Sprintf("HP %s: +%d", dict.Get("word.bonus"), int(d.HPBonus)))
	if d.SpeedPenalty != 0 {
		lines = append(lines, fmt.Sprintf("%s: -%d", dict.GetTitleCase("word.speed"), int(d.SpeedPenalty)))
	}
	return strings.Join(lines, "\n")
}

func (c *unitStatsController) Update(delta float64) {
	if c.input.ActionIsJustPressed(ActionExit) || c.input.ActionIsJustPressed(ActionConfirm) || c.input.ActionIsJustPressed(ActionOpenMenu) {
		c.scene.Context().ChangeScene(newMenuController(c.gameState))
	}
}
