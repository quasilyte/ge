package ge

import (
	"image"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

type unsafeImage struct {
	addr *unsafeImage

	image *byte

	bounds   image.Rectangle
	original *unsafeImage

	setVerticesCache map[[2]int][4]byte
}

func toUnsafeImage(img *ebiten.Image) *unsafeImage {
	return (*unsafeImage)(unsafe.Pointer(img))
}

func toEbitenImage(img *unsafeImage) *ebiten.Image {
	return (*ebiten.Image)(unsafe.Pointer(img))
}
