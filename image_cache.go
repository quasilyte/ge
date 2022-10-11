package ge

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type imageCacheKey struct {
	imageWidth  int
	imageHeight int
}

type imageCache struct {
	m map[imageCacheKey]*ebiten.Image
}

func (cache *imageCache) Init() {
	cache.m = make(map[imageCacheKey]*ebiten.Image)
}

func (cache *imageCache) NewTempImage(width, height int) *ebiten.Image {
	key := imageCacheKey{imageWidth: width, imageHeight: height}
	if cachedImage, ok := cache.m[key]; ok {
		cachedImage.Clear()
		return cachedImage
	}
	img := ebiten.NewImage(width, height)
	cache.m[key] = img
	return img
}
