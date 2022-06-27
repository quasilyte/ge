package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gemath"
	"github.com/quasilyte/ge/input"
)

const harvestDelaySeconds = 5.0

type teamsMode int

const (
	teams2vs2 teamsMode = iota
	teams1vs3
	teamsDeathmatch
	teamsLeader
)

var teamsModeNames = []string{
	"2 VS 2",
	"1 VS 3",
	"DEATHMATCH",
	"VS LEADER",
}

func (t teamsMode) String() string { return teamsModeNames[t] }

type playerKind int

const (
	pkEmpty playerKind = iota
	pkLocalPlayer1keyboard
	pkLocalPlayer1
	pkLocalPlayer2
	pkLocalPlayer3
	pkLocalPlayer4
	pkEasyBot
	pkBot
)

var playerKindNames = []string{
	"EMPTY",
	"PLAYER 1",
	"PLAYER 1 gamepad",
	"PLAYER 2 gamepad",
	"PLAYER 3 gamepad",
	"PLAYER 4 gamepad",
	"EASY BOT",
	"BOT",
}

func (pk playerKind) String() string { return playerKindNames[pk] }

type battleConfig struct {
	teamsMode teamsMode
	players   [4]playerKind
	rules     map[string]bool
}

var battleRules = []string{
	// Makes starting locations closer to each other.
	"Close Combat",

	// Central sectors have no resources.
	"Barren Center",

	// All sector incomes are doubled (2 instead of 1).
	"Doubled Income",

	// Generate fair resources layout instead of a purely random one.
	"Balanced Resources",

	// Players start with two bases instead of one.
	"Quick Start",

	// Losing HQ causes the immediate loss.
	"HQ Siege",

	// Building battle post fortifications (gauss turrets) is prohibited.
	"No Fortifications",

	// Causes tanks to move and turn slower.
	"Mud Terrain",
}

type battleController struct {
	input *input.MultiHandler

	config battleConfig

	gameState   *gameState
	battleState *battleState
	scene       *ge.Scene

	harvestDelay    float64
	reallianceDelay float64
}

func newBattleController(state *gameState, config battleConfig) *battleController {
	return &battleController{
		gameState:   state,
		input:       state.MenuInput,
		battleState: newBattleState(),
		config:      config,
	}
}

func (c *battleController) Init(scene *ge.Scene) {
	ctx := scene.Context()

	c.battleState.DoubledIncome = c.config.rules["Doubled Income"]
	c.battleState.FortificationsAllowed = !c.config.rules["No Fortifications"]
	c.battleState.MudTerrain = c.config.rules["Mud Terrain"]
	c.battleState.HQDefeat = c.config.rules["HQ Siege"]

	bg := ge.NewTiledBackground()
	if c.battleState.MudTerrain {
		bg.Hue = gemath.DegToRad(160)
	}
	bg.LoadTileset(scene.Context(), 1920, 1080, ImageBackgroundTiles, RawTilesJSON)
	scene.AddGraphicsBelow(bg, 1)

	c.scene = scene

	{
		var alliances [4]int
		switch c.config.teamsMode {
		case teams2vs2:
			alliances = [4]int{0, 0, 1, 1}
		case teams1vs3:
			alliances = [4]int{0, 1, 1, 1}
		case teamsDeathmatch:
			alliances = [4]int{0, 1, 2, 3}
		case teamsLeader:
			c.battleState.DynamicAlliances = true
			alliances = [4]int{0, 1, 2, 3}
		default:
			panic("unexpected option")
		}
		for i := range c.battleState.Players {
			c.battleState.Players[i].Alliance = alliances[i]
		}
	}

	c.battleState.DeploySectors(c.scene.Rand(), c.config.rules["Balanced Resources"])

	if c.config.rules["Barren Center"] {
		emptySectors := []int{8, 9, 14, 15}
		for _, id := range emptySectors {
			c.battleState.Sectors[id].Resource = resNone
		}
	}

	startingSectors := []int{0, 18, 0 + 5, 18 + 5}
	extraBaseSectors := []int{1, 19, 4, 22}
	if c.config.rules["Close Combat"] {
		startingSectors = []int{7, 20, 16, 3}
		extraBaseSectors = []int{13, 21, 10, 2}
	}
	quickStartEnabled := c.config.rules["Quick Start"]
	for i, pk := range c.config.players {
		if pk == pkEmpty {
			continue
		}
		c.battleState.Sectors[startingSectors[i]].Resource = resCombined
		if quickStartEnabled {
			p := &c.battleState.Players[i]
			s := c.battleState.Sectors[extraBaseSectors[i]]
			bp := c.battleState.NewBattlePost(p, s.Center(), nil)
			s.AssignBase(bp)
			scene.AddObject(bp)

			guard := c.battleState.NewBattleTank(bp.Player, tankDesign{
				Hull:   hullDesigns["scout"],
				Turret: turretDesigns["gatling gun"],
			})
			guard.Body.Pos = s.Center().Add(gemath.Vec{Y: 64})
			guard.Body.Rotation = guard.Body.Pos.AngleToPoint(gemath.Vec{X: 1920 / 2, Y: 1080 / 2})
			s.AddTank(guard)
			scene.AddObject(guard)
		}
	}
	for _, s := range c.battleState.Sectors {
		c.scene.AddObjectBelow(s, 1)
	}

	for i, pk := range c.config.players {
		s := c.battleState.Sectors[startingSectors[i]]
		p := &c.battleState.Players[i]
		var object ge.SceneObject
		switch pk {
		case pkEmpty:
			// Do nothing.
		case pkLocalPlayer1keyboard:
			object = newLocalPlayer(p, c.gameState.Player1keyboard, s)
		case pkLocalPlayer1:
			object = newLocalPlayer(p, c.gameState.Player1gamepad, s)
		case pkLocalPlayer2:
			object = newLocalPlayer(p, c.gameState.Player2gamepad, s)
		case pkLocalPlayer3:
			object = newLocalPlayer(p, c.gameState.Player3gamepad, s)
		case pkLocalPlayer4:
			object = newLocalPlayer(p, c.gameState.Player4gamepad, s)
		case pkBot:
			object = newComputerPlayer(p, c.battleState, s, false)
		case pkEasyBot:
			object = newComputerPlayer(p, c.battleState, s, true)
		}
		if object != nil {
			scene.AddObject(object)
		}
	}

	ctx.Audio.ContinueMusic(AudioMusic)
}

func (c *battleController) recalculateAlliances() {
	var numSectors [4]int
	for _, s := range c.battleState.Sectors {
		if s.Base == nil {
			continue
		}
		numSectors[s.Base.Player.ID]++
	}
	maxSectors := 0
	for _, num := range &numSectors {
		if num > maxSectors {
			maxSectors = num
		}
	}
	alliance := 1
	for i := range c.battleState.Players {
		p := &c.battleState.Players[i]
		if numSectors[p.ID] == maxSectors {
			p.Alliance = alliance
			alliance++
		} else {
			p.Alliance = 0
		}
	}
}

func (c *battleController) Update(delta float64) {
	if c.battleState.DynamicAlliances {
		c.reallianceDelay = gemath.ClampMin(c.reallianceDelay-delta, 0)
		if c.reallianceDelay == 0 {
			c.recalculateAlliances()
			c.reallianceDelay = c.scene.Rand().FloatRange(15, 30)
		}
	}

	c.harvestDelay -= delta
	if c.harvestDelay < 0 {
		for i := range c.battleState.Players {
			c.battleState.Players[i].Income = resourceContainer{}
		}
		c.harvestDelay = harvestDelaySeconds - c.harvestDelay
		c.calculateIncome()
		for i := range c.battleState.Players {
			p := &c.battleState.Players[i]
			p.Resources.Add(p.Income)
		}
	}

	if c.input.ActionIsJustPressed(ActionExit) {
		c.scene.Audio().PauseCurrentMusic()
		c.scene.Context().ChangeScene("game", newGameController(c.gameState, c.config))
	}
}

func (c *battleController) calculateIncome() {
	amount := 1
	if c.config.rules["Doubled Income"] {
		amount = 2
	}
	for _, s := range c.battleState.Sectors {
		if s.Base == nil {
			continue
		}
		s.Base.Player.Income.AddOfKind(s.Resource, amount)
	}
}

type playerData struct {
	ID       int
	Alliance int

	Resources resourceContainer
	Income    resourceContainer

	BattleState *battleState
}
