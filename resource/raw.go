package resource

type Raw struct {
	Path string
}

type RawID int

type RawRegistry struct {
	mapping map[RawID]Raw
}

func (r *RawRegistry) Set(id RawID, info Raw) {
	r.mapping[id] = info
}
