package ge

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type imageCacheKey struct {
	imageWidth  int
	imageHeight int
}

type imageCache struct {
	m map[imageCacheKey]*ebiten.Image

	tmp unsafeImage
}

func (cache *imageCache) UnsafeSubImage(i *ebiten.Image, bounds image.Rectangle) *ebiten.Image {
	// Basically, we're doing this:
	// > subImage = i.SubImage(bounds)
	// But without redundant allocation.
	unsafeImg := toUnsafeImage(i)
	unsafeSubImage := cache.UnsafeImageForSubImage()
	unsafeSubImage.original = unsafeImg
	unsafeSubImage.bounds = bounds
	unsafeSubImage.image = unsafeImg.image
	return toEbitenImage(unsafeSubImage)
}

func (cache *imageCache) UnsafeImageForSubImage() *unsafeImage {
	return &cache.tmp
}

func (cache *imageCache) Init() {
	cache.m = make(map[imageCacheKey]*ebiten.Image)

	cache.tmp.addr = &cache.tmp
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
