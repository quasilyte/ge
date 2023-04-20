package ge

import (
	"image"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

type unsafeImage struct {
	addr *unsafeImage

	image *byte

	original *unsafeImage
	bounds   image.Rectangle

	tmpVertices []float32
	tmpUniforms []uint32
}

func toUnsafeImage(img *ebiten.Image) *unsafeImage {
	return (*unsafeImage)(unsafe.Pointer(img))
}

func toEbitenImage(img *unsafeImage) *ebiten.Image {
	return (*ebiten.Image)(unsafe.Pointer(img))
}
