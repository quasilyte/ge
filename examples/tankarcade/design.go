package main

import (
	"fmt"
	"math"

	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/ge/xmaps"
	"github.com/quasilyte/gmath"
)

func init() {
	hurricaneBouncedDesign = &weaponDesign{}
	*hurricaneBouncedDesign = *hurricaneDesign
	hurricaneBouncedDesign.minRange = 64 * 1
	hurricaneBouncedDesign.fireSound = AudioNone

	*improvedRegeneratorDesign = *regeneratorDesign.extra
	improvedRegeneratorDesign.damage = -25

	allBoosts = append(allBoosts, "none")
	allBoosts = append(allBoosts, xmaps.KeysSortedByValue(boostDesigns, func(d1, d2 *boostDesign) bool {
		return d1.cost < d2.cost
	})...)
	for k, d := range boostDesigns {
		d.name = k
	}
}

type standardBodyDesign struct {
	image          resource.ImageID
	maxHP          float64
	speed          float64
	rotationTime   float64
	productionTime float64
	reloadModified float64
	techLevel      int
}

var standardBodyDesignList = []*standardBodyDesign{
	{
		image:          ImageVessel4,
		maxHP:          5,
		speed:          180,
		rotationTime:   0.4,
		productionTime: 10,
		reloadModified: 1,
		techLevel:      1,
	},
	{
		image:          ImageVessel1,
		maxHP:          20,
		speed:          150,
		rotationTime:   0.25,
		productionTime: 20,
		reloadModified: 1.1,
		techLevel:      2,
	},
	{
		image:          ImageVessel2,
		maxHP:          10,
		speed:          160,
		rotationTime:   0.2,
		productionTime: 15,
		reloadModified: 0.9,
		techLevel:      2,
	},
	{
		image:          ImageVessel3,
		maxHP:          25,
		speed:          100,
		rotationTime:   0.4,
		productionTime: 25,
		reloadModified: 1.1,
		techLevel:      3,
	},
	{
		image:          ImageVessel5,
		maxHP:          45,
		speed:          70,
		rotationTime:   0.7,
		productionTime: 30,
		reloadModified: 1,
		techLevel:      4,
	},
	{
		image:          ImageVessel6,
		maxHP:          40,
		speed:          130,
		rotationTime:   0.4,
		productionTime: 30,
		reloadModified: 0.8,
		techLevel:      6,
	},
}

type specialWeaponDesign struct {
	ammo int

	extra *weaponDesign
}

type weaponDesign struct {
	cost                       int
	fireSound                  resource.AudioID
	hitSound                   resource.AudioID
	projectileImage            resource.ImageID
	projectileExplosion        resource.ImageID
	projectileExplosionScale   float64
	projectileExplosionHue     gmath.Rad
	projectileExplosionRotates bool
	projectileRotationSpeed    gmath.Rad
	projectileSpeed            float64
	maxRange                   float64
	minRange                   float64
	reload                     float64
	damage                     float64
	canReflect                 bool
}

var simpleWeaponDesignList = []*weaponDesign{
	plasmaCannonDesign,
	gaussGunDesign,
	assaultCannonDesign,
}

var advancedWeaponDesignList = []*weaponDesign{
	sonicCannonDesign,
	howitzerDesign,
	hurricaneDesign,
}

var allWeapons = []string{
	"plasma",
	"gauss",
	"assault",
	"hurricane",
	"sonic",
	"howitzer",
	"railgun",
}

var allSpecials = []string{
	"none",
	"regenerator",
	"flamer",
	"launcher",
	"minelayer",
	"shield",
}

var allBoosts []string

type boostDesign struct {
	name string
	cost int
}

var boostDesigns = map[string]*boostDesign{
	"special": {cost: 10},
	"primary": {cost: 20},
	"both":    {cost: 30},
}

func specialWeaponDesignByName(name string) *specialWeaponDesign {
	switch name {
	case "none":
		return nil
	case "flamer":
		return flamethrowerDesign
	case "launcher":
		return rocketLauncherDesign
	case "minelayer":
		return mineLayerDesign
	case "regenerator":
		return regeneratorDesign
	case "shield":
		return shieldDesign
	default:
		panic(fmt.Sprintf("unexpected special weapon design name: %s", name))
	}
}

func weaponDesignByName(name string) *weaponDesign {
	switch name {
	case "plasma":
		return plasmaCannonDesign
	case "gauss":
		return gaussGunDesign
	case "assault":
		return assaultCannonDesign
	case "sonic":
		return sonicCannonDesign
	case "howitzer":
		return howitzerDesign
	case "hurricane":
		return hurricaneDesign
	case "railgun":
		return railgunDesign
	default:
		panic(fmt.Sprintf("unexpected weapon design name: %s", name))
	}
}

var plasmaCannonDesign = &weaponDesign{
	cost:                     0,
	fireSound:                AudioLaser1,
	hitSound:                 AudioExplosion1,
	projectileImage:          ImageBullet1,
	projectileExplosion:      ImageExplosion1,
	projectileExplosionScale: 1,
	projectileSpeed:          400,
	maxRange:                 64 * 4,
	reload:                   1.5,
	damage:                   6,
	canReflect:               true,
}

var gaussGunDesign = &weaponDesign{
	cost:                     5,
	fireSound:                AudioLaser2,
	hitSound:                 AudioExplosion2,
	projectileImage:          ImageBullet2,
	projectileExplosion:      ImageExplosion2,
	projectileExplosionScale: 1,
	projectileSpeed:          300,
	maxRange:                 64 * 6,
	reload:                   1,
	damage:                   5,
	canReflect:               true,
}

var assaultCannonDesign = &weaponDesign{
	cost:                     5,
	fireSound:                AudioLaser3,
	hitSound:                 AudioExplosion3,
	projectileImage:          ImageBullet3,
	projectileExplosion:      ImageExplosion3,
	projectileExplosionScale: 0.8,
	projectileSpeed:          500,
	maxRange:                 64 * 3,
	reload:                   0.5,
	damage:                   4,
	canReflect:               true,
}

var sonicCannonDesign = &weaponDesign{
	cost:                     20,
	fireSound:                AudioLaser5,
	hitSound:                 AudioExplosion6,
	projectileImage:          ImageBullet5,
	projectileExplosion:      ImageExplosion2,
	projectileExplosionScale: 1.2,
	projectileExplosionHue:   2.8,
	projectileSpeed:          550,
	maxRange:                 64 * 5,
	reload:                   1.7,
	damage:                   13,
	canReflect:               true,
}

var howitzerDesign = &weaponDesign{
	cost:                       25,
	fireSound:                  AudioHowitz,
	hitSound:                   AudioExplosion7,
	projectileImage:            ImageBullet6,
	projectileExplosion:        ImageExplosion3,
	projectileExplosionScale:   1,
	projectileExplosionHue:     math.Pi,
	projectileExplosionRotates: true,
	projectileSpeed:            350,
	maxRange:                   64 * 7,
	minRange:                   64 * 2,
	reload:                     2,
	damage:                     10,
	canReflect:                 true,
}

var hurricaneDesign = &weaponDesign{
	cost:                       10,
	fireSound:                  AudioHurricane,
	hitSound:                   AudioExplosion1,
	projectileImage:            ImageBullet7,
	projectileExplosion:        ImageExplosion1,
	projectileExplosionHue:     -1.4,
	projectileExplosionScale:   1,
	projectileExplosionRotates: true,
	projectileRotationSpeed:    36,
	projectileSpeed:            550,
	maxRange:                   64 * 2,
	reload:                     0.6,
	damage:                     4,
}

var hurricaneBouncedDesign *weaponDesign

var railgunDesign = &weaponDesign{
	cost:                       40,
	fireSound:                  AudioRailgun,
	projectileExplosion:        ImageExplosion2,
	projectileExplosionScale:   1,
	projectileExplosionHue:     -2.1,
	projectileExplosionRotates: true,
	projectileRotationSpeed:    36,
	maxRange:                   64 * 4,
	reload:                     1.7,
	damage:                     15,
}

var turretWeaponDesign = &weaponDesign{
	fireSound:                AudioLaser4,
	hitSound:                 AudioExplosion3,
	projectileImage:          ImageBullet4,
	projectileExplosion:      ImageExplosion3,
	projectileExplosionScale: 1,
	projectileSpeed:          450,
	maxRange:                 64 * 6,
	reload:                   2.5,
	damage:                   9,
	canReflect:               true,
}

var flameTurretWeaponDesign = &weaponDesign{
	fireSound:                AudioFlamethrower,
	hitSound:                 AudioExplosion5,
	projectileImage:          ImageFlamethrower,
	projectileExplosion:      ImageExplosion3,
	projectileExplosionScale: 1,
	projectileSpeed:          300,
	maxRange:                 64 * 1,
	reload:                   2.5,
	damage:                   10,
}

var shieldDesign = &specialWeaponDesign{
	ammo: 6,
	extra: &weaponDesign{
		cost:      10,
		fireSound: AudioShieldActivated,
		reload:    1.5,
	},
}

var regeneratorDesign = &specialWeaponDesign{
	ammo: 4,
	extra: &weaponDesign{
		cost:                     5,
		fireSound:                AudioRegenerator,
		hitSound:                 AudioRegeneratorExplosion,
		projectileImage:          ImageRegenerator,
		projectileExplosion:      ImageRegeneratorExplosion,
		projectileExplosionScale: 1,
		projectileSpeed:          250,
		projectileRotationSpeed:  10,
		maxRange:                 64 * 3,
		reload:                   3,
		damage:                   -20,
	},
}

var improvedRegeneratorDesign = &weaponDesign{}

var flamethrowerDesign = &specialWeaponDesign{
	ammo: 4,
	extra: &weaponDesign{
		cost:                     10,
		fireSound:                AudioFlamethrower,
		hitSound:                 AudioExplosion5,
		projectileImage:          ImageFlamethrower,
		projectileExplosion:      ImageExplosion3,
		projectileExplosionScale: 1,
		projectileSpeed:          250,
		maxRange:                 64 * 1,
		reload:                   2.5,
		damage:                   20,
	},
}

var rocketLauncherDesign = &specialWeaponDesign{
	ammo: 2,
	extra: &weaponDesign{
		cost:                     10,
		fireSound:                AudioMissile,
		hitSound:                 AudioExplosion4,
		projectileImage:          ImageMissile,
		projectileExplosion:      ImageExplosion3,
		projectileExplosionScale: 1,
		projectileSpeed:          450,
		maxRange:                 64 * 7,
		reload:                   2,
		damage:                   25,
	},
}

var mineLayerDesign = &specialWeaponDesign{
	ammo: 4,
	extra: &weaponDesign{
		cost:      10,
		fireSound: AudioMine,
		hitSound:  AudioExplosion4,
		reload:    1.5,
		damage:    30,
	},
}
