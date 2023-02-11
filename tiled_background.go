package ge

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

type TiledBackground struct {
	Pos Pos

	Visible bool

	ColorScale ColorScale

	Hue gmath.Rad

	disposed bool

	combined *ebiten.Image
}

type tileInfo struct {
	offset gmath.Vec
	frame  int
}

func NewTiledBackground() *TiledBackground {
	return &TiledBackground{
		Visible:    true,
		ColorScale: defaultColorScale,
	}
}

func (bg *TiledBackground) LoadTileset(ctx *Context, width, height float64, source resource.ImageID, tileset resource.RawID) {
	ts, err := tiled.UnmarshalTileset(ctx.Loader.LoadRaw(tileset).Data)
	if err != nil {
		panic(err)
	}

	spriteSheet := ctx.Loader.LoadImage(source)
	frames := make([]*ebiten.Image, 0, ts.NumTiles)
	for i := 0; i < ts.NumTiles; i++ {
		x := i * int(ts.TileWidth)
		frameRect := image.Rect(x, 0, x+int(ts.TileWidth), int(ts.TileHeight))
		frameImage := spriteSheet.Data.SubImage(frameRect).(*ebiten.Image)
		frames = append(frames, frameImage)
	}

	framePicker := gmath.NewRandPicker[int](&ctx.Rand)
	for i := 0; i < ts.NumTiles; i++ {
		framePicker.AddOption(i, *ts.Tiles[i].Probability)
	}

	combined := ebiten.NewImage(int(width), int(height))
	var op ebiten.DrawImageOptions
	applyColorScale(bg.ColorScale, &op.ColorM)
	if bg.Hue != 0 {
		op.ColorM.RotateHue(float64(bg.Hue))
	}
	for y := float64(0); y < height; y += ts.TileHeight {
		for x := float64(0); x < width; x += ts.TileWidth {
			offset := gmath.Vec{X: x, Y: y}
			frameIndex := framePicker.Pick()
			img := frames[frameIndex]
			op.GeoM.Reset()
			op.GeoM.Translate(offset.X, offset.Y)
			combined.DrawImage(img, &op)
		}
	}
	bg.combined = combined
}

func (bg *TiledBackground) IsDisposed() bool {
	return bg.disposed
}

func (bg *TiledBackground) Dispose() {
	bg.disposed = true
}

func (bg *TiledBackground) Draw(screen *ebiten.Image) {
	if !bg.Visible {
		return
	}

	var op ebiten.DrawImageOptions
	screen.DrawImage(bg.combined, &op)
}
