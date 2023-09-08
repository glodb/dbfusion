package hooks

type PreUpdate interface {
	PreUpdate()
}

type PostUpdate interface {
	PostUpdate()
}
