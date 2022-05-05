package resource

type Image struct {
	Path string
}

type ImageID int

type ImageRegistry struct {
	mapping map[ImageID]Image
}

func (r *ImageRegistry) Set(id ImageID, info Image) {
	r.mapping[id] = info
}
