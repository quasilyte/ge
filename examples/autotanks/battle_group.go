package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type battleGroup struct {
	scene *ge.Scene

	player *playerData

	tanks       []*battleTank
	baseBuilder *battleTank

	numBuilders int

	dstSector *sector

	disposed bool

	EventDisposed gesignal.Event[*battleGroup]
}

func newBattleGroup(p *playerData, dst *sector, tanks []*battleTank) *battleGroup {
	return &battleGroup{player: p, dstSector: dst, tanks: tanks}
}

func removeTankFromList(list []*battleTank, bt *battleTank) []*battleTank {
	for i := range list {
		if list[i] != bt {
			continue
		}
		list[i] = list[len(list)-1]
		return list[:len(list)-1]
	}
	return list
}

func (g *battleGroup) Init(scene *ge.Scene) {
	for _, bt := range g.tanks {
		if bt.Turret.IsBuilder() {
			g.numBuilders++
		}
		bt.EventDestroyed.Connect(g, g.onTankDestroyed)
		bt.EventWaypointReached.Connect(g, g.onTankWaypointReached)
	}

	g.scene = scene
	g.setWaypoint(g.dstSector.Center())
}

func (g *battleGroup) IsDisposed() bool {
	return g.disposed
}

func (g *battleGroup) Update(delta float64) {}

func (g *battleGroup) setWaypoint(pos gmath.Vec) {
	if len(g.tanks) == 0 {
		panic("empty group?")
	}
	occupied := [len(sectorLocations)]bool{}
	for i := 0; i < len(g.tanks); i++ {
		probe := g.scene.Rand().IntRange(0, len(sectorLocations)-1)
		for occupied[probe] {
			probe++
			if probe >= len(occupied) {
				probe = 0
			}
		}
		g.tanks[i].Waypoint = sectorLocations[probe].Add(pos)
		occupied[probe] = true
	}
}

func (g *battleGroup) checkDisposed() bool {
	if len(g.tanks) == 0 {
		g.disposed = true
		g.EventDisposed.Emit(g)
	}
	return g.disposed
}

func (g *battleGroup) onTankDestroyed(bt *battleTank) {
	g.tanks = removeTankFromList(g.tanks, bt)
	if g.checkDisposed() {
		return
	}
	if bt.Turret.IsBuilder() {
		g.numBuilders--
	}

	if g.baseBuilder == bt {
		// Failed to build a base.
		g.baseBuilder = nil
		g.pickNextAction()
		return
	}
}

func (g *battleGroup) onTankWaypointReached(bt *battleTank) {
	if bt == g.baseBuilder {
		// While the builder was moving to the center, someone
		// could already build a base.
		if g.dstSector.Base == nil {
			g.buildBase()
			return
		}
		g.baseBuilder = nil
	}

	bt.Waypoint = gmath.Vec{}

	allReady := xslices.All(g.tanks, func(bt *battleTank) bool {
		return bt.Waypoint.IsZero()
	})
	if allReady {
		g.pickNextAction()
	}
}

func (g *battleGroup) pickNextAction() {
	if g.numBuilders > 0 && g.dstSector.Base == nil {
		for _, bt := range g.tanks {
			if bt.Turret.IsBuilder() {
				g.baseBuilder = bt
				bt.Waypoint = g.dstSector.Center()
				return
			}
		}
	}

	if g.numBuilders > 0 && g.dstSector.Base != nil {
		var closest *sector
		closestDist := math.MaxFloat64
		for _, s := range g.player.BattleState.Sectors {
			if s.Base != nil || s.Resource != g.dstSector.Resource {
				continue
			}
			dist := g.dstSector.Center().DistanceTo(s.Center())
			if closest == nil || dist < closestDist {
				closestDist = dist
				closest = s
			}
		}
		if closest != nil {
			g.dstSector = closest
			g.setWaypoint(closest.Center())
			return
		}
	}

	g.pickNewWaypoint()
}

func (g *battleGroup) pickNewWaypoint() {
	candidates := make([]*sector, 0, 8)
	for _, s := range g.player.BattleState.Sectors {
		if s.Base == nil || s.Base.Player.Alliance == g.player.Alliance {
			continue
		}
		candidates = append(candidates, s)
	}
	if len(candidates) == 0 {
		return
	}
	s := gmath.RandElem(g.scene.Rand(), candidates)
	g.dstSector = s
	g.setWaypoint(s.Center())
}

func (g *battleGroup) buildBase() {
	// The base builder tank is transformed into a base.
	tanks := removeTankFromList(g.tanks, g.baseBuilder)
	g.baseBuilder.EventDestroyed.Disconnect(g)
	g.baseBuilder.EventDestroyed.Emit(g.baseBuilder)
	g.baseBuilder.Dispose()

	bp := g.player.BattleState.NewBattlePost(g.player, g.dstSector.Center(), nil)
	g.dstSector.AssignBase(bp)
	for _, bt := range tanks {
		g.dstSector.AddTank(bt)
	}
	g.scene.AddObject(bp)

	if g.player == g.player.BattleState.SingleLocalPlayer {
		g.scene.Audio().EnqueueSound(AudioCueConstructionCompleted)
	}

	g.tanks = nil
	g.checkDisposed()
}

var sectorLocations = [...]gmath.Vec{
	{X: -64 * 2, Y: -64},
	{X: -64 * 1, Y: -64},
	{X: 0, Y: -64},
	{X: 64 * 1, Y: -64},
	{X: 64 * 2, Y: -64},

	{X: -64 * 2, Y: 0},
	{X: -64 * 1, Y: 0},
	{X: 64 * 1, Y: 0},
	{X: 64 * 2, Y: 0},

	{X: -64 * 2, Y: 64},
	{X: -64 * 1, Y: 64},
	{X: 0, Y: 64},
	{X: 64 * 1, Y: 64},
	{X: 64 * 2, Y: 64},
}
