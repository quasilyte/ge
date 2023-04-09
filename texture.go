package ge

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Texture struct {
	image  *ebiten.Image
	width  float64
	height float64
}

func NewHorizontallyRepeatedTexture(img resource.Image, maxLen float64) *Texture {
	tex := &Texture{}

	w, h := img.Data.Size()
	tex.width = math.Round(maxLen)
	tex.height = float64(h)

	repeatedImage := ebiten.NewImage(int(tex.width), h)
	x := 0.0
	var drawOptions ebiten.DrawImageOptions
	for x < maxLen {
		segmentWidth := float64(w)
		srcImage := img.Data
		if x+segmentWidth > maxLen {
			segmentWidth = maxLen - x
			srcImage = img.Data.SubImage(image.Rectangle{
				Max: image.Point{X: int(math.Round(segmentWidth)), Y: h},
			}).(*ebiten.Image)
		}
		repeatedImage.DrawImage(srcImage, &drawOptions)
		drawOptions.GeoM.Translate(segmentWidth, 0)
		x += segmentWidth
	}

	tex.image = repeatedImage
	return tex
}
