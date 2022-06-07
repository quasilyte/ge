package main

import (
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/resource"
)

type resourceKind int

const (
	resIron resourceKind = iota
	resGold
	resOil
	resCombined
	resNone
)

type resourceContainer struct {
	Iron int
	Gold int
	Oil  int
}

type damageKind int

const (
	damageKinetic damageKind = iota
	// Energy does 60% less damage against buildings.
	damageEnergy
	// Thermal does 20% less damage against vehicles,
	// but 20% more damage to buildings.
	damageThermal
)

func (d damageKind) String() string {
	switch d {
	case damageKinetic:
		return "kinetic"
	case damageEnergy:
		return "energy"
	default:
		return "thermal"
	}
}

func (c *resourceContainer) AddOfKind(kind resourceKind, amount int) {
	switch kind {
	case resIron:
		c.Iron = gemath.ClampMax(c.Iron+amount, 99)
	case resGold:
		c.Gold = gemath.ClampMax(c.Gold+amount, 99)
	case resOil:
		c.Oil = gemath.ClampMax(c.Oil+amount, 99)
	case resCombined:
		c.AddOfKind(resIron, amount)
		c.AddOfKind(resGold, amount)
		c.AddOfKind(resOil, amount)
	}
}

func (c *resourceContainer) Add(other resourceContainer) {
	c.Iron = gemath.ClampMax(c.Iron+other.Iron, 99)
	c.Gold = gemath.ClampMax(c.Gold+other.Gold, 99)
	c.Oil = gemath.ClampMax(c.Oil+other.Oil, 99)
}

func (c *resourceContainer) Sub(other resourceContainer) {
	c.Iron -= other.Iron
	c.Gold -= other.Gold
	c.Oil -= other.Oil
}

func (c resourceContainer) Contains(other resourceContainer) bool {
	return c.Iron >= other.Iron && c.Gold >= other.Gold && c.Oil >= other.Oil
}

type hullSize int

const (
	hullSmall hullSize = iota
	hullMedium
	hullLarge
)

type tankDesign struct {
	Hull   *hullDesign
	Turret *turretDesign
}

func (d *tankDesign) IsEmpty() bool {
	return d.Hull == nil
}

func (d tankDesign) Price() resourceContainer {
	var total resourceContainer
	total.Add(d.Hull.Price)
	total.Add(d.Turret.Price)
	return total
}

func (d tankDesign) ProductionTime() float64 {
	return d.Hull.Production + d.Turret.Production
}

type hullDesign struct {
	ID uint8

	Name string

	Size    hullSize
	OriginX float64

	HP float64

	Speed         float64
	RotationSpeed gemath.Rad

	Price      resourceContainer
	Production float64

	Image resource.ImageID
}

type turretDesign struct {
	ID uint8

	Name string

	Price      resourceContainer
	Production float64

	// Only for battle post installations.
	HP float64

	HPBonus      float64
	SpeedPenalty float64

	FireRange  float64
	Reload     float64
	Damage     float64
	DamageKind damageKind

	ProjectileSpeed float64

	RotationSpeed gemath.Rad

	Image     resource.ImageID
	AmmoImage resource.ImageID
	Sound     resource.AudioID
}

func init() {
	for i := range hullDesignList {
		hullDesignList[i].ID = uint8(i)
	}
	for i := range turretDesignList {
		turretDesignList[i].ID = uint8(i)
	}

	for _, d := range hullDesignList {
		hullDesigns[d.Name] = d
	}
	for _, d := range turretDesignList {
		turretDesigns[d.Name] = d
	}
	for _, d := range turretDesignList {
		if d.Name == "builder" {
			continue
		}
		turretDesignListNoBuilder = append(turretDesignListNoBuilder, d)
	}
}

var hullDesigns = map[string]*hullDesign{}

var hullDesignList = []*hullDesign{
	{
		Name:          "viper",
		HP:            140,
		Size:          hullSmall,
		Speed:         130,
		RotationSpeed: 1.3,
		Price:         resourceContainer{Iron: 3, Gold: 2, Oil: 3},
		Production:    2.5,
		Image:         ImageHullViper,
		OriginX:       -2,
	},

	{
		Name:          "scout",
		HP:            90,
		Size:          hullSmall,
		Speed:         155,
		RotationSpeed: 1.9,
		Price:         resourceContainer{Iron: 3, Oil: 5},
		Production:    1.5,
		Image:         ImageHullScout,
		OriginX:       -4,
	},

	{
		Name:          "hunter",
		HP:            185,
		Size:          hullMedium,
		Speed:         120,
		RotationSpeed: 3.25,
		Price:         resourceContainer{Iron: 4, Gold: 4, Oil: 4},
		Production:    3,
		Image:         ImageHullHunter,
	},

	{
		Name:          "scorpion",
		HP:            160,
		Size:          hullSmall,
		Speed:         140,
		RotationSpeed: 2.75,
		Price:         resourceContainer{Iron: 4, Gold: 7, Oil: 2},
		Production:    7,
		Image:         ImageHullScorpion,
	},

	{
		Name:          "fighter",
		HP:            230,
		Size:          hullMedium,
		Speed:         110,
		RotationSpeed: 1.1,
		Price:         resourceContainer{Iron: 11, Gold: 3},
		Production:    6.5,
		Image:         ImageHullFighter,
	},

	{
		Name:          "mammoth",
		HP:            600,
		Size:          hullLarge,
		Speed:         95,
		RotationSpeed: 0.8,
		Price:         resourceContainer{Iron: 15, Gold: 4, Oil: 5},
		Production:    16,
		Image:         ImageHullMammoth,
		OriginX:       4,
	},
}

var gaussTurretDesign = &turretDesign{
	Name:            "gauss",
	HP:              500,
	FireRange:       295,
	ProjectileSpeed: 550,
	RotationSpeed:   2,
	Reload:          2.2,
	Damage:          20,
	DamageKind:      damageEnergy,
	Price:           resourceContainer{Iron: 2, Gold: 4, Oil: 7},
	AmmoImage:       ImageAmmoGauss,
	Image:           ImageTurretGauss,
	Sound:           AudioGauss,
}

var ionTurretDesign = &turretDesign{
	Name:            "ion",
	HP:              700,
	FireRange:       350,
	ProjectileSpeed: 500,
	RotationSpeed:   0.5,
	Reload:          0.4,
	Damage:          15,
	DamageKind:      damageEnergy,
	AmmoImage:       ImageAmmoIon,
	Image:           ImageTurretIon,
	Sound:           AudioIon,
}

var turretDesigns = map[string]*turretDesign{}

var turretDesignListNoBuilder []*turretDesign

var turretDesignList = []*turretDesign{
	{
		Name:         "builder",
		SpeedPenalty: 25,
		Price:        resourceContainer{Gold: 4, Oil: 5},
		Production:   10,
		Image:        ImageTurretBuilder,
	},

	{
		Name:            "gatling gun",
		HPBonus:         60,
		FireRange:       200,
		ProjectileSpeed: 600,
		RotationSpeed:   1.3,
		Reload:          0.9,
		Damage:          8,
		DamageKind:      damageKinetic,
		Price:           resourceContainer{Iron: 1, Gold: 2},
		Production:      2.5,
		AmmoImage:       ImageAmmoGatlingGun,
		Image:           ImageTurretGatlingGun,
		Sound:           AudioGatlingGun,
	},

	{
		Name:            "light cannon",
		HPBonus:         40,
		FireRange:       270,
		ProjectileSpeed: 520,
		RotationSpeed:   2,
		Reload:          0.9,
		Damage:          12,
		DamageKind:      damageKinetic,
		Price:           resourceContainer{Iron: 3, Oil: 1},
		Production:      6.5,
		AmmoImage:       ImageAmmoMediumCannon,
		Image:           ImageTurretLightCannon,
		Sound:           AudioLightCannon,
	},

	{
		Name:            "dual cannon",
		HPBonus:         140,
		FireRange:       260,
		SpeedPenalty:    20,
		ProjectileSpeed: 500,
		RotationSpeed:   1.4,
		Reload:          1.5,
		Damage:          40,
		DamageKind:      damageKinetic,
		Price:           resourceContainer{Iron: 4, Gold: 3, Oil: 2},
		Production:      15,
		AmmoImage:       ImageAmmoDualCannon,
		Image:           ImageTurretDualCannon,
		Sound:           AudioDualCannon,
	},

	{
		Name:            "heavy cannon",
		HPBonus:         200,
		SpeedPenalty:    20,
		ProjectileSpeed: 500,
		FireRange:       320,
		RotationSpeed:   1.5,
		Reload:          1.4,
		Damage:          35,
		DamageKind:      damageKinetic,
		Price:           resourceContainer{Iron: 7, Oil: 3},
		Production:      8,
		AmmoImage:       ImageAmmoMediumCannon,
		Image:           ImageTurretHeavyCannon,
		Sound:           AudioHeavyCannon,
	},

	{
		Name:          "railgun",
		HPBonus:       45,
		SpeedPenalty:  5,
		FireRange:     290,
		RotationSpeed: 2.2,
		Reload:        1.8,
		Damage:        50,
		DamageKind:    damageEnergy,
		Price:         resourceContainer{Gold: 6, Oil: 1},
		Production:    6,
		Image:         ImageTurretRailgun,
		Sound:         AudioRailgun,
	},

	{
		Name:            "lancer",
		SpeedPenalty:    45,
		ProjectileSpeed: 600,
		FireRange:       370,
		RotationSpeed:   0.7,
		Reload:          2.6,
		Damage:          55,
		DamageKind:      damageThermal,
		Price:           resourceContainer{Iron: 1, Gold: 2, Oil: 9},
		Production:      15,
		AmmoImage:       ImageAmmoLancer,
		Image:           ImageTurretLancer,
		Sound:           AudioLancer,
	},
}
