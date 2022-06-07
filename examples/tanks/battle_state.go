package main

import "github.com/quasilyte/ge/gemath"

const maxPlayers = 4

type battleState struct {
	numRows int
	numCols int

	Sectors []*sector
	Groups  map[*battleGroup]struct{}
	Tanks   [maxPlayers]map[*battleTank]struct{}
	Players [maxPlayers]playerData

	FortificationsAllowed bool
	MudTerrain            bool
	HQDefeat              bool
	DynamicAlliances      bool
}

func newBattleState() *battleState {
	state := &battleState{
		numRows: 1080 / 270,
		numCols: 1920 / 320,
		Groups:  make(map[*battleGroup]struct{}),
	}
	for i := 0; i < maxPlayers; i++ {
		state.Tanks[i] = make(map[*battleTank]struct{})
		state.Players[i] = playerData{
			ID: i,
			Resources: resourceContainer{
				Iron: 20,
				Gold: 20,
				Oil:  20,
			},
			BattleState: state,
		}
	}
	return state
}

func (state *battleState) AddGroup(g *battleGroup) {
	g.EventDisposed.Connect(nil, state.onGroupDisposed)
	state.Groups[g] = struct{}{}
}

func (state *battleState) onGroupDisposed(g *battleGroup) {
	delete(state.Groups, g)
}

func (state *battleState) NewBattlePost(p *playerData, pos gemath.Vec, turret *turretDesign) *battlePost {
	bp := newBattlePost(p, pos, turret)
	bp.EventDestroyed.Connect(nil, state.onBaseDestroyed)
	return bp
}

func (state *battleState) NewBattleTank(p *playerData, design tankDesign) *battleTank {
	bt := newBattleTank(p, design, state.MudTerrain)
	state.addTank(bt)
	return bt
}

func (state *battleState) addTank(bt *battleTank) {
	state.Tanks[bt.Player.ID][bt] = struct{}{}
	bt.EventDestroyed.Connect(nil, state.onTankDestroyed)
}

func (state *battleState) DeploySectors(r *gemath.Rand) {
	for y := 0; y < state.numRows; y++ {
		for x := 0; x < state.numCols; x++ {
			resource := resourceKind(r.IntRange(0, 2))
			id := len(state.Sectors)
			s := newSector(resource, id, x, y)
			state.Sectors = append(state.Sectors, s)
		}
	}
}

func (state *battleState) onTankDestroyed(bt *battleTank) {
	delete(state.Tanks[bt.Player.ID], bt)
}

func (state *battleState) onBaseDestroyed(bp *battlePost) {
	if state.HQDefeat && bp.HQ {
		state.playerDefeat(bp.Player)
	}
}

func (state *battleState) playerDefeat(p *playerData) {
	for _, s := range state.Sectors {
		if s.Base == nil || s.Base.IsDisposed() || s.Base.Player.ID != p.ID {
			continue
		}
		s.Base.Destroy()
	}
	for bt := range state.Tanks[p.ID] {
		bt.Destroy()
	}
}

func (state *battleState) GetSector(x, y int) *sector {
	return state.Sectors[y*state.numCols+x]
}

func (state *battleState) WrapXY(x, y int) (int, int) {
	if x >= state.numCols {
		x = 0
	} else if x < 0 {
		x = state.numCols - 1
	}
	if y >= state.numRows {
		y = 0
	} else if y < 0 {
		y = state.numRows - 1
	}
	return x, y
}
