package main

import (
	"strconv"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

type popupBuildTank struct {
	Pos     gemath.Vec
	visible bool

	player *playerData

	cachedResources resourceContainer

	scene *ge.Scene

	priceDisplay *priceDisplay
	priceLabel   *ge.Label
	stashLabel   *ge.Label
	sprite       *ge.Sprite
	hullSprite   *ge.Sprite
	hullName     *ge.Label
	turretSprite *ge.Sprite
	turretName   *ge.Label
	ironStash    *ge.Label
	goldStash    *ge.Label
	oilStash     *ge.Label
	ironIncome   *ge.Label
	goldIncome   *ge.Label
	oilIncome    *ge.Label

	selectingTurret bool

	selectedTurret int
	selectedHull   int
}

func newPopupBuildTank(p *playerData) *popupBuildTank {
	return &popupBuildTank{player: p}
}

func (popup *popupBuildTank) Init(scene *ge.Scene) {
	popup.scene = scene

	popup.sprite = scene.NewSprite(ImagePopupBuildTank)
	popup.sprite.Centered = false
	popup.sprite.Pos.Base = &popup.Pos
	scene.AddGraphics(popup.sprite)

	popup.priceDisplay = newPriceDisplay()
	popup.priceDisplay.Pos.Set(&popup.Pos, 180, 14)
	scene.AddObject(popup.priceDisplay)

	popup.priceLabel = scene.NewLabel(FontSmall)
	popup.priceLabel.Text = scene.Dict().GetTitleCase("word.cost") + ":"
	popup.priceLabel.Pos.Set(&popup.Pos, 180-54, 52)
	scene.AddGraphics(popup.priceLabel)

	popup.stashLabel = scene.NewLabel(FontSmall)
	popup.stashLabel.Text = scene.Dict().GetTitleCase("word.resources_amount") + ":"
	popup.stashLabel.Pos.Set(&popup.Pos, 180-54, 52+24)
	scene.AddGraphics(popup.stashLabel)

	popup.ironStash = scene.NewLabel(FontSmall)
	popup.ironStash.AlignHorizontal = ge.AlignHorizontalCenter
	popup.ironStash.Width = 28
	popup.ironStash.Pos.Set(&popup.Pos, 180, 52+24)
	scene.AddGraphics(popup.ironStash)

	popup.ironIncome = scene.NewLabel(FontSmall)
	popup.ironIncome.AlignHorizontal = ge.AlignHorizontalCenter
	popup.ironIncome.Width = 28
	popup.ironIncome.Pos.Set(&popup.Pos, 180, 52+48)
	scene.AddGraphics(popup.ironIncome)

	popup.goldStash = scene.NewLabel(FontSmall)
	popup.goldStash.AlignHorizontal = ge.AlignHorizontalCenter
	popup.goldStash.Width = 28
	popup.goldStash.Pos.Set(&popup.Pos, 180+40, 52+24)
	scene.AddGraphics(popup.goldStash)

	popup.goldIncome = scene.NewLabel(FontSmall)
	popup.goldIncome.AlignHorizontal = ge.AlignHorizontalCenter
	popup.goldIncome.Width = 28
	popup.goldIncome.Pos.Set(&popup.Pos, 180+40, 52+48)
	scene.AddGraphics(popup.goldIncome)

	popup.oilStash = scene.NewLabel(FontSmall)
	popup.oilStash.AlignHorizontal = ge.AlignHorizontalCenter
	popup.oilStash.Width = 28
	popup.oilStash.Pos.Set(&popup.Pos, 180+80, 52+24)
	scene.AddGraphics(popup.oilStash)

	popup.oilIncome = scene.NewLabel(FontSmall)
	popup.oilIncome.AlignHorizontal = ge.AlignHorizontalCenter
	popup.oilIncome.Width = 28
	popup.oilIncome.Pos.Set(&popup.Pos, 180+80, 52+48)
	scene.AddGraphics(popup.oilIncome)

	{
		hullSprite := ge.NewSprite()
		hullSprite.Pos.Base = &popup.Pos
		hullSprite.Pos.Offset = gemath.Vec{X: 64, Y: 64}
		popup.hullSprite = hullSprite
		applyPlayerColor(popup.player.ID, popup.hullSprite)
		scene.AddGraphics(hullSprite)

		hullName := scene.NewLabel(FontSmall)
		hullName.Pos.Set(&popup.Pos, 16+32, 128)
		popup.hullName = hullName
		scene.AddGraphics(hullName)

		popup.updateHull()
	}
	{
		turretSprite := ge.NewSprite()
		turretSprite.Pos.Base = &popup.Pos
		turretSprite.Pos.Offset = gemath.Vec{X: 64, Y: 64}
		popup.turretSprite = turretSprite
		applyPlayerColor(popup.player.ID, popup.turretSprite)
		scene.AddGraphics(popup.turretSprite)

		turretName := scene.NewLabel(FontSmall)
		turretName.Pos.Set(&popup.Pos, 16+32, 128+24)
		popup.turretName = turretName
		scene.AddGraphics(turretName)

		popup.updateTurret()
	}

	popup.updateCategory()
	popup.updatePrice()
	popup.SetVisibility(false)
}

func (popup *popupBuildTank) IsDisposed() bool {
	return false
}

func (popup *popupBuildTank) Update(delta float64) {
	if !popup.visible {
		return
	}
	if popup.cachedResources != popup.player.Resources {
		popup.updateStashLabels()
	}
}

func (popup *popupBuildTank) updateStashLabels() {
	popup.cachedResources = popup.player.Resources
	popup.ironStash.Text = strconv.Itoa(popup.cachedResources.Iron)
	popup.goldStash.Text = strconv.Itoa(popup.cachedResources.Gold)
	popup.oilStash.Text = strconv.Itoa(popup.cachedResources.Oil)
	popup.ironIncome.Text = "+" + strconv.Itoa(popup.player.Income.Iron)
	popup.goldIncome.Text = "+" + strconv.Itoa(popup.player.Income.Gold)
	popup.oilIncome.Text = "+" + strconv.Itoa(popup.player.Income.Oil)
	popup.updatePrice()
}

func (popup *popupBuildTank) Confirm() tankDesign {
	d := popup.tankDesign()
	if !popup.cachedResources.Contains(d.Price()) {
		return tankDesign{}
	}
	popup.SetVisibility(false)
	return d
}

func (popup *popupBuildTank) SetVisibility(v bool) {
	popup.visible = v
	popup.sprite.Visible = v
	popup.stashLabel.Visible = v
	popup.priceLabel.Visible = v
	popup.priceDisplay.SetVisibility(v)
	popup.ironStash.Visible = v
	popup.goldStash.Visible = v
	popup.oilStash.Visible = v
	popup.ironIncome.Visible = v
	popup.goldIncome.Visible = v
	popup.oilIncome.Visible = v
	popup.hullSprite.Visible = v
	popup.hullName.Visible = v
	popup.turretSprite.Visible = v
	popup.turretName.Visible = v
}

func (popup *popupBuildTank) SelectNextCategory() {
	popup.selectingTurret = !popup.selectingTurret
	popup.updateCategory()
}

func (popup *popupBuildTank) SelectPrevCategory() {
	popup.SelectNextCategory()
}

func (popup *popupBuildTank) SelectNextItem() {
	if popup.selectingTurret {
		popup.selectNextTurret()
	} else {
		popup.selectNextHull()
	}
	popup.updatePrice()
}

func (popup *popupBuildTank) SelectPrevItem() {
	if popup.selectingTurret {
		popup.selectPrevTurret()
	} else {
		popup.selectPrevHull()
	}
	popup.updatePrice()
}

func (popup *popupBuildTank) selectPrevTurret() {
	popup.selectedTurret--
	if popup.selectedTurret < 0 {
		popup.selectedTurret = len(turretDesignList) - 1
	}
	popup.updateTurret()
}

func (popup *popupBuildTank) selectNextTurret() {
	popup.selectedTurret++
	if popup.selectedTurret >= len(turretDesignList) {
		popup.selectedTurret = 0
	}
	popup.updateTurret()
}

func (popup *popupBuildTank) selectPrevHull() {
	popup.selectedHull--
	if popup.selectedHull < 0 {
		popup.selectedHull = len(hullDesignList) - 1
	}
	popup.updateHull()
}

func (popup *popupBuildTank) selectNextHull() {
	popup.selectedHull++
	if popup.selectedHull >= len(hullDesignList) {
		popup.selectedHull = 0
	}
	popup.updateHull()
}

func (popup *popupBuildTank) updatePrice() {
	total := popup.tankDesign().Price()
	popup.priceDisplay.SetPrice(total)
	popup.priceDisplay.SetAvailable(
		total.Iron <= popup.cachedResources.Iron,
		total.Gold <= popup.cachedResources.Gold,
		total.Oil <= popup.cachedResources.Oil,
	)
}

func (popup *popupBuildTank) tankDesign() tankDesign {
	var d tankDesign
	d.Hull = hullDesignList[popup.selectedHull]
	d.Turret = turretDesignList[popup.selectedTurret]
	return d
}

func (popup *popupBuildTank) updateCategory() {
	if popup.selectingTurret {
		popup.turretName.ColorScale = ge.ColorScale{R: 0.8, G: 1, B: 0.8, A: 1}
		popup.hullName.ColorScale = ge.ColorScale{R: 0.6, G: 0.6, B: 0.6, A: 1}
	} else {
		popup.hullName.ColorScale = ge.ColorScale{R: 0.8, G: 1, B: 0.8, A: 1}
		popup.turretName.ColorScale = ge.ColorScale{R: 0.6, G: 0.6, B: 0.6, A: 1}
	}
}

func (popup *popupBuildTank) updateHull() {
	d := hullDesignList[popup.selectedHull]
	popup.hullSprite.SetImage(popup.scene.LoadImage(d.Image))
	popup.hullSprite.Pos.Offset.X = 64 + d.OriginX
	popup.hullName.Text = popup.scene.Dict().GetTitleCase("word.hull") + ": " + popup.scene.Dict().Get("design.hull."+d.Name)
}

func (popup *popupBuildTank) updateTurret() {
	d := turretDesignList[popup.selectedTurret]
	popup.turretSprite.SetImage(popup.scene.LoadImage(d.Image))
	popup.turretName.Text = popup.scene.Dict().GetTitleCase("word.turret") + ": " + popup.scene.Dict().Get("design.turret."+d.Name)
}
