package resource

type ShaderInfo struct {
	Path string
}

type ShaderID int

type ShaderRegistry struct {
	mapping map[ShaderID]ShaderInfo
}

func (r *ShaderRegistry) Set(id ShaderID, info ShaderInfo) {
	r.mapping[id] = info
}
