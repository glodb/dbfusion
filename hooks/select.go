package hooks

type PreSelect interface {
	PreSelect()
}

type PostSelect interface {
	PostSelect()
}
