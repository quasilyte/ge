package input

type Action uint32

type Keymap struct {
	mapping map[Action]Key
}

func (m *Keymap) Set(action Action, key Key) {
	if m.mapping == nil {
		m.mapping = make(map[Action]Key)
	}
	m.mapping[action] = key
}
