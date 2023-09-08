package hooks

type PreDelete interface {
	PreDelete()
}

type PostDelete interface {
	PostDelete()
}
