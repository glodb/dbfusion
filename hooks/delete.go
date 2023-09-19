package hooks

type PreDelete interface {
	PreDelete() PreDelete
}

type PostDelete interface {
	PostDelete() PostDelete
}
