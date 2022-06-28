package main

import (
	"embed"
	"io"
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/resource"

	_ "image/png"
)

//go:embed assets/*
var gameAssets embed.FS

const (
	ActionSectorLeft input.Action = iota
	ActionSectorRight
	ActionSectorUp
	ActionSectorDown
	ActionCancel
	ActionOpenMenu
	ActionConfirm
	ActionNextItem
	ActionPrevItem
	ActionNextCategory
	ActionPrevCategory
	ActionFortify
	ActionExit
)

const (
	ImageHullViper resource.ImageID = iota
	ImageHullScout
	ImageHullHunter
	ImageHullFighter
	ImageHullScorpion
	ImageHullMammoth
	ImageTurretBuilder
	ImageTurretGatlingGun
	ImageTurretLightCannon
	ImageTurretDualCannon
	ImageTurretHeavyCannon
	ImageTurretRailgun
	ImageTurretLancer
	ImageTurretGauss
	ImageTurretIon
	ImageBattlePost
	ImageAmmoGatlingGun
	ImageAmmoLightCannon
	ImageAmmoMediumCannon
	ImageAmmoDualCannon
	ImageAmmoLancer
	ImageAmmoGauss
	ImageAmmoIon
	ImageExplosion
	ImageBackgroundTiles
	ImageSectorSelector
	ImageUnitSelector
	ImageGrid
	ImageIronResourceIcon
	ImageGoldResourceIcon
	ImageOilResourceIcon
	ImageCombinedResourceIcon
	ImageResourceRow
	ImagePopupBuildTank
	ImageMenuBackground
	ImageMenuButton
	ImageMenuSelectButton
	ImageMenuCheckboxButton
	ImageMenuSlideLeft
)

const (
	AudioGatlingGun resource.AudioID = iota
	AudioLightCannon
	AudioDualCannon
	AudioHeavyCannon
	AudioRailgun
	AudioLancer
	AudioGauss
	AudioIon
	AudioCueError
	AudioCueSendUnits
	AudioCueProductionStarted
	AudioCueProductionCompleted
	AudioCueConstructionCompleted
	AudioCueScreenReset
	AudioMusic
)

const (
	FontSmall resource.FontID = iota
	FontDescription
	FontBig
)

const (
	RawTilesJSON resource.RawID = iota
)

func main() {
	ctx := ge.NewContext()
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.WindowTitle = "Tanks"
	ctx.WindowWidth = 1920
	ctx.WindowHeight = 1080
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = func(path string) io.ReadCloser {
		f, err := gameAssets.Open("assets/" + path)
		if err != nil {
			ctx.OnCriticalError(err)
		}
		return f
	}

	state := &gameState{}

	// Bind controls.
	gamepadKeymap := input.MakeKeymap(map[input.Action]input.Key{
		ActionSectorLeft:   input.KeyGamepadLeft,
		ActionSectorRight:  input.KeyGamepadRight,
		ActionSectorDown:   input.KeyGamepadDown,
		ActionSectorUp:     input.KeyGamepadUp,
		ActionConfirm:      input.KeyGamepadA,
		ActionOpenMenu:     input.KeyGamepadX,
		ActionCancel:       input.KeyGamepadB,
		ActionPrevItem:     input.KeyGamepadLeft,
		ActionNextItem:     input.KeyGamepadRight,
		ActionNextCategory: input.KeyGamepadDown,
		ActionPrevCategory: input.KeyGamepadUp,
		ActionFortify:      input.KeyGamepadY,
		ActionExit:         input.KeyGamepadStart,
	})
	keyboardKeymap := input.MakeKeymap(map[input.Action]input.Key{
		ActionSectorLeft:   input.KeyA,
		ActionSectorRight:  input.KeyD,
		ActionSectorDown:   input.KeyS,
		ActionSectorUp:     input.KeyW,
		ActionConfirm:      input.KeySpace,
		ActionOpenMenu:     input.KeyEnter,
		ActionCancel:       input.KeyQ,
		ActionPrevItem:     input.KeyA,
		ActionNextItem:     input.KeyD,
		ActionNextCategory: input.KeyS,
		ActionPrevCategory: input.KeyW,
		ActionFortify:      input.KeyE,
		ActionExit:         input.KeyEscape,
	})
	state.Player1keyboard = ctx.Input.NewHandler(0, keyboardKeymap)
	state.Player1gamepad = ctx.Input.NewHandler(0, gamepadKeymap)
	state.Player2gamepad = ctx.Input.NewHandler(1, gamepadKeymap)
	state.Player3gamepad = ctx.Input.NewHandler(2, gamepadKeymap)
	state.Player4gamepad = ctx.Input.NewHandler(3, gamepadKeymap)

	state.MenuInput = ctx.Input.NewMultiHandler()
	state.MenuInput.AddHandler(state.Player1gamepad)
	state.MenuInput.AddHandler(state.Player1keyboard)

	// Associate audio resources.
	audioResources := map[resource.AudioID]resource.Audio{
		AudioGatlingGun:               {Path: "sounds/gatling_gun.wav", Volume: -0.5},
		AudioLightCannon:              {Path: "sounds/light_cannon.wav", Volume: -0.4},
		AudioDualCannon:               {Path: "sounds/dual_cannon.wav", Volume: -0.3},
		AudioHeavyCannon:              {Path: "sounds/heavy_cannon.wav", Volume: -0.75},
		AudioRailgun:                  {Path: "sounds/railgun.wav", Volume: -0.5},
		AudioLancer:                   {Path: "sounds/lancer.wav", Volume: -0.75},
		AudioGauss:                    {Path: "sounds/gauss.wav", Volume: -0.5},
		AudioIon:                      {Path: "sounds/ion.wav", Volume: -0.5},
		AudioCueError:                 {Path: "sounds/error.wav", Volume: -0.45},
		AudioCueSendUnits:             {Path: "sounds/send_units.wav", Volume: -0.3},
		AudioCueProductionStarted:     {Path: "sounds/production_started.wav", Volume: -0.1},
		AudioCueProductionCompleted:   {Path: "sounds/production_completed.wav", Volume: +0.1},
		AudioCueConstructionCompleted: {Path: "sounds/construction_completed.wav"},
		AudioCueScreenReset:           {Path: "sounds/screen_reset.wav"},
		AudioMusic:                    {Path: "sounds/music.ogg"},
	}
	for id, res := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, res)
		ctx.Loader.PreloadAudio(id)
	}

	// Associate image resources.
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageHullViper:            {Path: "hull_viper.png"},
		ImageHullScout:            {Path: "hull_scout.png"},
		ImageHullHunter:           {Path: "hull_hunter.png"},
		ImageHullFighter:          {Path: "hull_fighter.png"},
		ImageHullScorpion:         {Path: "hull_scorpion.png"},
		ImageHullMammoth:          {Path: "hull_mammoth.png"},
		ImageTurretBuilder:        {Path: "turret_builder.png"},
		ImageTurretGatlingGun:     {Path: "turret_gatling_gun.png"},
		ImageTurretLightCannon:    {Path: "turret_light_cannon.png"},
		ImageTurretDualCannon:     {Path: "turret_dual_cannon.png"},
		ImageTurretHeavyCannon:    {Path: "turret_heavy_cannon.png"},
		ImageTurretRailgun:        {Path: "turret_railgun.png"},
		ImageTurretLancer:         {Path: "turret_lancer.png"},
		ImageTurretGauss:          {Path: "turret_gauss.png"},
		ImageTurretIon:            {Path: "turret_ion.png"},
		ImageBattlePost:           {Path: "battle_post.png"},
		ImageAmmoGatlingGun:       {Path: "ammo_gatling_gun.png"},
		ImageAmmoLightCannon:      {Path: "ammo_light_cannon.png"},
		ImageAmmoMediumCannon:     {Path: "ammo_medium_cannon.png"},
		ImageAmmoDualCannon:       {Path: "ammo_dual_cannon.png"},
		ImageAmmoLancer:           {Path: "ammo_lancer.png"},
		ImageAmmoGauss:            {Path: "ammo_gauss.png"},
		ImageAmmoIon:              {Path: "ammo_ion.png"},
		ImageExplosion:            {Path: "explosion.png"},
		ImageBackgroundTiles:      {Path: "tiles.png"},
		ImageSectorSelector:       {Path: "sector_selector.png"},
		ImageUnitSelector:         {Path: "unit_selector.png"},
		ImageGrid:                 {Path: "grid.png"},
		ImageIronResourceIcon:     {Path: "resource_iron.png"},
		ImageGoldResourceIcon:     {Path: "resource_gold.png"},
		ImageOilResourceIcon:      {Path: "resource_oil.png"},
		ImageCombinedResourceIcon: {Path: "resource_combined.png"},
		ImageResourceRow:          {Path: "resource_row.png"},
		ImagePopupBuildTank:       {Path: "popup_build_tank.png"},
		ImageMenuButton:           {Path: "menu_button.png"},
		ImageMenuSelectButton:     {Path: "menu_select_button.png"},
		ImageMenuCheckboxButton:   {Path: "menu_checkbox_button.png"},
		ImageMenuSlideLeft:        {Path: "menu_slide_left.png"},
		ImageMenuBackground:       {Path: "menu_bg.png"},
	}
	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
		ctx.Loader.PreloadImage(id)
	}

	// Associate font resources.
	fontResources := map[resource.FontID]resource.Font{
		FontSmall:       {Path: "DejavuSansMono.ttf", Size: 12},
		FontDescription: {Path: "DejavuSansMono.ttf", Size: 14, LineSpacing: 1.15},
		FontBig:         {Path: "DejavuSansMono.ttf", Size: 20},
	}
	for id, res := range fontResources {
		ctx.Loader.FontRegistry.Set(id, res)
		ctx.Loader.PreloadFont(id)
	}

	// Associate other resources.
	rawResources := map[resource.RawID]resource.Raw{
		RawTilesJSON: {Path: "tiles.json"},
	}
	for id, res := range rawResources {
		ctx.Loader.RawRegistry.Set(id, res)
		ctx.Loader.PreloadRaw(id)
	}

	ctx.CurrentScene = ctx.NewRootScene("menu", newMenuController(state))

	if err := ge.RunGame(ctx); err != nil {
		panic(err)
	}
}

type gameState struct {
	MenuInput       *input.MultiHandler
	Player1keyboard *input.Handler
	Player1gamepad  *input.Handler
	Player2gamepad  *input.Handler
	Player3gamepad  *input.Handler
	Player4gamepad  *input.Handler
}
