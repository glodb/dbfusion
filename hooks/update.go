package hooks

type PreUpdate interface {
	PreUpdate() PreUpdate
}

type PostUpdate interface {
	PostUpdate() PostUpdate
}
