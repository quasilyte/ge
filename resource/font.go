package resource

type Font struct {
	Path string

	Size int
}

type FontRegistry struct {
	mapping map[ID]Font
}

func (r *FontRegistry) Set(id ID, info Font) {
	r.mapping[id] = info
}
