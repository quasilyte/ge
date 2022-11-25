package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

type battleController struct {
	gameState   *gameState
	battleState *battleState
	scene       *ge.Scene
	config      battleConfig

	tutorialText *ge.Label
}

type battleConfig struct {
	levelData  []byte
	numPlayers int
}

func newBattleController(s *gameState, config battleConfig) *battleController {
	return &battleController{
		gameState: s,
		config:    config,
	}
}

func (c *battleController) Init(scene *ge.Scene) {
	c.scene = scene

	scene.Audio().ContinueMusic(AudioMusic)

	bg := ge.NewTiledBackground()
	bg.LoadTileset(scene.Context(), 1920, 896, ImageTiles, RawTilesJSON)
	scene.AddGraphicsBelow(bg, 1)

	gridBackground := scene.NewSprite(ImageGradientGrid)
	gridBackground.Centered = false
	scene.AddGraphics(gridBackground)

	state := newBattleState(c.gameState, c.config.numPlayers)
	state.rand = scene.Rand()
	state.scene = scene
	c.battleState = state

	scene.AddObject(newBattlePanel(state))

	c.initLevel(scene, state, c.config.levelData)
	if c.gameState.gameLevel == 0 {
		c.initTutorial()
	}

	state.updateWalls()
}

func (c *battleController) Update(delta float64) {
	if c.gameState.PlayerInput[0].ActionIsJustPressed(ActionEscape) {
		c.scene.Audio().PauseCurrentMusic()
		if c.battleState.gameState.gameLevel == 0 {
			c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
		} else {
			c.scene.Context().ChangeScene(newPrebattleController(c.gameState))
		}
	}
}

func (c *battleController) initTutorial() {
	bg := ge.NewRect(c.scene.Context(), 1980, 64*4)
	bg.FillColorScale = ge.ColorScale{A: 1}
	bg.Centered = false
	c.scene.AddGraphics(bg)

	h := c.gameState.PlayerInput[0]
	isGamepad := h.GamepadConnected()
	c.tutorialText = c.scene.NewLabel(FontSmall)
	c.tutorialText.Pos.Offset = gmath.Vec{X: 64, Y: 28}
	textLines := []string{
		"welcome to the tutorial level!",
		"",
	}
	if isGamepad {
		textLines = append(textLines, "use dpad or left stick to move.")
		textLines = append(textLines, "(a) fires a primary weapon, (b) activates a special weapon.")
	} else {
		textLines = append(textLines, "use arrows to move.")
		textLines = append(textLines, "[q] fires a primary weapon, [w] activates a special weapon.")
	}
	textLines = append(textLines, "note that your special weapon has limited ammo.")
	textLines = append(textLines, "")
	textLines = append(textLines, "get to the other side of the wall to continue.")

	c.tutorialText.Text = strings.Join(textLines, "\n")
	c.scene.AddGraphics(c.tutorialText)
}

func (c *battleController) onTrigger(name string) {
	switch name {
	case "turret_hint":
		c.scene.Audio().PlaySound(AudioBonus)
		textLines := []string{
			"oh, and another hint for you.",
			"",
			"usually, you can avoid fighting the turrets.",
			"but there are efficient ways to defeat them if you have to.",
			"your flamer special weapon can handle them, just try a correct angle.",
			"",
			"in order to win, destroy the hq, the turret is optional.",
		}
		c.tutorialText.Text = strings.Join(textLines, "\n")

	case "attackers_defeated":
		c.scene.Audio().PlaySound(AudioBonus)
		textLines := []string{
			"nice work!",
			"",
			"you lose when all of your hq facilities are destroyed.",
			"the same holds for your enemy.",
			"it is as simple as that.",
			"",
			"now go and claim your victory!",
		}
		c.tutorialText.Text = strings.Join(textLines, "\n")
		c.scene.DelayedCall(15, func() {
			c.onTrigger("turret_hint")
		})

	case "waypoint1":
		c.scene.Audio().PlaySound(AudioBonus)
		textLines := []string{
			"good!",
			"",
			"the two enemy units will commence attack now.",
			"they will try to destroy your hq.",
			"if they manage to do so, you will lose.",
			"",
			"elliminate the threat.",
		}
		c.tutorialText.Text = strings.Join(textLines, "\n")
		numAttackersLeft := 2
		c.battleState.walkUnitsOfGroup("attackers", func(u *battleUnit) {
			u.bot.program = botHQAttack
			u.bot.delayedAttackCountdown = 10
			u.EventDestroyed.Connect(nil, func(_ *battleUnit) {
				numAttackersLeft--
				if numAttackersLeft == 0 {
					c.onTrigger("attackers_defeated")
				}
			})
		})

	default:
		panic(fmt.Sprintf("unexpected trigger %q", name))
	}
}

func (c *battleController) initLevel(scene *ge.Scene, state *battleState, levelData []byte) {
	mapdata, err := tiled.UnmarshalTileset(scene.LoadRaw(RawMapDateTilesetJSON))
	if err != nil {
		panic(err)
	}
	m, err := tiled.UnmarshalMap(levelData)
	if err != nil {
		panic(err)
	}
	ref := m.Tilesets[0]
	layer := m.Layers[0]
	getAlliance := func(key string) int {
		if strings.HasPrefix(key, "enemy_") {
			return 1
		}
		return 0
	}
	numPlayers := 0
	for _, p := range state.players {
		if p.active {
			numPlayers++
		}
	}
	for _, o := range layer.Objects {
		id := o.GID - ref.FirstGID
		t := mapdata.TileByID(id)
		pos := gmath.Vec{X: float64(o.X + 32), Y: float64(o.Y - 32)}
		switch o.Rotation {
		case 90:
			pos.Y += 64
		case 180:
			pos.X -= 64
			pos.Y += 64
		case 270:
			pos.X -= 64
		}

		switch t.Class {
		case "trigger":
			name := o.GetStringProp("name", "")
			if name == "" {
				panic("trigger name can't be empty")
			}
			t := newTriggerNode(pos, name)
			t.EventActivated.Connect(nil, func(t *triggerNode) {
				c.onTrigger(t.name)
			})
			scene.AddObject(t)

		case "prop":
			mapFlag := o.GetStringProp("flag", "")
			switch mapFlag {
			case "force_base_attack":
				state.forceBaseAttack = true
			default:
				panic(fmt.Sprintf("unexpected map flag: %q", mapFlag))
			}

		case "spawn":
			playerID := o.GetIntProp("player_id", 1)
			playerInfo := state.players[playerID-1]
			if !playerInfo.active {
				break
			}
			pConfig := c.gameState.playerConfig[playerID-1]
			if c.gameState.gameLevel == 0 {
				pConfig = &playerConfig{
					weaponID:      0,
					specialID:     2,
					armorLevel:    2,
					speedLevel:    1,
					rotationLevel: 1,
				}
			}
			playerInfo.spawnPos = pos
			unitConfig := battleUnitConfig{
				battleState:            state,
				playerID:               playerID,
				maxHP:                  (20) + float64(pConfig.armorLevel)*5,
				alliance:               0,
				pos:                    pos,
				image:                  ImagePlayerVessel,
				speed:                  (130) + float64(pConfig.speedLevel)*25,
				rotationTime:           (0.4) - float64(pConfig.rotationLevel)*0.06,
				weapon:                 weaponDesignByName(allWeapons[pConfig.weaponID]),
				weaponReloadMultiplier: 0.8,
				rotation:               gmath.DegToRad(float64(o.Rotation)),
				special:                specialWeaponDesignByName(allSpecials[pConfig.specialID]),
			}
			if boost := boostDesigns[allBoosts[pConfig.boostID]]; boost != nil {
				switch boost.name {
				case "primary":
					unitConfig.weaponReloadMultiplier -= 0.2
				case "special":
					unitConfig.improvedSpecial = true
				case "both":
					unitConfig.weaponReloadMultiplier -= 0.2
					unitConfig.improvedSpecial = true
				}
			}
			u := state.newBattleUnit(unitConfig)
			p := newLocalPlayer(c.gameState.PlayerInput[playerID-1], u)
			scene.AddObject(p)
			scene.AddObject(u)
			if numPlayers != 1 {
				label := newPlayerUnitLayer(ge.Pos{Base: &u.Body.Pos}, strconv.Itoa(playerID))
				scene.AddObject(label)
				u.EventDestroyed.Connect(nil, func(_ *battleUnit) {
					label.Dispose()
				})
			}

		case "enemy_unit", "unit":
			if minNumPlayers := o.GetIntProp("min_players_num", 1); minNumPlayers > numPlayers {
				break
			}
			bodyDesign := standardBodyDesignList[o.GetIntProp("body_design_id", 0)]
			u := state.newBattleUnit(battleUnitConfig{
				group:                  o.GetStringProp("group", ""),
				battleState:            state,
				maxHP:                  bodyDesign.maxHP,
				alliance:               getAlliance(t.Class),
				pos:                    pos,
				image:                  bodyDesign.image,
				rotation:               gmath.DegToRad(float64(o.Rotation)),
				speed:                  bodyDesign.speed,
				rotationTime:           bodyDesign.rotationTime,
				weapon:                 weaponDesignByName(o.GetStringProp("weapon_design", "plasma")),
				weaponReloadMultiplier: 1,
			})
			var prog botProgramKind
			switch progString := o.GetStringProp("program", "delayed_attack"); progString {
			case "guard":
				prog = botGuard
			case "delayed_attack":
				prog = botDelayedAttack
			case "base_attack":
				prog = botBaseAttack
			case "hq_attack":
				prog = botHQAttack
			case "tank_hunt":
				prog = botTankHunt
			case "roam":
				prog = botRoam
			default:
				panic(fmt.Sprintf("unexpected bot program: %s", progString))
			}
			b := newLocalBot(localBotConfig{
				state:       state,
				unit:        u,
				program:     prog,
				attackDelay: o.GetFloatProp("attack_delay", 0.0),
			})
			u.bot = b
			scene.AddObject(b)
			scene.AddObject(u)

		case "enemy_base":
			if minNumPlayers := o.GetIntProp("min_players_num", 1); minNumPlayers > numPlayers {
				break
			}
			scene.AddObject(state.newBattleBase(pos, getAlliance(t.Class), o.GetIntProp("level", 3)))
		case "base":
			scene.AddObject(state.newBattleBase(pos, getAlliance(t.Class), o.GetIntProp("level", 3)))

		case "wall", "enemy_wall":
			scene.AddObject(state.newBattleWall(pos, getAlliance(t.Class)))

		case "factory", "enemy_factory":
			scene.AddObject(state.newBattleFactory(battleFactoryConfig{
				pos:         pos,
				alliance:    getAlliance(t.Class),
				rotation:    gmath.DegToRad(float64(o.Rotation)),
				techLevel:   o.GetIntProp("tech_level", 1),
				maxUnits:    o.GetIntProp("max_units", 1),
				battleState: state,
			}))

		case "bonus_generator":
			scene.AddObject(state.newBattleBonusGenerator(battleBonusGeneratorConfig{
				pos:      pos,
				alliance: getAlliance(t.Class),
				rotation: gmath.DegToRad(float64(o.Rotation)),
			}))

		case "pickup_hp":
			scene.AddObject(newPickupBonus(pos, pickupHP))
		case "pickup_ammo":
			scene.AddObject(newPickupBonus(pos, pickupAmmo))

		case "roam_stop":
			state.addCellInfoAt(pos, cellRoamStop)
		case "dark_tile":
			scene.AddObject(newDarkTile(pos))
			state.addCellInfoAt(pos, cellDarkTile)

		case "mine", "enemy_mine":
			scene.AddObject(newBattleMine(pos, getAlliance(t.Class), -1))

		case "enemy_turret", "enemy_flamer":
			scene.AddObject(state.newBattleFort(battleFortConfig{
				pos:         pos,
				alliance:    getAlliance(t.Class),
				battleState: state,
				rotation:    gmath.DegToRad(float64(o.Rotation)),
				flamer:      t.Class == "enemy_flamer",
			}))
		case "turret", "flamer":
			if maxNumPlayers := o.GetIntProp("max_players_num", 3); maxNumPlayers < numPlayers {
				break
			}
			scene.AddObject(state.newBattleFort(battleFortConfig{
				pos:         pos,
				alliance:    getAlliance(t.Class),
				battleState: state,
				rotation:    gmath.DegToRad(float64(o.Rotation)),
				flamer:      t.Class == "flamer" || t.Class == "enemy_flamer",
			}))
		}
	}
}
