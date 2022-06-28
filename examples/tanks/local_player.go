package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/input"
)

type localPlayer struct {
	*playerData

	input *input.Handler

	sectorSelector *sectorSelector
	scene          *ge.Scene

	popupOpened    bool
	buildTank      *popupBuildTank
	selectedSector *sector
}

func newLocalPlayer(data *playerData, input *input.Handler, startingSector *sector) *localPlayer {
	p := &localPlayer{
		playerData:     data,
		input:          input,
		sectorSelector: newSectorSelector(data.ID),
	}
	p.sectorSelector.SetSector(startingSector)
	return p
}

func (p *localPlayer) Init(scene *ge.Scene) {
	p.scene = scene
	scene.AddObject(p.sectorSelector)

	bp := p.BattleState.NewBattlePost(p.playerData, p.sector().Center(), ionTurretDesign)
	bp.HQ = true
	p.sector().AssignBase(bp)
	scene.AddObject(bp)

	firstBuilder := p.BattleState.NewBattleTank(bp.Player, tankDesign{
		Hull:   hullDesigns["fighter"],
		Turret: turretDesigns["builder"],
	})
	firstBuilder.Body.Pos = p.sector().Center().Add(gemath.Vec{Y: 64})
	firstBuilder.Body.Rotation = firstBuilder.Body.Pos.AngleToPoint(gemath.Vec{X: 1920 / 2, Y: 1080 / 2})
	p.sector().AddTank(firstBuilder)
	scene.AddObject(firstBuilder)

	p.buildTank = newPopupBuildTank(p.playerData)
	scene.AddObjectAbove(p.buildTank, 1)
}

func (p *localPlayer) IsDisposed() bool { return false }

func (p *localPlayer) Update(delta float64) {
	if p.popupOpened {
		p.handlePopupControls()
	} else {
		p.handleControls()
	}
}

func (p *localPlayer) handlePopupControls() {
	if p.sector().Base == nil {
		p.buildTank.SetVisibility(false)
		p.popupOpened = false
		return
	}

	if p.input.ActionIsJustPressed(ActionOpenMenu) || p.input.ActionIsJustPressed(ActionCancel) {
		p.buildTank.SetVisibility(false)
		p.popupOpened = false
	} else if p.input.ActionIsJustPressed(ActionConfirm) {
		d := p.buildTank.Confirm()
		if d.IsEmpty() {
			if p.BattleState.SingleLocalPlayer == p.playerData {
				p.scene.Audio().PlaySound(AudioCueError)
			}
		} else {
			p.Resources.Sub(d.Price())
			p.sector().Base.StartProduction(d)
			p.popupOpened = false
			if p.BattleState.SingleLocalPlayer == p.playerData {
				p.scene.Audio().EnqueueSound(AudioCueProductionStarted)
			}
		}
	} else if p.input.ActionIsJustPressed(ActionNextItem) {
		p.buildTank.SelectNextItem()
	} else if p.input.ActionIsJustPressed(ActionPrevItem) {
		p.buildTank.SelectPrevItem()
	} else if p.input.ActionIsJustPressed(ActionNextCategory) {
		p.buildTank.SelectNextCategory()
	} else if p.input.ActionIsJustPressed(ActionPrevCategory) {
		p.buildTank.SelectPrevCategory()
	}
}

func (p *localPlayer) handleControls() {
	base := p.sector().Base

	if p.input.ActionIsJustPressed(ActionOpenMenu) {
		if base != nil && base.Player.ID == p.ID && !base.IsBusy() {
			p.buildTank.SetVisibility(true)
			p.buildTank.Pos = p.sector().Pos.Add(gemath.Vec{X: 10, Y: 40})
			if p.selectedSector != nil {
				p.selectedSector.UnselectUnits()
				p.selectedSector = nil
			}
			p.popupOpened = true
			return
		}
	}

	if p.BattleState.FortificationsAllowed && p.input.ActionIsJustPressed(ActionFortify) {
		if base != nil && base.Player.ID == p.ID && base.Turret == nil && p.Resources.Contains(gaussTurretDesign.Price) {
			base.InstallTurret(gaussTurretDesign)
			p.Resources.Sub(gaussTurretDesign.Price)
		}
	}

	if p.input.ActionIsJustPressed(ActionCancel) {
		if p.selectedSector != nil {
			p.selectedSector.UnselectUnits()
			p.selectedSector = nil
		}
	}

	if p.input.ActionIsJustPressed(ActionConfirm) {
		if base != nil && base.Player.ID == p.ID && p.selectedSector == nil {
			p.selectedSector = p.sector()
			p.sector().SelectUnits()
		} else if p.selectedSector != nil && p.selectedSector.NumDefenders() != 0 {
			p.selectedSector.UnselectUnits()
			p.selectedSector.SendUnits(p.sector())
			p.selectedSector = nil
			if p.BattleState.SingleLocalPlayer == p.playerData {
				p.scene.Audio().PlaySound(AudioCueSendUnits)
			}
		}
	}

	state := p.BattleState
	if p.input.ActionIsJustPressed(ActionSectorRight) {
		p.sectorSelector.SetSector(state.GetSector(state.WrapXY(p.sector().X+1, p.sector().Y)))
	} else if p.input.ActionIsJustPressed(ActionSectorLeft) {
		p.sectorSelector.SetSector(state.GetSector(state.WrapXY(p.sector().X-1, p.sector().Y)))
	} else if p.input.ActionIsJustPressed(ActionSectorDown) {
		p.sectorSelector.SetSector(state.GetSector(state.WrapXY(p.sector().X, p.sector().Y+1)))
	} else if p.input.ActionIsJustPressed(ActionSectorUp) {
		p.sectorSelector.SetSector(state.GetSector(state.WrapXY(p.sector().X, p.sector().Y-1)))
	}
}

func (p *localPlayer) sector() *sector {
	return p.sectorSelector.Sector()
}
