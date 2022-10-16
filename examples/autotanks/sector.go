package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/xslices"
)

type sector struct {
	ID       int
	X        int
	Y        int
	Resource resourceKind

	Base *battlePost

	scene *ge.Scene

	Pos gemath.Vec

	defenders []*battleTank

	unitsSelected bool

	resourceIcon  *ge.Sprite
	bordersSprite *ge.Sprite
}

func newSector(resource resourceKind, id, x, y int) *sector {
	return &sector{
		ID:       id,
		X:        x,
		Y:        y,
		Resource: resource,
		Pos:      gemath.Vec{X: float64(x) * 320, Y: float64(y) * 270},
	}
}

func (s *sector) Init(scene *ge.Scene) {
	s.scene = scene

	switch s.Resource {
	case resIron:
		s.resourceIcon = scene.NewSprite(ImageIronResourceIcon)
	case resGold:
		s.resourceIcon = scene.NewSprite(ImageGoldResourceIcon)
	case resOil:
		s.resourceIcon = scene.NewSprite(ImageOilResourceIcon)
	case resCombined:
		s.resourceIcon = scene.NewSprite(ImageCombinedResourceIcon)
	}
	if s.resourceIcon != nil {
		s.resourceIcon.Pos.Base = &s.Pos
		s.resourceIcon.SetAlpha(0.6)
		s.resourceIcon.Centered = false
		scene.AddGraphics(s.resourceIcon)
	}

	s.bordersSprite = scene.NewSprite(ImageGrid)
	s.bordersSprite.Centered = false
	s.bordersSprite.SetColorScale(ge.ColorScale{A: 0.2})
	s.bordersSprite.Pos.Base = &s.Pos
	scene.AddGraphics(s.bordersSprite)
}

func (s *sector) IsDisposed() bool { return false }

func (s *sector) Update(delta float64) {}

func (s *sector) Center() gemath.Vec {
	return s.Pos.Add(gemath.Vec{X: 320 / 2, Y: 270 / 2})
}

func (s *sector) AddTank(bt *battleTank) {
	s.defenders = append(s.defenders, bt)
	bt.Selected = s.unitsSelected
}

func (s *sector) NumDefenders() int {
	s.updateDefenders()
	return len(s.defenders)
}

func (s *sector) HasLancer() bool {
	s.updateDefenders()
	return xslices.Any(s.defenders, func(bt *battleTank) bool {
		return bt.Turret.IsLancer()
	})
}

func (s *sector) HasBuilder() bool {
	s.updateDefenders()
	return xslices.Any(s.defenders, func(bt *battleTank) bool {
		return bt.Turret.IsBuilder()
	})
}

func (s *sector) updateDefenders() {
	s.defenders = xslices.RemoveIf(s.defenders, func(bt *battleTank) bool {
		return bt.IsDisposed()
	})
}

func (s *sector) SendUnits(target *sector) {
	s.updateDefenders()
	if len(s.defenders) == 0 {
		return
	}
	tanks := make([]*battleTank, len(s.defenders))
	copy(tanks, s.defenders)
	g := newBattleGroup(s.Base.Player, target, tanks)
	s.Base.Player.BattleState.AddGroup(g)
	s.scene.AddObject(g)
	s.defenders = s.defenders[:0]
}

func (s *sector) AssignBase(bp *battlePost) {
	s.Base = bp
	bp.EventDestroyed.Connect(nil, s.onBaseDestroyed)
	bp.EventProductionCompleted.Connect(nil, s.onProductionCompleted)
}

func (s *sector) onBaseDestroyed(bp *battlePost) {
	s.SendUnits(s)
	s.Base = nil
}

func (s *sector) onProductionCompleted(design tankDesign) {
	st := s.Base.Player.BattleState
	bt := st.NewBattleTank(s.Base.Player, design)
	bt.Body.Pos = s.Center().Add(gemath.Vec{Y: 48})
	bt.Waypoint = s.Center().Add(gemath.Vec{Y: 64})
	s.AddTank(bt)
	s.scene.AddObjectAbove(bt, 1)
	if st.SingleLocalPlayer == s.Base.Player {
		s.scene.Audio().EnqueueSound(AudioCueProductionCompleted)
	}
}

func (s *sector) SelectUnits() {
	s.updateDefenders()
	s.unitsSelected = true
	for _, bt := range s.defenders {
		bt.Selected = true
	}
}

func (s *sector) UnselectUnits() {
	s.updateDefenders()
	s.unitsSelected = false
	for _, bt := range s.defenders {
		bt.Selected = false
	}
}
