package resource

import "github.com/hajimehoshi/ebiten/v2"

type Image struct {
	ID   ImageID
	Data *ebiten.Image

	DefaultFrameWidth  float64
	DefaultFrameHeight float64
}

type ImageInfo struct {
	Path string

	FrameWidth  float64
	FrameHeight float64
}

type ImageID int

type ImageRegistry struct {
	mapping map[ImageID]ImageInfo
}

func (r *ImageRegistry) Set(id ImageID, info ImageInfo) {
	r.mapping[id] = info
}
