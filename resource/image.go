package resource

type Image struct {
	Path string
}

type ImageRegistry struct {
	mapping map[ID]Image
}

func (r *ImageRegistry) Set(id ID, info Image) {
	r.mapping[id] = info
}
