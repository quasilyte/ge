package main

import (
	"embed"
	"io"
	"time"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"

	_ "image/png"
)

//go:embed all:_assets
var gameAssets embed.FS

const buildVersion = 1

const (
	ImageNone resource.ImageID = iota
	ImageInterfacePanel
	ImageGradientGrid
	ImageTiles
	ImageTurret1
	ImageTurret2
	ImageFort
	ImageBase
	ImageWall
	ImageFactory
	ImageBonusGenerator
	ImagePlayerVessel
	ImageVessel1
	ImageVessel2
	ImageVessel3
	ImageVessel4
	ImageVessel5
	ImageVessel6
	ImageBullet1
	ImageBullet2
	ImageBullet3
	ImageBullet4
	ImageBullet5
	ImageBullet6
	ImageBullet7
	ImageRailgunRay
	ImageMissile
	ImagePickupHP
	ImagePickupAmmo
	ImageMine
	ImageFlamethrower
	ImageRegenerator
	ImageShield
	ImageExplosion1
	ImageExplosion2
	ImageExplosion3
	ImageRegeneratorExplosion
	ImageDamageMask1
	ImageDamageMask2
	ImageDamageMask3
	ImageDamageMask4
	ImageDarkTile
	ImageUIArrow
)

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
	ActionFire
	ActionSpecial
	ActionPressButton
	ActionConfirm
	ActionEscape
	ActionNextItem
	ActionPrevItem
)

const (
	AudioNone resource.AudioID = iota
	AudioLaser1
	AudioLaser2
	AudioLaser3
	AudioLaser4
	AudioLaser5
	AudioRailgun
	AudioHowitz
	AudioHurricane
	AudioMissile
	AudioFlamethrower
	AudioRegenerator
	AudioMine
	AudioShieldActivated
	AudioShieldAbsorb
	AudioExplosion1
	AudioExplosion2
	AudioExplosion3
	AudioExplosion4
	AudioExplosion5
	AudioExplosion6
	AudioExplosion7
	AudioRegeneratorExplosion
	AudioBonus
	AudioMusic
	AudioScreenReset
)

const (
	RawTilesJSON resource.RawID = iota
	RawMapDateTilesetJSON
	RawLevel0JSON
	RawLevel1JSON
	RawLevel2JSON
	RawLevel3JSON
	RawLevel4JSON
	RawLevel5JSON
	RawLevel6JSON
	RawLevel7JSON
	RawLevel8JSON
	RawLevel9JSON
	RawLevel10JSON
	RawLevel11JSON
	RawLevel12JSON
	RawLevel13JSON
	RawLevel14JSON
)

const lastLevel = 14

const (
	ShaderBuildingDamage resource.ShaderID = iota
)

const (
	FontSmall resource.FontID = iota
	FontMedium
	FontBig
)

func gameMain() {
	ctx := ge.NewContext(ge.ContextConfig{})
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "retrowave_city"
	ctx.WindowTitle = "Retrowave City"
	ctx.WindowWidth = 1920
	ctx.WindowHeight = 1080
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = func(path string) io.ReadCloser {
		f, err := gameAssets.Open("_assets/" + path)
		if err != nil {
			ctx.OnCriticalError(err)
		}
		return f
	}

	// Associate image resources.
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageGradientGrid:         {Path: "grid.png"},
		ImageInterfacePanel:       {Path: "ui_background.png"},
		ImageTiles:                {Path: "tiles.png"},
		ImageTurret1:              {Path: "turret1.png"},
		ImageTurret2:              {Path: "turret2.png"},
		ImageBase:                 {Path: "base.png"},
		ImageFactory:              {Path: "factory.png"},
		ImageBonusGenerator:       {Path: "bonus_generator.png"},
		ImageFort:                 {Path: "fort.png"},
		ImageWall:                 {Path: "wall.png", FrameWidth: 64},
		ImagePlayerVessel:         {Path: "player_vessel.png", FrameWidth: 64},
		ImageVessel1:              {Path: "vessel1.png", FrameWidth: 64},
		ImageVessel2:              {Path: "vessel2.png", FrameWidth: 64},
		ImageVessel3:              {Path: "vessel3.png", FrameWidth: 64},
		ImageVessel4:              {Path: "vessel4.png", FrameWidth: 64},
		ImageVessel5:              {Path: "vessel5.png", FrameWidth: 64},
		ImageVessel6:              {Path: "vessel6.png", FrameWidth: 64},
		ImageBullet1:              {Path: "bullet1.png"},
		ImageBullet2:              {Path: "bullet2.png"},
		ImageBullet3:              {Path: "bullet3.png"},
		ImageBullet4:              {Path: "bullet4.png"},
		ImageBullet5:              {Path: "bullet5.png"},
		ImageBullet6:              {Path: "bullet6.png"},
		ImageBullet7:              {Path: "bullet7.png"},
		ImageRailgunRay:           {Path: "railgun_ray.png"},
		ImageMissile:              {Path: "missile.png"},
		ImagePickupHP:             {Path: "pickup_hp.png"},
		ImagePickupAmmo:           {Path: "pickup_ammo.png"},
		ImageMine:                 {Path: "mine.png"},
		ImageFlamethrower:         {Path: "flamethrower.png"},
		ImageRegenerator:          {Path: "regenerator.png"},
		ImageShield:               {Path: "shield.png"},
		ImageExplosion1:           {Path: "explosion1.png", FrameWidth: 64},
		ImageExplosion2:           {Path: "explosion2.png", FrameWidth: 64},
		ImageExplosion3:           {Path: "explosion3.png", FrameWidth: 64},
		ImageRegeneratorExplosion: {Path: "regenerator_explosion.png", FrameWidth: 64},
		ImageDamageMask1:          {Path: "damage_mask1.png"},
		ImageDamageMask2:          {Path: "damage_mask2.png"},
		ImageDamageMask3:          {Path: "damage_mask3.png"},
		ImageDamageMask4:          {Path: "damage_mask4.png"},
		ImageDarkTile:             {Path: "dark_tile.png"},
		ImageUIArrow:              {Path: "ui_arrow.png"},
	}
	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
	}

	// Associate other resources.
	rawResources := map[resource.RawID]resource.RawInfo{
		RawTilesJSON:          {Path: "tiles.json"},
		RawMapDateTilesetJSON: {Path: "tankarcade.tsj"},
		RawLevel0JSON:         {Path: "maps/level0.json"},
		RawLevel1JSON:         {Path: "maps/level1.json"},
		RawLevel2JSON:         {Path: "maps/level2.json"},
		RawLevel3JSON:         {Path: "maps/level3.json"},
		RawLevel4JSON:         {Path: "maps/level4.json"},
		RawLevel5JSON:         {Path: "maps/level5.json"},
		RawLevel6JSON:         {Path: "maps/level6.json"},
		RawLevel7JSON:         {Path: "maps/level7.json"},
		RawLevel8JSON:         {Path: "maps/level8.json"},
		RawLevel9JSON:         {Path: "maps/level9.json"},
		RawLevel10JSON:        {Path: "maps/level10.json"},
		RawLevel11JSON:        {Path: "maps/level11.json"},
		RawLevel12JSON:        {Path: "maps/level12.json"},
		RawLevel13JSON:        {Path: "maps/level13.json"},
		RawLevel14JSON:        {Path: "maps/level14.json"},
	}
	for id, res := range rawResources {
		ctx.Loader.RawRegistry.Set(id, res)
		ctx.Loader.LoadRaw(id)
	}

	// Associate shader resources.
	shaderResources := map[resource.ShaderID]resource.ShaderInfo{
		ShaderBuildingDamage: {Path: "shaders/building_damage.go"},
	}
	for id, res := range shaderResources {
		ctx.Loader.ShaderRegistry.Set(id, res)
		ctx.Loader.LoadShader(id)
	}

	// Associate audio resources.
	audioResources := map[resource.AudioID]resource.AudioInfo{
		AudioLaser1:               {Path: "laser1.wav", Volume: -0.6},
		AudioLaser2:               {Path: "laser2.wav", Volume: -0.7},
		AudioLaser3:               {Path: "laser3.wav", Volume: -0.65},
		AudioLaser4:               {Path: "laser4.wav", Volume: -0.7},
		AudioLaser5:               {Path: "laser5.wav", Volume: -0.5},
		AudioRailgun:              {Path: "railgun.wav", Volume: -0.6},
		AudioHowitz:               {Path: "howitz.wav", Volume: -0.65},
		AudioHurricane:            {Path: "hurricane.wav", Volume: -0.6},
		AudioMissile:              {Path: "missile.wav", Volume: -0.6},
		AudioFlamethrower:         {Path: "flamethrower.wav", Volume: -0.6},
		AudioRegenerator:          {Path: "regenerator.wav"},
		AudioMine:                 {Path: "laymine.wav", Volume: -0.2},
		AudioShieldActivated:      {Path: "shield_activated.wav", Volume: -0.4},
		AudioShieldAbsorb:         {Path: "shield_absorb.wav", Volume: -0.3},
		AudioExplosion1:           {Path: "explosion1.wav", Volume: -0.8},
		AudioExplosion2:           {Path: "explosion2.wav", Volume: -0.75},
		AudioExplosion3:           {Path: "explosion3.wav", Volume: -0.75},
		AudioExplosion4:           {Path: "explosion4.wav", Volume: -0.7},
		AudioExplosion5:           {Path: "explosion5.wav", Volume: -0.7},
		AudioExplosion6:           {Path: "explosion6.wav", Volume: -0.75},
		AudioExplosion7:           {Path: "explosion7.wav", Volume: -0.75},
		AudioRegeneratorExplosion: {Path: "regenerator_explosion.wav", Volume: -0.7},
		AudioBonus:                {Path: "bonus.wav", Volume: -0.5},
		AudioMusic:                {Path: "music.ogg", Volume: -0.7},
		AudioScreenReset:          {Path: "screen_reset.wav", Volume: -0.5},
	}
	for id, res := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, res)
		ctx.Loader.LoadAudio(id)
	}

	// Associate font resources.
	fontResources := map[resource.FontID]resource.FontInfo{
		FontSmall:  {Path: "font.otf", Size: 20},
		FontMedium: {Path: "font.otf", Size: 26},
		FontBig:    {Path: "font.otf", Size: 48},
	}
	for id, res := range fontResources {
		ctx.Loader.FontRegistry.Set(id, res)
		ctx.Loader.LoadFont(id)
	}

	state := newGameState()

	// Bind controls.
	keymap := input.Keymap{
		ActionMoveLeft:    {input.KeyGamepadLeft, input.KeyGamepadLStickLeft, input.KeyLeft},
		ActionMoveRight:   {input.KeyGamepadRight, input.KeyGamepadLStickRight, input.KeyRight},
		ActionMoveDown:    {input.KeyGamepadDown, input.KeyGamepadLStickDown, input.KeyDown},
		ActionMoveUp:      {input.KeyGamepadUp, input.KeyGamepadLStickUp, input.KeyUp},
		ActionFire:        {input.KeyGamepadA, input.KeyQ},
		ActionSpecial:     {input.KeyGamepadB, input.KeyW},
		ActionPressButton: {input.KeyGamepadA, input.KeyEnter, input.KeyMouseLeft},
		ActionConfirm:     {input.KeyGamepadA, input.KeyEnter},
		ActionEscape:      {input.KeyGamepadStart, input.KeyEscape},
	}
	state.PlayerInput[0] = ctx.Input.NewHandler(0, keymap)
	state.PlayerInput[1] = ctx.Input.NewHandler(1, keymap)
	state.PlayerInput[2] = ctx.Input.NewHandler(2, keymap)

	if err := ge.RunGame(ctx, newMainMenuController(state)); err != nil {
		panic(err)
	}
}

func main() {
	gameMain()
}
