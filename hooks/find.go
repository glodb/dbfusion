package hooks

type PreFind interface {
	PreFind() PreFind
}

type PostFind interface {
	PostFind() PostFind
}
