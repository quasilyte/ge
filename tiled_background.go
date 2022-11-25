package ge

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

type TiledBackground struct {
	Pos Pos

	Visible bool

	ColorScale ColorScale

	Hue gmath.Rad

	disposed bool

	tiles  []tileInfo
	frames []*ebiten.Image
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
	ts, err := tiled.UnmarshalTileset(ctx.Loader.LoadRaw(tileset))
	if err != nil {
		panic(err)
	}

	spriteSheet := ctx.Loader.LoadImage(source)
	bg.frames = bg.frames[:0]
	for i := 0; i < ts.NumTiles; i++ {
		x := i * int(ts.TileWidth)
		frameRect := image.Rect(x, 0, x+int(ts.TileWidth), int(ts.TileHeight))
		frameImage := spriteSheet.Data.SubImage(frameRect).(*ebiten.Image)
		bg.frames = append(bg.frames, frameImage)
	}

	framePicker := gmath.NewRandPicker[int](&ctx.Rand)
	for i := 0; i < ts.NumTiles; i++ {
		framePicker.AddOption(i, *ts.Tiles[i].Probability)
	}

	bg.tiles = bg.tiles[:0]
	for y := float64(0); y < height; y += ts.TileHeight {
		for x := float64(0); x < width; x += ts.TileWidth {
			offset := gmath.Vec{X: x, Y: y}
			frame := framePicker.Pick()
			bg.tiles = append(bg.tiles, tileInfo{
				offset: offset,
				frame:  frame,
			})
		}
	}
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
	applyColorScale(bg.ColorScale, &op)
	if bg.Hue != 0 {
		op.ColorM.RotateHue(float64(bg.Hue))
	}
	for _, t := range bg.tiles {
		img := bg.frames[t.frame]
		op.GeoM.Reset()
		op.GeoM.Translate(t.offset.X, t.offset.Y)
		screen.DrawImage(img, &op)
	}
}
