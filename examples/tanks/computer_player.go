package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
)

type computerPlayer struct {
	*playerData

	scene *ge.Scene
	easy  bool

	startingSector *sector

	orderBuilderDelay float64
	orderTankDelay    float64
	actionDelay       float64

	enemySectors   []*sector
	ownSectors     []*sector
	numFreeSectors int

	favDesigns []computerDesignOption
}

func newComputerPlayer(data *playerData, state *battleState, startingSector *sector, easy bool) *computerPlayer {
	p := &computerPlayer{
		playerData:     data,
		startingSector: startingSector,
		easy:           easy,
	}
	p.initFavDesigns()
	return p
}

func (p *computerPlayer) Init(scene *ge.Scene) {
	p.scene = scene

	if p.easy {
		p.orderBuilderDelay = scene.Rand().FloatRange(15, 25)
	} else {
		p.orderBuilderDelay = scene.Rand().FloatRange(0, 5)
	}

	bp := p.BattleState.NewBattlePost(p.playerData, p.startingSector.Center(), ionTurretDesign)
	bp.HQ = true
	p.startingSector.AssignBase(bp)
	scene.AddObject(bp)

	firstBuilder := p.BattleState.NewBattleTank(bp.Player, tankDesign{
		Hull:   hullDesigns["fighter"],
		Turret: turretDesigns["builder"],
	})
	firstBuilder.Body.Pos = p.startingSector.Center().Add(gemath.Vec{Y: 64})
	firstBuilder.Body.Rotation = firstBuilder.Body.Pos.AngleToPoint(gemath.Vec{X: 1920 / 2, Y: 1080 / 2})
	p.startingSector.AddTank(firstBuilder)
	scene.AddObject(firstBuilder)
}

func (p *computerPlayer) IsDisposed() bool { return false }

func (p *computerPlayer) Update(delta float64) {
	p.actionDelay = gemath.ClampMin(p.actionDelay-delta, 0)
	p.orderBuilderDelay = gemath.ClampMin(p.orderBuilderDelay-delta, 0)
	p.orderTankDelay = gemath.ClampMin(p.orderTankDelay-delta, 0)
	if p.actionDelay != 0 {
		return
	}

	actionRoll := p.scene.Rand().Float()
	switch {
	case actionRoll < 0.6:
		p.orderProduction()
	case actionRoll < 0.7:
		p.sendBuilder()
	case actionRoll < 0.8:
		p.sendUnits()
	case actionRoll < 0.9:
		p.defendBase()
	case actionRoll < 0.95:
		p.buildDefence()
	default:
		// Idle.
	}

	if p.easy {
		p.actionDelay = p.scene.Rand().FloatRange(0.6, 1.2)
		if p.scene.Rand().Chance(0.4) {
			p.actionDelay += 6
		}
	} else {
		delay := p.scene.Rand().FloatRange(0.1, 0.2)
		if p.BattleState.DoubledIncome {
			delay *= 0.5
		}
		p.actionDelay = delay
	}
}

func (p *computerPlayer) buildDefence() {
	if !p.BattleState.FortificationsAllowed {
		return // Fortifications are disabled
	}
	if !p.Resources.Contains(gaussTurretDesign.Price) {
		return // Can't build anyway
	}
	if p.Resources.Iron > 10 && p.Resources.Gold > 10 && p.Resources.Oil > 10 {
		// Keep resources for now. No one resource is scarce at the moment.
		// The bot can build a turret when the base is attacked.
		return
	}

	p.updateSectorsInfo()

	scarceResource := resIron
	switch {
	case p.Resources.Gold < p.Resources.Iron && p.Resources.Gold < p.Resources.Oil:
		scarceResource = resGold
	case p.Resources.Oil < p.Resources.Iron && p.Resources.Oil < p.Resources.Gold:
		scarceResource = resOil
	}

	s := p.selectOwnedSector(func(s *sector) bool {
		if s.Base.Turret != nil {
			return false
		}
		return s.Resource == scarceResource
	})
	if s == nil {
		return
	}

	s.Base.InstallTurret(gaussTurretDesign)
	p.Resources.Sub(gaussTurretDesign.Price)
}

func (p *computerPlayer) defendBase() {
	p.updateSectorsInfo()

	for _, s := range p.ownSectors {
		if s.Base.UnderAttack == 0 {
			continue
		}
		if p.BattleState.FortificationsAllowed {
			// If we're rich, building a turret could be a good idea.
			turretPriceX2 := gaussTurretDesign.Price
			turretPriceX2.Add(turretPriceX2)
			if s.Base.Turret == nil && p.Resources.Contains(turretPriceX2) {
				s.Base.InstallTurret(gaussTurretDesign)
				p.Resources.Sub(gaussTurretDesign.Price)
				return
			}
		}
		// Can use the base defenders to repell the attack.
		if s.NumDefenders() != 0 {
			s.SendUnits(s)
			return
		}
		// Need reinforcements.
		otherSector := p.selectOwnedSector(func(candidate *sector) bool {
			if candidate == s || candidate.NumDefenders() == 0 {
				return false
			}
			return !candidate.HasBuilder()
		})
		if otherSector != nil {
			otherSector.SendUnits(s)
			return
		}
		// As a last resort, order a small tank to distract the enemy?
		if s.Base.Turret == nil && s.Base.HP >= 300 {
			cheapestDesign := p.pickCheapDesign()
			if p.Resources.Contains(cheapestDesign.Price()) {
				s.Base.StartProduction(cheapestDesign)
				p.Resources.Sub(cheapestDesign.Price())
			}
		}
		return
	}
}

func (p *computerPlayer) sendBuilder() {
	p.updateSectorsInfo()

	s := p.selectOwnedSector(func(s *sector) bool {
		return s.HasBuilder()
	})
	if s == nil {
		return
	}

	if best := p.pickSectorToOccupy(s); best != nil {
		s.SendUnits(best)
	}
}

func (p *computerPlayer) sendUnits() {
	p.updateSectorsInfo()

	squadSizeRoll := p.scene.Rand().Float()
	var minSquadSize int
	switch {
	case squadSizeRoll < 0.05:
		minSquadSize = 1
	case squadSizeRoll < 0.6:
		minSquadSize = 2
	default:
		minSquadSize = 3
	}
	s := p.selectOwnedSector(func(s *sector) bool {
		return s.NumDefenders() >= minSquadSize
	})
	if s == nil {
		return // Can't assemble an attack group
	}

	// Pick the initial waypoint.
	if s.HasBuilder() {
		if best := p.pickSectorToOccupy(s); best != nil {
			s.SendUnits(best)
			return
		}
	}

	// Choose a sector to attack.
	var attackTarget *sector
	if !p.easy && p.scene.Rand().Bool() {
		// Pick a wise attack target.
		// Prefer bases that:
		// * Don't have a defensive turret
		// * Are located closer to the attack squad
		// * Have less defenders
		var bestScore float64
		for _, candidate := range p.enemySectors {
			numDefendersScore := 1.5 / float64(s.NumDefenders())
			turretScore := 2.0
			if s.Base.Turret != nil {
				turretScore = 0
			}
			healthScore := 0.0
			if !p.easy {
				switch {
				case s.Base.HP <= 100:
					healthScore = 1.25
				case s.Base.HP <= 400:
					healthScore = 0.75
				}
			}
			resourceScore := 1.5
			if s.Resource == resNone {
				resourceScore = 0.5
			}
			rollScore := p.scene.Rand().FloatRange(0, 0.6)
			score := (3000 - candidate.Center().DistanceTo(s.Center())) * (rollScore + healthScore + numDefendersScore + turretScore + resourceScore)
			if score > bestScore {
				attackTarget = candidate
			}
		}
	} else {
		// Choose a random target.
		attackTarget = gemath.RandElem(p.scene.Rand(), p.enemySectors)
	}
	if attackTarget == nil {
		return
	}
	// Do not try to attack a base with insufficient forces.
	if !p.easy && attackTarget.Base.Turret != nil && s.NumDefenders() < 3 && !s.HasLancer() {
		if p.scene.Rand().Chance(0.9) {
			return
		}
	}
	s.SendUnits(attackTarget)
}

func (p *computerPlayer) resourceScore(kind resourceKind) float64 {
	if p.easy {
		return p.scene.Rand().FloatRange(0.5, 4.0)
	}

	stock := p.Resources.Iron
	switch kind {
	case resGold:
		stock = p.Resources.Gold
	case resOil:
		stock = p.Resources.Oil
	}
	stockMultiplier := 0.5
	switch {
	case stock < 10:
		stockMultiplier = 3
	case stock < 20:
		stockMultiplier = 2.5
	case stock < 30:
		stockMultiplier = 1.75
	case stock < 40:
		stockMultiplier = 1
	}
	numSectorsMultiplier := 0.0
	numSectors := 0
	for _, s := range p.ownSectors {
		if s.Resource == kind {
			numSectors++
		}
	}
	switch numSectors {
	case 0:
		numSectorsMultiplier = 5.0
	case 1:
		numSectorsMultiplier = 1
	case 2:
		numSectorsMultiplier = 0.25
	}
	return stockMultiplier + numSectorsMultiplier
}

func (p *computerPlayer) orderProduction() {
	p.updateSectorsInfo()
	if len(p.ownSectors) == 0 {
		return
	}

	// Pick a base where a new unit will be produced.
	var s *sector
	pick := gemath.RandIndex(p.scene.Rand(), p.ownSectors)
	for i := 0; i < len(p.ownSectors); i++ {
		if !p.ownSectors[pick].Base.IsBusy() && p.ownSectors[pick].NumDefenders() < 4 {
			s = p.ownSectors[pick]
			break
		}
		pick++
		if pick >= len(p.ownSectors) {
			pick = 0
		}
	}
	if s == nil {
		return // All bases are busy
	}

	// When already have a good amount of sectors, order a lancer
	// instead of a builder to add some pressure to the opponent.
	// 5 bases => 50% chance
	// 10 bases => 75% chance
	lancerChance := 0.25 + float64(len(p.ownSectors))*0.05
	if !p.easy && p.orderBuilderDelay == 0 && len(p.ownSectors) >= 5 && p.scene.Rand().Chance(lancerChance) {
		lancerDesign := tankDesign{
			Hull:   hullDesigns["viper"],
			Turret: turretDesigns["lancer"],
		}
		if p.Resources.Contains(lancerDesign.Price()) {
			s.Base.StartProduction(lancerDesign)
			p.Resources.Sub(lancerDesign.Price())
		}
	}

	if p.orderBuilderDelay == 0 && (p.numFreeSectors != 0 || p.scene.Rand().Bool()) {
		var builderDesign tankDesign
		builderDesign.Turret = turretDesigns["builder"]
		if p.Resources.Iron >= 14 && p.Resources.Gold >= 12 && p.Resources.Oil >= 12 && p.scene.Rand().Chance(0.85) {
			builderDesign.Hull = hullDesigns["hunter"]
		} else if p.Resources.Gold > p.Resources.Oil {
			builderDesign.Hull = hullDesigns["viper"]
		} else {
			builderDesign.Hull = hullDesigns["scout"]
		}
		if p.Resources.Contains(builderDesign.Price()) {
			p.updateBuilderDelay()
			s.Base.StartProduction(builderDesign)
			p.Resources.Sub(builderDesign.Price())
			return
		}
		if p.scene.Rand().Bool() {
			return
		}
		p.orderBuilderDelay = p.scene.Rand().FloatRange(2, 5)
		if len(p.ownSectors) < 8 {
			p.orderTankDelay = p.scene.Rand().FloatRange(4, 8)
		}
	}

	if !p.easy && p.orderTankDelay == 0 {
		waitForResources := p.Resources.Iron <= 6 &&
			p.Resources.Gold <= 5 &&
			p.Resources.Oil <= 5 &&
			len(p.ownSectors) < 5
		if waitForResources {
			p.orderTankDelay = p.scene.Rand().FloatRange(4, 10)
			return
		}
	}

	if p.orderTankDelay != 0 {
		return
	}

	cheapestDesign := p.pickCheapDesign()
	if !p.Resources.Contains(cheapestDesign.Price()) {
		// Can't order anything.
		return
	}

	var selectedDesign tankDesign
	if !p.easy && p.scene.Rand().Bool() {
		selectedDesign = p.selectFavDesign()
	}
	if selectedDesign.IsEmpty() {
		for i := 0; i < 4; i++ {
			var turret *turretDesign
			turret = gemath.RandElem(p.scene.Rand(), turretDesignListNoBuilder)
			hull := gemath.RandElem(p.scene.Rand(), hullDesignList)
			design := tankDesign{Hull: hull, Turret: turret}
			if p.Resources.Contains(design.Price()) {
				selectedDesign = design
				break
			}
		}
	}
	if selectedDesign.IsEmpty() {
		if p.scene.Rand().Chance(0.75) {
			// No resources to build units.
			// Consider preserving resources for a while and/or order a builder sooner.
			p.orderTankDelay = p.scene.Rand().FloatRange(0.5, 3)
			if len(p.ownSectors) < 8 {
				p.orderBuilderDelay = gemath.ClampMin(p.orderBuilderDelay*0.8, 0)
			}
			return
		}
		selectedDesign = cheapestDesign
	}
	s.Base.StartProduction(selectedDesign)
	p.Resources.Sub(selectedDesign.Price())
}

func (p *computerPlayer) updateBuilderDelay() {
	// Easy bots only consider the number of sectors they have.
	// More sectors cause them to be less agressive.
	if p.easy {
		// This calculation method is more complicated than it should be really...
		delay := float64((len(p.ownSectors)-1)*5 + 5)
		p.orderBuilderDelay = p.scene.Rand().FloatRange(delay, delay*1.05)
		p.orderBuilderDelay = gemath.ClampMax(p.orderBuilderDelay, 30)
		p.orderBuilderDelay += 10
		p.orderBuilderDelay *= 2
		return
	}

	// For harder bots, they will focus on sectors capturing as long
	// as there is something to capture.
	minBaseDelay := 2.0
	maxBaseDelay := 6.0
	delay := p.scene.Rand().FloatRange(minBaseDelay, maxBaseDelay)
	delay += (24.0 - float64(p.numFreeSectors)) * 1.5
	switch {
	case len(p.ownSectors) < 4:
		delay *= 0.7
	case len(p.ownSectors) < 8:
		delay *= 0.9
	}
	if p.BattleState.DoubledIncome {
		delay *= 0.4
	}
	p.orderBuilderDelay = delay
}

func (p *computerPlayer) updateSectorsInfo() {
	p.ownSectors = p.ownSectors[:0]
	p.enemySectors = p.enemySectors[:0]
	p.numFreeSectors = 0
	for _, s := range p.BattleState.Sectors {
		if s.Base == nil {
			p.numFreeSectors++
		}
		if s.Base != nil && s.Base.Player.Alliance != p.Alliance {
			p.enemySectors = append(p.enemySectors, s)
		} else if s.Base != nil && s.Base.Player.ID == p.ID {
			p.ownSectors = append(p.ownSectors, s)
		}
	}
}

func (p *computerPlayer) pickSectorToOccupy(from *sector) *sector {
	// Choose a sector to occupy.
	// Two criterias: proximity and resource kinds (their score).
	ironScore := p.resourceScore(resIron) + 0.1
	goldScore := p.resourceScore(resGold)
	oilScore := p.resourceScore(resOil) + 0.2
	combinedScore := ironScore + goldScore + oilScore
	var best *sector
	bestScore := float64(0)
	for _, candidate := range p.BattleState.Sectors {
		if candidate.Base != nil {
			continue
		}
		score := 3000 - candidate.Center().DistanceTo(from.Center())
		switch candidate.Resource {
		case resIron:
			score *= ironScore
		case resGold:
			score *= goldScore
		case resOil:
			score *= oilScore
		case resCombined:
			score *= combinedScore
		case resNone:
			score *= 0.8 // Some very low score multiplier
		}
		if score > bestScore {
			bestScore = score
			best = candidate
		}
	}
	return best
}

type computerDesignOption struct {
	hull   string
	turret string
	cond   func() bool
}

func (p *computerPlayer) initFavDesigns() {
	var options = []computerDesignOption{
		{
			hull:   "scout",
			turret: "gatling gun",
			cond: func() bool {
				return p.Resources.Iron < 8 && p.Resources.Gold < 8
			},
		},

		{
			hull:   "fighter",
			turret: "light cannon",
			cond: func() bool {
				return p.resourceScore(resIron) < p.resourceScore(resGold) &&
					p.resourceScore(resIron) < p.resourceScore(resOil) &&
					p.Resources.Iron >= 14
			},
		},

		{
			hull:   "hunter",
			turret: "dual cannon",
		},

		{
			hull:   "scorpion",
			turret: "dual cannon",
			cond: func() bool {
				return p.resourceScore(resGold) > p.resourceScore(resOil)
			},
		},

		{
			hull:   "fighter",
			turret: "dual cannon",
			cond: func() bool {
				return p.resourceScore(resIron) < p.resourceScore(resOil) &&
					p.Resources.Iron >= 16
			},
		},

		{
			hull:   "hunter",
			turret: "railgun",
		},

		{
			hull:   "mammoth",
			turret: "railgun",
		},

		{
			hull:   "scout",
			turret: "railgun",
			cond: func() bool {
				return p.resourceScore(resIron) > p.resourceScore(resGold) &&
					p.resourceScore(resOil) > p.resourceScore(resGold)
			},
		},

		{
			hull:   "scorpion",
			turret: "heavy cannon",
		},

		{
			hull:   "scout",
			turret: "heavy cannon",
			cond: func() bool {
				return p.resourceScore(resGold) > p.resourceScore(resIron) &&
					p.resourceScore(resGold) > p.resourceScore(resOil) &&
					p.Resources.Iron >= 14 &&
					p.Resources.Oil >= 14
			},
		},

		{
			hull:   "fighter",
			turret: "heavy cannon",
			cond: func() bool {
				return p.resourceScore(resIron) < p.resourceScore(resGold) &&
					p.resourceScore(resIron) < p.resourceScore(resOil) &&
					p.Resources.Iron >= 16
			},
		},

		{
			hull:   "viper",
			turret: "lancer",
		},

		{
			hull:   "fighter",
			turret: "lancer",
			cond: func() bool {
				return p.resourceScore(resIron) > p.resourceScore(resOil) &&
					p.Resources.Iron >= 16
			},
		},
	}

	p.favDesigns = options
}

func (p *computerPlayer) selectFavDesign() tankDesign {
	for i := 0; i < 3; i++ {
		option := gemath.RandElem(p.scene.Rand(), p.favDesigns)
		if option.cond != nil && !option.cond() {
			continue
		}
		d := tankDesign{
			Hull:   hullDesigns[option.hull],
			Turret: turretDesigns[option.turret],
		}
		if p.Resources.Contains(d.Price()) {
			return d
		}
	}
	return tankDesign{}
}

func (p *computerPlayer) pickCheapDesign() tankDesign {
	var cheapestDesign tankDesign
	switch {
	case p.Resources.Oil > p.Resources.Iron && p.Resources.Oil > p.Resources.Gold:
		// 4,1,6
		cheapestDesign.Hull = hullDesigns["scout"]
		cheapestDesign.Turret = turretDesigns["gatling gun"]
	case p.Resources.Gold > p.Resources.Oil:
		// 5,4,2
		cheapestDesign.Hull = hullDesigns["viper"]
		cheapestDesign.Turret = turretDesigns["gatling gun"]
	case p.Resources.Oil < p.Resources.Iron && p.Resources.Gold < p.Resources.Iron:
		// 6,3,3
		cheapestDesign.Hull = hullDesigns["viper"]
		cheapestDesign.Turret = turretDesigns["light cannon"]
	default:
		// 5,0,7
		cheapestDesign.Hull = hullDesigns["scout"]
		cheapestDesign.Turret = turretDesigns["light cannon"]
	}
	return cheapestDesign
}

func (p *computerPlayer) selectOwnedSector(pred func(*sector) bool) *sector {
	pick := gemath.RandIndex(p.scene.Rand(), p.ownSectors)
	for i := 0; i < len(p.ownSectors); i++ {
		if pred(p.ownSectors[pick]) {
			return p.ownSectors[pick]
		}
		pick++
		if pick >= len(p.ownSectors) {
			pick = 0
		}
	}
	return nil
}
