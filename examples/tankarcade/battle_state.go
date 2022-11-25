package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type battlePlayer struct {
	unit     *battleUnit
	active   bool
	spawnPos gmath.Vec
}

type cellInfo uint8

const (
	cellRoamStop cellInfo = 1 << iota
	cellDarkTile
)

type battleState struct {
	scene *ge.Scene

	gameState *gameState

	rand *gmath.Rand

	bases           []*battleBase
	forts           []*battleFort
	factories       []*battleFactory
	bonusGenerators []*battleBonusGenerator
	units           []*battleUnit
	walls           []*battleWall

	players []*battlePlayer

	wallMap  [14][30]uint8
	cellInfo [14][30]cellInfo

	forceBaseAttack bool

	tmpTargetSlice []targetable
}

func newBattleState(s *gameState, numPlayers int) *battleState {
	players := make([]*battlePlayer, 3)
	for i := range players {
		players[i] = &battlePlayer{active: i < numPlayers}
	}
	return &battleState{
		gameState:      s,
		tmpTargetSlice: make([]targetable, 0, 32),
		players:        players,
	}
}

func (s *battleState) numPlayers(alliance int) int {
	return xslices.CountIf(s.units, func(u *battleUnit) bool {
		return u.config.alliance == alliance && u.config.playerID != 0
	})
}

func (s *battleState) walkUnitsOfGroup(group string, f func(u *battleUnit)) {
	for _, u := range s.units {
		if u.config.group == group {
			f(u)
		}
	}
}

func (s *battleState) newBattleUnit(config battleUnitConfig) *battleUnit {
	u := newBattleUnit(config)
	s.units = append(s.units, u)
	if u.config.playerID != 0 {
		s.players[u.config.playerID-1].unit = u
	}
	u.EventDestroyed.Connect(nil, func(u *battleUnit) {
		s.units = xslices.Remove(s.units, u)
		if u.config.playerID != 0 {
			s.players[u.config.playerID-1].unit = nil
		}
		if u.config.alliance != 0 {
			return
		}
		numPlayers := s.numPlayers(u.config.alliance)
		if numPlayers != 0 {
			return
		}
		s.scene.DelayedCall(2, func() {
			s.endBattle(false)
		})
	})
	return u
}

func (s *battleState) endBattle(victory bool) {
	numActive := xslices.CountIf(s.players, func(p *battlePlayer) bool {
		return p.active
	})
	s.scene.Audio().PauseCurrentMusic()
	s.scene.Context().ChangeScene(newLevelTransitionController(s.gameState, numActive, victory))
}

func (s *battleState) newBattleBase(pos gmath.Vec, alliance, level int) *battleBase {
	b := newBattleBase(pos, alliance, level)
	s.bases = append(s.bases, b)
	b.EventDestroyed.Connect(nil, func(b *battleBase) {
		s.bases = xslices.Remove(s.bases, b)
		numBases := xslices.CountIf(s.bases, func(b2 *battleBase) bool {
			return b2.alliance == b.alliance
		})
		if numBases != 0 {
			return
		}
		s.scene.DelayedCall(2, func() {
			s.endBattle(b.alliance != 0)
		})
	})
	return b
}

func (s *battleState) newBattleFactory(config battleFactoryConfig) *battleFactory {
	f := newBattleFactory(config)
	s.factories = append(s.factories, f)
	f.EventDestroyed.Connect(nil, func(f *battleFactory) {
		s.factories = xslices.Remove(s.factories, f)
	})
	return f
}

func (s *battleState) newBattleBonusGenerator(config battleBonusGeneratorConfig) *battleBonusGenerator {
	b := newBattleBonusGenerator(config)
	s.bonusGenerators = append(s.bonusGenerators, b)
	b.EventDestroyed.Connect(nil, func(f *battleBonusGenerator) {
		s.bonusGenerators = xslices.Remove(s.bonusGenerators, b)
	})
	return b
}

func (s *battleState) newBattleFort(config battleFortConfig) *battleFort {
	f := newBattleFort(config)
	s.forts = append(s.forts, f)
	f.EventDestroyed.Connect(nil, func(f *battleFort) {
		s.forts = xslices.Remove(s.forts, f)
	})
	return f
}

func (s *battleState) newBattleWall(pos gmath.Vec, alliance int) *battleWall {
	w := newBattleWall(pos, alliance)
	s.walls = append(s.walls, w)
	w.EventDestroyed.Connect(nil, func(w *battleWall) {
		s.walls = xslices.Remove(s.walls, w)
		s.updateWalls()
	})
	return w
}

func (s *battleState) getWallAt(pos gmath.Vec) *battleWall {
	x := int(pos.X) / 64
	y := int(pos.Y) / 64
	return s.getWallAtXY(x, y)
}

func (s *battleState) getCellInfoAt(pos gmath.Vec) cellInfo {
	x := int(pos.X) / 64
	y := int(pos.Y) / 64
	return s.cellInfo[y][x]
}

func (s *battleState) addCellInfoAt(pos gmath.Vec, bits cellInfo) {
	x := int(pos.X) / 64
	y := int(pos.Y) / 64
	s.cellInfo[y][x] |= bits
}

func (s *battleState) getWallAtXY(x, y int) *battleWall {
	index := s.wallMap[y][x]
	if index == 0 {
		return nil
	}
	return s.walls[index-1]
}

func (s *battleState) updateWalls() {
	if len(s.walls) > (255 - 1) {
		panic("too many walls!")
	}
	s.wallMap = [14][30]uint8{}
	for i, w := range s.walls {
		id := uint8(i)
		s.wallMap[w.gridY][w.gridX] = (id + 1)
	}
	type checkOption struct {
		dx int8
		dy int8
	}
	checkList := [4]checkOption{
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
	}
	for _, w := range s.walls {
		connectionsMask := uint8(0)
		for bitIndex, option := range checkList {
			x := w.gridX + int(option.dx)
			y := w.gridY + int(option.dy)
			if x < 0 || x >= 30 || y < 0 || y >= 14 {
				continue // This grid cell is out of bounds
			}
			w2 := s.getWallAtXY(x, y)
			if w2 == nil {
				continue // No wall here
			}
			if w.alliance != w2.alliance {
				continue // Only friendly walls can connect
			}
			bitMask := 1 << bitIndex
			connectionsMask |= uint8(bitMask)
		}
		w.SetFrame(connectionsMask)
	}
}

func (s *battleState) findTargetUnit(pos gmath.Vec, maxDist float64, alliance int) targetable {
	s.tmpTargetSlice = s.tmpTargetSlice[:0]
	for _, x := range s.units {
		if x.config.alliance == alliance {
			continue
		}
		if maxDist != 0 && pos.DistanceTo(x.Body.Pos) > maxDist {
			continue
		}
		s.tmpTargetSlice = append(s.tmpTargetSlice, x)
	}
	if len(s.tmpTargetSlice) == 0 {
		return nil
	}
	return gmath.RandElem(s.rand, s.tmpTargetSlice)
}

func (s *battleState) findClosestBase(pos gmath.Vec, alliance int) targetable {
	minDist := math.MaxFloat64
	var target targetable
	for _, x := range s.bases {
		if x.alliance == alliance {
			continue
		}
		dist := pos.DistanceTo(x.Body.Pos)
		if dist < minDist {
			minDist = dist
			target = x
		}
	}
	return target
}

func (s *battleState) findTargetBuilding(alliance int) targetable {
	s.tmpTargetSlice = s.tmpTargetSlice[:0]
	for _, x := range s.forts {
		if x.alliance == alliance {
			continue
		}
		s.tmpTargetSlice = append(s.tmpTargetSlice, x)
	}
	for _, x := range s.factories {
		if x.alliance == alliance {
			continue
		}
		s.tmpTargetSlice = append(s.tmpTargetSlice, x)
	}
	for _, x := range s.bonusGenerators {
		if x.alliance == alliance {
			continue
		}
		s.tmpTargetSlice = append(s.tmpTargetSlice, x)
	}
	for _, x := range s.bases {
		if x.alliance == alliance {
			continue
		}
		s.tmpTargetSlice = append(s.tmpTargetSlice, x)
	}
	if len(s.tmpTargetSlice) == 0 {
		return nil
	}
	return gmath.RandElem(s.rand, s.tmpTargetSlice)
}
