package main

import (
	"github.com/quasilyte/gmath"

	"github.com/quasilyte/ge"
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
	"menu.team_2vs2",
	"menu.team_1vs3",
	"menu.team_deathmatch",
	"menu.team_vs_leader",
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
	"menu.slot_empty",
	"menu.slot_player1",
	"menu.slot_player1_gamepad",
	"menu.slot_player2_gamepad",
	"menu.slot_player3_gamepad",
	"menu.slot_player4_gamepad",
	"menu.slot_easy_bot",
	"menu.slot_bot",
}

func (pk playerKind) String() string { return playerKindNames[pk] }

type battleConfig struct {
	teamsMode teamsMode
	players   [4]playerKind
	rules     map[string]bool
}

var battleRules = []string{

	// Bots get unfair bonuses.
	"mode.unfair_bots",

	// Makes starting locations closer to each other.
	"mode.close_combat",

	// Central sectors have no resources.
	"mode.barren_center",

	// All sector incomes are doubled (2 instead of 1).
	"mode.doubled_income",

	// Generate fair resources layout instead of a purely random one.
	"mode.balanced_resources",

	// Players start with two bases instead of one.
	"mode.quick_start",

	// Losing HQ causes the immediate loss.
	"mode.hq_siege",

	// Building battle post fortifications (gauss turrets) is prohibited.
	"mode.no_fortifications",

	// Causes tanks to move and turn slower.
	"mode.mud_terrain",
}

type battleController struct {
	input *MultiHandler

	config battleConfig

	gameState   *gameState
	battleState *battleState
	scene       *ge.Scene

	harvestDelay    float64
	reallianceDelay float64

	winnerAlliance int
	victoryState   bool

	curtainLeft  *ge.Sprite
	curtainRight *ge.Sprite
}

func newBattleController(state *gameState, config battleConfig) *battleController {
	return &battleController{
		gameState:      state,
		input:          state.MenuInput,
		battleState:    newBattleState(),
		config:         config,
		winnerAlliance: -1,
	}
}

func (c *battleController) Init(scene *ge.Scene) {
	c.battleState.DoubledIncome = c.config.rules["mode.doubled_income"]
	c.battleState.FortificationsAllowed = !c.config.rules["mode.no_fortifications"]
	c.battleState.MudTerrain = c.config.rules["mode.mud_terrain"]
	c.battleState.HQDefeat = c.config.rules["mode.hq_siege"]
	c.battleState.UnfairBots = c.config.rules["mode.unfair_bots"]

	bg := ge.NewTiledBackground(scene.Context())
	if c.battleState.MudTerrain {
		bg.Hue = gmath.DegToRad(160)
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
			c.battleState.Players[i].Stats.ID = i
		}
	}

	c.battleState.DeploySectors(c.scene.Rand(), c.config.rules["mode.balanced_resources"])

	if c.config.rules["mode.barren_center"] {
		emptySectors := []int{8, 9, 14, 15}
		for _, id := range emptySectors {
			c.battleState.Sectors[id].Resource = resNone
		}
	}

	startingSectors := []int{0, 18, 0 + 5, 18 + 5}
	extraBaseSectors := []int{1, 19, 4, 22}
	if c.config.rules["mode.close_combat"] {
		startingSectors = []int{7, 20, 16, 3}
		extraBaseSectors = []int{13, 21, 10, 2}
	}
	quickStartEnabled := c.config.rules["mode.quick_start"]
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
				Turret: turretDesigns["gatling_gun"],
			})
			guard.Body.Pos = s.Center().Add(gmath.Vec{Y: 64})
			guard.Body.Rotation = guard.Body.Pos.AngleToPoint(gmath.Vec{X: 1920 / 2, Y: 1080 / 2})
			s.AddTank(guard)
			scene.AddObject(guard)
		}
	}
	for _, s := range c.battleState.Sectors {
		c.scene.AddObjectBelow(s, 1)
	}

	numLocalPlayers := 0
	var singleLocalPlayer *playerData
	for i, pk := range c.config.players {
		s := c.battleState.Sectors[startingSectors[i]]
		p := &c.battleState.Players[i]
		p.IsBot = false
		var object ge.SceneObject
		switch pk {
		case pkEmpty:
			p.Alliance = -1 // Mark as inactive
		case pkLocalPlayer1keyboard:
			numLocalPlayers++
			singleLocalPlayer = p
			object = newLocalPlayer(p, c.gameState.Player1keyboard, s)
		case pkLocalPlayer1:
			numLocalPlayers++
			singleLocalPlayer = p
			object = newLocalPlayer(p, c.gameState.Player1gamepad, s)
		case pkLocalPlayer2:
			numLocalPlayers++
			singleLocalPlayer = p
			object = newLocalPlayer(p, c.gameState.Player2gamepad, s)
		case pkLocalPlayer3:
			numLocalPlayers++
			singleLocalPlayer = p
			object = newLocalPlayer(p, c.gameState.Player3gamepad, s)
		case pkLocalPlayer4:
			numLocalPlayers++
			singleLocalPlayer = p
			object = newLocalPlayer(p, c.gameState.Player4gamepad, s)
		case pkBot:
			object = newComputerPlayer(p, c.battleState, s, false)
			p.IsBot = true
		case pkEasyBot:
			object = newComputerPlayer(p, c.battleState, s, true)
			p.IsBot = true
		}
		if object != nil {
			scene.AddObject(object)
		}
	}
	if numLocalPlayers == 1 {
		c.battleState.SingleLocalPlayer = singleLocalPlayer
	} else {
		c.battleState.SingleLocalPlayer = nil
	}

	scene.Audio().ContinueMusic(AudioMusic)

	{
		window := scene.Context().WindowRect()
		c.curtainLeft = scene.NewRepeatedSprite(ImageMenuBackground, window.Width()/2, window.Height())
		c.curtainLeft.Centered = false
		scene.AddGraphics(c.curtainLeft)
		c.curtainRight = scene.NewRepeatedSprite(ImageMenuBackground, window.Width()/2, window.Height())
		c.curtainRight.Centered = false
		c.curtainRight.Pos.Offset.X += window.Width() / 2
		scene.AddGraphics(c.curtainRight)
	}
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
	if c.curtainLeft != nil {
		c.curtainLeft.Pos.Offset.X -= delta * 512 * 5
		c.curtainRight.Pos.Offset.X += delta * 512 * 5
		if c.curtainLeft.Pos.Offset.X <= -c.curtainLeft.FrameWidth {
			c.curtainLeft.Dispose()
			c.curtainRight.Dispose()
			c.curtainLeft = nil
			c.curtainRight = nil
		}
	}

	if c.winnerAlliance == -1 {
		if winner := c.checkVictory(); winner != -1 {
			c.winnerAlliance = winner
			c.scene.DelayedCall(3, c.onVictoryState)
		}
	}

	if c.battleState.DynamicAlliances {
		c.reallianceDelay = gmath.ClampMin(c.reallianceDelay-delta, 0)
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
			p.Stats.Iron += p.Income.Iron
			p.Stats.Gold += p.Income.Gold
			p.Stats.Oil += p.Income.Oil
			if p.BattleState.UnfairBots && p.IsBot {
				// This unfair bot income is not recorded in the statistics.
				var unfairIncome resourceContainer
				multiplier := 1
				if c.battleState.DoubledIncome {
					multiplier = 2
				}
				kind := resourceKind(c.scene.Rand().IntRange(0, 2))
				unfairIncome.AddOfKind(kind, c.scene.Rand().IntRange(1, 2)*multiplier)
				p.Resources.Add(unfairIncome)
			}
		}
	}

	if c.input.ActionIsJustPressed(ActionExit) {
		c.scene.Audio().PauseCurrentMusic()
		c.scene.Context().ChangeScene(newGameController(c.gameState, c.config))
	}
}

func (c *battleController) calculateIncome() {
	amount := 1
	if c.battleState.DoubledIncome {
		amount = 2
	}
	for _, s := range c.battleState.Sectors {
		if s.Base == nil {
			continue
		}
		s.Base.Player.Income.AddOfKind(s.Resource, amount)
	}
}

func (c *battleController) checkVictory() int {
	winner := -1
	for i := range c.battleState.Players {
		p := &c.battleState.Players[i]
		if p.Alliance == -1 || p.NumBases == 0 {
			continue
		}
		if winner == -1 {
			winner = p.Alliance
			continue
		}
		if winner != p.Alliance {
			return -1
		}
	}
	return winner
}

func (c *battleController) onVictoryState() {
	c.scene.Audio().PauseCurrentMusic()
	c.scene.Audio().PlaySound(AudioCueScreenReset)
	results := battleResult{
		alliance: c.winnerAlliance,
	}
	for i := range c.battleState.Players {
		p := &c.battleState.Players[i]
		if p.Alliance == -1 {
			continue
		}
		stats := p.Stats
		stats.ID = p.ID
		stats.Alliance = p.Alliance
		results.players = append(results.players, stats)
	}
	c.scene.Context().ChangeScene(newResultsController(c.gameState, c.config, results))
}

type playerData struct {
	ID       int
	Alliance int

	NumBases int

	Resources resourceContainer
	Income    resourceContainer

	IsBot bool

	BattleState *battleState

	Stats playerResult
}
